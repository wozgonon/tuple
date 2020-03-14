package main

import "fmt"

import (
	"tuple"
	"bufio"
	"io"
	"log"
	"os"
	"unicode"
)

//	Open rune
//	Close rune
//	return ParserContext{sourceName, 0, 0, '{', '}'}

func next(parser * tuple.ParserContext) (interface{}, error) {

	for {
		ch, err := parser.ReadRune()
		switch {
		case err != nil: return "", err
		case err == io.EOF: return "", nil
		case unicode.IsSpace(ch): break
		case ch == '(' :  return "(", nil
		case ch == ')' :  return ")", nil
		case ch == '"' :  return tuple.ReadCLanguageString(parser)
		case ch == '.' || unicode.IsNumber(ch): return tuple.ReadNumber(parser, string(ch))    // TODO minus
		case tuple.IsArithmetic(ch): return tuple.ReadAtom(parser, string(ch), func(r rune) bool { return tuple.IsArithmetic(r) })
		case unicode.IsLetter(ch):  return tuple.ReadAtom(parser, string(ch), func(r rune) bool { return unicode.IsLetter(r) })
		case unicode.IsGraphic(ch): parser.Error("Error graphic character not recognised '%s'", string(ch))
		case unicode.IsControl(ch): parser.Error("Error control character not recognised '%d'", ch)
		default: parser.Error("Error character not recognised '%d'", ch)
		}
	}
}

func parse(parser * tuple.ParserContext) (tuple.Tuple, error) {

	tuple := tuple.NewTuple()
	for {
		token, err := next(parser)
		switch {
		case err != nil: return tuple, err
		case token == ")": return tuple, nil
		case token == "(":
			subTuple, err := parse(parser)
			if err == io.EOF {
				parser.Error ("Missing close bracket")
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

func Run(parser * tuple.ParserContext) {
	tuple, err := parse(parser)
	if err != io.EOF && err != nil {
		parser.Error("Failed while parsing: %s", err)
	} else {
		fmt.Printf ("%s\n", tuple.PrettyPrint(""))
	}

}

func main() {

	
	if len(os.Args) == 1 {

		reader := bufio.NewReader(os.Stdin)
		parser := tuple.NewParserContext("<stdin>", reader)
		Run(&parser)
	} else {
		for _, fileName := range os.Args[1:] {
			file, err := os.Open(fileName)
			if err != nil {
				log.Fatal(err)
			}
			reader := bufio.NewReader(file)
			parser := tuple.NewParserContext(fileName, reader)
			Run(&parser)
			file.Close()
		}
	}
}
