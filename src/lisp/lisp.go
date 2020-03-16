package main

import (
	"tuple"
	"io"
	"os"
	"unicode"
	//"errors"
	"fmt"
	"log"
	"unicode/utf8"
	//"reflect"
)

//func check(e error) {
//    if e != nil {
//        panic(e)
//    }
//}

type Parser struct {
	style tuple.Style
	outputStyle tuple.Style
	openChar rune
	closeChar rune
	next tuple.Next
}

func NewParser(style tuple.Style, outputStyle tuple.Style, next tuple.Next) Parser {

	openChar, _ := utf8.DecodeRuneInString(style.Open)
	closeChar, _ := utf8.DecodeRuneInString(style.Close)
	return Parser{style,outputStyle,openChar,closeChar,next}
}

func (parser Parser) getNext(context * tuple.ParserContext) (interface{}, error) {

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
		case ch == '"' :  return tuple.ReadCLanguageString(context)
		case ch == '.' || unicode.IsNumber(ch): return tuple.ReadNumber(context, string(ch))    // TODO minus
		case tuple.IsArithmetic(ch): return tuple.ReadAtom(context, string(ch), func(r rune) bool { return tuple.IsArithmetic(r) })
		case tuple.IsCompare(ch): return tuple.ReadAtom(context, string(ch), func(r rune) bool { return tuple.IsCompare(r) })
		case unicode.IsLetter(ch):  return tuple.ReadAtom(context, string(ch), func(r rune) bool { return unicode.IsLetter(r) })
		case unicode.IsGraphic(ch): context.Error("Error graphic character not recognised '%s'", string(ch))
		case unicode.IsControl(ch): context.Error("Error control character not recognised '%d'", ch)
		default: context.Error("Error character not recognised '%d'", ch)
		}
	}
}

func (parser Parser) parseTuple(context * tuple.ParserContext) (tuple.Tuple, error) {
	tuple := tuple.NewTuple()
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

func (parser Parser) parseSExpression(context * tuple.ParserContext) {

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

func (parser Parser) parseTCL(context * tuple.ParserContext) {

	resultTuple := tuple.NewTuple()
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
				if _, ok := first.(tuple.Atom); ok {
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

func help() {
	fmt.Printf("%s <options> <files...>\n", os.Args[0])
	fmt.Printf("-o|--output <...>\n")
	fmt.Printf("-p|--pretty\n")
	fmt.Printf("-h|--help - prints this message\n")
	fmt.Printf("-v|--version\n")
}

func main() {

	//if len(os.Args) == 1 {
	//	help()
	//	return
	//}

	tclStyle := tuple.Style{"  ", "{", "}", "", "\n", "true", "false", '#'}
	jmlStyle := tuple.Style{"  ", "{", "}", "", "\n", "true", "false", '#'}
	tupleStyle := tuple.Style{"  ", "(", ")", ",", "\n", "true", "false", '%'} // prolog, sql '--' for 
	lispStyle := tuple.Style{"  ", "(", ")", "", "\n", "true", "false", ';'}

	logStyle := lispStyle
	outputStyle := lispStyle

	var eval bool
	files := make([]string, 0)
	l := len(os.Args[1:])
	for k, v := range os.Args[1:] {
		switch v {
		case "-o", "--output":
			if k < l-1 {
				
			}
		case "-e", "--eval":
			eval = true
		case "-p", "--pretty":
			outputStyle = lispStyle
		case "--jml":
			outputStyle= jmlStyle
		case "--tcl":
			outputStyle= tclStyle
		case "--tuple":
			outputStyle= tupleStyle
		case "--l_tcl":
			logStyle= tclStyle
		case "-h", "--help":
			help()
			return
		case "-v", "--version":
			fmt.Print("%s\n", 0.1)
			return
		default:
			log.Print(v)
			files = append(files, v)
		}
	}

	nextFunction := func(outputStyle tuple.Style) tuple.Next {
		pretty := func(tuple interface{}) {
			outputStyle.PrettyPrint(tuple, func(value string) {
				fmt.Printf ("%s", value)
			})
		}
		if eval {
			return func(value interface{}) {
				tuple.SimpleEval(value, pretty)
			}
		} else {
			return pretty
		}
	}
	tuple.RunParser(files, logStyle,
		func (context * tuple.ParserContext) {
			suffix := context.Suffix()
			switch suffix {
			case ".l":
				parser := NewParser(lispStyle, outputStyle, nextFunction(outputStyle))
				parser.parseSExpression (context)
				fmt.Printf("\n")
			case ".jml":
				parser := NewParser(jmlStyle, outputStyle, nextFunction(outputStyle))
				parser.parseSExpression (context)
				fmt.Printf("\n")
			case ".tcl": 
				parser := NewParser(tclStyle, outputStyle, nextFunction(outputStyle))
				parser.parseTCL (context)
				fmt.Printf("\n")
			case ".yaml": fallthrough
			case ".json": fallthrough
			case ".xml": fallthrough
			case ".jpost": fallthrough
			case ".tsv": fallthrough
			case ".csv":
				context.Error("Not implemented file suffix: '%s'", suffix)
			default:
				context.Error("Unsupported file suffix: '%s', source: '%s'", suffix, context.SourceName)
			}
		})
}
