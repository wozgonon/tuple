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
import "tuple/lexer"
import "bufio"
import "strings"

type Grammar = tuple.Grammar
type Context = tuple.Context
type Tag = tuple.Tag
type Value = tuple.Value
type StringFunction = tuple.StringFunction
type String = tuple.String
type Tuple = tuple.Tuple
type Next = tuple.Next
type Lexer = tuple.Lexer
type Float64 = tuple.Float64
type Int64 = tuple.Int64
type Style = lexer.Style

var CONS_ATOM = tuple.CONS_ATOM
var IsAtom = tuple.IsAtom
var Head = tuple.Head
var PrintTuple = tuple.PrintTuple
var PrintExpression = tuple.PrintExpression
var PrintExpression1 = tuple.PrintExpression1
var PrintScalar = tuple.PrintScalar
var NewTuple = tuple.NewTuple
var NewStyle = lexer.NewStyle
//var NewScalar = tuple.NewScalar
var Error = tuple.Error
var Verbose = tuple.Verbose

const OPEN_BRACKET = lexer.OPEN_BRACKET
const CLOSE_BRACKET = lexer.CLOSE_BRACKET
const OPEN_SQUARE_BRACKET = lexer.OPEN_SQUARE_BRACKET
const CLOSE_SQUARE_BRACKET = lexer.CLOSE_SQUARE_BRACKET
const OPEN_BRACE = lexer.OPEN_BRACE
const CLOSE_BRACE = lexer.CLOSE_BRACE
const NEWLINE = lexer.NEWLINE
const DOUBLE_QUOTE = lexer.DOUBLE_QUOTE
const CONS_OPERATOR = lexer.CONS_OPERATOR
var SPACE_ATOM = lexer.SPACE_ATOM


func UnexpectedCloseBracketError(context Context, token string) {
	Error(context,"Unexpected close bracket '%s'", token)
}

func UnexpectedEndOfInputErrorBracketError(context Context) {
	Error(context,"Unexpected end of input")
}

func AddStandardCOperators(operators *Operators) {
	
	operators.AddBracket(OPEN_BRACKET, CLOSE_BRACKET)
	operators.AddBracket(OPEN_SQUARE_BRACKET, CLOSE_SQUARE_BRACKET)
	operators.AddBracket(OPEN_BRACE, CLOSE_BRACE)
	operators.AddPrefix("+", 110)
	operators.AddPrefix("-", 110)
	operators.AddPostfix("++", 105)
	// TODO operators.AddPostfix("%", 105)
	// TODO operators.AddPostfix("--", 105)
	operators.AddInfix("**", 100)
	operators.AddInfix("*", 90)
	operators.AddInfix("/", 90)
	operators.AddInfix("%", 90)
	operators.AddInfix("+", 80)
	operators.AddInfix("-", 80)
	operators.AddInfix("..", 70)  // Range operator
	operators.AddInfix("<", 60)
	operators.AddInfix(">", 60)
	operators.AddInfix("<=", 60)
	operators.AddInfix(">=", 60)
	operators.AddInfix("==", 60)
	operators.AddInfix("!=", 60)
	operators.AddPrefix("!", 55) // TODO check
	operators.AddInfix("|", 55)  // Pipe, what about redirect
	operators.AddInfix("&&", 50)
	operators.AddInfix("||", 50)
	operators.AddInfix("=", 40)
	//operators.AddInfix(",", 30)
	operators.AddInfix(";", 10)
	operators.AddInfix(SPACE_ATOM.Name, 20)  // TODO space???

	//operators.AddInfix3(":", 30, CONS_OPERATOR)

}


// TODO can this be removed
func quote(value string, out func(value string)) {
	out(DOUBLE_QUOTE)
	out(value)   // TODO Escape
	out(DOUBLE_QUOTE)
}


func ParseString(logger LocationLogger, grammar Grammar, expression string) Value {
	var result Value = tuple.EMPTY // TODO Void?
	pipeline := func(value Value) {
		result = value
	}

	reader := bufio.NewReader(strings.NewReader(expression))
	context := NewParserContext("<parse>", reader, logger)
	grammar.Parse(&context, pipeline)
	return result
}


