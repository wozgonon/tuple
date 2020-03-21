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
import "strconv"
import "math"

const NEWLINE = '\n'
const DOUBLE_QUOTE = "\""
const UNKNOWN = "<???>"
const WORLD = "世界"

/////////////////////////////////////////////////////////////////////////////


// "Code point" Go introduces a shorter term for the concept: rune and means exactly the same as "code point", with one interesting addition.
//  Go language defines the word rune as an alias for the type int32, so programs can be clear when an integer value represents a code point.
//  Moreover, what you might think of as a character constant is called a rune constant in Go. 

func ReadString (context * ParserContext, token string, unReadLast bool, test func(r rune) bool) (string, error) {
	for {
		ch, err := context.ReadRune()
		if err == io.EOF {
			//context.Error("ERROR missing close quote: '%s'", DOUBLE_QUOTE)
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

func ReadAtom(context * ParserContext, prefix string, test func(rune) bool) (interface{}, error) {
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

func ReadNumber(context * ParserContext, token string) (interface{}, error) {  // Number
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
	//return Number{dots=true,token}
	switch dots {
	case 0: return strconv.ParseInt(token, 10, 0)
	default: return strconv.ParseFloat(token, 64)
	} 
}

func ReadUntilEndOfLine(context * ParserContext) (Comment, error) {
	token := ""
	for {
		ch, err := context.ReadRune()
		switch {
		case err == io.EOF:
			return NewComment(*context, token), nil
		case err != nil:
			return NewComment(*context, token), err
		case ch == NEWLINE:
			context.UnreadRune()
			return NewComment(*context, token), err
		default:
			token = token + string(ch)
		}
	}
}

func ReadUntilSpace(context * ParserContext, token string) (string, error) {
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

func ReadCLanguageString(context * ParserContext) (string, error) {
	token := ""
	for {
		ch, err := context.ReadRune()
		switch {
		case err == io.EOF:
			context.Error("ERROR missing close quote: '%s'", DOUBLE_QUOTE)
			return token, nil
		case err != nil: return "", err
		case ch == '"': return token, nil
		case ch == '\\':
			ch, err := context.ReadRune()
			if err == io.EOF {
				context.Error("ERROR missing close quote: '%s'", DOUBLE_QUOTE)
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


