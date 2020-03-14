package main

import (
	"tuple"
	"io"
	"os"
	"unicode"
	"errors"
)

type Parser struct {
	Open rune
	Close rune
}

func (parser Parser) next(context * tuple.ParserContext) (interface{}, error) {

	for {
		ch, err := context.ReadRune()
		switch {
		case err != nil: return "", err
		case err == io.EOF: return "", nil
		case unicode.IsSpace(ch): break
		case ch == '(' :  return "(", nil
		case ch == ')' :  return ")", nil
		case ch == '"' :  return tuple.ReadCLanguageString(context)
		case ch == '.' || unicode.IsNumber(ch): return tuple.ReadNumber(context, string(ch))    // TODO minus
		case tuple.IsArithmetic(ch): return tuple.ReadAtom(context, string(ch), func(r rune) bool { return tuple.IsArithmetic(r) })
		case unicode.IsLetter(ch):  return tuple.ReadAtom(context, string(ch), func(r rune) bool { return unicode.IsLetter(r) })
		case unicode.IsGraphic(ch): context.Error("Error graphic character not recognised '%s'", string(ch))
		case unicode.IsControl(ch): context.Error("Error control character not recognised '%d'", ch)
		default: context.Error("Error character not recognised '%d'", ch)
		}
	}
}

func (parser Parser) parse(context * tuple.ParserContext) (tuple.Tuple, error) {

	tuple := tuple.NewTuple()
	for {
		token, err := parser.next(context)
		switch {
		case err != nil: return tuple, err
		case token == ")":
			return tuple, nil
			return tuple, errors.New("Unexpected ')'")
		case token == "(":
			subTuple, err := parser.parse(context)
			if err == io.EOF {
				context.Error ("Missing close bracket")
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

func main() {
	
	parser := Parser{'(', ')'}
	tuple.RunParser(os.Args[1:], func (context * tuple.ParserContext) (tuple.Tuple, error) { return parser.parse (context) })
}
