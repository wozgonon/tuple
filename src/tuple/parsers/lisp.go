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

func parse(context Context, operators Operators, style Style, next Next) error {
	operatorGrammar := NewOperatorGrammar(context, &operators)

	flush := func() bool {
		if context.Location().Depth() == 0 && operatorGrammar.Values.Arity() == 1 && len(operatorGrammar.operatorStack) == 0  {
			err := operatorGrammar.EndOfInput(next)
			if err != nil {
				Error(context, "%s", err)
			}
			return true
		}
		return false
	}
	for {
		err := style.GetNext(context,
			func() {
				flush()
				if context.Location().Depth() == 0 {
					err := operatorGrammar.EndOfInput(next)
					if err != nil {
						Error(context, "%s", err)
					}
				}
			},
			func (open string) {
				flush()
				operatorGrammar.OpenBracket(Tag{open})
			},
			func (close string) {
				flush()
				operatorGrammar.CloseBracket(Tag{close})
			},
			func (tag Tag) {
				flush()
				if operators.Precedence(tag) != -1 {
					operatorGrammar.PushOperator(tag)
				} else {
					operatorGrammar.PushValue(tag)
				}
			},
			func (literal Value) {
				flush()
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

/////////////////////////////////////////////////////////////////////////////
// Lisp with conventional Prefix Grammar
/////////////////////////////////////////////////////////////////////////////

type LispGrammar struct {
	style Style
	operators Operators
	//parser SExpressionParser
}

func (grammar LispGrammar) Name() string {
	return "Lisp"
}

func (grammar LispGrammar) FileSuffix() string {
	return ".l"
}

func (grammar LispGrammar) Parse(context Context, next Next) error {
	return parse(context, grammar.operators, grammar.style, next)
}

func (grammar LispGrammar) Print(object Value, next func(value string)) {
	//PrintExpression(grammar.parser.lexer, "", object, out)
	PrintExpression(&(grammar.operators), "", object, next)
}

func NewLispGrammar() Grammar {
	//style := NewStyle("", "", "  ",
	//	OPEN_BRACKET, CLOSE_BRACKET, "", "", LISP_CONS_OPERATOR,
	//	"", "\n", "true", "false", ';', "")


	style := LispWithInfixStyle()
	style.RecognizeNegative = true
	//return LispGrammar{NewSExpressionParser(style)}

	operators := NewOperators(style)
	operators.AddBracket(OPEN_BRACKET, CLOSE_BRACKET)
	operators.AddInfix(CONS_ATOM.Name, 30)
	operators.AddInfix(LISP_CONS_OPERATOR, 105) // CONS Operator
	operators.AddInfix(SPACE_ATOM.Name, 20)  // TODO space???
	return LispGrammar{style, operators}

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
	return parse(context, grammar.operators, grammar.style, next)
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
	// TODO style.RecognizeNegative = false
	operators := NewOperators(style)
	AddStandardCOperators(&operators)
	operators.AddInfix(CONS_ATOM.Name, 30)
	operators.AddInfix(LISP_CONS_OPERATOR, 105) // CONS Operator
	return LispWithInfixGrammar{style, operators}
}
