package tuple

/////////////////////////////////////////////////////////////////////////////
// Parsing S-Expressions
/////////////////////////////////////////////////////////////////////////////


import (
	"io"
	"unicode"
	//"errors"
	"fmt"
	"unicode/utf8"
	//"reflect"
)

//func check(e error) {
//    if e != nil {
//        panic(e)
//    }
//}

type SExpressionParser struct {
	style Style
	outputStyle Style
	openChar rune
	closeChar rune
	next Next
}

func NewSExpressionParser(style Style, outputStyle Style, next Next) SExpressionParser {

	openChar, _ := utf8.DecodeRuneInString(style.Open)
	closeChar, _ := utf8.DecodeRuneInString(style.Close)
	return SExpressionParser{style,outputStyle,openChar,closeChar,next}
}

func (parser SExpressionParser) getNext(context * ParserContext) (interface{}, error) {

	for {
		ch, err := context.ReadRune()
		switch {
		case err != nil: return "", err
		case err == io.EOF: return "", nil
		case unicode.IsSpace(ch): break
		// TODO case ch == parser.style.OneLineComment:
			//comment := make(
			// TODO return 
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

func (parser SExpressionParser) parseTuple(context * ParserContext) (Tuple, error) {
	tuple := NewTuple()
	for {
		token, err := parser.getNext(context)
		switch {
		case err != nil:
			context.Error("parsing %s", err);
			return tuple, err /// ??? Any need to return
		case token == parser.style.Close:
			return tuple, nil
		case token == parser.style.Open:
			subTuple, err := parser.parseTuple(context)
			if err == io.EOF {
				context.Error ("Missing close bracket")
				return tuple, err
			}
			if err != nil {
				return tuple, err
			}
			tuple.Append(subTuple)
		default:
			tuple.Append(token)
		}
	}

}

func (parser SExpressionParser) ParseSExpression(context * ParserContext) {

	for {
		token, err := parser.getNext(context)
		switch {
		case err == io.EOF:
			return
		case err != nil:
			context.Error ("'%s'", err)
			return
		case token == parser.style.Close:
			context.Error ("Unexpected close bracket '%s'", parser.style.Close)
		case token == parser.style.Open:
			tuple, err := parser.parseTuple(context)
			if err != nil {
				return // tuple, err
			}
			parser.next(tuple)
		default:
			parser.next(token)
		}
		fmt.Print(" ")
	}
}

func (parser SExpressionParser) ParseTCL(context * ParserContext) {

	resultTuple := NewTuple()
	for {
		token, err := parser.getNext(context)
		switch {
		case err == io.EOF:
			return
		case err != nil:
			context.Error ("'%s'", err)
			return
		case token == "\n":
			l := len(resultTuple.List)
			if l == 1 {
				first := resultTuple.List[0]
				if _, ok := first.(Atom); ok {
					parser.next(resultTuple)
				} else {
					parser.next(token)
				}
			} else {
				parser.next(resultTuple)
			}
		case token == parser.style.Close:
			context.Error ("Unexpected close bracket '%s'", parser.style.Close)
		case token == parser.style.Open:
			resultTuple, err := parser.parseTuple(context)
			if err != nil {
				return // tuple, err
			}
			parser.next(resultTuple)
		default:
			parser.next(token)
		}
		fmt.Print(" ")
	}
}
