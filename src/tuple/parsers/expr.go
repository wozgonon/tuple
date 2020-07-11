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
//import "fmt"

var EXPR_CONS_OPERATOR = ":"

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

func handleTag(tag Tag, style Style, context Context, operatorGrammar * OperatorGrammar) error {
	open := style.Open
	operators := operatorGrammar.operators

	//fmt.Printf("tag=%s precedence %d\n", tag, operators.Precedence(tag))
	if operators.Precedence(tag) != -1 {
		operatorGrammar.PushOperator(tag)
	} else {
		ch := context.LookAhead()
		if ch == style.OpenChar { //  || ch == style.openChar2
			_, err := context.ReadRune()
			if err == io.EOF {
				operatorGrammar.PushValue(tag)
				return nil
			}
			if err != nil {
				return err
			}
			context.Open()
			operatorGrammar.OpenBracket(Tag{open})
		}
		operatorGrammar.PushValue(tag)
	}
	return nil
}

func (grammar InfixExpressionGrammar) Parse(context Context, next Next) error {

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
				err := operatorGrammar.CloseBracket(Tag{close})
				if err != nil {
					// TODO handle error
				}
			},
			func(tag Tag) {
				err := handleTag(tag, grammar.style, context, &operatorGrammar)
				if err != nil {
					// TODO handle error
				}
			},
			func (literal Value) {
				operatorGrammar.PushValue(literal)  // TODO WithoutInsertingMissingSepator
			})
	
		if err == io.EOF {
			return operatorGrammar.EndOfInput(next)
		}
		if err != nil {
			return err
		}
	}
}

func (grammar InfixExpressionGrammar) Print(token Value, next func(value string)) {
	PrintExpression(&(grammar.operators), "", token, next)
}

func NewInfixExpressionGrammar() Grammar {
	style := NewStyle("", "", "  ",
		OPEN_BRACKET, CLOSE_BRACKET, OPEN_BRACE, CLOSE_BRACE, EXPR_CONS_OPERATOR,
		",", "\n", "true", "false", '%', "") // prolog, sql '--' for 

	operators := NewOperators(style)
	AddStandardCOperators(&operators)
	operators.AddInfix(CONS_ATOM.Name, 30)

	return InfixExpressionGrammar{style, operators}
}

/////////////////////////////////////////////////////////////////////////////
// The basic syntax of command shell parsers is typically very simple:
// everything is a line of strings separated by spaces terminate by a newline.
// Mostly no need for quotes or double quotes or semi-colons.
// This makes it very easy to type command with a few parameters on a command line interface (CLI).
//
// Examples include [TCL](https://en.wikipedia.org/wiki/Tcl), [Bash](https://en.wikipedia.org/wiki/Bash_(Unix_shell)), DOS cmd shell
// TCL and bash use braces { ... } for nesting.
//
// Bash and DOS of course have lots of extra syntax for working with files but the
// basic syntax typically does not understand arithmetic, with infix notation, one has to use a special tool:
// * TCL has a 'expr' function that understand arithmetic with infix notation
// * Bash one can use an external 'expr(1)' tool:  'expr 8.3 + 6'
// * DOS has a special version of the SET command 'SET /a c=a+b'
//
// It would be nice not to have to put quotes around strings, in particular working interactice on the command line (CLI) with a shell such as bash or
// dos command line.  It is very nice not to have to include any boiler plate when just entering commands but this quickly becomes very awkward to add any more complex syntax.
/////////////////////////////////////////////////////////////////////////////

type ShellGrammar struct {
	style Style
	operators Operators
}

func (grammar ShellGrammar) Name() string {
	return "Shell Expression with infix operators"
}

func (grammar ShellGrammar) FileSuffix() string {
	return ".wsh"
}

func (grammar ShellGrammar) Parse(context Context, next Next) error {

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
				} else if ! operatorGrammar.wasOperator {
					operatorGrammar.PushOperator(Tag{";"})
				}
			},
			func (open string) {
				operatorGrammar.OpenBracket(Tag{open})
			},
			func (close string) {
				err := operatorGrammar.CloseBracket(Tag{close})
				if err != nil {
					// TODO handle error
				}
			},
			func(tag Tag) {
				handleTag(tag, grammar.style, context, &operatorGrammar)
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

func (grammar ShellGrammar) Print(token Value, next func(value string)) {
	PrintExpression(&(grammar.operators), "", token, next)
}

func NewShellGrammar() Grammar {
	style := NewStyle("", "", "  ",
		OPEN_BRACKET, CLOSE_BRACKET, OPEN_BRACE, CLOSE_BRACE, EXPR_CONS_OPERATOR,
		",", "\n", "true", "false", '#', "")

	// **
	// /*
	// +-
	// == <= >=     < >
	// ' '
	// &
	// | > >>
	// && ||
	// ;
	
	operators := NewOperators(style)
	AddStandardCOperators(&operators)
	operators.AddInfix(CONS_ATOM.Name, 30)
	operators.AddPrefix("$", 150)
	operators.AddPostfix("&", 20)

	return ShellGrammar{style, operators}
}

