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
	inputGrammar := tuple.NewShellGrammar()
	outputGrammar := inputGrammar

	table := tuple.NewUnSafeSymbolTable()
	table.Add("|", tuple.Pipe)
	table.Add("=", tuple.Assign)

	tuple.ParseAndEval(inputGrammar, table, "func count t { progn (c=0) (for v t { c=c+1 }) c }")

	var symbols * tuple.SymbolTable = nil
	if !*ast {
		symbols = &table
	}


	
	pipeline := tuple.SimplePipeline (symbols, *queryPattern, outputGrammar, tuple.PrintString)


	grammars := tuple.NewGrammars()
	grammars.Add(inputGrammar)
	files := tuple.GetRemainingNonFlagOsArgs()
	errors := tuple.RunFiles(files, tuple.GetLogger(nil), *verbose, inputGrammar, &grammars, pipeline)

	if errors > 0 {
		os.Exit(1)
	}
}
