package main

import (
	"tuple"
	"os"
	"fmt"
	"log"
)

//func check(e error) {
//    if e != nil {
//        panic(e)
//    }
//}

func help() {
	fmt.Printf("%s <options> <files...>\n", os.Args[0])
	fmt.Printf("-o|--output <...>\n")
	fmt.Printf("-p|--pretty\n")
	fmt.Printf("-h|--help - prints this message\n")
	fmt.Printf("-v|--version\n")
}

func main() {

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
				parser := tuple.NewSExpressionParser(lispStyle, outputStyle, nextFunction(outputStyle))
				parser.ParseSExpression (context)
				fmt.Printf("\n")
			case ".jml":
				parser := tuple.NewSExpressionParser(jmlStyle, outputStyle, nextFunction(outputStyle))
				parser.ParseSExpression (context)
				fmt.Printf("\n")
			case ".tcl": 
				parser := tuple.NewSExpressionParser(tclStyle, outputStyle, nextFunction(outputStyle))
				parser.ParseTCL (context)
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
