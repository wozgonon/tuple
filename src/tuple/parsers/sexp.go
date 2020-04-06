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

// A [S-Expression](https://en.wikipedia.org/wiki/S-expression) or symbolic expression is a very old and general notation.
//
// See https://en.wikipedia.org/wiki/S-expression
//
//  Not quite since these are not CONS cells
//
// A nested structure of scalars (tags and numbers), lists and key-values pairs (called cons cells).
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
	lexer Lexer
	keyValueSeparator string
	//lexer Lexer
}

func NewSExpressionParser(style Style) SExpressionParser {
	return SExpressionParser{style,style.KeyValueSeparator}
}

// Deal with a binary operator  a : b or cons cell a.b
// TODO- should be done best by operator grammar
func (parser SExpressionParser) parserKeyValueOperator(context Context, tuple *Tuple) {
	if tuple.Length() == 0 {
		Error(context,"Unexpected operator '%s'", parser.keyValueSeparator)
		return // errors.New("Unexpected")
	}
	left := tuple.List[tuple.Length()-1]
	Verbose(context,"CONS %s : ... ", left)
	var right Value = nil
	for {
		err := parser.lexer.GetNext(context,
			func() {
				// EOL do nothing 
			},
			func (open string) {
				Verbose(context,"** OPEN")
				tuple1 := NewTuple()
				parser.parseSExpressionTuple(context, &tuple1)
				right = tuple1
			},
			func (close string) {
				Error(context,"Unexpected close bracket '%s'", close)
			},
			func (tag Tag) {
				Verbose(context,"parse tag: %s", tag)
				right = tag
			},
			func (literal Value) {
				right = literal
			})
		if err != nil {
			Verbose(context,"** ERR")
			return // err
		}
		if right == nil {
			Verbose(context,"RIGHT is NIL")
		} else {
			tuple.List[tuple.Length() -1] = NewTuple(CONS_ATOM, left, right)
			return
		}
	}
}

func (parser SExpressionParser) parseSExpressionTuple(context Context, tuple *Tuple) error {

	closeBracketFound := false
	for {
		err := parser.lexer.GetNext(context,
			func() {},
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
			func (tag Tag) {
				if tag.Name == parser.keyValueSeparator {
					parser.parserKeyValueOperator(context, tuple)
				} else {
					tuple.Append(tag)
				}
			},
			func (literal Value) {
				tuple.Append(literal)
			})
		
		if err != nil {
			if err != io.EOF {
				Error(context,"parsing %s", err);
			}
			if ! closeBracketFound {
				Error(context,"Found EOF but expected a close bracket")
			}
			return err /// ??? Any need to return
		}
		if closeBracketFound {
			return nil
		}
	}
}

func (parser SExpressionParser) parse(context Context, next Next) (error) {

	err := parser.lexer.GetNext(context,
		func() {},
		func (open string) {
			tuple := NewTuple()
			parser.parseSExpressionTuple(context, &tuple)
			next(tuple)
		},
		func (close string) {
			Error(context,"Unexpected close bracket '%s'", close)
		},
		func (tag Tag) {
			next(tag)
		},
		func (literal Value) {
			Verbose(context,"parse literal: %s", literal)
			next(literal)
		})
	return err
}

// Reads a given text and produces an Asbstract Syntax Tree (AST)
// See the Grammar interface
func (parser SExpressionParser) Parse(context Context, next Next) {

	for {
		err := parser.parse(context, func (value Value) {
			next(value)
		})
		if err != nil {
			return
		}
	}
}

