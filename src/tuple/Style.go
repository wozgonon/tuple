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
//import "log"
import "io"
import "unicode"
import "strconv"
import "math"
import "unicode/utf8"

const NEWLINE = '\n'
const DOUBLE_QUOTE = "\""
const UNKNOWN = "<???>"
const WORLD = "世界"

const OPEN_BRACKET = "("
const CLOSE_BRACKET = ")"
const OPEN_SQUARE_BRACKET = "["
const CLOSE_SQUARE_BRACKET = "]"
const OPEN_BRACE = "{"
const CLOSE_BRACE = "}"

// LISP cons operator (https://en.wikipedia.org/wiki/Cons)
const CONS_OPERATOR = "."

var (
	SPACE_ATOM = Atom{" "}
	CONS_ATOM = Atom{"cons"}
)

/////////////////////////////////////////////////////////////////////////////
//  Style
/////////////////////////////////////////////////////////////////////////////

//  Style is a configurable 'lexer' (https://en.wikipedia.org/wiki/Lexical_analysis).
//  It is not a general purpose lexer but can manage most C-like operators and tokens.
//  It implements the 'Lexer' interface
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

	//INF string
	//NAN string

	openChar rune
	closeChar rune
	openChar2 rune
	closeChar2 rune
	KeyValueSeparatorRune rune
}

func NewStyle(
	StartDoc string,
	EndDoc string,
	Indent string,

	Open string,
	Close string,
	Open2 string,
	Close2 string,
	KeyValueSeparator string,
	
	Separator string,
	LineBreak string,
	True string,
	False string,
	OneLineComment rune,
	ScalarPrefix string) Style {

	openChar, _ := utf8.DecodeRuneInString(Open)
	closeChar, _ := utf8.DecodeRuneInString(Close)
	openChar2, _ := utf8.DecodeRuneInString(Open2)
	closeChar2, _ := utf8.DecodeRuneInString(Close2)
	KeyValueSeparatorRune, _ := utf8.DecodeRuneInString(KeyValueSeparator)

	return Style{StartDoc,EndDoc,Indent, Open,Close,Open2,Close2,KeyValueSeparator,Separator,LineBreak,True,False,OneLineComment,ScalarPrefix,
		openChar,closeChar,openChar2,closeChar2,KeyValueSeparatorRune}
}

/////////////////////////////////////////////////////////////////////////////
// Lexer
/////////////////////////////////////////////////////////////////////////////

func readRune(context Context, style Style) (rune, error) {
	ch, err := context.ReadRune()
	switch {
	case err != nil: return ch, err
	case ch == style.openChar, ch == style.openChar2  :
		context.Open()
	case ch == style.closeChar, ch == style.closeChar2 :
		context.Close()
	} 
	return ch, nil
}

func (style Style) GetNext(context Context, open func(open string), close func(close string), nextAtom func(atom Atom), nextLiteral func (literal interface{})) error {

	ch, err := readRune(context, style)
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
	case ch == style.openChar : open(style.Open)
	case ch == style.closeChar : close(style.Close)
	case ch == style.openChar2 : open(style.Open2)
	case ch == style.closeChar2 : close(style.Close2)
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
	case ch == style.KeyValueSeparatorRune:
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
		
	case unicode.IsGraphic(ch): Error(context,"Error graphic character not recognised '%s'", string(ch))
	case unicode.IsControl(ch): Error(context,"Error control character not recognised '%d'", ch)
	default: Error(context,"Error character not recognised '%d'", ch)
	}
	return nil
}

/////////////////////////////////////////////////////////////////////////////

// "Code point" Go introduces a shorter term for the concept: rune and means exactly the same as "code point", with one interesting addition.
//  Go language defines the word rune as an alias for the type int32, so programs can be clear when an integer value represents a code point.
//  Moreover, what you might think of as a character constant is called a rune constant in Go. 

func ReadString (context Context, token string, unReadLast bool, test func(r rune) bool) (string, error) {
	for {
		ch, err := context.ReadRune()
		if err == io.EOF {
			//Error(context,"ERROR missing close quote: '%s'", DOUBLE_QUOTE)
			return token, nil
		} else if err != nil {
			//log.Printf("ERROR nil")
			//return ""
		} else if ! test(ch) {
			if unReadLast {
				context.UnreadRune()
			}
			return token, nil
		} else {
			// TODO not efficient
			token = token + string(ch)
		}
	}
}

func ReadAtom(context Context, prefix string, test func(rune) bool) (interface{}, error) {
	atom, err := ReadString(context, prefix, true, test)
	if err != nil {
		return Atom{""}, err
	}
	switch atom {
	case "NaN": return math.NaN(), nil
	case "Inf": return math.Inf(1), nil // TODO "+Inf", and "-Inf" 
	default: return Atom{atom}, err
	}
}

func ReadNumber(context Context, token string) (interface{}, error) {  // Number
	var dots int
	if token == "." {
		dots = 1
	} else {
		dots = 0
	}
	for {
		ch, err := context.ReadRune()
		if err == io.EOF {
			break
		} else if err != nil {
			return "", err
		} else if ch == '.' && dots == 0 {
			dots += 1
			token = token + "." // TODO not efficient
		} else if unicode.IsNumber(ch) {
			// TODO ought to be much more efficient to build up a number dynamically
			token = token + string(ch) // TODO not efficient
		} else {
			context.UnreadRune()
			//if token == "." {
			//	context.UnreadRune()
			//}
			break
		}
	}
	if token == "." {
		return Atom{"."}, nil
	}
	//return Number{dots=true,token}
	switch dots {
	case 0: return strconv.ParseInt(token, 10, 0)
	default: return strconv.ParseFloat(token, 64)
	} 
}

func ReadUntilEndOfLine(context Context) (Comment, error) {
	token := ""
	for {
		ch, err := context.ReadRune()
		switch {
		case err == io.EOF:
			return NewComment(context, token), nil
		case err != nil:
			return NewComment(context, token), err
		case ch == NEWLINE:
			context.UnreadRune()
			return NewComment(context, token), err
		default:
			token = token + string(ch)
		}
	}
}

func ReadUntilSpace(context Context, token string) (string, error) {
	for {
		ch, err := context.ReadRune()
		switch {
		case err == io.EOF:
			return token, nil
		case err != nil:
			return token, err
		case unicode.IsSpace(ch), ch == NEWLINE:
			context.UnreadRune()
			return token, nil
		default:
			token = token + string(ch)
		}
	}
}

func ReadCLanguageString(context Context) (string, error) {
	token := ""
	for {
		ch, err := context.ReadRune()
		switch {
		case err == io.EOF:
			Error(context,"ERROR missing close quote: '%s'", DOUBLE_QUOTE)
			return token, nil
		case err != nil: return "", err
		case ch == '"': return token, nil
		case ch == '\\':
			ch, err := context.ReadRune()
			if err == io.EOF {
				Error(context,"ERROR missing close quote: '%s'", DOUBLE_QUOTE)
				return token, nil
			}
			token = token + string(cLanguageEscapeCharacters(ch))
		default:
			// TODO not efficient
			token = token + string(ch)
		}
	}
}

func cLanguageEscapeCharacters(ch rune) rune {
	switch ch {
	case 'n': return NEWLINE
	case 'r': return '\r'
	case 't': return '\t'
	// TODO 
	default:
		return ch;
	}
}

func IsArithmetic(ch rune) bool {
	switch ch {
		case '+': return true
		case '-': return true
		case '/': return true
		case '*': return true
		case '^': return true
		default: return false
	}
}

func IsCompare(ch rune) bool {
	switch ch {
		case '=': return true
		case '!': return true
		case '<': return true
		case '>': return true
		default: return false
	}
}

func AddStandardCOperators(operators *Operators) {
	operators.unary["-"] = Atom{"_unary_minus"}
	operators.unary["+"] = Atom{"_unary_plus"}
	operators.AddBracket(OPEN_BRACKET, CLOSE_BRACKET)
	operators.AddBracket(OPEN_SQUARE_BRACKET, CLOSE_SQUARE_BRACKET)
	operators.AddBracket(OPEN_BRACE, CLOSE_BRACE)
	operators.Add("_unary_plus", 110)
	operators.Add("_unary_minus", 110)
	operators.Add("^", 100)
	operators.Add("*", 90)
	operators.Add("/", 90)
	operators.Add("+", 80)
	operators.Add("-", 80)
	operators.Add("<", 60)
	operators.Add(">", 60)
	operators.Add("<=", 60)
	operators.Add(">=", 60)
	operators.Add("==", 60)
	operators.Add("!=", 60)
	operators.Add("&&", 50)
	operators.Add("||", 50)
	//operators.Add(",", 40)
	//operators.Add(";", 30)
	operators.Add(SPACE_ATOM.Name, 10)  // TODO space???
}

/////////////////////////////////////////////////////////////////////////////
// Printer
/////////////////////////////////////////////////////////////////////////////

func (printer Style) PrintIndent(depth string, out StringFunction) {
	out(depth)
}

func (printer Style) PrintSuffix(depth string, out StringFunction) {
	out(string(NEWLINE))
}

func (printer Style) PrintSeparator(depth string, out StringFunction) {
	//out(printer.Separator)
}

func (printer Style) PrintEmptyTuple(depth string, out StringFunction) {
	out(printer.Open)
	out(printer.Close)
}

func (printer Style) PrintNullaryOperator(depth string, atom Atom, out StringFunction) {
	PrintTuple(&printer, depth, NewTuple(atom), out)
}

func (printer Style) PrintUnaryOperator(depth string, atom Atom, value interface{}, out StringFunction) {
	PrintTuple(&printer, depth, NewTuple(atom, value), out)
}

func (printer Style) PrintBinaryOperator(depth string, atom Atom, value1 interface{}, value2 interface{}, out StringFunction) {
	PrintTuple(&printer, depth, NewTuple(atom, value1, value2), out)
}

func isCons(tuple Tuple) bool {
	cons := false
	if tuple.Length() > 0 {
		t, ok := tuple.List[0].(Tuple)
		cons = ok && t.IsCons()
	}
	return cons
}

func (printer Style) PrintOpenTuple(depth string, tuple Tuple, out StringFunction) string {
	if isCons(tuple) {
		out(printer.Open2)
	} else {
		out(printer.Open)
	}
	return depth + "  "
}

func (printer Style) PrintCloseTuple(depth string, tuple Tuple, out StringFunction) {
	printer.PrintIndent(depth, out)
	if isCons(tuple) {
		out(printer.Close2)
	} else {
		out(printer.Close)
	}
}

func (printer Style) PrintHeadAtom(atom Atom, out StringFunction) {
	out(atom.Name)
}

func (printer Style) PrintAtom(depth string, atom Atom, out StringFunction) {
	out(atom.Name)
}

func (printer Style) PrintInt64(depth string, value int64, out StringFunction) {
	out(strconv.FormatInt(value, 10))
}

func (printer Style) PrintFloat64(depth string, value float64, out StringFunction) {
	if math.IsInf(value, 64) {
		out("Inf") // style.INF)  // Do not print +Inf
	} else {
		out(fmt.Sprint(value))
	}
}

// TODO can this be removed
func quote(value string, out func(value string)) {
	out(DOUBLE_QUOTE)
	out(value)   // TODO Escape
	out(DOUBLE_QUOTE)
}


func (printer Style) PrintString(depth string, value string, out StringFunction) {
	out(DOUBLE_QUOTE)
	out(value)   // TODO Escape
	out(DOUBLE_QUOTE)
}

func (printer Style) PrintBool(depth string, value bool, out StringFunction) {
	if value {
		out(printer.True)
	} else {
		out(printer.False)
	}				
}

func (printer Style) PrintComment(depth string, value Comment, out StringFunction) {
	out(string(printer.OneLineComment))
	out(value.Comment)
}

func (printer Style) PrintScalarPrefix(depth string, out StringFunction) {}
