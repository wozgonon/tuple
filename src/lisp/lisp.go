package main

import (
	"tuple"
	"io"
	"os"
	"unicode"
	//"errors"
	"fmt"
	"log"
)

type Parser struct {
	style tuple.Style
	outputStyle tuple.Style
}

func (parser Parser) getNext(context * tuple.ParserContext) (interface{}, error) {

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

func (parser Parser) next(tuple interface{}) {
	result := ""
	parser.outputStyle.PrettyPrint(tuple, func(value string) { result = result + value })
	fmt.Printf ("%s\n", result)
}

func (parser Parser) parseTuple(context * tuple.ParserContext) (tuple.Tuple, error) {
	tuple := tuple.NewTuple()
	for {
		token, err := parser.getNext(context)
		switch {
		case err != nil:
			context.Error("parsing %s", err);
			return tuple, err /// ??? Any need to return
		case token == ")":
			return tuple, nil
		case token == "(":
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

func (parser Parser) parse(context * tuple.ParserContext) {

	for {
		token, err := parser.getNext(context)
		switch {
		case err == io.EOF:
			return
		case err != nil:
			context.Error ("'%s'", err)
			return
		case token == ")":
			context.Error ("Unexpected close bracket '%s'", ")")
		case token == "(":
			tuple, err := parser.parseTuple(context)
			if err != nil {
				return // tuple, err
			}
			parser.next(tuple)
		default:
			parser.next(token)
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

	tclStyle := tuple.Style{"  ", "{", "}", "", "\n"}
	jmlStyle := tuple.Style{"  ", "{", "}", "", "\n"}
	tupleStyle := tuple.Style{"  ", "(", ")", ",", "\n"} // prolog, sql
	lispStyle := tuple.Style{"  ", "(", ")", "", "\n"}

	outputStyle := lispStyle

	files := make([]string, 0)
	for _, v := range os.Args[1:] {
		switch v {
		case "-o", "--output":
		case "-p", "--pretty":
			outputStyle = lispStyle
		case "--jml":
			outputStyle= jmlStyle
		case "--tcl":
			outputStyle= tclStyle
		case "--tuple":
			outputStyle= tupleStyle
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
	tuple.RunParser(files,
		func (context * tuple.ParserContext) {
			suffix := context.Suffix()
			switch suffix {
			case ".l":
				parser := Parser{lispStyle, outputStyle}
				parser.parse (context)
			case ".jml":
				parser := Parser{jmlStyle, outputStyle}
				parser.parse (context)
			case ".tcl": 
				parser := Parser{tclStyle, outputStyle}
				parser.parse (context)
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
