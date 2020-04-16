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
	logger := tuple.NewVerboseFilterLogger(*verbose, tuple.NewDefaultLocationLogger())
	runner1 := runner.NewSafeEvalContext(logger)
	
	grammars := runner.NewGrammars(parsers.NewInfixExpressionGrammar())
	pipeline := runner.SimplePipeline (runner1, !*ast, *queryPattern, grammars.Default(), runner.PrintString)
	reader := bufio.NewReader(strings.NewReader(expression))
	context := runner.NewParserContext("<cli>", reader, logger)

	//
	//  Set up the translator pipeline.
	//
	grammars.Default().Parse(&context, pipeline)

	if context.Errors() > 0 {
		os.Exit(1)
	}
}
