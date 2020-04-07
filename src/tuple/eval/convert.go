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
import "strconv"
import "errors"

// TODO could populate this just from the function
var conversions = map[string]reflect.Value{
	"Int64 int64": reflect.ValueOf(func(value Int64) int64 { return int64(value) }),
	"Float64 float64": reflect.ValueOf(func(value Float64) float64 { return float64(value) }),
	"Bool bool": reflect.ValueOf(func(value Bool) bool { return bool(value) }),
	"String string": reflect.ValueOf(func(value String) string { return string(value) }),
	"Float64 bool":reflect.ValueOf(func (value Float64) bool { return value!=0. }),
	"Int64 bool": reflect.ValueOf(func (value Int64) bool { return value!=0 }),
	"Bool float64":reflect.ValueOf(boolToFloat),
	"Bool int64": reflect.ValueOf(boolToInt),
	"Int64 float64": reflect.ValueOf(func(value Int64) float64 { return float64(int64(value)) }),
	"Float64 int64": reflect.ValueOf(func(value Float64) int64 { return int64(float64(value)) }),
	"Float64 string": reflect.ValueOf(fmt.Sprint),  // TODO Inf rather than +Inf
	"Int64 string": reflect.ValueOf(Int64ToString),
}

func Convert (context EvalContext, evaluated Value, expectedType reflect.Type) (interface{}, error) {

	if expectedType == reflect.TypeOf(evaluated) || expectedType == ValueType {
		return evaluated, nil
	}
	key := reflect.TypeOf(evaluated).Name() + " " + expectedType.Name()
	Verbose(context, "key=%s", key)
	convert, ok := conversions[key]
	if ok {
		in := convert.Call([]reflect.Value{reflect.ValueOf(evaluated)})
		reflectValue := in [0]
		switch reflectValue.Type() {
		case IntType: return reflectValue.Int(), nil
		case FloatType: return reflectValue.Float(), nil
		case BoolType: return reflectValue.Bool(), nil
		case StringType: return reflectValue.String(), nil
		case TupleType: return reflectValue.Interface().(Tuple), nil
		case TagType: return reflectValue.Interface().(Tag), nil
		case ValueType: return reflectValue.Interface().(Value), nil
		default:
		}
	}
	return nil, errors.New("No conversion")
}

func boolToFloat(value Bool) float64 {
	if bool(value) {
		return 1.
	}
	return 0.0
}
func boolToInt(value Bool) int64 {
	if bool(value) {
		return 1
	}
	return 0
}
func Int64ToString(value int64) string {
	return strconv.FormatInt(int64(value), 10)
}

func convertCallResult(table * SymbolTable, in []reflect.Value) Value {
	if len(in) == 0  {
		return tuple.EMPTY  // TODO VOID
	}
	result:= in[0]
	switch result.Type() {
	case IntType: return tuple.Int64(result.Int())
	case FloatType: return tuple.Float64(result.Float())
	case BoolType: return tuple.Bool(result.Bool())
	case StringType: return tuple.String(result.String())
	case TupleType: return result.Interface().(Tuple)
	case TagType: return result.Interface().(Tag)
	case ValueType: return result.Interface().(Value)
	default:
		table.Error(tuple.EMPTY, "Cannot find type of '%s'", result) // TODO EMPTY
		return tuple.NAN
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
