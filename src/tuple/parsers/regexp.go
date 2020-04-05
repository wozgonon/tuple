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
import "io"

// To provide an example of how the OperatorGrammar can be used for a variety of purposes,
// in this case to parse regular expressions.
// TODO finish off and provide a full set of test cases.
func NewRegexpGrammar(context Context) OperatorGrammar {

	style := NewStyle("", "", "  ",
		OPEN_BRACKET, CLOSE_BRACKET, OPEN_BRACE, CLOSE_BRACE, ":",
		",", "\n", "true", "false", '#', "")

	operators := NewOperators(style)

	operators.AddBracket(OPEN_BRACKET, CLOSE_BRACKET)
	operators.AddBracket(OPEN_SQUARE_BRACKET, CLOSE_SQUARE_BRACKET)
	operators.AddBracket(OPEN_BRACE, CLOSE_BRACE)
	operators.AddPostfix("*", 105)
	operators.AddPostfix("+", 105)
	operators.AddPostfix("?", 105)
	operators.AddInfix("-", 100)
	operators.AddInfix("|", 5)

	operatorGrammar := NewOperatorGrammar(context, &operators)
	return operatorGrammar;
}

//  Parses a stream of 'runes' that represent an Regular Expression
//  and return a parse tree/AST
func ParseRegexp(context Context) Value {

	grammar := NewRegexpGrammar(context)

	var result Value = tuple.EMPTY
	for {  // TODO read runes
		v, err := context.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			Error(context, "on reading: %s", err)
			break
		}
		s := string(v)
		switch v {
		case '[', '(', '{': grammar.OpenBracket(Atom{s})
		case ']', ')', '}': grammar.CloseBracket(Atom{s})
		case '-', '|', '*', '+', '?': grammar.PushOperator(Atom{s})
		default: grammar.PushValue(String(s))
		}
	}
	grammar.EndOfInput(func (value Value) {
		result = value
	})

	return result
}



