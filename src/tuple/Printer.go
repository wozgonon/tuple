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

type Printer interface {

	PrintIndent(depth string, out StringFunction)
	PrintSuffix(depth string, out StringFunction)
	PrintEmptyTuple(depth string, out StringFunction)
	PrintNullaryOperator(depth string, atom Atom, out StringFunction)
	PrintUnaryOperator(depth string, atom Atom, value interface{}, out StringFunction)
	PrintBinaryOperator(depth string, atom Atom, value1 interface{}, value2 interface{}, out StringFunction)
	PrintOpenTuple(depth string, out StringFunction) string
	PrintCloseTuple(depth string, out StringFunction)
	PrintAtom(depth string, value Atom, out StringFunction)
	PrintInt64(depth string, value int64, out StringFunction)
	PrintFloat64(depth string, value float64, out StringFunction)
	PrintString(depth string, value string, out StringFunction)
	PrintBool(depth string, value bool, out StringFunction)
	PrintComment(depth string, value Comment, out StringFunction)
}

type StringFunction func(value string)

func PrintScalar(printer Printer, depth string, token interface{}, out StringFunction) {
	switch token.(type) {
	case Atom:
		printer.PrintAtom(depth, token.(Atom), out)
	case string:
		printer.PrintString(depth, token.(string), out)
	case bool:
		printer.PrintBool(depth, token.(bool), out)
	case Comment:
		printer.PrintComment(depth, token.(Comment), out)
	case int64:
		printer.PrintInt64(depth, token.(int64), out)
	case float64:
		printer.PrintFloat64(depth, token.(float64), out)
	default:
		log.Printf("ERROR type '%s' not recognised: %s", reflect.TypeOf(token), token);
	}
}

func PrintTuple(printer Printer, depth string, tuple Tuple, out StringFunction) {
	newDepth := printer.PrintOpenTuple(depth, out)
	printer.PrintSuffix(depth, out)
	for _, value := range tuple.List {
		PrintExpression(printer, newDepth, value, out)
	}
	printer.PrintIndent(depth, out)
	printer.PrintCloseTuple(depth, out)
}

func PrintExpression(printer Printer, depth string, token interface{}, out StringFunction) {

	switch token.(type) {
	case Tuple:
		tuple := token.(Tuple)
		len := len(tuple.List)
		if len == 0 {
			printer.PrintIndent(depth, out)
			printer.PrintEmptyTuple(depth, out)
			printer.PrintSuffix(depth, out)
		} else {
			head := tuple.List[0]
			atom, ok := head.(Atom)
			//log.Printf("Tuple [%s] %d\n", atom, len)
			if ok {  // TODO and head in a (binary) operator
				switch len {
				case 1:
					printer.PrintIndent(depth, out)
					printer.PrintNullaryOperator(depth, atom, out)
					printer.PrintSuffix(depth, out)
				case 2:
					printer.PrintIndent(depth, out)
					printer.PrintUnaryOperator(depth, atom, tuple.List[1], out)
					printer.PrintSuffix(depth, out)
				case 3:
					printer.PrintIndent(depth, out)
					printer.PrintBinaryOperator(depth, atom, tuple.List[1], tuple.List[2], out)
					printer.PrintSuffix(depth, out)
				default:
					printer.PrintIndent(depth, out)
					PrintTuple(printer, depth, tuple, out)
					printer.PrintSuffix(depth, out)
				}
			} else {
				printer.PrintIndent(depth, out)
				PrintTuple(printer, depth, tuple, out)
				printer.PrintSuffix(depth, out)
			}
		}
	default:
		printer.PrintIndent(depth, out)
		PrintScalar(printer, depth, token, out)
		printer.PrintSuffix(depth, out)
	}
}
