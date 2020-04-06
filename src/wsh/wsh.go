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
	inputGrammar := parsers.NewShellGrammar()
	outputGrammar := inputGrammar

	table := eval.NewLessSafeSymbolTable()
	table.Add("|", eval.Pipe)
	table.Add("=", eval.AssignLocal)

	//func reduce f t { progn c=1 accumulator=first(t) (for v t { accumulator = f(accumulator v))  accumulator}
		
	runner.ParseAndEval(&table, inputGrammar, "func count  t { progn (c=0) (for v t { c=c+1 }) c }")
	runner.ParseAndEval(&table, inputGrammar, "func first  t { nth 0 t }")
	runner.ParseAndEval(&table, inputGrammar, "func second t { nth 1 t }")
	runner.ParseAndEval(&table, inputGrammar, "func third  t { nth 2 t }")


	var symbols * eval.SymbolTable = nil
	if !*ast {
		symbols = &table
	}
	
	pipeline := runner.SimplePipeline (symbols, *queryPattern, outputGrammar, runner.PrintString)


	grammars := runner.NewGrammars()
	runner.AddAllKnownGrammars(&grammars)
	runner1 := runner.NewRunner(grammars, &table, runner.GetLogger(nil, *verbose), inputGrammar)
	runner.AddSafeGrammarFunctions(&table, &runner1.Grammars)

	files := runner.GetRemainingNonFlagOsArgs()
	errors := runner1.RunFiles(files, pipeline)

	if errors > 0 {
		os.Exit(1)
	}
}
