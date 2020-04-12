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
import 	"fmt"
import 	"bufio"
import 	"os"
import "errors"
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

/////////////////////////////////////////////////////////////////////////////

func NewSafeEvalContext(logger LocationLogger) eval.EvalContext {
	ifNotFound := eval.NewErrorIfFunctionNotFound()
	runner := eval.NewRunner(ifNotFound, logger)
	eval.AddSafeFunctions(&runner)
	return &runner
}

func NewHarmlessEvalContext(logger LocationLogger) eval.EvalContext {
	ifNotFound := eval.NewErrorIfFunctionNotFound()
	runner := eval.NewRunner(ifNotFound, logger)
	eval.AddHarmlessFunctions(&runner)
	return &runner
}

/////////////////////////////////////////////////////////////////////////////

func LocationForValue(value Value) tuple.Location {
	// TODO get the location associate with a Value
	return tuple.NewLocation("<eval>", 0, 0, 0) // TODO
}

/////////////////////////////////////////////////////////////////////////////

func ParseAndEval(context eval.EvalContext, grammar Grammar, expression string) (Value, error) {

	var result Value = tuple.NAN  // TODO ought to be EMPTY
	pipeline := func(value Value) error {
		evaluated, err := eval.Eval(context, value)
		if err != nil {
			return err
		}
		result = evaluated
		return nil
	}
	ctx, err := parsers.RunParser(grammar, expression, context.LocationLogger(), pipeline)
	if ctx.Errors() > 0 {
		return nil, errors.New("Errors during parse")
	}
	return result, err
}

func RunStdin(logger LocationLogger, inputGrammar Grammar, next Next) (Context, error) {
	reader := bufio.NewReader(os.Stdin)
	context := parsers.NewParserContext2(STDIN, reader, logger, promptOnEOL)
	context.EOL() // prompt
	err := inputGrammar.Parse(&context, next)
	return &context, err
}

func RunFile(grammars * Grammars, locationLogger LocationLogger, fileName string, next Next) (Context, error) {
	suffix := path.Ext(fileName)
	grammar, ok := grammars.FindBySuffix(suffix)
	if ! ok {
		return nil, errors.New("Unsupported file suffix: " + suffix)
	}
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	context := NewParserContext(fileName, reader, locationLogger)
	err = grammar.Parse(&context, next)
	file.Close()
	return &context, err
}

func RunFiles(grammars * Grammars, locationLogger LocationLogger, args []string, next Next) (int64) {
	errors := int64(0)
	if len(args) == 0 {
		 context, err := RunStdin(locationLogger, grammars.Default(), next)
		if err != nil {
			tuple.Error(context, "%s", err)
		}
		errors += context.Errors()
	} else {
		for _, fileName := range args {
			context, err := RunFile(grammars, locationLogger, fileName, next)
			if err != nil {
				tuple.Error(context, "%s", err)
				break
			}
			errors += context.Errors()
		}
	}
	return errors
}

func PrintString(value string) {
	fmt.Printf ("%s", value)
}

//  Set up the translator pipeline.
func SimplePipeline (context eval.EvalContext, runEval bool, queryPattern string, outputGrammar Grammar, out func(value string)) Next {

	prettyPrint := func(tuple Value) error {
		outputGrammar.Print(tuple, out)
		return nil
	}
	pipeline := prettyPrint
	if runEval {
		next := pipeline
		pipeline = func(value Value) error {
			evaluated, err := eval.Eval(context, value)
			if err != nil {
				return err
			}
			next(evaluated)
			return nil
		}
	}
	if queryPattern != "" {
		next := pipeline
		query := NewQuery(queryPattern)
		pipeline = func(value Value) error {
			query.Match(context, value, next)
			return nil
		}
	}
	return pipeline
}

func GetRemainingNonFlagOsArgs() []string {
	args := len(os.Args)
	numberOfFiles := flag.NArg()
	return os.Args[args-numberOfFiles:]
}


