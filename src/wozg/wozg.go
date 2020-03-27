/*
    This file is part of WOZG.

    WOZG is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    WOZG is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with WOZG.  If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"tuple"
	"os"
	"fmt"
	"flag"
	"strings"
	"bufio"
)

func main() {

	//
	//  Set up the command line arguments
	//
	
	var in = flag.String("in", ".l", "The format of the input.")
	var out = flag.String("out", ".l", "The format of the output.")
	var loggerGrammarSuffix = flag.String("log", "", "The format of the error logging.")
	var verbose = flag.Bool("verbose", false, "Verbose logging.")
	var eval = flag.Bool("eval", false, "Run 'eval' interpretter.")
	var queryPattern = flag.String("query", "", "Select parts of the AST matching a query pattern.")
	var version = flag.Bool("version", false, "Print version of this software.")
	var command = flag.Bool("command", false, "Execute command lines arguments rather than files.")
	var listGrammars = flag.Bool("list-grammars", false, "List supported grammars.")

	flag.Parse()

	//
	//  Print the version message.
	//
 
	if *version {
		fmt.Printf("%s %s %s %s\n", os.Args[0], VERSION, COMMIT, BUILT)
		return
	}

	//
	// Set up and then look up the set of supported grammars.
	//
	grammars := tuple.NewGrammars()
	tuple.AddAllKnownGrammars(&grammars)

	if *listGrammars {
		grammars.Forall(func (grammar tuple.Grammar) {
			fmt.Printf("%s\t%s\n", grammar.FileSuffix(), grammar.Name())
			})
		return
	}

	
	outputGrammar := grammars.FindBySuffixOrPanic(*out)
	loggerGrammar, _ := grammars.FindBySuffix(*loggerGrammarSuffix)
	logger := tuple.GetLogger(loggerGrammar)
	var inputGrammar tuple.Grammar = nil
	if *in != "" {
		inputGrammar = grammars.FindBySuffixOrPanic(*in)
	}
	
	//
	//  Set up the translator pipeline.
	//
	var symbols * tuple.SymbolTable = nil
	table := tuple.NewSymbolTable()
	if *eval {
		symbols = &table
	}
	pipeline := tuple.SimplePipeline (symbols, *queryPattern, outputGrammar)

	if *command {
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
		reader := bufio.NewReader(strings.NewReader(expression))
		context := tuple.NewRunnerContext("<cli>", reader, logger, *verbose)
		grammar := grammars.FindBySuffixOrPanic(*in)

		grammar.Parse(&context, pipeline)
		if context.Errors() > 0 {
			os.Exit(1)
		}

	} else {
		//
		//  Run the translators over all the input files.
		//
		args := len(os.Args)
		numberOfFiles := flag.NArg()
		files := os.Args[args-numberOfFiles:]
		errors := tuple.RunFiles(files, logger, *verbose, inputGrammar, &grammars, pipeline)
		//
		//  Exit with non-zero response code if any errors occurred.
		//
		if errors > 0 {
			os.Exit(1)
		}
	}

}
