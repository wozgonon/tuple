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

import "tuple"
import "tuple/eval"
import "tuple/runner"
import "tuple/parsers"
import "os"
import "fmt"
import "flag"

/////////////////////////////////////////////////////////////////////////////
//  Functions specific to command shell
/////////////////////////////////////////////////////////////////////////////


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
	logger := tuple.GetLogger(nil, *verbose)
	ifNotFound := eval.NewExecIfNotFound()

	grammars := runner.NewGrammars(parsers.NewShellGrammar())
	grammars.AddAllKnownGrammars()
	runner1 := eval.NewRunner(ifNotFound, logger)

	eval.AddSafeFunctions(&runner1)
	runner.AddSafeQueryFunctions(&runner1)
	runner.AddTranslatedSafeFunctions(&runner1)
	grammars.AddSafeGrammarFunctions(&runner1)

	eval.AddLessSafeFunctions(&runner1, &runner1)
	runner1.Add("|", eval.Pipe)
	runner1.Add("=", eval.AssignLocal)

	outputGrammar := grammars.Default()
	pipeline := runner.SimplePipeline (&runner1, !*ast, *queryPattern, outputGrammar, runner.PrintString)

	files := runner.GetRemainingNonFlagOsArgs()
	errors := grammars.RunFiles(logger, files, pipeline)

	if errors > 0 {
		os.Exit(1)
	}
}
