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
	"fmt"
	"net/http"
	"bufio"
	"flag"
	"os"
)

func do(w http.ResponseWriter, req *http.Request) {

	fmt.Printf("Request: %s %s\n", req.URL, req.Method)
	//strings.NewReader("1+2")
	r := req.Body

	var grammarName string
	grammarNames, ok := req.Header["Content-type"]
	if ! ok {
		grammarName = "expr"
	} else {
		grammarName = grammarNames[0]
	}
	grammars := tuple.NewGrammars()
	tuple.AddAllKnownGrammars(&grammars)
	grammar, ok := grammars.FindBySuffix(grammarName)
	if ok {
		reader := bufio.NewReader(r) // TODO read from request
		context := tuple.NewRunnerContext("<http>", reader, tuple.GetLogger(nil), false)
		grammar.Parse(&context, func (value tuple.Value) {
			fmt.Printf("GOT '%s'\n", value)
			fmt.Fprintf(w, "%s", value)
		})
	} else {
		fmt.Printf("Unkown grammar '%s'\n", grammarName)
	}
	fmt.Fprintf(w, "\n")
	fmt.Printf("DONE\n")
}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func main() {

	//var verbose = flag.Bool("verbose", false, "Verbose logging.")
	var listenPort = flag.String("port", ":8888", "Listen port.")
	var version = flag.Bool("version", false, "Print version of this software.")
	flag.Parse()
	
	if *version {
		fmt.Printf("%s version 0.1", os.Args[0])
		return
	}

	fmt.Printf("Listen on: '%s'\n", *listenPort)

	http.HandleFunc("/do", do)
	http.HandleFunc("/headers", headers)
	http.ListenAndServe(*listenPort, nil)
}
