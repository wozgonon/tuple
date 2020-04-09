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

import "log"
import "reflect"

/////////////////////////////////////////////////////////////////////////////
//  Printer
/////////////////////////////////////////////////////////////////////////////

type Printer interface {

	PrintIndent(depth string, out StringFunction)
	PrintSuffix(depth string, out StringFunction)
	PrintScalarPrefix(depth string, out StringFunction)
	PrintSeparator(depth string, out StringFunction)
	PrintEmptyTuple(depth string, out StringFunction)
	PrintNullaryOperator(depth string, tag Tag, out StringFunction)
	PrintUnaryOperator(depth string, tag Tag, value Value, out StringFunction)
	PrintBinaryOperator(depth string, tag Tag, value1 Value, value2 Value, out StringFunction)
	PrintOpenTuple(depth string, tuple Value, out StringFunction) string
	PrintCloseTuple(depth string, tuple Value, out StringFunction)
	PrintHeadTag(value Tag, out StringFunction)
	PrintScalar(depth string, token Value, out StringFunction)
}

func PrintScalar(printer Printer, depth string, value Value, out StringFunction) {
	printer.PrintScalarPrefix(depth, out)

	switch value.(type) {
	case Tag: out(value.(Tag).Name)
	case String: Quote(string(value.(String)), out)
	case Bool: out(BoolToString(bool(value.(Bool))))
	case Int64: out(Int64ToString(value.(Int64)))
	case Float64: out(Float64ToString(value.(Float64)))
	default:
		if value.Arity() == 0 {
			printer.PrintEmptyTuple(depth, out)
		} else {
			log.Printf("ERROR type '%s' not recognised: %s", reflect.TypeOf(value), value);  // TODO return error or prevent from ever happening
			//log.Printf("ERROR unexpected tuple '%s", value);  // TODO return error or prevent from ever happening
		}
	}
}


func PrintTuple(printer Printer, depth string, tuple Array, out StringFunction) {
	newDepth := printer.PrintOpenTuple(depth, tuple, out)
	printer.PrintSuffix(depth, out)
	ll := tuple.Arity()
	first := false
	if ll > 0 {
		_, first = tuple.Get(0).(Tag)
	}
	if mapp, ok := tuple.(Map); ok {
		mapp.ForallKeyValue(func (k Tag, value Value) {
			printer.PrintHeadTag(k, out)
			out (":") // TODO
			PrintExpression1(printer, newDepth, value, out)
			printer.PrintSeparator(newDepth, out)
			printer.PrintSuffix(depth, out)
		})
	} else {
		k := 0
		tuple.ForallValues(func (value Value) error {
			printer.PrintIndent(newDepth, out)
			if first && k == 0 {
				printer.PrintHeadTag(value.(Tag), out)
			} else {
				PrintExpression1(printer, newDepth, value, out)
			}
			if k < ll-1 {
				printer.PrintSeparator(newDepth, out)
			}
			printer.PrintSuffix(depth, out)
			k += 1
			return nil
		})
	}
	printer.PrintCloseTuple(depth, tuple, out)
}

func PrintExpression(printer Printer, depth string, token Value, out StringFunction) {
	printer.PrintIndent(depth, out)
	PrintExpression1(printer, depth, token, out)
	printer.PrintSuffix(depth, out)
}

func PrintExpression1(printer Printer, depth string, token Value, out StringFunction) {

	if IsAtom(token) {
		printer.PrintScalar(depth, token, out)
		return
	}
	ll := token.Arity()
	if array, ok := token.(Array); ok {
		head := array.Get(0)
		tag, ok := head.(Tag)
		//log.Printf("Array [%s] %d\n", tag, len)
		if ok {  // TODO and head in a (binary) operator
			switch ll {
			case 1:
				printer.PrintNullaryOperator(depth, tag, out)
			case 2:
				printer.PrintUnaryOperator(depth, tag, array.Get(1), out)
			case 3:
				printer.PrintBinaryOperator(depth, tag, array.Get(1), array.Get(2), out)
			default:
				PrintTuple(printer, depth, array, out)
			}
		} else {
			PrintTuple(printer, depth, array, out)
		}
		return
	}
	newDepth := printer.PrintOpenTuple(depth, token, out)
	printer.PrintSuffix(depth, out)
	if mapp, ok := token.(Map); ok {
		mapp.ForallKeyValue(func (k Tag, value Value) {
			printer.PrintHeadTag(k, out)
			out (":") // TODO
			PrintExpression1(printer, newDepth, value, out)
			printer.PrintSeparator(newDepth, out)
			printer.PrintSuffix(depth, out)
		})
	} else {
		k := 0
		token.ForallValues(func (value Value) error {
			printer.PrintIndent(newDepth, out)
			PrintExpression1(printer, newDepth, value, out)
			if k < ll-1 {
				printer.PrintSeparator(newDepth, out)
			}
			printer.PrintSuffix(depth, out)
			k += 1
			return nil
		})
	}
	printer.PrintCloseTuple(depth, token, out)
}
