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
import 	"io"
import 	"log"
import 	"fmt"
import 	"bufio"
import 	"os"
import 	"strings"
import "flag"
import "tuple/eval"
//import "tuple/parsers"

type Logger = tuple.Logger
type Grammar = tuple.Grammar
type Grammars = tuple.Grammars
type Value = tuple.Value
type Next = tuple.Next
type Context = tuple.Context
type String = tuple.String
type Int64 = tuple.Int64

/////////////////////////////////////////////////////////////////////////////
// For running the language translations
/////////////////////////////////////////////////////////////////////////////

const STDIN = "<stdin>"
const PROMPT = "$ "

func IsInteractive(context Context) bool {
	return context.SourceName() == STDIN
}

type RunnerContext struct {
	sourceName string
	line int64
	column int64
	depth int
	errors int64
	scanner io.RuneScanner
	logger Logger
	verbose bool
}

func NewRunnerContext(sourceName string, scanner io.RuneScanner, logger Logger, verbose bool) RunnerContext {
	context :=  RunnerContext{sourceName, 1, 0, 0, 0, scanner, logger, verbose}
	tuple.Verbose(&context,"Parsing file [%s] suffix [%s]", sourceName, tuple.Suffix(&context))
	return context
}

func (context * RunnerContext) Line() int64 {
	return context.line
}
func (context * RunnerContext) Column() int64 {
	return context.column
}
func (context * RunnerContext) Depth() int {
	return context.depth
}

func (context * RunnerContext) Errors() int64 {
	return context.errors
}

func (context * RunnerContext) SourceName() string {
	return context.sourceName
}

func (context * RunnerContext) Open() {
	tuple.Verbose(context, "*OPEN")
	context.depth += 1
}

func (context * RunnerContext) Close() {
	if context.depth > 0 {
		context.depth -= 1
	}
	tuple.Verbose(context, "*CLOSE")
}

func (context * RunnerContext) EOL() {
	if IsInteractive(context) {

		// TODO display grammar name: fmt.Printf("%s (%s) ", os.Args[0], tuple.Suffix(context))
		fmt.Printf("%s", os.Args[0])
		if context.depth > 0 {
			fmt.Printf (" %d%s", context.depth, PROMPT)
		} else {
			fmt.Print (PROMPT)
		}
	}
}

func (context * RunnerContext) ReadRune() (rune, error) {
	ch, _, err := context.scanner.ReadRune()
	switch {
	case err != nil: return ch, err
	case ch == '\n':
		context.line ++
		context.column = 0
		tuple.Verbose(context,"New line")
	default:
		context.column ++
	}
	return ch, nil
}

func (context * RunnerContext) LookAhead() rune {
	ch, _, err := context.scanner.ReadRune()
	if err != nil {
		// TODO Is this okay to just return false rather than an error
		context.scanner.UnreadRune()
		return ' '
	}
	context.scanner.UnreadRune()
	return ch

}

func (context * RunnerContext) Log(level string, format string, args ...interface{}) {

	switch level {
	case "VERBOSE":
		if ! context.verbose {
			return
		}
	case "ERROR": context.errors += 1
	default:
	}
	suffix := fmt.Sprintf(format, args...)
	context.logger(context, level, suffix)
}

func GetLogger(logGrammar Grammar) Logger {
	if logGrammar == nil {
		return func (context Context, level string, message string) {
			prefix := fmt.Sprintf("%s at %d, %d depth=%d in '%s': %s", level, context.Line(), context.Column(), context.Depth(), context.SourceName(), message)
			log.Print(prefix)
		}
	} else {
		return func(context Context, level string, message string) {
			record := tuple.NewTuple()
			record.Append(String(level))
			record.Append(Int64(context.Line()))
			record.Append(Int64(context.Column()))
			record.Append(Int64(context.Depth()))
			record.Append(String(context.SourceName()))
			record.Append(String(message))
			logGrammar.Print(record, func (value string) { fmt.Print(value) })
		}
	}
}

func ParseAndEval(grammar Grammar, symbols eval.SymbolTable, expression string) Value {

	var result Value = tuple.NAN
	pipeline := func(value Value) {
		result = eval.Eval(&symbols, value)
	}
	reader := bufio.NewReader(strings.NewReader(expression))
	context := NewRunnerContext("<eval>", reader, GetLogger(nil), false)
	//fmt.Printf("*** Eval: '%s'\n", expression)
	grammar.Parse(&context, pipeline)
	return result
}

func ParseString(grammar Grammar, expression string) Value {
	var result Value = tuple.NAN
	pipeline := func(value Value) {
		result = value
	}

	reader := bufio.NewReader(strings.NewReader(expression))
	context := NewRunnerContext("<parse>", reader, GetLogger(nil), false)
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
		context := NewRunnerContext(STDIN, reader, logger, verbose)
		tuple.Verbose(&context, "STDIN isinteractive: %s", IsInteractive(&context))
		context.EOL() // prompt
		parse(&context)
		errors += context.errors
		
	} else {
		for _, fileName := range args {
			file, err := os.Open(fileName)
			if err != nil {
				log.Fatal(err)
			}
			reader := bufio.NewReader(file)
			context := NewRunnerContext(fileName, reader, logger, verbose)
			parse(&context)
			errors += context.errors
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


func AddSafeGrammarFunctions(table * eval.SymbolTable) {

//	table.Add("expr", func (context eval.EvalContext, value Value) Value {
//		grammar := parsers.NewShellGrammar()
//		return ParseAndEval(grammar, context, value)
//	})

	//table.Add("grammars", func (context eval.EvalContext, value Value) Value {
	//	return ParseAndEval(grammar, context, value)
	//})

}
