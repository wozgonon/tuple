package main

import (
	"tuple"
	"io"
	"os"
	"unicode"
	"errors"
	"fmt"
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

func help() {
	fmt.Printf("%s <options> <files...>\n", os.Args[0])
	fmt.Printf("-o|--output <...>\n")
	fmt.Printf("-p|--pretty\n")
	fmt.Printf("-h|--help - prints this message\n")
	fmt.Printf("-v|--version\n")
}

func main() {

	if len(os.Args) == 1 {
		help()
		return
	}

	prettyPrint := func (tuple tuple.Tuple) {
		fmt.Printf ("%s\n", tuple.PrettyPrint(""))
	}
	processTuple := prettyPrint

	files := make([]string, 0)
	for _, v := range os.Args {
		switch v {
		case "-o", "--output":
		case "=p", "--pretty":
			processTuple = prettyPrint 
		case "-h", "--help":
			help()
			return
		case "-v", "--version":
			fmt.Print("%s\n", 0.1)
			return
		default:
			files = append(files, v)
		}
	}
	parser := Parser{'(', ')'}
	tuple.RunParser(files,
		func (context * tuple.ParserContext) (tuple.Tuple, error) {
			suffix := context.Suffix()
			switch suffix {
			case ".l":
				return parser.parse (context)
			case ".jml":
				return parser.parse (context)
			case ".yaml": fallthrough
			case ".json": fallthrough
			case ".tcl": fallthrough
			case ".xml": fallthrough
			case ".jpost": fallthrough
			case ".tsv": fallthrough
			case ".csv":
				context.Error("Not implemented file suffix: '%s'", suffix)
				return tuple.NewTuple(), errors.New("Not implemented file suffix: " + suffix)
			default:
				context.Error("Unsupported file suffix: '%s'", suffix)
				return tuple.NewTuple(), errors.New("Unsupported file suffix: " + suffix)
			}
			
		},
		processTuple)
}
