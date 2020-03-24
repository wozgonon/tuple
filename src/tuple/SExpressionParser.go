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
//import "unicode"
//import "unicode/utf8"
//import "errors"

// A [S-Expression](https://en.wikipedia.org/wiki/S-expression) or symbolic expression is a very old and general notation.
//
// See https://en.wikipedia.org/wiki/S-expression
//
//  Not quite since these are not CONS cells
//
// A nested structure of scalars (atoms and numbers), lists and key-values pairs (called cons cells).
// These are used for the syntax of LISP but also any other language can typically be converted to an S-Expression,
// it is in particular a very useful format for debugging a parser by printing out the Abstract Syntaxt Tree (AST) created by parsing.
//
//  TODO Cons cells
//  https://www.gnu.org/software/emacs/manual/html_node/elisp/Dotted-Pair-Notation.html#Dotted-Pair-Notation
//

/////////////////////////////////////////////////////////////////////////////
//  Parsing
/////////////////////////////////////////////////////////////////////////////

type SExpressionParser struct {
	style Style
	//lexer Lexer
}

func NewSExpressionParser(lexer Style) SExpressionParser {
	return SExpressionParser{lexer}
}

func (parser SExpressionParser) parseSExpressionTuple(context * ParserContext, tuple *Tuple) error {

	style := parser.style // lexer
	closeBracketFound := false
	for {
		err := style.GetNext(context,
			func (atom Atom) {
				token := atom.Name
				switch {
				case token == style.Close, token == style.Close2:
					closeBracketFound = true
					return
				case token == style.Open || token == style.Open2:
					subTuple := NewTuple()
					err := parser.parseSExpressionTuple(context, &subTuple)
					if err == io.EOF {
						// TODO context.Error ("Found EOF but expected a bracket to close '%s'", token)
						tuple.Append(subTuple)
						//return err
						return
					}
					if err != nil {
						return // err
					}
					tuple.Append(subTuple)
				case token == style.KeyValueSeparator:  // TODO check if it is an operator
					if tuple.Length() == 0 {
						context.Error("Unexpected operator '%s'", style.KeyValueSeparator)
						return // errors.New("Unexpected")
					}
					left := tuple.List[tuple.Length()-1]
					context.Verbose("CONS %s : ... ", left)
					var right interface{} = nil
					for {
						err := style.GetNext(context,
							func (atom Atom) {
								context.Verbose("parse atom: %s", atom)
								token := atom.Name
								switch {
								case token == style.Close, token == style.Close2:
									context.Error ("Unexpected close bracket '%s'", token)
									return
									//return nil, errors.New("Unexpected")
								case token == style.Open || token == style.Open2:
									context.Verbose("** OPEN")
									tuple1 := NewTuple()
									parser.parseSExpressionTuple(context, &tuple1)
									right = tuple1
									return
								default:
									right = atom
									return
								}
							},
							func (literal interface{}) {
								right = literal
							})
						if err != nil {
							context.Verbose("** ERR")
							return // err
						}
						if right == nil {
							context.Verbose("RIGHT is NIL")
						} else {
							tuple.List[tuple.Length() -1] = NewTuple(CONS_ATOM, left, right)
							return
						}
					}
				default:
					tuple.Append(atom)
				}
			},
			func (literal interface{}) {
				tuple.Append(literal)
			})
		
		if err != nil {
			if err != io.EOF {
				context.Error("parsing %s", err);
			}
			if ! closeBracketFound {
				context.Error ("Found EOF but expected a close bracket")
			}
			return err /// ??? Any need to return
		}
		if closeBracketFound {
			return nil
		}
	}
}

func (parser SExpressionParser) parse(context * ParserContext, next Next) (error) {

	context.Verbose("*** parse:")
	style := parser.style
	err := style.GetNext(context,
		func (atom Atom) {
			context.Verbose("parse atom: %s", atom)
			token := atom.Name
			switch {
			case token == style.Close, token == style.Close2:
				context.Error ("Unexpected close bracket '%s'", token)
				//return nil, errors.New("Unexpected")
			case token == style.Open || token == style.Open2:
				tuple := NewTuple()
				parser.parseSExpressionTuple(context, &tuple)
				next(tuple)
			default:
				next(atom)
			}
		},
		func (literal interface{}) {
			context.Verbose("parse literal: %s", literal)
			next(literal)
		})
	context.Verbose("*** parse: err=%s", err)
	return err
}

// Reads a given text and produces an Asbstract Syntax Tree (AST)
// See the Grammar interface
func (parser SExpressionParser) Parse(context * ParserContext) {

	for {
		err := parser.parse(context, func (value interface{}) {
			context.next(value)
		})
		if err != nil {
			return
		}
	}
}

