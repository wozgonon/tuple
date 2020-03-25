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
package tuple

import 	"io"
import 	"log"
import 	"fmt"
import 	"bufio"
import 	"os"

const PROMPT = "$ "

type ParserContext struct {
	sourceName string
	line int64
	column int64
	depth int
	errors int64
	scanner io.RuneScanner
	logger Logger
	verbose bool
}

func NewParserContext(sourceName string, scanner io.RuneScanner, logger Logger, verbose bool) ParserContext {
	context :=  ParserContext{sourceName, 1, 0, 0, 0, scanner, logger, verbose}
	Verbose(context,"Parsing file [%s] suffix [%s]", sourceName, Suffix(context))
	return context
}

func (context ParserContext) Line() int64 {
	return context.line
}
func (context ParserContext) Column() int64 {
	return context.column
}
func (context ParserContext) Depth() int {
	return context.depth
}

func (context ParserContext) Errors() int64 {
	return context.errors
}

func (context ParserContext) SourceName() string {
	return context.sourceName
}

func (context ParserContext) Open() {
	context.depth += 1
}

func (context ParserContext) Close() {
	context.depth -= 1
}

func (context ParserContext) prompt() {
	if IsInteractive(context) {
		fmt.Print (context.SourceName)
		if context.depth > 0 {
			fmt.Printf (" %d%s", context.depth, PROMPT)
		} else {
			fmt.Print (PROMPT)
		}
	}
}

func (context ParserContext) ReadRune() (rune, error) {
	ch, _, err := context.scanner.ReadRune()
	switch {
	case err != nil: return ch, err
	case ch == '\n':
		context.line ++
		context.column = 0
		Verbose(context,"New line")
		context.prompt()
	default:
		context.column ++
	}
	return ch, nil
}

func (context ParserContext) UnreadRune() {
	context.scanner.UnreadRune()
	if context.column == 0 {
		context.line --
	} else {
		context.column --
	}
}

func (context ParserContext) Log(format string, level string, args ...interface{}) {

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
			prefix := fmt.Sprintf("%s at %d, %d depth=%d in '%s': ", level, context.Line(), context.Column(), context.Depth(), context.SourceName(), message)
			log.Print(prefix)
		}
	} else {
		return func(context Context, level string, message string) {
			tuple := NewTuple()
			tuple.Append(level)
			tuple.Append(context.Line())
			tuple.Append(context.Column())
			tuple.Append(int64(context.Depth()))
			tuple.Append(context.SourceName())
			tuple.Append(message)
			logGrammar.Print(tuple, func (value string) { fmt.Print(value) })
		}
	}
}

func RunParser(args []string, logger Logger, verbose bool, inputGrammar Grammar, grammars *Grammars, next Next) int64 {

	errors := int64(0)
	// TODO this can be improved
	parse := func (context Context) {
		suffix := Suffix(context)
		var grammar Grammar
		if suffix == "" {
			if inputGrammar == nil {
				panic("Input grammar for '" + context.SourceName() + "' not given, use -in ...")
			}
			grammar = inputGrammar
		} else {
			grammar = grammars.FindBySuffixOrPanic(suffix)
		}
		Verbose(context,"source [%s] suffix [%s]", context.SourceName(), grammar.FileSuffix ())
		grammar.Parse(context, next)
	}

	if len(args) == 0 {
		reader := bufio.NewReader(os.Stdin)
		context := NewParserContext(STDIN, reader, logger, verbose)
		context.prompt()
		parse(context)
		errors += context.errors
		
	} else {
		for _, fileName := range args {
			file, err := os.Open(fileName)
			if err != nil {
				log.Fatal(err)
			}
			reader := bufio.NewReader(file)
			context := NewParserContext(fileName, reader, logger, verbose)
			parse(context)
			errors += context.errors
			file.Close()
		}
	}
	return errors
}

//
//  Set up the translator pipeline.
//
func SimplePipeline (eval bool, queryPattern string, outputGrammar Grammar) Next {

	prettyPrint := func(tuple interface{}) {
		outputGrammar.Print(tuple, func(value string) {
			fmt.Printf ("%s", value)
		})
	}
	pipeline := prettyPrint
	if eval {
		next := pipeline
		pipeline = func(value interface{}) {
			SimpleEval(value, next)
		}
	}
	if queryPattern != "" {
		next := pipeline
		query := NewQuery(queryPattern)
		pipeline = func(value interface{}) {
			query.Match(value, next)
		}
	}
	return pipeline
}
