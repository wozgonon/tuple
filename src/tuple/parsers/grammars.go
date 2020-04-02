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

type Grammars = tuple.Grammars
type Grammar = tuple.Grammar
type Context = tuple.Context
type Atom = tuple.Atom
type Value = tuple.Value
type Comment = tuple.Comment
type StringFunction = tuple.StringFunction
type String = tuple.String
type Tuple = tuple.Tuple
type Next = tuple.Next
type Lexer = tuple.Lexer
type Float64 = tuple.Float64
type Int64 = tuple.Int64

var CONS_ATOM = tuple.CONS_ATOM
var PrintTuple = tuple.PrintTuple
var PrintExpression = tuple.PrintExpression
var PrintExpression1 = tuple.PrintExpression1
var PrintScalar = tuple.PrintScalar
var NewComment = tuple.NewComment
var NewTuple = tuple.NewTuple
//var NewScalar = tuple.NewScalar
var Error = tuple.Error
var Verbose = tuple.Verbose



func AddAllKnownGrammars(grammars * Grammars) {
	grammars.Add(NewLispWithInfixGrammar())
	grammars.Add(NewLispGrammar())
	grammars.Add(NewInfixExpressionGrammar())
	grammars.Add(NewYamlGrammar())
	grammars.Add(NewIniGrammar())
	grammars.Add(NewPropertyGrammar())
	grammars.Add(NewJSONGrammar())
	grammars.Add(NewShellGrammar())
}

func UnexpectedCloseBracketError(context Context, token string) {
	Error(context,"Unexpected close bracket '%s'", token)
}

func UnexpectedEndOfInputErrorBracketError(context Context) {
	Error(context,"Unexpected end of input")
}
