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

func AddStringFunctions(table SymbolTable) SymbolTable {
	table.Add("len", func(value string) int64 { return int64(len(value)) })
	table.Add("lower", strings.ToLower)
	table.Add("join", func (separator string, tuple Tuple) string { return strings.Join(ValuesToStrings(tuple.List), separator) })
	table.Add("concat", func (aa string, bb string) string { return aa + bb })
	table.Add("upper", strings.ToUpper)

	return table
}

// These functions are harmless, one can execute them from a script and not cause any damage.
func AddBooleanAndArithmeticFunctions(table SymbolTable) SymbolTable {

	AddStringFunctions(table)
	
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

	return table
}


// These functions are harmless, one can execute them from a script and not cause any damage.
func NewSafeSymbolTable(notFound CallHandler) SymbolTable {
	table := NewSymbolTable(notFound)
	AddBooleanAndArithmeticFunctions(table)
	return table
}

/////////////////////////////////////////////////////////////////////////////
// Declare variables and functions
/////////////////////////////////////////////////////////////////////////////

func AddDeclareFunctions(table SymbolTable) {

	table.Add("get", func(atom Atom) Value {
		return table.Call(atom, []Value{})
	})
	table.Add("set", func(atom Atom, value Value) Value {
		table.Add(atom.Name, func () Value {
			return value
		})
		return value
	})
	table.Add("func", func(atom Atom, arg Atom, code Value) Value {
		table.Add(atom.Name, func (argValue Value) Value {
			newScope := NewSymbolTable(&table)
			newScope.Add(arg.Name, func () Value {
				return argValue
			})
			return newScope.Eval(code)
		})
		return atom
	})
	table.Add("if", func(condition Value, trueCode Value, falseCode Value) Value {
		if toBool(table.Eval(condition)) {
			return table.Eval(trueCode)
		} else {
			return table.Eval(falseCode)
		}
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
func AddOperatinSystemFunctions(table SymbolTable) {
	//
	//  Add shell specific commands
	//  These are typically not a 'safe' in that they allow access to the file system
	// 
	table.Add("exec", executeProcess0)
	table.Add("spawn", spawnProcess)
	table.Add("pipe", Pipe)

	table.Add("echo", func (values... Value) bool {
		for k,_:= range values {
			fmt.Print(toString(values[k]))
		}
		return true
	})
}

type ErrorIfFunctionNotFound struct {}

func (function * ErrorIfFunctionNotFound) Find(name Atom, args [] Value) reflect.Value {
	return reflect.ValueOf(func(args... Value) bool {
		fmt.Printf("ERROR: function not found: '%s' %s\n", name.Name, args)  // TODO ought to use context logger
		return false
	})
}

type ExecIfNotFound struct {}

func (exec * ExecIfNotFound) Find (name Atom, args [] Value) reflect.Value {

	return reflect.ValueOf(func(args... Value) bool {
		return executeProcess(name.Name, ValuesToStrings(args)...)
	})
}

// These functions are potentially not harmless
func NewUnSafeSymbolTable() SymbolTable {
	table := SymbolTable{map[string]reflect.Value{},&ExecIfNotFound{}}
	AddBooleanAndArithmeticFunctions(table)
	AddDeclareFunctions(table)
	AddOperatinSystemFunctions(table)
	return table
}
