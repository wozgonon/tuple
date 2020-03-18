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
)

func main() {

	//
	//  Set up the command line arguments
	//
	
	var in = flag.String("in", ".l", "The format of the input.")
	var out = flag.String("out", ".l", "The format of the output.")
	var logger = flag.String("log", ".l", "The format of the error logging.")
	var verbose = flag.Bool("verbose", false, "Verbose logging.")
	var eval = flag.Bool("eval", false, "Run 'eval' interpretter.")
	var query = flag.String("query", "", "Select parts of the AST matching a query pattern.")
	var version = flag.Bool("version", false, "Print version of this software.")

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
	grammars.Add((tuple.NewLispGrammar()))
	grammars.Add((tuple.NewTclGrammar()))
	grammars.Add((tuple.NewJmlGrammar()))
	grammars.Add((tuple.NewTupleGrammar()))
	grammars.Add((tuple.NewYamlGrammar()))
	grammars.Add((tuple.NewIniGrammar()))
	grammars.Add((tuple.NewPropertyGrammar()))

	outputGrammar := grammars.FindBySuffixOrPanic(*out)
	loggerGrammar := grammars.FindBySuffixOrPanic(*logger)
	var inputGrammar *tuple.Grammar = nil
	if *in != "" {
		inputGrammar = grammars.FindBySuffixOrPanic(*in)
	}
	
	//
	//  Set up the translator pipeline.
	//

	prettyPrint := func(tuple interface{}) {
		(*outputGrammar).Print(tuple, func(value string) {
			fmt.Printf ("%s", value)
		})
	}
	pipeline := prettyPrint
	if *eval {
		next := pipeline
		pipeline = func(value interface{}) {
			tuple.SimpleEval(value, next)
		}
	}
	if *query != "" {
		next := pipeline
		pipeline = func(value interface{}) {
			tuple.Query(*query, value, next)
		}
	}
	
	//
	//  Run the translators over all the input files.
	//

	args := len(os.Args)
	numberOfFiles := flag.NArg()
	files := os.Args[args-numberOfFiles:]
 
	tuple.RunParser(files, *loggerGrammar, *verbose, inputGrammar, &grammars, pipeline)
}
