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
package runner

import "tuple"
import 	"log"
import 	"fmt"
import 	"bufio"
import 	"os"
import 	"strings"
import "path"
import "flag"
import "tuple/eval"
import "tuple/parsers"

type Grammar = tuple.Grammar
type Value = tuple.Value
type Next = tuple.Next
type Context = tuple.Context
type String = tuple.String
type Int64 = tuple.Int64
type LocationLogger = tuple.LocationLogger

var NewParserContext = parsers.NewParserContext
var GetLogger = tuple.GetLogger
const STDIN = "<stdin>"
const PROMPT = "$ "

/////////////////////////////////////////////////////////////////////////////
// For running the language translations
/////////////////////////////////////////////////////////////////////////////

func promptOnEOL(context Context) {
	// TODO display grammar name: fmt.Printf("%s (%s) ", os.Args[0], tuple.Suffix(context))
	fmt.Printf("%s", os.Args[0])
	depth := context.Location().Depth()
	if depth > 0 {
		fmt.Printf (" %d%s", depth, PROMPT)
	} else {
		fmt.Print (PROMPT)
	}
}

type Runner struct {
	Grammars Grammars
	symbols * eval.SymbolTable
	logger LocationLogger
	inputGrammar Grammar
}

func NewRunner(grammars Grammars, symbols * eval.SymbolTable, logger LocationLogger, inputGrammar Grammar) Runner {
	
	return Runner{grammars, symbols, logger, inputGrammar}
}

func ParseAndEval(context eval.EvalContext, grammar Grammar, expression string) Value {

	var result Value = tuple.NAN
	pipeline := func(value Value) {
		result = eval.Eval(context, value)
	}
	RunParser(grammar, expression, GetLogger(nil, false), pipeline)  // TODO
	return result
}

func RunParser(grammar Grammar, expression string, logger LocationLogger, next Next) Context {

	reader := bufio.NewReader(strings.NewReader(expression))
	context := NewParserContext("<eval>", reader, logger)
	grammar.Parse(&context, next)
	return &context
}

func RunStdin(logger LocationLogger, inputGrammar Grammar, next Next) int64 {

	reader := bufio.NewReader(os.Stdin)
	context := parsers.NewParserContext2(STDIN, reader, logger, promptOnEOL)
	context.EOL() // prompt
	inputGrammar.Parse(&context, next)
	return context.Errors()
}

func (runner * Runner) RunFiles(args []string, next Next) int64 {

	if len(args) == 0 {
		return RunStdin(runner.logger, runner.inputGrammar, next)
	}
	errors := int64(0)
	for _, fileName := range args {
		suffix := path.Ext(fileName)
		grammar, ok := runner.Grammars.FindBySuffix(suffix)
		if ok {
			file, err := os.Open(fileName)
			if err != nil {
				log.Fatal(err)
			}
			reader := bufio.NewReader(file)
			context := NewParserContext(fileName, reader, runner.logger)
			grammar.Parse(&context, next)
			errors += context.Errors()
			file.Close()
		} else {
			panic("Unsupported file suffix: " + suffix)  // TODO should not be fatal
		}
	}
	return errors
}

func PrintString(value string) {
	fmt.Printf ("%s", value)
}

//
//  Set up the translator pipeline.
//
func SimplePipeline (symbols * eval.SymbolTable, queryPattern string, outputGrammar Grammar, out func(value string)) Next {

	prettyPrint := func(tuple Value) {
		outputGrammar.Print(tuple, out)
	}
	pipeline := prettyPrint
	if symbols != nil {
		next := pipeline
		pipeline = func(value Value) {
			next(eval.Eval(symbols, value))
		}
	}
	if queryPattern != "" {
		next := pipeline
		query := NewQuery(queryPattern)
		pipeline = func(value Value) {
			query.Match(value, next)
		}
	}
	return pipeline
}

func GetRemainingNonFlagOsArgs() []string {
	args := len(os.Args)
	numberOfFiles := flag.NArg()
	return os.Args[args-numberOfFiles:]
}


