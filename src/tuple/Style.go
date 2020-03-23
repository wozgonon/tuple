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

import "fmt"
import "strconv"
import "math"

/////////////////////////////////////////////////////////////////////////////
//
/////////////////////////////////////////////////////////////////////////////

type Style struct {
	StartDoc string
	EndDoc string
	Indent string

	Open string
	Close string
	Open2 string
	Close2 string
	KeyValueSeparator string
	
	Separator string
	LineBreak string
	True string
	False string
	OneLineComment rune
	ScalarPrefix string

	//INF string
	//NAN string
}

/////////////////////////////////////////////////////////////////////////////
// Printer
/////////////////////////////////////////////////////////////////////////////

func (style Style) PrintIndent(depth string, out StringFunction) {
	out(depth)
}

func (style Style) PrintSuffix(depth string, out StringFunction) {
	out(string(NEWLINE))
}

func (style Style) PrintSeparator(depth string, out StringFunction) {
	//out(style.Separator)
}

func (style Style) PrintEmptyTuple(depth string, out StringFunction) {
	out(style.Open)
	out(style.Close)
}

func (style Style) PrintNullaryOperator(depth string, atom Atom, out StringFunction) {
	PrintTuple(&style, depth, NewTuple(atom), out)
}

func (style Style) PrintUnaryOperator(depth string, atom Atom, value interface{}, out StringFunction) {
	PrintTuple(&style, depth, NewTuple(atom, value), out)
}

func (style Style) PrintBinaryOperator(depth string, atom Atom, value1 interface{}, value2 interface{}, out StringFunction) {
	PrintTuple(&style, depth, NewTuple(atom, value1, value2), out)
}

func (style Style) PrintOpenTuple(depth string, tuple Tuple, out StringFunction) string {
	if tuple.IsCons() {
		out(style.Open2)
	} else {
		out(style.Open)
	}
	return depth + "  "
}

func (style Style) PrintCloseTuple(depth string, tuple Tuple, out StringFunction) {
	if tuple.IsCons() {
		out(style.Close2)
	} else {
		out(style.Close)
	}
}

func (style Style) PrintAtom(depth string, atom Atom, out StringFunction) {
	out(atom.Name)
}

func (style Style) PrintInt64(depth string, value int64, out StringFunction) {
	out(strconv.FormatInt(value, 10))
}

func (style Style) PrintFloat64(depth string, value float64, out StringFunction) {
	if math.IsInf(value, 64) {
		out("Inf") // style.INF)  // Do not print +Inf
	} else {
		out(fmt.Sprint(value))
	}
}

func (style Style) PrintString(depth string, value string, out StringFunction) {
	out(DOUBLE_QUOTE)
	out(value)   // TODO Escape
	out(DOUBLE_QUOTE)
}

func (style Style) PrintBool(depth string, value bool, out StringFunction) {
	if value {
		out(style.True)
	} else {
		out(style.False)
	}				
}

func (style Style) PrintComment(depth string, value Comment, out StringFunction) {
	out(string(style.OneLineComment))
	out(value.Comment)
}

