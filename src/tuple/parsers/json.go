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

/////////////////////////////////////////////////////////////////////////////
// JSON Grammar
/////////////////////////////////////////////////////////////////////////////

type JSONGrammar struct {
	Style
}

func (grammar JSONGrammar) Name() string {
	return "JSON"
}

func (grammar JSONGrammar) FileSuffix() string {
	return ".json"
}

func (grammar JSONGrammar) Parse(context Context, next Next) {
	parser := NewSExpressionParser(grammar.Style)
	parser.Parse(context, next)
}

func (grammar JSONGrammar) Print(object Value, out func(value string)) {
	PrintExpression(grammar, "", object, out)  // TODO Use Printer
}

func NewJSONGrammar() Grammar {
	style := NewStyle("", "", "  ",
		OPEN_SQUARE_BRACKET, CLOSE_SQUARE_BRACKET, OPEN_BRACE, CLOSE_BRACE, ":",
		",", "\n", "true", "false", '%', "") // prolog, sql '--' for 
	return JSONGrammar{style}
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

func (printer JSONGrammar) PrintBinaryOperator(depth string, tag Tag, value1 Value, value2 Value, out StringFunction) {

	if tag == CONS_ATOM {
		newDepth := depth + "  "
		printer.PrintIndent(depth, out)
		PrintExpression1(printer, newDepth, value1, out)
		out(" :")
		if _, ok := value2.(Tuple); ok {
			printer.PrintSuffix(newDepth, out)
			printer.PrintIndent(newDepth, out)
		} else {
			out (" ")
		}
		PrintExpression1(printer, newDepth, value2, out)
	} else {
		PrintTuple(&printer, depth, NewTuple(tag, value1, value2), out)
	}
}
