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
import "errors"
import "tuple"


// A table of basic conversion functions
var conversions = NewConversions(
	func(value Int64) int64 { return int64(value) },
	func(value Float64) float64 { return float64(value) },
	func(value Bool) bool { return bool(value) },
	func(value String) string { return string(value) },
	func (value Float64) bool { return value!=0. },
	func (value Int64) bool { return value!=0 },
	tuple.BoolToFloat,
	tuple.BoolToInt,
	func(value Int64) float64 { return float64(int64(value)) },
	func(value Float64) int64 { return int64(float64(value)) },
	func(value tuple.TagValueMap) tuple.Map { return value },
	fmt.Sprint,  // TODO Inf rather than +If
	tuple.Int64ToString)

/////////////////////////////////////////////////////////////////////////////
//  Conversions using reflection
/////////////////////////////////////////////////////////////////////////////

type Conversions struct {
	functions map[string]reflect.Value
}

func NewConversions(functions ... interface{}) Conversions {

	result := Conversions{make(map[string]reflect.Value)}
	for _, function := range functions {
		value := reflect.ValueOf(function)
		typ := value.Type()
		in := typ.In(0).Name()
		out := typ.Out(0).Name()
		key := in + " " + out
		result.functions[key] = value
	}
	return result
}

func (conversions Conversions) Convert(evaluated Value, expectedType reflect.Type) (interface{}, error) {
	key := reflect.TypeOf(evaluated).Name() + " " + expectedType.Name()
	//Verbose(context, "key=%s", key)
	convert, ok := conversions.functions[key]
	if ok {
		reflectValues := convert.Call([]reflect.Value{reflect.ValueOf(evaluated)})
		reflectValue := reflectValues [0]
		//return reflectValueToValue(in[0])
		switch reflectValue.Type() {
		case IntType: return reflectValue.Int(), nil
		case FloatType: return reflectValue.Float(), nil
		case BoolType: return reflectValue.Bool(), nil
		case StringType: return reflectValue.String(), nil
		case TupleType: return reflectValue.Interface().(Tuple), nil  // TODO is this needed
		case TagType: return reflectValue.Interface().(Tag), nil
		case ValueType: return reflectValue.Interface().(Value), nil
		case MapType: return reflectValue.Interface().(tuple.Map), nil
		case ArrayType: return reflectValue.Interface().(tuple.Array), nil
		default:
		}
	}
	message := fmt.Sprintf("No conversion '%s', cannot convert '%s' to '%s'", key, evaluated, expectedType)  // TODO Log v location
	return nil, errors.New(message)
}

func Convert (context EvalContext, evaluated Value, expectedType reflect.Type) (interface{}, error) {

	if expectedType == reflect.TypeOf(evaluated) || expectedType == ValueType {
		return evaluated, nil
	}
	if array, isArray := evaluated.(tuple.Array); isArray && expectedType == ArrayType {
		return array, nil
	}
	return conversions.Convert(evaluated, expectedType)
}
