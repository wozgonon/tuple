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

import "tuple"
import "reflect"
import "errors"

// TODO type constants
var IntType = reflect.TypeOf(int64(1))
var FloatType = reflect.TypeOf(float64(1.0))
var BoolType = reflect.TypeOf(true)
var StringType = reflect.TypeOf("")
var TupleType = reflect.TypeOf(tuple.NewTuple())
var TagType = reflect.TypeOf(Tag{""})
var ValueType = reflect.TypeOf(func (_ Value) {}).In(0)
var MapType = reflect.TypeOf(func (_ tuple.Map) {}).In(0)
var ArrayType = reflect.TypeOf(func (_ tuple.Array) {}).In(0)
var EvalContextType = reflect.TypeOf(func (_ EvalContext) {}).In(0)
var QuotedType = reflect.TypeOf(func (_ Quoted) {}).In(0)


// Represents a call to a function using the golang reflect API.
type ReflectCall struct {
	function reflect.Value
	functionType reflect.Type
	start int
	args []reflect.Value
}

func NewReflectCall(context EvalContext, f reflect.Value, nn int) ReflectCall {
	t := f.Type()
	start := 0
	if takesContext(t) {
		start = 1
	}
	//nn := t.NumIn()
	reflectedArgs := make([]reflect.Value, nn + start)
	if start == 1 {
		reflectedArgs[0] = reflect.ValueOf(context)
	}
	return ReflectCall{f, t, start, reflectedArgs}
}

func takesContext(typ reflect.Type) bool {
	nn := typ.NumIn()
	return nn > 0 && typ.In(0) == EvalContextType
}

func (call * ReflectCall) expectedTypeOfArg(k int) reflect.Type {
	return call.functionType.In(k + call.start)
}

func (call * ReflectCall) setArg(key int, value interface{}) {
	call.args[key  + call.start] = reflect.ValueOf(value)
}

func (call * ReflectCall) Call(context EvalContext, head string) (Value, error) {
	Trace(context, "Call '%s' (%s)", head, call.args)
	reflectValues := call.function.Call(call.args)
	Trace(context, "  Call '%s' (%s)   f=%s -> %s", head, call.args, call.function, reflectValues)

	if len(reflectValues) == 0  {
		return tuple.EMPTY, nil  // TODO VOID
	}
	if len(reflectValues) == 2 {
		err := reflectValues[1].Interface()
		if err != nil {
			return nil, err.(error)
		}
	}
	return reflectValueToValue(reflectValues[0])
}

func reflectValueToValue(result reflect.Value) (Value, error) {
	switch result.Type() {
	case IntType: return tuple.Int64(result.Int()), nil
	case FloatType: return tuple.Float64(result.Float()), nil
	case BoolType: return tuple.Bool(result.Bool()), nil
	case StringType: return tuple.String(result.String()), nil
	case TupleType: return result.Interface().(Tuple), nil
	case TagType: return result.Interface().(Tag), nil
	case ValueType: return result.Interface().(Value), nil
	default:
		return nil, errors.New("Cannot find type of: " + result.Type().Name())
	}
}

