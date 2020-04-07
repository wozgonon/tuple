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
}

type EvalContext interface {
	Global
	Logger

	Add(name string, function interface{})
	//Eval(expression Value) Value
	Call(head Tag, args []Value) (Value, error)  // Reduce
}

/////////////////////////////////////////////////////////////////////////////

type SymbolTable struct {
	symbols map[string]reflect.Value
	global Global
}

func NewSymbolTable(notFound Global) SymbolTable {
	if notFound.Logger() == nil {
		panic("nil logger")
	}
	return SymbolTable{map[string]reflect.Value{},notFound}
}

func (context * SymbolTable) Logger() LocationLogger {
	return context.global.Logger()
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


/*func ValuesToStrings(values []Value) []string {
	result := make([]string, len(values))
	for k,_:= range values {
		result[k] = toString(values[k])
	}
	return result
}*/

func EvalToStrings(context EvalContext, values []Value) []string {

	result := make([]string, len(values))
	for k,_:= range values {
		value, _ := Eval(context, values[k])
		result[k] = toString(context, value)
	}
	return result
}

/*
var b bytes.Buffer
	for _,value := range values {
		evaluated := Eval(context, value)
		str := toString(context, evaluated)
		b.WriteString(str)
	}
	return b.String()
}*/

func makeKey(name string, arity int, variadic bool) string {
	//return name  // TODO do exact match then do a general match
	if variadic {
		return fmt.Sprintf("*_%s", name)
	}
	return fmt.Sprintf("%d_%s", arity, name)
}


func (table * SymbolTable) Count() int {
	return len(table.symbols)
}

func TakesContext(typ reflect.Type) bool {
	nn := typ.NumIn()
	return nn > 0 && typ.In(0) == EvalContextType
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
	key := makeKey(name, nn, typ.IsVariadic())
	table.symbols[key] = reflectValue
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

func (context * SymbolTable) Call(head Tag, args []Value) (Value, error) {
	return context.call3(context, head, args)
}

func (table * SymbolTable) call3(context EvalContext, head Tag, args []Value) (Value, error) {  // Reduce

	//name := head.Name
	nn := len(args)

	_, f := table.Find(context, head, args)
	t := f.Type()

	start := 0
	if TakesContext(t) {
		start = 1
	}
	reflectedArgs := make([]reflect.Value, nn + start)
	if start == 1 {
		reflectedArgs[0] = reflect.ValueOf(context)
	}
	
	for key,v:= range args {
		k := start + key
		var result interface{}
		if t.IsVariadic() {
			Verbose(context, "isvariadic='%s' ", t) 
			result = v
		} else {
			_, isTuple := v.(Tuple)
			_, isTag := v.(Tag)
			
			expectedType := t.In(k)
			Verbose(context, "expected type='%s' got '%s'", expectedType, v)

			switch  {
			case expectedType == TagType && isTag: result = v
			case expectedType == TupleType && isTuple: result = v.(Tuple)
			case expectedType == ValueType: result = v
			default:
				evaluated, err := Eval(context, v)
				if err != nil {
					return tuple.EMPTY, err
				}
				Trace(context, "** EVAL head=%s  v=%s-> evaluated=%s type=(%s) expectedType=%s", head, v, evaluated, reflect.TypeOf(evaluated), expectedType)
				converted, err := Convert(context, evaluated, expectedType)
				if err != nil {
					table.Error(v, "Cannot convert '%s' to '%s'", evaluated, expectedType)
					return tuple.EMPTY, err
				}
				result = converted
			}
		}
		if result == nil {
			table.Error(v, "MUST not be nil v=%s head=%s", v, head)
			// TODO return nil, errors.New
		}
		Trace(context, "Call '%s' arg=%d value=%s", head, k, result)
		reflectedArgs[k] = reflect.ValueOf(result)
	}
	Trace(context, "Call '%s' (%s)", head, reflectedArgs)
	reflectValue := f.Call(reflectedArgs)
	Trace(context, "  Call '%s' (%s)   f=%s -> %s", head, reflectedArgs, f, reflectValue)
	return convertCallResult(table, reflectValue), nil
}

func (table * SymbolTable) Find(context EvalContext, head Tag, args []Value) (*SymbolTable, reflect.Value) {  // Reduce

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
	context.Log ("TRACE", "FIND Could not find '%s' in this symbol table (%d entries), forwarding",  head, len(table.symbols))
	// TODO look up variatic functions
	return table.global.Find(context, head, args)
}

func Eval(context EvalContext, expression Value) (Value, error) {

	switch val := expression.(type) {
	case Tuple:
		ll := len(val.List)
		if ll == 0 {
			return val, nil
		}
		head := val.List[0]
		if ll == 1 {
			return Eval(context, head)
		}
		tag, ok := head.(Tag)
		if ! ok {
			return evalTuple(context, val)
			//return val // TODO Handle case of list: (1 2 3)
		}
		return context.Call(tag, val.List[1:])
	case Tag:
		return context.Call(val, []Value{})
	default:
		return val, nil
	}
}
