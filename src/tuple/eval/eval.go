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

type EvalContext interface {
	Logger
	LocalScope
	GlobalScope() GlobalScope
	NewLocalScope() EvalContext
}

type Finder interface {
	// TODO return a single Callable object not a reflect.Value
	Find(context EvalContext, name Tag, args [] Value) (LocalScope, reflect.Value)
}

type GlobalScope interface {
	Logger
	LocationLogger() LocationLogger
	
	// The root is analageous to the root of a data hierarchy such as a directory, file system or registry.
	// Rather than provide individual functions to return contextual information it is much more flexible
	// to provide a searchable directory structure.
	Root() Value
	AddToRoot(string Tag, value Value)  // Is this needed
}

type LocalScope interface {
	Finder
	Add(name string, function interface{})
}

type Quoted struct {
	value Value
}
func (quoted Quoted) Value() Value { return quoted.value }

func Eval(context EvalContext, expression Value) (Value, error) {

	context.Log("VERBOSE", "eval: '%s'", expression)
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
		return Call(context, tag, val.List[1:])
	case Tag:
		return Call(context, val, []Value{})
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

func Call(context EvalContext, head Tag, args []Value) (Value, error) {  // Reduce

	_, f := context.Find(context, head, args)
	call := NewReflectCall(context, f, len(args))
	for key,v:= range args {
		var result interface{} = v
		if ! f.Type().IsVariadic() {
			_, isTag := v.(Tag)
			expectedType := call.expectedTypeOfArg(key)
			switch  {
			case expectedType == TagType && isTag:
			case expectedType == QuotedType: result = Quoted{v}
			default:
				evaluated, err := Eval(context, v)
				if err != nil {
					return tuple.EMPTY, err
				}
				converted, err := Convert(context, evaluated, expectedType)
				if err != nil {
					context.Log("ERROR", "Cannot convert '%s' to '%s'", evaluated, expectedType)  // TODO Log v location
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

/////////////////////////////////////////////////////////////////////////////

type SymbolTable struct {
	symbols map[string]reflect.Value
	notFound Finder
}

func NewSymbolTable(notFound Finder) SymbolTable {
	return SymbolTable{map[string]reflect.Value{},notFound}
}

func (context * SymbolTable) Arity() int {
	return len(context.symbols)
}

func  (context * SymbolTable) ForallValues(next func(value Value) error) error {
	for k, v := range context.symbols {
		key := signatureOfFunction(k, v)
		err := next(Tag{key})
		if err != nil {
			return err
		}
	}
	return nil
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

func (table * SymbolTable) Find(context EvalContext, head Tag, args []Value) (LocalScope, reflect.Value) {  // Reduce

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
	return table.notFound.Find(context, head, args)
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
