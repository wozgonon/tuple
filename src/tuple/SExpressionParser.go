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

// Deal with a binary operator  a : b or cons cell a.b
// TODO- should be done best by operator grammar
func (parser SExpressionParser) parserKeyValueOperator(context * ParserContext, tuple *Tuple) {
	if tuple.Length() == 0 {
		context.Error("Unexpected operator '%s'", parser.style.KeyValueSeparator)
		return // errors.New("Unexpected")
	}
	left := tuple.List[tuple.Length()-1]
	context.Verbose("CONS %s : ... ", left)
	var right interface{} = nil
	for {
		err := parser.style.GetNext(context,
			func (open string) {
				context.Verbose("** OPEN")
				tuple1 := NewTuple()
				parser.parseSExpressionTuple(context, &tuple1)
				right = tuple1
			},
			func (close string) {
				context.Error ("Unexpected close bracket '%s'", close)
			},
			func (atom Atom) {
				context.Verbose("parse atom: %s", atom)
				right = atom
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
}

func (parser SExpressionParser) parseSExpressionTuple(context * ParserContext, tuple *Tuple) error {

	closeBracketFound := false
	for {
		err := parser.style.GetNext(context,
			func (open string) {
				subTuple := NewTuple()
				err := parser.parseSExpressionTuple(context, &subTuple)
				if err == io.EOF {
					tuple.Append(subTuple)
					return
				}
				if err != nil {
					return
				}
				tuple.Append(subTuple)
			},
			func (close string) {
				closeBracketFound = true
			},
			func (atom Atom) {
				if atom.Name == parser.style.KeyValueSeparator {
					parser.parserKeyValueOperator(context, tuple)
				} else {
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

	err := parser.style.GetNext(context,
		func (open string) {
			tuple := NewTuple()
			parser.parseSExpressionTuple(context, &tuple)
			next(tuple)
		},
		func (close string) {
			context.Error ("Unexpected close bracket '%s'", close)
		},
		func (atom Atom) {
			next(atom)
		},
		func (literal interface{}) {
			context.Verbose("parse literal: %s", literal)
			next(literal)
		})
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

