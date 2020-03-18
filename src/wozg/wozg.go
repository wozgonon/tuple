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
	// Set up and then look up the set of supported syntaxes.
	//

	syntaxes := tuple.NewGrammars()
	syntaxes.Add((tuple.NewLispGrammar()))
	syntaxes.Add((tuple.NewTclGrammar()))
	syntaxes.Add((tuple.NewJmlGrammar()))
	syntaxes.Add((tuple.NewTupleGrammar()))
	syntaxes.Add((tuple.NewYamlGrammar()))
	syntaxes.Add((tuple.NewIniGrammar()))
	syntaxes.Add((tuple.NewPropertyGrammar()))

	outputGrammar := syntaxes.FindBySuffixOrPanic(*out)
	loggerGrammar := syntaxes.FindBySuffixOrPanic(*logger)
	var inputGrammar *tuple.Grammar = nil
	if *in != "" {
		inputGrammar = syntaxes.FindBySuffixOrPanic(*in)
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
 
	tuple.RunParser(files, *loggerGrammar, *verbose,
		pipeline,
		func (context * tuple.ParserContext) {
			suffix := context.Suffix()
			var syntax *tuple.Grammar
			if suffix == "" {
				if inputGrammar == nil {
					panic("Input syntax for '" + context.SourceName + "' not given, use -in ...")
				}
				syntax = inputGrammar
			} else {
				syntax = syntaxes.FindBySuffixOrPanic(suffix)
			}
			context.Verbose("source [%s] suffix [%s]", context.SourceName, (*syntax).FileSuffix ())
			(*syntax).Parse(context)
		})
}
