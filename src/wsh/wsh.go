package main

import "tuple"
import "os"
import "os/exec"
//import "io"
import "fmt"
import "flag"
import "log"
import "bytes"

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
	table := tuple.NewSymbolTable()

	//
	//  Add shell specific commands
	//  These are typically not a 'safe' in that they allow access to the file system
	// 
	table.Add("exec", func (arg string) bool {
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
	})
	
	pipeline := tuple.SimplePipeline (&table, *queryPattern, outputGrammar)

	grammars := tuple.NewGrammars()
	errors := tuple.RunFiles([]string{}, tuple.GetLogger(nil), *verbose, inputGrammar, &grammars, pipeline)

	if errors > 0 {
		os.Exit(1)
	}
}
