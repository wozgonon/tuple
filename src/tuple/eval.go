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
	ifFunctionNotFound FunctionNotFound
}

type FunctionNotFound func(name Atom, args [] Value) Value

func ErrorIfFunctionNotFound(name Atom, args [] Value) Value {
	fmt.Printf("ERROR: function not found: '%s' %s\n", name.Name, args)  // TODO ought to use context logger
	return Bool(false)
}

func NewSymbolTable(notFound FunctionNotFound) SymbolTable {

	table := SymbolTable{map[string]reflect.Value{},notFound}
	table.Add("len", func(value string) int64 { return int64(len(value)) })
	table.Add("lower", strings.ToLower)
	table.Add("upper", strings.ToUpper)
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
	table.Add("^", math.Pow)
	table.Add("+", func (aa float64) float64 { return aa })
	table.Add("-", func (aa float64) float64 { return -aa })
	table.Add("+", func (aa float64, bb float64) float64 { return aa+bb })
	table.Add("-", func (aa float64, bb float64) float64 { return aa-bb })
	table.Add("*", func (aa float64, bb float64) float64 { return aa*bb })
	table.Add("/", func (aa float64, bb float64) float64 { return aa/bb })
	table.Add("==", func (aa float64, bb float64) bool { return aa==bb })
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

func (table * SymbolTable) Add(name string, function interface{}) {
	reflectValue := reflect.ValueOf(function)
	typ := reflectValue.Type()
	key := makeKey(name, typ.NumIn())
	table.symbols[key] = reflectValue
}

func (table SymbolTable) Call(head Atom, args []Value) Value {

	name := head.Name
	nn := len(args)

	key := makeKey(name, len(args))
	f, ok := table.symbols[key]
	if ok {
		t := f.Type()
		//fmt.Printf("FUNC %s %d - %s %s\n", name, nn, key, t)
		//fmt.Printf("  FUNC %s %d - %s %d\n", name, nn, key, t.NumIn())
		//fmt.Printf("   FUNC %s %d\n", name, nn)
		reflectedArgs := make([]reflect.Value, nn)
		for k,_:= range args {
			reflectedArgs[k] = mapToReflectValue(args[k], t.In(k))
		}
		reflectValue := f.Call(reflectedArgs)
		return mapFromReflectValue(reflectValue)
	}
	return table.ifFunctionNotFound(head, args)
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

func (table SymbolTable) Eval(expression Value) Value {

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
		evaluated := make([]Value, ll-1)
		for k, v := range val.List[1:] {
			evaluated[k] = table.Eval(v)
		}
		return table.Call(atom, evaluated)
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

func mapToReflectValue (v Value, expected reflect.Type) reflect.Value {

	_, isTuple := v.(Tuple)

	var result interface{}
	
	/// TODO should this take a Scalar or a Value?
	switch expected {
	case IntType: result = toInt64(v)
	case FloatType: result = toFloat64(v)
	case BoolType: result = toBool(v)
	case StringType: result = toString(v)
	case TupleType:
		if isTuple {
			result = v
		}
	default: result = float64(v.(Float64)) // TODO???
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
	default: return Float64(in[0].Float()) // TODO
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

