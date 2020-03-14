package tuple

import "io"
import "log"
import "unicode"
import "strconv"
import "math"

/////////////////////////////////////////////////////////////////////////////


// "Code point" Go introduces a shorter term for the concept: rune and means exactly the same as "code point", with one interesting addition.
//  Go language defines the word rune as an alias for the type int32, so programs can be clear when an integer value represents a code point.
//  Moreover, what you might think of as a character constant is called a rune constant in Go. 

func ReadString (context * ParserContext, token string, keepLast bool, test func(r rune) bool) (string, error) {
	for {
		ch, err := context.ReadRune()
		if err == io.EOF {
			log.Printf("ERROR missing close \"")
			return token, nil
		} else if err != nil {
			//log.Printf("ERROR nil")
			//return ""
		} else if ! test(ch) {
			if keepLast {
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

func ReadNumber(context * ParserContext, token string) (interface{}, error) {
	var dots int
	if token == "." {
		dots = 1
	} else {
		dots = 0
	}
	for {
		ch, err := context.ReadRune()
		if err == io.EOF {
			return token, nil
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
			switch dots {
			case 0: return strconv.ParseInt(token, 10, 0)
			case 1:	return strconv.ParseFloat(token, 64)
			} 
		}
	}
}

func ReadCLanguageString(context * ParserContext) (string, error) {
	token := ""
	for {
		ch, err := context.ReadRune()
		switch {
		case err == io.EOF:
			log.Printf("ERROR missing lose \"")
			return token, nil
		case err != nil: return "", err
		case ch == '"': return token, nil
		case ch == '\\':
			ch, err := context.ReadRune()
			if err == io.EOF {
				log.Printf("ERROR missing lose \"")
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
	case 'n': return '\n'
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
