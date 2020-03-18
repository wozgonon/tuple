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

/////////////////////////////////////////////////////////////////////////////
// Pretty printer
/////////////////////////////////////////////////////////////////////////////

type Style struct {
	StartDoc string
	EndDoc string
	Indent string
	Open string
	Close string
	Separator string
	LineBreak string
	True string
	False string
	OneLineComment rune
	ScalarPrefix string
}

type SExpressionParser struct {
	style Style
	openChar rune
	closeChar rune
}

func NewSExpressionParser(style Style) SExpressionParser {

	openChar, _ := utf8.DecodeRuneInString(style.Open)
	closeChar, _ := utf8.DecodeRuneInString(style.Close)
	return SExpressionParser{style,openChar,closeChar}
}

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
		case ch == '"' :  return ReadCLanguageString(context)
		case ch == '.' || unicode.IsNumber(ch): return ReadNumber(context, string(ch))    // TODO minus
		case IsArithmetic(ch): return ReadAtom(context, string(ch), func(r rune) bool { return IsArithmetic(r) })
		case IsCompare(ch): return ReadAtom(context, string(ch), func(r rune) bool { return IsCompare(r) })
		case unicode.IsLetter(ch):  return ReadAtom(context, string(ch), func(r rune) bool { return unicode.IsLetter(r) })
		case unicode.IsGraphic(ch): context.Error("Error graphic character not recognised '%s'", string(ch))
		case unicode.IsControl(ch): context.Error("Error control character not recognised '%d'", ch)
		default: context.Error("Error character not recognised '%d'", ch)
		}
	}
}

func quote(value string, out func(value string)) {
	out(DOUBLE_QUOTE)
	out(value)   // TODO Escape
	out(DOUBLE_QUOTE)
}

func (style Style) indentOnly() bool { // TODO Remove this
       return style.Close == "" && style.Indent != "" // TODO
}

func (style Style) printScalar(token interface{}, out func(value string)) {
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

func (style Style) printToken(depth string, token interface{}, out func(value string)) {
	if tuple, ok := token.(Tuple); ok {

		len := len(tuple.List)
		out(depth)
		if len == 0 {
			out(style.Open)
			out(style.Close)
			return
		}
		head := tuple.List[0]
		atom, ok := head.(Atom)
		first := ok && style.indentOnly()
		if first {
			out(atom.Name)
		}			
		out(style.Open)
		out(style.LineBreak)
		newDepth := depth + style.Indent
		for k, token := range tuple.List {
			style.printToken(newDepth, token, out)
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
		style.printScalar(token, out)
	}
}

func (style Style) PrettyPrint(token interface{}, out func(value string)) {

	style.printToken("", token, out)
	out (string(NEWLINE))
}

