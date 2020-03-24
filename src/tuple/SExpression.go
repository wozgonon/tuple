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
import "unicode"
import "unicode/utf8"
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
type SExpressionParser struct {
	style Style
	openChar rune
	closeChar rune
	openChar2 rune
	closeChar2 rune
	KeyValueSeparator rune
}

func NewSExpressionParser(style Style) SExpressionParser {

	openChar, _ := utf8.DecodeRuneInString(style.Open)
	closeChar, _ := utf8.DecodeRuneInString(style.Close)
	openChar2, _ := utf8.DecodeRuneInString(style.Open2)
	closeChar2, _ := utf8.DecodeRuneInString(style.Close2)
	KeyValueSeparator, _ := utf8.DecodeRuneInString(style.KeyValueSeparator)
	return SExpressionParser{style,openChar,closeChar,openChar2,closeChar2,KeyValueSeparator}
}

/////////////////////////////////////////////////////////////////////////////
// Lexer
/////////////////////////////////////////////////////////////////////////////

func readRune(context * ParserContext, parser SExpressionParser) (rune, error) {
	ch, err := context.ReadRune()
	switch {
	case err != nil: return ch, err
	case ch == parser.openChar, ch == parser.openChar2  :
		context.Open()
	case ch == parser.closeChar, ch == parser.closeChar2 :
		context.Close()
	}
	return ch, nil

}

func (parser SExpressionParser) GetNext(context * ParserContext, nextAtom func(atom Atom), nextLiteral func (literal interface{})) error {

	style := parser.style
//	for {
		ch, err := readRune(context, parser)
		switch {
		case err != nil: return err
		case err == io.EOF:
			//next.NextEOF()
			return err
		case ch == ',' || unicode.IsSpace(ch): break // TODO fix comma
		case ch == style.OneLineComment:
			_, err = ReadUntilEndOfLine(context)
			if err != nil {
				return err
			}
			// TODO next.NextComment
		case ch == parser.openChar : nextAtom(Atom{style.Open})
		case ch == parser.closeChar : nextAtom(Atom{style.Close})
		case ch == parser.openChar2 : nextAtom(Atom{style.Open2})
		case ch == parser.closeChar2 : nextAtom(Atom{style.Close2})
		//case ch == '+', ch== '*', ch == '-', ch== '/': return string(ch), nil
		case ch == '"' :
			value, err := ReadCLanguageString(context)
			if err != nil {
				return err
			}
			nextLiteral(value)
		case ch == '.' || unicode.IsNumber(ch):
			value, err := ReadNumber(context, string(ch))    // TODO minus
			if err != nil {
				return err
			}
			if atom, ok := value.(Atom); ok {
				nextAtom(atom)
			} else {
				nextLiteral(value)
			}
		case ch == parser.KeyValueSeparator :
			nextAtom(Atom{style.KeyValueSeparator})
		case IsArithmetic(ch): nextAtom(Atom{string(ch)}) // }, nil // ReadAtom(context, string(ch), func(r rune) bool { return IsArithmetic(r) })
		case IsCompare(ch):
			value, err := (ReadAtom(context, string(ch), func(r rune) bool { return IsCompare(r) }))
			if err != nil {
				return err
			}
			if atom, ok := value.(Atom); ok {
				nextAtom(atom)
			} else {
				nextLiteral(value)
			}
		case ch == '_' || unicode.IsLetter(ch):
			value, err :=(ReadAtom(context, string(ch), func(r rune) bool { return r == '_' || unicode.IsLetter(r) || unicode.IsNumber(r) }))
			if err != nil {
				return err
			}
			if atom, ok := value.(Atom); ok {
				nextAtom(atom)
			} else {
				nextLiteral(value)
			}
			
		case unicode.IsGraphic(ch): context.Error("Error graphic character not recognised '%s'", string(ch))
		case unicode.IsControl(ch): context.Error("Error control character not recognised '%d'", ch)
		default: context.Error("Error character not recognised '%d'", ch)
		}
	//}
	return nil
}

/////////////////////////////////////////////////////////////////////////////
//  Parsing
/////////////////////////////////////////////////////////////////////////////


func (parser SExpressionParser) parseSExpressionTuple(context * ParserContext, tuple *Tuple) error {

	style := parser.style
	closeBracketFound := false
	for {
		err := parser.GetNext(context,
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
						err := parser.GetNext(context,
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
	err := parser.GetNext(context,
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
	/*
	switch {
	case err == io.EOF:
		return nil, err
	case err != nil:
		context.Error ("'%s'", err)
		return nil, err
	case token == style.Close, token == style.Close2:
		context.Error ("Unexpected close bracket '%s'", style.Close)
		return nil, errors.New("Unexpected")
		return token, nil
	case token == style.Open || token == style.Open2:
		tuple := NewTuple()
		err := parser.parseSExpressionTuple(context, &tuple)
		if err != nil {
			return nil, err
		}
		return tuple, nil
	default:
		//if _,ok := token.(Comment); ok {
			// TODO Ignore ???
		//} else {
			return token, nil
		//}
	}*/
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

/////////////////////////////////////////////////////////////////////////////
//  Printing
/////////////////////////////////////////////////////////////////////////////

func quote(value string, out func(value string)) {
	out(DOUBLE_QUOTE)
	out(value)   // TODO Escape
	out(DOUBLE_QUOTE)
}

