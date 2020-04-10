package main

import (
	"tuple"
	"tuple/runner"
	"tuple/eval"
	"tuple/parsers"
)
import "os"
import "strings"
import "bufio"
import "fmt"
import "flag"

type SymbolTable = eval.SymbolTable

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
	args := runner.GetRemainingNonFlagOsArgs()
	expression := strings.Join(args, " ")

	//
	//  Set up the translator pipeline.
	//
	logger := tuple.GetLogger(nil, *verbose)
	outputGrammar := parsers.NewInfixExpressionGrammar()
	var symbols * SymbolTable = nil
	global := eval.NewErrorIfFunctionNotFound(logger)
	table := eval.NewSafeSymbolTable(global)
	if !*ast {
		symbols = &table
	}
	//table.Add("func", func(name string, body tuple.Tuple) { fmt.Printf("TODO Implement 'func' '%s' '%s'", name, body) })

	pipeline := runner.SimplePipeline (symbols, *queryPattern, outputGrammar, runner.PrintString)
	reader := bufio.NewReader(strings.NewReader(expression))
	context := runner.NewParserContext("<cli>", reader, logger)
	grammar := parsers.NewInfixExpressionGrammar()

	//
	//  Set up the translator pipeline.
	//
	grammar.Parse(&context, pipeline)

	if context.Errors() > 0 {
		os.Exit(1)
	}
}
