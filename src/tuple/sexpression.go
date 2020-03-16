package tuple

/////////////////////////////////////////////////////////////////////////////
// Parsing S-Expressions
/////////////////////////////////////////////////////////////////////////////


import (
	"io"
	"unicode"
	"fmt"
	"unicode/utf8"
)


// A [S-Expression](https://en.wikipedia.org/wiki/S-expression) or symbolic expression is a very old and general notation.
// A nested structure of scalars (atoms and numbers), lists and key-values pairs (called cons cells).
// These are used for the syntax of LISP but also any other language can typically be converted to an S-Expression,
// it is in particular a very useful format for debugging a parser by printing out the Abstract Syntaxt Tree (AST) created by parsing.
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

func (parser SExpressionParser) parseCommaTuple(context * ParserContext, tuple *Tuple) (error) {

	// TODO comma and semi-colon
	for {
		token, err := parser.getNext(context)
		switch {
		case err != nil:
			context.Error("parsing %s", err);
			return err /// ??? Any need to return
		case token == parser.style.Close:
			return nil
		case token == parser.style.Open:
			subTuple := NewTuple()
			err := parser.parseCommaTuple(context, &subTuple)
			if err == io.EOF {
				context.Error ("Missing close bracket")
				return err
			}
			if err != nil {
				return err
			}
			tuple.Append(subTuple)
		default:
			if _,ok := token.(Comment); ok {
				// TODO Ignore ???
			} else {
				tuple.Append(token)
			}
		}
	}

}

func (parser SExpressionParser) ParseTuple(context * ParserContext) {

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
		default:
			if atom,ok := token.(Atom); ok {
				bracket, err := parser.getNext(context)
				if err != nil {
					context.Error ("'%s'", err)
					return
				}
				if bracket != parser.style.Open {
					context.Error ("Expected open bracket '%s' after '%s', not '%s'", parser.style.Open, token, bracket)
				} else {
					subTuple := NewTuple()
					subTuple.Append(atom)
					err := parser.parseCommaTuple(context, &subTuple)
					if err != nil {
						return
					}
					parser.next(subTuple)
				}
			} else {
				parser.next(token)
			}
		}
		fmt.Print("\n")
	}
}

func (parser SExpressionParser) parseSExpressionTuple(context * ParserContext, tuple *Tuple) (error) {
	for {
		token, err := parser.getNext(context)
		switch {
		case err != nil:
			context.Error("parsing %s", err);
			return err /// ??? Any need to return
		case token == parser.style.Close:
			return nil
		case token == parser.style.Open:
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
		default:
			if _,ok := token.(Comment); ok {
				// TODO Ignore ???
			} else {
				tuple.Append(token)
			}
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
			subTuple := NewTuple()
			err := parser.parseSExpressionTuple(context, &subTuple)
			if err != nil {
				return
			}
			parser.next(subTuple)
		default:
			parser.next(token)
		}
		fmt.Print("\n")
	}
}

func (parser SExpressionParser) readCommandString(context * ParserContext, token string) (string, error) {
	return ReadString(context, token, true, func (ch rune) bool {
		return ! unicode.IsSpace(ch) && string(ch) != parser.style.Close && string(ch) != parser.style.Open && ch != '$'
	})

}

func (parser SExpressionParser) getNextCommandShell(context * ParserContext) (interface{}, error) {
	for {
		ch, err := context.ReadRune()
		switch {
		case err != nil: return "", err
		case err == io.EOF: return "", nil
		case ch == NEWLINE: return string(NEWLINE), nil
		case unicode.IsSpace(ch): break
		case ch == parser.style.OneLineComment:
			// TODO ignore for now
			//return string(ch), nil
		case ch == parser.openChar :  return parser.style.Open, nil
		case ch == parser.closeChar : return parser.style.Close, nil
		case ch == '"' :  return ReadCLanguageString(context)
		case ch == '.' || unicode.IsNumber(ch): return ReadNumber(context, string(ch))    // TODO minus
		case ch == '$':
			value, err := parser.readCommandString(context, "")
			if err != nil {
				return nil, err
			}
			return Atom{value}, nil
		case unicode.IsGraphic(ch): return parser.readCommandString(context, string(ch))
		case unicode.IsControl(ch): context.Error("Error control character not recognised '%d'", ch)
		default: context.Error("Error character not recognised '%d'", ch)
		}
	}
}

func (parser SExpressionParser) parseCommandShellTuple(context * ParserContext, tuple *Tuple) (error) {
	for {
		token, err := parser.getNextCommandShell(context)
		switch {
		case err != nil:
			context.Error("parsing %s", err);
			return err /// ??? Any need to return
		case token == parser.style.Close:
			return nil
		case token == parser.style.Open:
			subTuple := NewTuple()
			err := parser.parseCommandShellTuple(context, &subTuple)
			if err == io.EOF {
				context.Error ("Missing close bracket")
				return err
			}
			if err != nil {
				return err
			}
			tuple.Append(subTuple)
		case token == string(NEWLINE):
		default:
			tuple.Append(token)
		}
	}
}

func (parser SExpressionParser) ParseCommandShell(context * ParserContext) {

	resultTuple := NewTuple()
	for {
		token, err := parser.getNextCommandShell(context)
		switch {
		case err == io.EOF:
			return
		case err != nil:
			context.Error ("'%s'", err)
			return
		case token == string(NEWLINE):
			l := resultTuple.Length()
			context.Verbose ("Newline length of tuple=%d", l)
			switch l {
			case 0: // Ignore
			case 1:
				first := resultTuple.List[0]
				if _, ok := first.(Atom); ok {
					parser.next(resultTuple)
				} else {
					parser.next(token)
				}
			default:
				parser.next(resultTuple)
			}
			resultTuple = NewTuple()
		case token == parser.style.OneLineComment:
			comment, err := ReadUntilEndOfLine(context)
			if err != nil {
				return
			}
			parser.next(comment)
		case token == parser.style.Close:
			context.Error ("Unexpected close bracket '%s'", parser.style.Close)
		case token == parser.style.Open:
			subTuple := NewTuple()
			err := parser.parseCommandShellTuple(context, &subTuple)
			if err != nil {
				return // tuple, err
			}
			resultTuple.Append(subTuple)
		default:
			context.Verbose("Add token: '%s'", token)
			resultTuple.Append(token)
		}
		//fmt.Print("% ")
	}
}
