package main

import "tuple"
import "os"
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
	//  Get the input expression from the command line
	//
	argsLength := len(os.Args)
	numberOfFiles := flag.NArg()
	args := os.Args[argsLength-numberOfFiles:] 
	expression := strings.Join(args, " ")

	//
	//  Set up the translator pipeline.
	//
	outputGrammar := tuple.NewInfixExpressionGrammar()
	var symbols * tuple.SymbolTable = nil
	table := tuple.NewSafeSymbolTable(&tuple.ErrorIfFunctionNotFound{})
	if !*ast {
		symbols = &table
	}
	//table.Add("func", func(name string, body tuple.Tuple) { fmt.Printf("TODO Implement 'func' '%s' '%s'", name, body) })

	pipeline := tuple.SimplePipeline (symbols, *queryPattern, outputGrammar, tuple.PrintString)
	reader := bufio.NewReader(strings.NewReader(expression))
	context := tuple.NewRunnerContext("<cli>", reader, tuple.GetLogger(nil), *verbose)
	grammar := tuple.NewInfixExpressionGrammar()

	//
	//  Set up the translator pipeline.
	//
	grammar.Parse(&context, pipeline)

	if context.Errors() > 0 {
		os.Exit(1)
	}
}
