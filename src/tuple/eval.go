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
import "reflect"
import "fmt"
import "strconv"

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

type SymbolTable struct {
	symbols map[string]reflect.Value
	ifFunctionNotFound CallHandler
}

func NewSymbolTable(notFound CallHandler) SymbolTable {
	return SymbolTable{map[string]reflect.Value{},notFound}
}

func ValuesToStrings(values []Value) []string {
	result := make([]string, len(values))
	for k,_:= range values {
		result[k] = toString(values[k])
	}
	return result
}

func makeKey(name string, arity int) string {
	return fmt.Sprintf("%d_%s", arity, name)
}


func (table * SymbolTable) Count() int {
	return len(table.symbols)
}

func (table * SymbolTable) Add(name string, function interface{}) {
	reflectValue := reflect.ValueOf(function)
	typ := reflectValue.Type()
	key := makeKey(name, typ.NumIn())
	table.symbols[key] = reflectValue
}

func (table * SymbolTable) Call(head Atom, args []Value) Value {  // Reduce

	//name := head.Name
	nn := len(args)

	f := table.Find(head, args)
	t := f.Type()
	
	reflectedArgs := make([]reflect.Value, nn)
	for k,v:= range args {
		if t.IsVariadic() {
			reflectedArgs[k] = reflect.ValueOf(v)
		} else {
			expectedType := t.In(k)
			if expectedType == AtomType {
				reflectedArgs[k] = reflect.ValueOf(v)
			} else if expectedType == TupleType {
				reflectedArgs[k] = reflect.ValueOf(v)
			} else if expectedType == ValueType {
				reflectedArgs[k] = reflect.ValueOf(v)
			} else {
				//fmt.Printf("Eval %s expectedType=%s\n", v, expectedType)
				evaluated := table.Eval(v)
				reflectedArgs[k] = mapToReflectValue(evaluated, expectedType)
			}
		}
	}
	//fmt.Printf("Call '%s' (%s)   f=%s\n", head, reflectedArgs, f)
	reflectValue := f.Call(reflectedArgs)
	return mapFromReflectValue(reflectValue)
}

func (table * SymbolTable) Find(head Atom, args []Value) reflect.Value {  // Reduce

	name := head.Name
	nn := len(args)

	//fmt.Printf("Call '%s' nn=%d count=%d\n", name, nn, table.Count())

	key := makeKey(name, nn)
	f, ok := table.symbols[key]
	if ok {
		return f
	}
	return table.ifFunctionNotFound.Find(head, args)
}

func (table * SymbolTable) Eval(expression Value) Value {

	switch val := expression.(type) {
	case Tuple:
		ll := len(val.List)
		if ll == 0 {
			return val
		}
		head := val.List[0]
		if ll == 1 {
			return table.Eval(head)
		}
		atom, ok := head.(Atom)
		if ! ok {
			return val // TODO Handle case of list: (1 2 3)
		}
		return table.Call(atom, val.List[1:])
	case Atom:
		return table.Call(val, []Value{})
	default:
		return val
	}
}

// TODO type constants
var IntType = reflect.TypeOf(int64(1))
var FloatType = reflect.TypeOf(float64(1.0))
var BoolType = reflect.TypeOf(true)
var StringType = reflect.TypeOf("")
var TupleType = reflect.TypeOf(NewTuple())


var AtomType = reflect.TypeOf(Atom{""})
var ValueType = reflect.TypeOf(func (_ Value) {}).In(0)


func mapToReflectValue (v Value, expected reflect.Type) reflect.Value {

	_, isTuple := v.(Tuple)
	_, isAtom := v.(Atom)

	var result interface{}
	
	/// TODO should this take a Scalar or a Value?
	switch expected {
	case IntType: result = toInt64(v)
	case FloatType: result = toFloat64(v)
	case BoolType: result = toBool(v)
	case StringType: result = toString(v)
	case AtomType:
		if isAtom {
			result = v
		}
	case TupleType:
		if isTuple {
			result = v
		}
	case ValueType:
		result = v
	default:
		fmt.Printf("ERROR should not get here Expected type: '%s' v=%s", expected, v) // TODO
		result = float64(v.(Float64)) // TODO???
	}
	return reflect.ValueOf(result)
}

func mapFromReflectValue (in []reflect.Value) Value {
	v:= in[0]
	switch v.Type() {
	case IntType: return Int64(v.Int())
	case FloatType: return Float64(v.Float())
	case BoolType: return Bool(v.Bool())
	case StringType: return String(v.String())
	case TupleType: return v.Interface().(Tuple)
	case AtomType: return v.Interface().(Atom)
	case ValueType: return v.Interface().(Value)
	default:
		fmt.Printf("Cannot find type of '%s'\n", v)
		return NAN
	}
}

func toString(value Value) string {
	switch val := value.(type) {
	case Atom: return val.Name
	case String: return string(val)  // Quote ???
	case Float64: return  fmt.Sprint(val)  // TODO Inf ???
	case Int64: return strconv.FormatInt(int64(val), 10)
	case Bool:
		if val {
			return "true"
		} else {
			return "false"
		}
	default: return "..." // TODO
	}
}

func toBool(value Value) bool {
	switch val := value.(type) {
	case Int64: return val != 0
	case Float64: return val != 0
	case Atom:
		if val.Name == "true" {
			return true
		}
		return false // TODO Nullary(val)
	case Bool: return bool(val)
	default: return false
	}
}

func toFloat64(value Value) float64 {
	switch val := value.(type) {
	case Int64: return float64(val)
	case Float64: return float64(val)
	case Atom: return math.NaN()
	case String: return math.NaN()
	case Bool:
		if val {
			return 1
		} else {
			return 0
		}
	default:
		return math.NaN()
	}
}

func toInt64(value Value) int64 {
	switch val := value.(type) {
	case Int64: return int64(val)
	case Float64: return int64(val)
	case Bool:
		if val {
			return 1
		} else {
			return 0
		}
	case Atom: return -1 // TODO
	default:
		return -1 //TODO
	}
}

