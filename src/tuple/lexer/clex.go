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
package lexer

import "io"
import "unicode"
import "strconv"
import "math"
import "unicode/utf8"
import "tuple"

type Grammar = tuple.Grammar
type Context = tuple.Context
type Tag = tuple.Tag
type Value = tuple.Value
type StringFunction = tuple.StringFunction
type String = tuple.String
type Tuple = tuple.Tuple
type Next = tuple.Next
type Lexer = tuple.Lexer
type Float64 = tuple.Float64
type Int64 = tuple.Int64

var PrintTuple = tuple.PrintTuple
var NewTuple = tuple.NewTuple
var Error = tuple.Error
var Verbose = tuple.Verbose

/////////////////////////////////////////////////////////////////////////////
//  A lexer similar to that used by UNIX/C based languages
//  such as C, C#, C++, Java, Go and also bash
//
//  TODO This is a first attempt to get something running, it can be much improed
/////////////////////////////////////////////////////////////////////////////

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

var SPACE_ATOM = Tag{" "}

/////////////////////////////////////////////////////////////////////////////
//  Style
/////////////////////////////////////////////////////////////////////////////

//  Style is a configurable 'lexer' (https://en.wikipedia.org/wiki/Lexical_analysis).
//  It is not a general purpose lexer but can manage most C-like operators and tokens.
//  Implements the 'Lexer' interface
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

	OpenChar rune
	closeChar rune
	OpenChar2 rune
	closeChar2 rune
	KeyValueSeparatorRune rune

	RecognizeNegative bool
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
	True string,  // Not used
	False string, // Not used
	OneLineComment rune,
	ScalarPrefix string) Style {

	openChar, _ := utf8.DecodeRuneInString(Open)
	closeChar, _ := utf8.DecodeRuneInString(Close)
	openChar2, _ := utf8.DecodeRuneInString(Open2)
	closeChar2, _ := utf8.DecodeRuneInString(Close2)
	KeyValueSeparatorRune, _ := utf8.DecodeRuneInString(KeyValueSeparator)

	return Style{StartDoc,EndDoc,Indent, Open,Close,Open2,Close2,KeyValueSeparator,Separator,LineBreak,True,False,OneLineComment,ScalarPrefix,
		openChar,closeChar,openChar2,closeChar2,KeyValueSeparatorRune, false}
}

/////////////////////////////////////////////////////////////////////////////
// Lexer
/////////////////////////////////////////////////////////////////////////////

func (style Style) GetNext(context Context, eol func(), open func(open string), close func(close string), nextTag func(tag Tag), nextLiteral func (literal Value)) error {

	ReadAndLookAhead := func(ch rune, expect1 rune, expect2 rune) bool {
		if ch == expect1 {
			if context.LookAhead() == expect2 {
				context.ReadRune()
				nextTag(Tag{string(expect1)+string(expect2)})
			} else {
				nextTag(Tag{string(expect1)})
			}
			return true
		}
		return false
	}

	ch, err := context.ReadRune()
	switch {
	case err != nil: return err
	case err == io.EOF:
		return err
	case ch == NEWLINE:
		eol()
		context.EOL()
	case unicode.IsSpace(ch) || ch == '\r': break // TODO fix comma
	case ch == style.OneLineComment:
		// TODO Comment is not part of parse tree, store elsewhere
		_, err = ReadUntilEndOfLine(context)
		if err != nil {
			return err
		}
	case ch == style.OpenChar : context.Open(); open(style.Open)
	case ch == style.closeChar : context.Close(); close(style.Close)
	case ch == style.OpenChar2 : context.Open(); open(style.Open2)
	case ch == style.closeChar2 : context.Close(); close(style.Close2)
	case ch == '"' :
		value, err := ReadCLanguageString(context)
		if err != nil {
			return err
		}
		nextLiteral(value)
	case ch == '.' && context.LookAhead() == '.': nextTag(Tag{".."})
	case ((ch == '.' || (ch== '-' && style.RecognizeNegative)) && unicode.IsNumber(context.LookAhead())) || unicode.IsNumber(ch): // TODO || ch == '+' 
		value, err := ReadNumber(context, string(ch))    // TODO minus
		if err != nil {
			return err
		}
		if tag, ok := value.(Tag); ok {
			nextTag(tag)
		} else {
			nextLiteral(value)
		}
	case ch == style.KeyValueSeparatorRune:		nextTag(tuple.CONS_ATOM)
	case ch == ',':
	case ch == ';':  nextTag(Tag{";"})
	case ReadAndLookAhead(ch, '.', '.'):
	case ReadAndLookAhead(ch, '>', '='):
	case ReadAndLookAhead(ch, '<', '='):
	case ReadAndLookAhead(ch, '!', '='):
	case ReadAndLookAhead(ch, '=', '='):
	case ReadAndLookAhead(ch, '*', '*'):
	case ReadAndLookAhead(ch, '+', '+'):
	case ReadAndLookAhead(ch, '|', '|'):
	case ReadAndLookAhead(ch, '&', '&'):
	case ch == '-' || ch== '/' || ch == '%': nextTag(Tag{string(ch)})
	case ch == '_' || unicode.IsLetter(ch):
		value, err :=(ReadTag(context, string(ch), func(r rune) bool { return r == '_' || unicode.IsLetter(r) || unicode.IsNumber(r) }))
		if err != nil {
			return err
		}
		if tag, ok := value.(Tag); ok {
			nextTag(tag)
		} else {
			nextLiteral(value)
		}		
	case unicode.IsGraphic(ch): Error(context,"Graphic character not recognised '%s'", string(ch))
	case unicode.IsControl(ch): Error(context,"Control character not recognised '%d'", ch)
	default: Error(context,"Character not recognised '%d'", ch)
	}
	return nil
}

/////////////////////////////////////////////////////////////////////////////

// "Code point" Go introduces a shorter term for the concept: rune and means exactly the same as "code point", with one interesting addition.
//  Go language defines the word rune as an alias for the type int32, so programs can be clear when an integer value represents a code point.
//  Moreover, what you might think of as a character constant is called a rune constant in Go. 

func ReadString (context Context, token string, unReadLast bool, test func(r rune) bool) (string, error) {
	for {
		if unReadLast {
			ch := context.LookAhead()
			if ! test(ch) {
				return token, nil
			}
		}
		ch, err := context.ReadRune()
		if err == io.EOF {
			//Error(context,"ERROR missing close quote: '%s'", DOUBLE_QUOTE)
			return token, nil
		} else if err != nil {
			//log.Printf("ERROR nil")
			return "", err
		} else {
			// TODO not efficient
			token = token + string(ch)
		}
	}
}

func ReadTag(context Context, prefix string, test func(rune) bool) (Value, error) {
	tag, err := ReadString(context, prefix, true, test)
	if err != nil {
		return Tag{""}, err
	}
	switch tag {
	case "NaN": return Float64(math.NaN()), nil
	case "Inf": return Float64(math.Inf(1)), nil // TODO "+Inf", and "-Inf" 
	default: return Tag{tag}, err
	}
}

func ReadNumber(context Context, token string) (Value, error) {  // Number
	var dots int
	if token == "." {
		dots = 1
	} else {
		dots = 0
	}
	for {
		ch := context.LookAhead()
		if (ch == '.' && dots == 0) || unicode.IsNumber(ch) {
			ch, err := context.ReadRune()
			if err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			} else if ch == '.' && dots == 0 {
				dots += 1
				token = token + "." // TODO not efficient
			} else if unicode.IsNumber(ch) {
				// TODO ought to be much more efficient to build up a number dynamically
				token = token + string(ch) // TODO not efficient
			}
		} else {
			break
		}
	}
	if token == "." {
		return Tag{"."}, nil
	}
	switch dots {
	case 0:
		value, err := strconv.ParseInt(token, 10, 0)  // TODO no need to parse as int64, could treat as bigint
		if err != nil {
			return Int64(0), err
		}
		return Int64(value), nil
	default:
		value, err := strconv.ParseFloat(token, 64)
		if err != nil {
			return Float64(0), err
		}
		return Float64(value), nil
	} 
}

func ReadUntilEndOfLine(context Context) (string, error) {
	token := ""
	for {
		ch := context.LookAhead()
		if ch == NEWLINE {
			return token, nil
		}
		ch, err := context.ReadRune()
		switch {
		case err == io.EOF:
			return token, nil
		case err != nil:
			return token, err
		default:
			token = token + string(ch)
		}
	}
}

func ReadCLanguageString(context Context) (String, error) {
	token := ""
	for {
		ch, err := context.ReadRune()
		switch {
		case err == io.EOF:
			Error(context,"ERROR missing close quote: '%s'", DOUBLE_QUOTE)
			return String(token), nil
		case err != nil: return String(""), err
		case ch == '"': return String(token), nil
		case ch == '\\':
			ch, err := context.ReadRune()
			if err == io.EOF {
				Error(context,"ERROR missing close quote: '%s'", DOUBLE_QUOTE)
				return String(token), nil
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

func IsCompare(ch rune) bool {
	switch ch {
		case '=': return true
		case '!': return true
		case '<': return true
		case '>': return true
		default: return false
	}
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

func (printer Style) PrintNullaryOperator(depth string, tag Tag, out StringFunction) {
	PrintTuple(&printer, depth, NewTuple(tag), out)
}

func (printer Style) PrintUnaryOperator(depth string, tag Tag, value Value, out StringFunction) {
	PrintTuple(&printer, depth, NewTuple(tag, value), out)
}

func (printer Style) PrintBinaryOperator(depth string, tag Tag, value1 Value, value2 Value, out StringFunction) {
	PrintTuple(&printer, depth, NewTuple(tag, value1, value2), out)
}

func (printer Style) PrintOpenTuple(depth string, value Value, out StringFunction) string {
	if tuple.IsConsInTuple(value) {
		out(printer.Open2)
	} else {
		out(printer.Open)
	}
	return depth + "  "
}

func (printer Style) PrintCloseTuple(depth string, value Value, out StringFunction) {
	printer.PrintIndent(depth, out)
	if tuple.IsConsInTuple(value) {
		out(printer.Close2)
	} else {
		out(printer.Close)
	}
}

func (printer Style) PrintHeadTag(tag Tag, out StringFunction) {
	out(tag.Name)
}

func (printer Style) PrintKey(tag Tag, out StringFunction) {
	out(tag.Name)
	out(printer.KeyValueSeparator)
}

func (printer Style) PrintScalar(depth string, value Value, out StringFunction) {
	tuple.PrintScalar(printer, depth, value, out)
}

func (printer Style) PrintScalarPrefix(depth string, out StringFunction) {}
