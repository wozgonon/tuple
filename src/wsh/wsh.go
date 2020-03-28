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
import "os/exec"
//import "io"
import "fmt"
import "flag"
import "log"
import "bytes"

/////////////////////////////////////////////////////////////////////////////
//  Functions specific to command shell
/////////////////////////////////////////////////////////////////////////////

func executeProcess(arg string) bool {
	cmd := exec.Command(arg)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
		return false
	}
	outStr := string(stdout.Bytes())
	fmt.Print(outStr)
	return true
}

func spawnProcess (arg string) bool {
		cmd := exec.Command(arg)
		var stdout bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = os.Stderr

		err := cmd.Start()
		if err != nil {
			log.Printf("Command finished with error: %v", err)
			return false
		}
		outStr := string(stdout.Bytes())
		fmt.Print(outStr)
		return true
}

func execIfNotFound(name tuple.Atom, args [] tuple.Value) tuple.Value {
	return tuple.Bool(executeProcess(name.Name))
}

func main () {

	var verbose = flag.Bool("verbose", false, "Verbose logging.")
	//var ast = flag.Bool("ast", false, "If set then returns the AST else runs the 'eval' interpretter.")
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
	outputGrammar := tuple.NewShellGrammar()
	inputGrammar := tuple.NewShellGrammar()
	table := tuple.NewSymbolTable(execIfNotFound)
	//
	//  Add shell specific commands
	//  These are typically not a 'safe' in that they allow access to the file system
	// 
	table.Add("exec", executeProcess)
	table.Add("spawn", spawnProcess)
	
	pipeline := tuple.SimplePipeline (&table, *queryPattern, outputGrammar)

	grammars := tuple.NewGrammars()
	errors := tuple.RunFiles([]string{}, tuple.GetLogger(nil), *verbose, inputGrammar, &grammars, pipeline)

	if errors > 0 {
		os.Exit(1)
	}
}
