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
	inputGrammar := parsers.NewShellGrammar()
	outputGrammar := inputGrammar

	table := eval.NewLessSafeSymbolTable()
	table.Add("|", eval.Pipe)
	table.Add("=", eval.Assign)
	table.Add("ast", func (expression string) tuple.Value { return runner.ParseString(inputGrammar, expression) })
	table.Add("expr", func (expression string) tuple.Value { return  runner.ParseAndEval(inputGrammar, table, expression) })

	runner.ParseAndEval(inputGrammar, table, "func count  t { progn (c=0) (for v t { c=c+1 }) c }")
	runner.ParseAndEval(inputGrammar, table, "func first  t { nth 0 t }")
	runner.ParseAndEval(inputGrammar, table, "func second t { nth 1 t }")
	runner.ParseAndEval(inputGrammar, table, "func third  t { nth 2 t }")

	var symbols * eval.SymbolTable = nil
	if !*ast {
		symbols = &table
	}
	
	pipeline := runner.SimplePipeline (symbols, *queryPattern, outputGrammar, runner.PrintString)

	grammars := tuple.NewGrammars()
	grammars.Add(inputGrammar)
	files := runner.GetRemainingNonFlagOsArgs()
	errors := runner.RunFiles(files, runner.GetLogger(nil), *verbose, inputGrammar, &grammars, pipeline)

	if errors > 0 {
		os.Exit(1)
	}
}
