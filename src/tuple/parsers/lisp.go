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

import "io"

var LISP_CONS_OPERATOR = "."

/////////////////////////////////////////////////////////////////////////////
// Lisp with conventional Prefix Grammar
/////////////////////////////////////////////////////////////////////////////

type LispGrammar struct {
	parser SExpressionParser
}

func (grammar LispGrammar) Name() string {
	return "Lisp"
}

func (grammar LispGrammar) FileSuffix() string {
	return ".l"
}

func (grammar LispGrammar) Parse(context Context, next Next) error {
	return grammar.parser.Parse(context, next)
}

func (grammar LispGrammar) Print(object Value, out func(value string)) {
	PrintExpression(grammar.parser.lexer, "", object, out)
}

func NewLispGrammar() Grammar {
	style := NewStyle("", "", "  ",
		OPEN_BRACKET, CLOSE_BRACKET, "", "", LISP_CONS_OPERATOR,
		"", "\n", "true", "false", ';', "")
	style.RecognizeNegative = true
	return LispGrammar{NewSExpressionParser(style)}
}

/////////////////////////////////////////////////////////////////////////////
// Lisp with an infix notation Grammar
/////////////////////////////////////////////////////////////////////////////

type LispWithInfixGrammar struct {
	style Style
	operators Operators
}

func (grammar LispWithInfixGrammar) Name() string {
	return "Lisp with infix"
}

func (grammar LispWithInfixGrammar) FileSuffix() string {
	return ".infix"
}

func (grammar LispWithInfixGrammar) Parse(context Context, next Next) error {

	operators := grammar.operators
	operatorGrammar := NewOperatorGrammar(context, &operators)
	for {
		err := grammar.style.GetNext(context,
			func() {
				if context.Location().Depth() == 0 {
					err := operatorGrammar.EndOfInput(next)
					if err != nil {
						Error(context, "%s", err)
					}
				}
			},
			func (open string) {
				operatorGrammar.OpenBracket(Tag{open})
			},
			func (close string) {
				operatorGrammar.CloseBracket(Tag{close})
			},
			func (tag Tag) {
				if operators.Precedence(tag) != -1 {
					operatorGrammar.PushOperator(tag)
				} else {
					operatorGrammar.PushValue(tag)
				}
			},
			func (literal Value) {
				operatorGrammar.PushValue(literal)
			})
		
		if err == io.EOF {
			return operatorGrammar.EndOfInput(next)
		}
		if err != nil {
			return err
		}			
	}
}

func (grammar LispWithInfixGrammar) Print(token Value, next func(value string)) {
	PrintExpression(&(grammar.operators), "", token, next)
}

func LispWithInfixStyle () Style {
	style := NewStyle("", "", "  ",
	OPEN_BRACKET, CLOSE_BRACKET, "", "", LISP_CONS_OPERATOR, 
		"", "\n", "true", "false", ';', "")

	// TODO infix should not have this
	style.RecognizeNegative = true
	return style
}

func NewLispWithInfixGrammar() Grammar {
	style := LispWithInfixStyle()
	
	operators := NewOperators(style)
	AddStandardCOperators(&operators)
	operators.AddInfix(LISP_CONS_OPERATOR, 105) // CONS Operator
	return LispWithInfixGrammar{style, operators}
}
