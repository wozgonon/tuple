/*
    This file is part of WOZG.

    WOZG is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    WOZG is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with WOZG.  If not, see <https://www.gnu.org/licenses/>.
*/
package tuple

import "math"
import "strings"
import "reflect"
import "fmt"
import "log"
import "os/exec"
import "os"
import "bytes"

func AddStringFunctions(table * SymbolTable) {
	table.Add("len", func(value string) int64 { return int64(len(value)) })
	table.Add("lower", strings.ToLower)
	table.Add("join", func (context EvalContext, separator string, tuple Tuple) string { return strings.Join(EvalToStrings(context, tuple.List), separator) })
	table.Add("concat", func (context EvalContext, values... Value) string  { return strings.Join(EvalToStrings(context, values), "") })
	table.Add("upper", strings.ToUpper)
}

func AddTupleFunctions(table * SymbolTable)  {

	table.Add("at", func(index int64, tuple Tuple) Value {
		if index < 0 || index >= int64(tuple.Length()) {
			return EMPTY
		}
		return tuple.List[index]
	})

	table.Add("list", func(_ EvalContext, values... Value) Value { return NewTuple(values...) })
	// TODO table.Add("quote", func(value Value) Value { return NewTuple("quote", value) })
}

// These functions are harmless, one can execute them from a script and not cause any damage.
func AddBooleanAndArithmeticFunctions(table * SymbolTable) {


	table.Add("exp", math.Exp)
	table.Add("log", math.Log)
	table.Add("sin", math.Sin)
	table.Add("cos", math.Cos)
	table.Add("tan", math.Tan)
	table.Add("acos", math.Acos)
	table.Add("asin", math.Asin)
	table.Add("atan", math.Atan)
	table.Add("atan2", math.Atan2)
	table.Add("round", math.Round)
	table.Add("round2", func(value float64) float64 { return math.Round(value*100)/100 }) // Not needed
	table.Add("**", math.Pow)
	table.Add("+", func (aa float64) float64 { return aa })
	table.Add("-", func (aa float64) float64 { return -aa })
	table.Add("+", func (aa float64, bb float64) float64 { return aa+bb })
	table.Add("-", func (aa float64, bb float64) float64 { return aa-bb })
	table.Add("*", func (aa float64, bb float64) float64 { return aa*bb })
	table.Add("/", func (aa float64, bb float64) float64 { return aa/bb })
	table.Add("eq", func (aa string, bb string) bool { return aa==bb })  // Should take Value as argument
	table.Add("==", func (aa float64, bb float64) bool { return aa==bb })  // Should take Value as argument
	table.Add("!=", func (aa float64, bb float64) bool { return aa!=bb })
	table.Add(">=", func (aa float64, bb float64) bool { return aa>=bb })
	table.Add("<=", func (aa float64, bb float64) bool { return aa<=bb })
	table.Add(">", func (aa float64, bb float64) bool { return aa>bb })
	table.Add("<", func (aa float64, bb float64) bool { return aa<bb })
	table.Add("&&", func (aa bool, bb bool) bool { return aa&&bb })
	table.Add("||", func (aa bool, bb bool) bool { return aa||bb })
	table.Add("!", func (aa bool) bool { return ! aa })
	table.Add("PI", func () float64 { return math.Pi })
	table.Add("PHI", func () float64 { return math.Phi })
	table.Add("E", func () float64 { return math.E })
	table.Add("true", func () bool { return true })
	table.Add("false", func () bool { return false })

}


// These functions are harmless, one can execute them from a script and not cause any damage.
func NewHarmlessSymbolTable(notFound CallHandler) SymbolTable {
	table := NewSymbolTable(notFound)
	AddBooleanAndArithmeticFunctions(&table)

	return table
}

func NewSafeSymbolTable(notFound CallHandler) SymbolTable {
	table := NewSymbolTable(notFound)

	AddBooleanAndArithmeticFunctions(&table)

	// These allocate memory and are not quite
	AddStringFunctions(&table)
	AddTupleFunctions(&table)
	return table
}

/////////////////////////////////////////////////////////////////////////////
// Declare variables and functions
/////////////////////////////////////////////////////////////////////////////


func Assign (context EvalContext, atom Atom, value Value) Value {
	context.Add(atom.Name, func () Value {
		return value
	})
	return value
}

func AddDeclareFunctions(table * SymbolTable) {

	table.Add("get", func(context EvalContext, atom Atom) Value {
		return context.Call(atom, []Value{})
	})
	table.Add("set", Assign)
	table.Add("func", func(context EvalContext, atom Atom, arg Atom, code Value) Value {
		context.Add(atom.Name, func (context1 EvalContext, argValue Value) Value {
			evaluated := Eval(context1, argValue)
			context.Log("TRACE", "** FUNC %s argName=%s argValue: %s evaluated:%s", atom.Name, arg.Name, argValue, evaluated)
			newScope := NewSymbolTable(context1)
			newScope.Add(arg.Name, func () Value {
				return evaluated
			})
			return Eval(&newScope, code)
		})
		return atom
	})
	table.Add("if", func(context EvalContext, condition Value, trueCode Value, falseCode Value) Value {
		if toBool(Eval(context, condition)) {
			return Eval(context, trueCode)
		} else {
			return Eval(context, falseCode)
		}
	})
	table.Add("for", func(context EvalContext, atom Atom, list Tuple, code Value) Tuple {
		var iterator Value = nil
		newScope := NewSymbolTable(context)
		newScope.Add(atom.Name, func () Value {
			return iterator
		})
		result := NewTuple()
		for _, v := range list.List {
			iterator = Eval(context, v)
			value := Eval(&newScope, code)
			result.Append(value)
		}
		return result
	})
}

/////////////////////////////////////////////////////////////////////////////
// Operating system access
/////////////////////////////////////////////////////////////////////////////

//
//  Add shell specific commands
//  These are typically not a 'safe' in that they allow access to the file system
// 
func Pipe(writer string, reader string) bool {

	log.Printf("pipe '%s', '%s'", writer, reader)
	cmdw := exec.Command(writer)
	stdoutw, err := cmdw.StdoutPipe()

	err = cmdw.Start()
	if err != nil {
		log.Printf("Start 1 finished with error: %v", err)
		return false
	}

	cmdr := exec.Command(reader)
	cmdr.Stdin = stdoutw
	cmdr.Stdout = os.Stdout

	if err != nil {
		log.Printf("StdoutPipe finished with error: %v", err)
		log.Fatal(err)
	}
	err = cmdr.Start()
	if err != nil {
		log.Printf("Start 2 finished with error: %v", err)
		log.Fatal(err)
	}
	err = cmdw.Wait()
	if err != nil {
		log.Printf("Wait 1 finished with error: %v", err)
		return false
	}
	err = cmdr.Wait()
	if ; err != nil {
		log.Printf("Wait 2 finished with error: %v", err)
		log.Fatal(err)
	}

	return true
}

func executeProcess(arg string, args... string) bool {
	cmd := exec.Command(arg, args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
		return false
	}
	outStr := string(stdout.Bytes())
	fmt.Print(outStr)
	return true
}

// TODO sort out variadic arguments
func executeProcess0(arg string) bool {
	return executeProcess(arg)
}


func spawnProcess (arg string) bool {
	cmd := exec.Command(arg)
	//var stdout bytes.Buffer
	cmd.Stdout = os.Stdout //&stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
		return false
	}
	//outStr := string(stdout.Bytes())
	//fmt.Print(outStr)
	return true
}


// These functions access the operating system and so an actor
// could use them to harm the computer.
func AddOperatingSystemFunctions(table * SymbolTable) {
	//
	//  Add shell specific commands
	//  These are typically not a 'safe' in that they allow access to the file system
	// 
	table.Add("exec", executeProcess0)
	table.Add("spawn", spawnProcess)
	table.Add("pipe", Pipe)

	table.Add("echo", func (context EvalContext, values... Value) bool {
		for k,_:= range values {
			evaluated := Eval(context, values[k])  // TODO This should use table from context so that it uses scope
			fmt.Print(toString(context, evaluated))
		}
		return true
	})
	table.Add("eval", func (context EvalContext, value Value) Value {
		return Eval(context, value)
	})
	//table.Add("expr", func (context EvalContext, value Value) Value {
	//	return ParseAndEval(context, value)
	//})
}

type ErrorIfFunctionNotFound struct {}

func (function * ErrorIfFunctionNotFound) Find(context EvalContext, name Atom, args [] Value) reflect.Value {
	return reflect.ValueOf(func(args... Value) bool {
		fmt.Printf("ERROR: function not found: '%s' %s\n", name.Name, args)  // TODO ought to use context logger
		return false
	})
}

type ExecIfNotFound struct {}

func (exec * ExecIfNotFound) Find (context EvalContext, name Atom, args [] Value) reflect.Value {

	return reflect.ValueOf(func(context EvalContext, args... Value) bool {
		return executeProcess(name.Name, EvalToStrings(context, args)...)
	})
}

// These functions are potentially not harmless
func NewUnSafeSymbolTable() SymbolTable {
	table := SymbolTable{map[string]reflect.Value{},&ExecIfNotFound{}}
	AddBooleanAndArithmeticFunctions(&table)
	AddStringFunctions(&table)
	AddTupleFunctions(&table)
	AddDeclareFunctions(&table)
	AddOperatingSystemFunctions(&table)
	return table
}
