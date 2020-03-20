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

import "fmt"
import "log"
import "strconv"
import "reflect"
import "math"
import "io"
import "unicode"
import "unicode/utf8"
import "errors"

/////////////////////////////////////////////////////////////////////////////
//
/////////////////////////////////////////////////////////////////////////////

type Style struct {
	StartDoc string
	EndDoc string
	Indent string

	Open string
	Close string
	Open2 string
	Close2 string
	KeyValueSeparator string
	
	Separator string
	LineBreak string
	True string
	False string
	OneLineComment rune
	ScalarPrefix string
}

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

func (parser SExpressionParser) GetNext(context * ParserContext) (interface{}, error) {

	for {
		ch, err := readRune(context, parser)
		switch {
		case err != nil: return "", err
		case err == io.EOF: return "", nil
		case ch == ',' || unicode.IsSpace(ch): break // TODO fix comma
		case ch == parser.style.OneLineComment:
			_, err = ReadUntilEndOfLine(context)
			if err != nil {
				return nil, err
			}
		case ch == parser.openChar : return parser.style.Open, nil
		case ch == parser.closeChar : return parser.style.Close, nil
		case ch == parser.openChar2 : return parser.style.Open2, nil
		case ch == parser.closeChar2 : return parser.style.Close2, nil
		//case ch == '+', ch== '*', ch == '-', ch== '/': return string(ch), nil
		case ch == '"' :  return ReadCLanguageString(context)
		case ch == '.' || unicode.IsNumber(ch): return ReadNumber(context, string(ch))    // TODO minus
		case ch == parser.KeyValueSeparator : return parser.style.KeyValueSeparator, nil
		case IsArithmetic(ch): return Atom{string(ch)}, nil // ReadAtom(context, string(ch), func(r rune) bool { return IsArithmetic(r) })
		case IsCompare(ch): return ReadAtom(context, string(ch), func(r rune) bool { return IsCompare(r) })
		case unicode.IsLetter(ch):  return ReadAtom(context, string(ch), func(r rune) bool { return unicode.IsLetter(r) })
		case unicode.IsGraphic(ch): context.Error("Error graphic character not recognised '%s'", string(ch))
		case unicode.IsControl(ch): context.Error("Error control character not recognised '%d'", ch)
		default: context.Error("Error character not recognised '%d'", ch)
		}
	}
}

/////////////////////////////////////////////////////////////////////////////
//  Parsing
/////////////////////////////////////////////////////////////////////////////

func (parser SExpressionParser) parseSExpressionTuple(context * ParserContext, tuple *Tuple) (error) {

	style := parser.style
	for {
		token, err := parser.GetNext(context)
		switch {
		case err != nil:
			context.Error("parsing %s", err);
			return err /// ??? Any need to return
		case token == style.Close, token == style.Close2:
			return nil
		case token == style.Open || token == style.Open2:
			subTuple := NewTuple()
			err := parser.parseSExpressionTuple(context, &subTuple)
			if err == io.EOF {
				context.Error ("Missing close bracket")
				return err
			}
			if err != nil {
				return err
			}
			tuple.Append(subTuple)
		case token == style.KeyValueSeparator:  // TODO check if it is an operator
			if tuple.Length() == 0 {
				context.Error("Unexpected operator '%s'", style.KeyValueSeparator)
				return errors.New("Unexpected")
			}
			left := tuple.List[tuple.Length()-1]
			right, err := parser.parse(context)
			if err != nil {
				return err
			}
			if right == style.Close || right == style.Close2 {
				context.Error ("Unexpected close bracket '%s'", token)
				return errors.New("Unexpected")

			}
			tuple.List[tuple.Length() -1] = NewTuple(Atom{"_cons"}, left, right)
		default:
			tuple.Append(token)
		}
	}
}

func (parser SExpressionParser) parse(context * ParserContext) (interface{}, error) {

	style := parser.style
	token, err := parser.GetNext(context)
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
	}
}

// Reads a given text and produces an Asbstract Syntax Tree (AST)
// See the Grammar interface
func (parser SExpressionParser) Parse(context * ParserContext) {

	for {
		value, err := parser.parse(context)
		if err == nil {
			context.next(value)
		} else {
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

func (style Style) indentOnly() bool { // TODO Remove this
       return style.Close == "" && style.Indent != "" // TODO
}

func (parser SExpressionParser) printScalar(token interface{}, out func(value string)) {
	style := parser.style
	switch token.(type) {
	case Atom:
		if style.indentOnly() {
			quote(token.(Atom).Name, out)
		} else {
			out(token.(Atom).Name)
		}
	case string:
		quote(token.(string), out)   // TODO Escape
	case bool:
		if token.(bool) {
			out(style.True)
		} else {
			out(style.False)
		}				
	case Comment:
		out(string(style.OneLineComment))
		out(token.(Comment).Comment)
	case int64:
		out(strconv.FormatInt(int64(token.(int64)), 10))
	case float64:
		float := token.(float64)
		if math.IsInf(float, 64) {
			out("Inf")  // Do not print +Inf
		} else {
			out(fmt.Sprint(token.(float64)))
		}
	default:
		log.Printf("ERROR type '%s' not recognised: %s", reflect.TypeOf(token), token);
		out(UNKNOWN)
	}
}

func (parser SExpressionParser) isCons(tuple Tuple) bool {
	if tuple.Length() > 0 {
		head := tuple.List[0]
		atom, ok := head.(Atom)
		if ok && atom.Name == "_cons" {
			return true
		}
	}
	return false
}

func (parser SExpressionParser) printObject(depth string, token interface{}, out func(value string)) {

	style := parser.style
	if tuple, ok := token.(Tuple); ok {

		len := len(tuple.List)
		out(depth)
		if len == 0 {
			out(style.Open)
			out(style.Close)
			return
		}
		newDepth := depth + style.Indent
		head := tuple.List[0]
		atom, ok := head.(Atom)
		first := ok && style.indentOnly()
		if first {
			out(atom.Name)
		} else if parser.isCons(tuple) {
			parser.printObject(depth, tuple.List[1], out)
			if _, ok = tuple.List[2].(Tuple); ok {
				out (" ")
				out(style.KeyValueSeparator)
				out (style.LineBreak)
				parser.printObject(newDepth, tuple.List[2], out)
			} else {
				out (" ")
				out(style.KeyValueSeparator)
				out (" ")
				parser.printScalar(tuple.List[2], out)
			}
			return
		}
		tuple1, ok := tuple.List[0].(Tuple)
		cons := ok && parser.isCons(tuple1)
		if cons {
			// TODO Need a way to differentiate between [] and {}
			out(style.Open2)
		} else {
			out(style.Open)
		}
		out(style.LineBreak)
		for k, token := range tuple.List {
			parser.printObject(newDepth, token, out)
			if ! first && k < len-1 {
				out(style.Separator)
				out(style.LineBreak)
			}
		}
		out(style.LineBreak)
		out(depth)
		if cons {
			out(style.Close2)
		} else {
			out(style.Close)
		}
	} else {
		out(depth)
		out(style.ScalarPrefix)
		parser.printScalar(token, out)
	}
}

/////////////////////////////////////////////////////////////////////////////

//  Converts the given object into a text string that can be parsed as an SExpression
//  See the Grammar interface.
func (parser SExpressionParser) Print(object interface{}, out func(value string)) {

	parser.printObject("", object, out)
	out (string(NEWLINE))
}

