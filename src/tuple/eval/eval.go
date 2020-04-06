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

import "math"
import "reflect"
import "fmt"
import "strconv"
import "tuple"

type Tag = tuple.Tag
type Value = tuple.Value
type Tuple = tuple.Tuple
type Int64 = tuple.Int64
type Float64 = tuple.Float64
type Bool = tuple.Bool
type String = tuple.String
type Logger = tuple.Logger


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

type CallHandler interface {
	Find(context EvalContext, name Tag, args [] Value) (*SymbolTable, reflect.Value)
}

type EvalContext interface {
	CallHandler
	Logger

	Add(name string, function interface{})
	//Eval(expression Value) Value
	Call(head Tag, args []Value) Value  // Reduce
}



type SymbolTable struct {
	symbols map[string]reflect.Value
	ifFunctionNotFound CallHandler
}

func (context * SymbolTable) Call(head Tag, args []Value) Value {
	return context.call3(context, head, args)
}

func (context * SymbolTable) Log(level string, format string, args ...interface{}) {
	if level == "VERBOSE" || level == "TRACE" {  // TODO
		return
	}
	fmt.Printf(level + " " + format + "\n", args...)
}


func NewSymbolTable(notFound CallHandler) SymbolTable {
	return SymbolTable{map[string]reflect.Value{},notFound}
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
		value := Eval(context, values[k])
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

func evalTuple(context EvalContext, value Tuple) Tuple {
	newTuple := tuple.NewTuple()
	for _,v:= range value.List {
		newTuple.Append(Eval(context, v))
	}
	context.Log("TRACE", "Eval tuple return '%s'", newTuple)
	return newTuple
}

func (table * SymbolTable) call3(context EvalContext, head Tag, args []Value) Value {  // Reduce

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
			context.Log("VERBOSE", "isvariadic='%s' ", t) 
			result = v
		} else {
			_, isTuple := v.(Tuple)
			_, isTag := v.(Tag)
			
			expectedType := t.In(k)
			context.Log("VERBOSE", "expected type='%s' got '%s'", expectedType, v)
			switch  {
			case expectedType == TagType && isTag: result = v
			case expectedType == TupleType && isTuple: result = v.(Tuple) //newTuple := evalTuple(context, v.(Tuple))
			case expectedType == ValueType: result = v //newTuple := evalTuple(context, v.(Tuple))
			default:
				evaluated := Eval(context, v)
				context.Log("TRACE", "** EVAL head=%s  v=%s-> evaluated=%s type=(%s) expectedType=%s", head, v, evaluated, reflect.TypeOf(evaluated), expectedType)
				/// TODO should this take a Scalar or a Value?
				switch expectedType {
				case IntType: result = toInt64(context, evaluated)
				case FloatType: result = toFloat64(context, evaluated)
				case BoolType: result = toBool(evaluated)
				case StringType: result = toString(context, evaluated)
				case TagType:
					if _, isTag := evaluated.(Tag); isTag {
						result = evaluated
					} else {
						context.Log("ERROR", "Expected tag but got: %s", evaluated)
						result = Tag{""}
					}
				case TupleType:
					context.Log("TRACE", "** TUPLE isTuple=%s evaluated=%s", isTuple, evaluated)
					if _, isTuple := evaluated.(Tuple); isTuple {
						result = evaluated
					} else {
						context.Log("ERROR", "Expected tuple but got: %s", evaluated)
						result = tuple.NewTuple()
					}
				case ValueType:
					result = evaluated
				default:
					context.Log("ERROR", "should not get here Expected type: '%s' v=%s", expectedType, evaluated) // TODO
					result = evaluated //float64(v.(Float64)) // TODO???
				}
			}
		}
		if result == nil {
			context.Log("ERROR", "MUST not be nil v=%s head=%s", v, head)
		}
		context.Log("TRACE", "Call '%s' arg=%d value=%s", head, k, result)
		reflectedArgs[k] = reflect.ValueOf(result)
	}
	context.Log("TRACE", "Call '%s' (%s)", head, reflectedArgs)
	reflectValue := f.Call(reflectedArgs)
	context.Log("TRACE", "  Call '%s' (%s)   f=%s -> %s", head, reflectedArgs, f, reflectValue)

	in := reflectValue
	if len(in) == 0  {
		return tuple.EMPTY  // TODO VOID
	}
	v:= in[0]
	switch v.Type() {
	case IntType: return tuple.Int64(v.Int())
	case FloatType: return tuple.Float64(v.Float())
	case BoolType: return tuple.Bool(v.Bool())
	case StringType: return tuple.String(v.String())
	case TupleType: return v.Interface().(Tuple)
	case TagType: return v.Interface().(Tag)
	case ValueType: return v.Interface().(Value)
	default:
		context.Log("ERROR", "Cannot find type of '%s'", v)
		return tuple.NAN
	}
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
	return table.ifFunctionNotFound.Find(context, head, args)
}

func Eval(context EvalContext, expression Value) Value {

	switch val := expression.(type) {
	case Tuple:
		ll := len(val.List)
		if ll == 0 {
			return val
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
		return val
	}
}

// TODO type constants
var IntType = reflect.TypeOf(int64(1))
var FloatType = reflect.TypeOf(float64(1.0))
var BoolType = reflect.TypeOf(true)
var StringType = reflect.TypeOf("")
var TupleType = reflect.TypeOf(tuple.NewTuple())


var TagType = reflect.TypeOf(Tag{""})
var ValueType = reflect.TypeOf(func (_ Value) {}).In(0)
var EvalContextType = reflect.TypeOf(func (_ EvalContext) {}).In(0)


func toString(context EvalContext, value Value) string {
	switch val := value.(type) {
	case Tag: return val.Name
	case String: return string(val)  // Quote ???
	case Float64: return  fmt.Sprint(val)  // TODO Inf ???
	case Int64: return strconv.FormatInt(int64(val), 10)
	case Bool:
		if val {
			return "true"
		} else {
			return "false"
		}
	default: 
		context.Log("ERROR", "cannot convert '%s' to string", value)
		return "..." // TODO
	}
}

func toBool(value Value) bool {
	switch val := value.(type) {
	case Int64: return val != 0
	case Float64: return val != 0
	case Tag:
		if val.Name == "true" {
			return true
		}
		return false // TODO Nullary(val)
	case Bool: return bool(val)
	default: return false
	}
}

func toFloat64(context EvalContext, value Value) float64 {
	switch val := value.(type) {
	case Int64: return float64(val)
	case Float64: return float64(val)
	case Bool:
		if val {
			return 1
		} else {
			return 0
		}
	default:
		context.Log("ERROR", "cannot convert '%s' to float", value)
		return math.NaN()
	}
}

func toInt64(context EvalContext, value Value) int64 {
	switch val := value.(type) {
	case Int64: return int64(val)
	case Float64: return int64(val)
	case Bool:
		if val {
			return 1
		} else {
			return 0
		}
	default:
		context.Log("ERROR", "cannot convert '%s' to int", value)
		return -1 //TODO
	}
}

