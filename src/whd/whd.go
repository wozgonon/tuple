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
	"tuple/runner"
	"tuple/parsers"
	"fmt"
	"net/http"
	"bufio"
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
	grammars := runner.NewGrammars(parsers.NewJSONGrammar())
	grammars.AddAllKnownGrammars()
	grammar, ok := grammars.FindBySuffix(grammarName)
	if ok {
		reader := bufio.NewReader(r) // TODO read from request
		context := runner.NewParserContext("<http>", reader, tuple.NewVerboseFilterLogger(false, tuple.NewDefaultLocationLogger()))
		grammar.Parse(&context, func (value tuple.Value) error {
			fmt.Fprintf(w, "%s", value)
			return nil
		})
	}
	fmt.Fprintf(w, "\n")
}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func main() {

	http.HandleFunc("/do", do)
	http.HandleFunc("/headers", headers)
	http.ListenAndServe(":8888", nil)
}
