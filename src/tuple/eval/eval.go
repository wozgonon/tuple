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
package eval

import "reflect"
import "fmt"
import "tuple"
import "errors"
import "strings"

type Tag = tuple.Tag
type Value = tuple.Value
type Tuple = tuple.Tuple
type Int64 = tuple.Int64
type Float64 = tuple.Float64
type Bool = tuple.Bool
type String = tuple.String
type Logger = tuple.Logger
type LocationLogger = tuple.LocationLogger
type Location = tuple.Location
var Error = tuple.Error
var Verbose = tuple.Verbose
var Trace = tuple.Trace


//  A simple toy evaluator.
//
//  Treats the data produced by parsing and treats it as an executable program
//  (code can be accessed and treated as if it is data and data can be treated as code).
//  This is normal for LISP but will work for any of the input grammars.
//
//  See:
//  * [Homoiconic](https://en.wikipedia.org/wiki/Homoiconicity) language treats "code as data".
//  * [Eval](https://en.wikipedia.org/wiki/Eval)
//  * [Meta-circular_evaluator](https://en.wikipedia.org/wiki/Meta-circular_evaluator)
//  

type Global interface {
	Logger() LocationLogger
	Find(context EvalContext, name Tag, args [] Value) (*SymbolTable, reflect.Value)

	// The root is analageous to the root of a data hierarchy such as a directory, file system or registry.
	// Rather than provide individual functions to return contextual information it is much more flexible
	// to provide a searchable directory structure.
	Root() Value
}

type EvalContext interface {
	Global
	Logger
	Add(name string, function interface{})
	Call(head Tag, args []Value) (Value, error)
	AllSymbols() Tuple
}

/////////////////////////////////////////////////////////////////////////////

type SymbolTable struct {
	symbols map[string]reflect.Value
	global Global
}

func NewSymbolTable(notFound Global) SymbolTable {
	if notFound.Logger() == nil {
		panic("nil logger")  // TODO not a good idea to have fatals in the code
	}
	return SymbolTable{map[string]reflect.Value{},notFound}
}

// TODO change to a Value
func (context * SymbolTable) AllSymbols() Tuple {
	tuple := tuple.NewTuple()
	for k,v := range context.symbols {
		key := signatureOfFunction(k, v)
		tuple.Append(String(key))
	}
	return tuple
}

func (context * SymbolTable) Logger() LocationLogger {
	return context.global.Logger()
}

func (context * SymbolTable) Root() Value {
	return context.global.Root()
}

func LocationForValue(value Value) Location {
	// TODO get the location associate with a Value
	return tuple.NewLocation("<eval>", 0, 0, 0) // TODO
}

// TODO use location from values from 
func (context * SymbolTable) Log(level string, format string, args ...interface{}) {
	location := tuple.NewLocation("<eval>", 0, 0, 0) // TODO
	message := fmt.Sprintf(format, args...)
	context.global.Logger()(location, level, message)
}

func (context * SymbolTable) Error(value Value, format string, args ...interface{}) {
	location := LocationForValue(value)
	message := fmt.Sprintf(format, args...)
	context.global.Logger()(location, "ERROR", message)
}

/////////////////////////////////////////////////////////////////////////////

func Eval(context EvalContext, expression Value) (Value, error) {

	switch val := expression.(type) {
	case Tuple:
		ll := val.Arity()
		if ll == 0 {
			return val, nil  // If IsAtom then return val, nul
		}
		head := val.Get(0)
		if ll == 1 {
			return Eval(context, head)
		}
		tag, ok := head.(Tag)
		if ! ok {
			return evalTuple(context, val)
		}
		return context.Call(tag, val.List[1:])
	case Tag:
		return context.Call(val, []Value{})
	default:
		return val, nil
	}
}

func evalTuple(context EvalContext, value Tuple) (Tuple, error) {
	newTuple := tuple.NewTuple()
	for _,v:= range value.List {
		evaluated, err := Eval(context, v)
		if err != nil {
			return tuple.EMPTY, err
		}
		newTuple.Append(evaluated)
	}
	Trace(context, "Eval tuple return '%s'", newTuple)
	return newTuple, nil
}

/////////////////////////////////////////////////////////////////////////////

func signatureOfCall(name string, args []Value) string {

	key := name
	for _,v := range args {
		key = key + " " + reflect.TypeOf(v).Name()
	}
	return key
}

func signatureOfFunction(name string, function reflect.Value) string {
	tt := function.Type()
	key := name
	numIn := tt.NumIn()
	if tt.IsVariadic() {
		return strings.ToLower(fmt.Sprintf("*_%s", name))
	} else {
		for nn := 0;  nn < numIn; nn += 1 {
			//Verbose(context, "name=%s nn=%d numin=%d", name, nn, numIn)
			argName := tt.In(nn).Name()
			if argName != "EvalContext" {
				key = key + " " + argName
			}
		}
	}
	return key
}


func makeKey(name string, arity int, variadic bool) string {
	if variadic {
		return fmt.Sprintf("*_%s", name)
	}
	return fmt.Sprintf("%d_%s", arity, name)
}

func (table * SymbolTable) Count() int {
	return len(table.symbols)
}

func (table * SymbolTable) Add(name string, function interface{}) {
	reflectValue := reflect.ValueOf(function)
	typ := reflectValue.Type()

	nn := typ.NumIn()
	if nn > 0 {
		if typ.In(0) == EvalContextType {
			nn -= 1
		}
	}
	key := signatureOfFunction(name, reflect.ValueOf(function))
	table.symbols[key] = reflectValue
	
	key = makeKey(name, nn, typ.IsVariadic())
	table.symbols[key] = reflectValue
}

func (table * SymbolTable) Find(context EvalContext, head Tag, args []Value) (*SymbolTable, reflect.Value) {  // Reduce

	/*key := signatureOfCall(head.Name, args)
	f, ok := table.symbols[key]
	if ok {
		return table, f
	}

	key = strings.Replace(key, "int64", "float64", 99999)
	f, ok = table.symbols[key]
	if ok {
		return table, f
	}*/
	
	name := head.Name
	nn := len(args)

	for _, variadic := range []bool{ false, true } {
		key := makeKey(name, nn, variadic)
		f, ok := table.symbols[key]
		if ok {
			context.Log ("TRACE", "FIND Found '%s' variadic=%s in this symbol table (%d entries), forwarding",  head, variadic, len(table.symbols))
			return table, f
		}
	}
	//context.Log ("TRACE1", "FIND Could not find '%s' in this symbol table (%d entries), forwarding",  head, len(table.symbols))
	// TODO look up variatic functions
	return table.global.Find(context, head, args)
}



func (context * SymbolTable) Call(head Tag, args []Value) (Value, error) {
	return context.call3(context, head, args)
}

func (table * SymbolTable) call3(context EvalContext, head Tag, args []Value) (Value, error) {  // Reduce

	_, f := table.Find(context, head, args)
	call := NewReflectCall(context, f, len(args))
	for key,v:= range args {
		var result interface{} = v
		if ! f.Type().IsVariadic() {
			_, isTag := v.(Tag)
			expectedType := call.expectedTypeOfArg(key)
			switch  {
			case expectedType == TagType && isTag:
			case expectedType == ValueType:
			default:
				evaluated, err := Eval(context, v)
				if err != nil {
					return tuple.EMPTY, err
				}
				converted, err := Convert(context, evaluated, expectedType)
				if err != nil {
					table.Error(v, "Cannot convert '%s' to '%s'", evaluated, expectedType)
					return tuple.EMPTY, err
				}
				result = converted
			}
		}
		if result == nil {
			return tuple.EMPTY, errors.New("Unexpected nil head=" + head.Name)
		}
		call.setArg(key, result)
	}
	return call.Call(context, head.Name)
}
