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
import "flag"
import "tuple/eval"
import "tuple/parsers"

type Logger = tuple.Logger
type Grammar = tuple.Grammar
type Grammars = tuple.Grammars
type Value = tuple.Value
type Next = tuple.Next
type Context = tuple.Context
type String = tuple.String
type Int64 = tuple.Int64

var NewParserContext = parsers.NewParserContext
var GetLogger = parsers.GetLogger
const STDIN = "<stdin>"
const PROMPT = "$ "

/////////////////////////////////////////////////////////////////////////////
// For running the language translations
/////////////////////////////////////////////////////////////////////////////

func promptOnEOL(context Context) {
	// TODO display grammar name: fmt.Printf("%s (%s) ", os.Args[0], tuple.Suffix(context))
	fmt.Printf("%s", os.Args[0])
	if context.Depth() > 0 {
		fmt.Printf (" %d%s", context.Depth(), PROMPT)
	} else {
		fmt.Print (PROMPT)
	}
}

type Runner struct {
	grammars Grammars
	
}

func ParseAndEval(grammar Grammar, symbols eval.SymbolTable, expression string) Value {

	var result Value = tuple.NAN
	pipeline := func(value Value) {
		result = eval.Eval(&symbols, value)
	}
	reader := bufio.NewReader(strings.NewReader(expression))
	context := NewParserContext("<eval>", reader, GetLogger(nil), false)
	//fmt.Printf("*** Eval: '%s'\n", expression)
	grammar.Parse(&context, pipeline)
	return result
}

func RunFiles(args []string, logger Logger, verbose bool, inputGrammar Grammar, grammars *Grammars, next Next) int64 {

	errors := int64(0)
	// TODO this can be improved
	parse := func (context Context) {
		suffix := tuple.Suffix(context)
		var grammar Grammar
		if suffix == "" {
			if inputGrammar == nil {
				panic("Input grammar for '" + context.SourceName() + "' not given, use -in ...")
			}
			grammar = inputGrammar
		} else {
			grammar = grammars.FindBySuffixOrPanic(suffix)
		}
		tuple.Verbose(context,"source [%s] suffix [%s]", context.SourceName(), grammar.FileSuffix ())
		grammar.Parse(context, next)
	}

	if len(args) == 0 {
		reader := bufio.NewReader(os.Stdin)
		context := parsers.NewParserContext2(STDIN, reader, logger, verbose, promptOnEOL)
		context.EOL() // prompt
		parse(&context)
		errors += context.Errors()
		
	} else {
		for _, fileName := range args {
			file, err := os.Open(fileName)
			if err != nil {
				log.Fatal(err)
			}
			reader := bufio.NewReader(file)
			context := NewParserContext(fileName, reader, logger, verbose)
			parse(&context)
			errors += context.Errors()
			file.Close()
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


func AddSafeGrammarFunctions(table * eval.SymbolTable, grammars * Grammars) {

	table.Add("grammars", func (context eval.EvalContext, value Value) Value {
		tuple := tuple.NewTuple()
		for _,v := range grammars.All {
			tuple.Append(String(v.FileSuffix()))
		}
		return tuple
	})

//	table.Add("expr", func (context eval.EvalContext, value Value) Value {
//		grammar := parsers.NewShellGrammar()
//		return ParseAndEval(grammar, context, value)
//	})

	//table.Add("grammars", func (context eval.EvalContext, value Value) Value {
	//	return ParseAndEval(grammar, context, value)
	//})

}
