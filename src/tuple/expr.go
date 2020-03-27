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

import "io"

/////////////////////////////////////////////////////////////////////////////
//
//  
//
// Conventional arithmetic expression grammar,
// with an infix notation and conventional function call syntax a(b, c, ...)
//
//  Also understands blocks of expressions surrounded by braces, for example:
//
// func abc {
//     1+2
//     3+4
//  }
/////////////////////////////////////////////////////////////////////////////

type InfixExpressionGrammar struct {
	style Style
	operators Operators
}

func (grammar InfixExpressionGrammar) Name() string {
	return "Expression with infix operators"
}

func (grammar InfixExpressionGrammar) FileSuffix() string {
	return ".expr"
}

func handleAtom(atom Atom, style Style, context Context, operatorGrammar * OperatorGrammar) {
	open := style.Open
	operators := (*operatorGrammar).operators
	if operators.Precedence(atom) != -1 {
		operatorGrammar.PushOperator(atom)
	} else {
		ch := context.LookAhead()
		if ch == style.openChar {
			_, err := context.ReadRune()
			if err == io.EOF {
				(*operatorGrammar).PushValue(atom)
				//operatorGrammar.EOF(next)
				return // TODO eof
			}
			if err != nil {
				// TODO
				return
			}
			(*operatorGrammar).OpenBracket(Atom{open})
		}
		(*operatorGrammar).PushValue(atom)
	}
}

func (grammar InfixExpressionGrammar) Parse(context Context, next Next) {

	operators := grammar.operators
	operatorGrammar := NewOperatorGrammar(context, &operators)
	for {
		err := grammar.style.GetNext(context,
			func() {
				operatorGrammar.EndOfInput(next)
			},
			func (open string) {
				operatorGrammar.OpenBracket(Atom{open})
			},
			func (close string) {
				operatorGrammar.CloseBracket(Atom{close})
			},
			func(atom Atom) {
				handleAtom(atom, grammar.style, context, &operatorGrammar)
			},
			func (literal Value) {
				operatorGrammar.PushValue(literal)
			})
	
		if err == io.EOF {
			operatorGrammar.EndOfInput(next)
			break
		}
		if err != nil {
			// TODO 
			break
		}
	}
}

func (grammar InfixExpressionGrammar) Print(token Value, next func(value string)) {
	PrintExpression(&(grammar.operators), "", token, next)
}

func NewInfixExpressionGrammar() Grammar {
	style := NewStyle("", "", "  ",
		OPEN_BRACKET, CLOSE_BRACKET, OPEN_BRACE, CLOSE_BRACE, ":",
		",", "\n", "true", "false", '%', "") // prolog, sql '--' for 

	operators := NewOperators(style)
	AddStandardCOperators(&operators)
	operators.Add(CONS_OPERATOR, 105) // CONS Operator

	return InfixExpressionGrammar{style, operators}
}

/////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////

type ShellGrammar struct {
	style Style
	operators Operators
}

func (grammar ShellGrammar) Name() string {
	return "Shell Expression with infix operators"
}

func (grammar ShellGrammar) FileSuffix() string {
	return ".sh"
}

func (grammar ShellGrammar) Parse(context Context, next Next) {

	operators := grammar.operators
	operatorGrammar := NewOperatorGrammar(context, &operators)
	for {
		err := grammar.style.GetNext(context,
			func() {
				operatorGrammar.EndOfInput(next)
			},
			func (open string) {
				operatorGrammar.OpenBracket(Atom{open})
				//context.Open()
			},
			func (close string) {
				//context.Close()
				operatorGrammar.CloseBracket(Atom{close})
			},
			func(atom Atom) {
				handleAtom(atom, grammar.style, context, &operatorGrammar)
			},
			func (literal Value) {
				operatorGrammar.PushValue(literal)
			})
	
		if err == io.EOF {
			operatorGrammar.EndOfInput(next)
			break
		}
		if err != nil {
			// TODO 
			break
		}
	}
}

func (grammar ShellGrammar) Print(token Value, next func(value string)) {
	PrintExpression(&(grammar.operators), "", token, next)
}

func NewShellGrammar() Grammar {
	style := NewStyle("", "", "  ",
		OPEN_BRACKET, CLOSE_BRACKET, OPEN_BRACE, CLOSE_BRACE, ":",
		",", "\n", "true", "false", '%', "") // prolog, sql '--' for 

	operators := NewOperators(style)
	AddStandardCOperators(&operators)
	operators.Add(CONS_OPERATOR, 105) // CONS Operator

	return ShellGrammar{style, operators}
}

