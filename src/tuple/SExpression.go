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

func (parser SExpressionParser) getNext(context * ParserContext) (interface{}, error) {

	for {
		ch, err := context.ReadRune()
		switch {
		case err != nil: return "", err
		case err == io.EOF: return "", nil
		case ch == ',' || unicode.IsSpace(ch): break // TODO fix comma
		case ch == parser.style.OneLineComment: return ReadUntilEndOfLine(context)
		case ch == parser.openChar :  return parser.style.Open, nil
		case ch == parser.closeChar : return parser.style.Close, nil
		case ch == parser.openChar2 :  return parser.style.Open2, nil
		case ch == parser.closeChar2 : return parser.style.Close2, nil
		case ch == '"' :  return ReadCLanguageString(context)
		case ch == '.' || unicode.IsNumber(ch): return ReadNumber(context, string(ch))    // TODO minus
		case ch == parser.KeyValueSeparator : return parser.style.KeyValueSeparator, nil
		case IsArithmetic(ch): return ReadAtom(context, string(ch), func(r rune) bool { return IsArithmetic(r) })
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

	//fmt.Printf("parseSExpressionTuple depth=%d, s\n", context.depth)
	
	style := parser.style
	for {
		token, err := parser.getNext(context)
		switch {
		case err != nil:
			context.Error("parsing %s", err);
			return err /// ??? Any need to return
		case token == style.Close:
			//fmt.Printf("*** close=%s\n", token)
			return nil
		case token == style.Close2:
			//fmt.Printf("*** close=%s\n", token)
			return nil
		case token == style.Open || token == style.Open2:
			context.Open()
			subTuple := NewTuple()
			err := parser.parseSExpressionTuple(context, &subTuple)
			context.Close()
			if err == io.EOF {
				context.Error ("Missing close bracket")
				return err
			}
			if err != nil {
				return err
			}
			//fmt.Printf("1. s=%s", subTuple)
			tuple.Append(subTuple)
		case token == style.KeyValueSeparator:  // TODO check if it is an operator
			//fmt.Printf("--------------------\n")
			if tuple.Length() == 0 {
				context.Error("Unexpected operator '%s'", style.KeyValueSeparator)
				return errors.New("Unexpected")
			}
			key := tuple.List[tuple.Length()-1]
			//fmt.Printf("** key=%s\n", key)
			value, err := parser.parse(context)
			if err != nil {
				return err
			}
			if value == style.Close || value == style.Close2 {
				context.Error ("Unexpected close bracket '%s'", token)
				return errors.New("Unexpected")

			}
			//fmt.Printf("depth=%d, key=%s value=%s\n", context.depth, key, value)
			tuple.List[tuple.Length() -1] = NewTuple(Atom{"_cons"}, key, value)
		default:
			if _,ok := token.(Comment); ok {
				// TODO Ignore ???
			} else {
				//fmt.Printf("depth=%d, append=%s\n", context.depth, token)
				tuple.Append(token)
			}
		}
	}
}

func (parser SExpressionParser) parse(context * ParserContext) (interface{}, error) {

	style := parser.style
	token, err := parser.getNext(context)
	switch {
	case err == io.EOF:
		return nil, err
	case err != nil:
		context.Error ("'%s'", err)
		return nil, err
	case token == style.Close:
		if context.depth == 0 {
			context.Error ("Unexpected close bracket '%s'", style.Close)
			return nil, errors.New("Unexpected")
		}
		//context.Close()
		return token, nil
	case token == style.Close2:
		context.Error ("Unexpected close bracket '%s'", style.Close2)
		return nil, errors.New("Unexpected")
	case token == style.Open || token == style.Open2:
		//fmt.Printf("!!!Open\n")
		context.Open()
		tuple := NewTuple()
		err := parser.parseSExpressionTuple(context, &tuple)
		context.Close()
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
		} else if ok && atom.Name == "_cons" {
			parser.printObject(depth, tuple.List[1], out)
			out (" ")
			out(style.KeyValueSeparator)
			out (" ")
			if _, ok = tuple.List[2].(Tuple); ok {
				parser.printObject(depth, tuple.List[2], out)
			} else {
				parser.printScalar(tuple.List[2], out)
			}
			return
		}
		out(style.Open)
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
		out(style.Close)

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

