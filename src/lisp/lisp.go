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

type Parser struct {
	//	Open rune
	//	Close rune
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
		case token == ")": return tuple, nil
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

func (parser Parser) Run(context * tuple.ParserContext) {
	tuple, err := parser.parse(context)
	if err != io.EOF && err != nil {
		context.Error("Failed while parsing: %s", err)
	} else {
		fmt.Printf ("%s\n", tuple.PrettyPrint(""))
	}

}

func main() {

	
	if len(os.Args) == 1 {

		reader := bufio.NewReader(os.Stdin)
		context := tuple.NewParserContext("<stdin>", reader)
		parser := Parser{}
		parser.Run(&context)
	} else {
		for _, fileName := range os.Args[1:] {
			file, err := os.Open(fileName)
			if err != nil {
				log.Fatal(err)
			}
			reader := bufio.NewReader(file)
			context := tuple.NewParserContext(fileName, reader)
			parser := Parser{}
			parser.Run(&context)
			file.Close()
		}
	}
}
