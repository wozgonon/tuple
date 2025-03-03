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
package parsers

import "tuple"

/////////////////////////////////////////////////////////////////////////////
// JSON Grammar
/////////////////////////////////////////////////////////////////////////////

type JSONGrammar struct {
	Style
	operators Operators
}

func (grammar JSONGrammar) Name() string {
	return "JSON"
}

func (grammar JSONGrammar) FileSuffix() string {
	return ".json"
}

func (grammar JSONGrammar) Parse(context Context, next Next) error {
	return parse(context, grammar.operators, grammar.Style, next)
}

func (grammar JSONGrammar) Print(object Value, next func(value string)) {
	PrintExpression(grammar, "", object, next)  // TODO Use Printer
}

var JSON_CONS_OPERATOR = ":"

func NewJSONGrammar() Grammar {
	style := NewStyle("", "", "  ",
		OPEN_SQUARE_BRACKET, CLOSE_SQUARE_BRACKET, OPEN_BRACE, CLOSE_BRACE, JSON_CONS_OPERATOR,
		",", "\n", "true", "false", '%', "") // prolog, sql '--' for   // TODO remove comment %
	style.RecognizeNegative = true
	operators := NewOperators(style)
	operators.AddBracket(OPEN_SQUARE_BRACKET, CLOSE_SQUARE_BRACKET)
	operators.AddBracket(OPEN_BRACE, CLOSE_BRACE)
	//operators.AddBracket(OPEN_BRACKET, CLOSE_BRACKET)
	operators.AddInfix(CONS_ATOM.Name, 30)
	operators.AddInfix(";", 10)
	operators.AddInfix(SPACE_ATOM.Name, 20)  // TODO space???
	return JSONGrammar{style, operators}

}

func (printer JSONGrammar) PrintKey(tag Tag, out StringFunction) {
	out(tuple.DoubleQuotedString(tag.Name))
	out(printer.KeyValueSeparator)
	out(" ")
}

func (printer JSONGrammar) PrintNullaryOperator(depth string, tag Tag, out StringFunction) {
	PrintTuple(&printer, depth, NewTuple(tag), out)
}

func (printer JSONGrammar) PrintUnaryOperator(depth string, tag Tag, value Value, out StringFunction) {
	PrintTuple(&printer, depth, NewTuple(tag, value), out)
}

func (printer JSONGrammar) PrintSeparator(depth string, out StringFunction) {
	out(printer.Style.Separator)
}
