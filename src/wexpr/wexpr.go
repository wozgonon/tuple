package main

import "tuple"
import "os"
import "io"
import "strings"
import "bufio"
import "fmt"
import "flag"

func main () {

	var verbose = flag.Bool("verbose", false, "Verbose logging.")
	var ast = flag.Bool("ast", false, "If set then returns the AST else runs the 'eval' interpretter.")
	var queryPattern = flag.String("query", "", "Select parts of the AST matching a query pattern.")
	var version = flag.Bool("version", false, "Print version of this software.")
	flag.Parse()
	
	if *version {
		fmt.Printf("%s version 0.1", os.Args[0])
		return
	}

	//
	//  Set up the translator pipeline.
	//
	outputGrammar := tuple.NewLispGrammar()
	pipeline := tuple.SimplePipeline (!*ast, *queryPattern, &outputGrammar)

	argsLength := len(os.Args)
	numberOfFiles := flag.NArg()
	args := os.Args[argsLength-numberOfFiles:]
	expression := strings.Join(args, " ")

	operators := tuple.NewOperators()
	operators.AddStandardCOperators()

	reader := bufio.NewReader(strings.NewReader(expression))
	logGrammar := tuple.NewLispGrammar()

	style := tuple.Style{"", "", "  ",
		"(", ")", "", "", ".",
		"", "\n", "true", "false", ';', ""}
	sexp := tuple.NewSExpressionParser(style)

	context := tuple.NewParserContext("<cli>", reader, logGrammar, *verbose, pipeline)
	grammar := tuple.NewOperatorGrammar(&context, &operators)

	context.Verbose("Inout: [%s]", expression)

	for {
		token, err := sexp.GetNext(&context)
		if err == io.EOF {
			grammar.EOF(pipeline)
			break
		}
		//context.Verbose("-- %s\n", token)
		if token == "(" {
			grammar.OpenBracket(tuple.Atom{"("})
		} else if token == ")" {
			grammar.CloseBracket(tuple.Atom{")"})

		} else if atom, ok := token.(tuple.Atom); ok {
			if operators.Precedence(atom) != -1 {
				grammar.PushOperator(atom)
			} else if operators.IsOpenBracket(atom) {
				grammar.OpenBracket(atom)
			} else if operators.IsCloseBracket(atom) {
				grammar.CloseBracket(atom)
			} else {
				
				grammar.PushValue(atom)
			}
		} else {
			grammar.PushValue(token)
		}
	}
	//fmt.Printf("Errors: %d\n", context.Errors)
	if context.Errors > 0 {
		os.Exit(1)
	}
}
