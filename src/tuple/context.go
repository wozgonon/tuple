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
import 	"path"

const PROMPT = "$ "
const STDIN = "<stdin>"

type Next func(value interface{})

type ParserContext struct {
	SourceName string
	line int64
	column int64
	depth int
	scanner io.RuneScanner
	logGrammar Grammar
	verbose bool
	next Next

}

func NewParserContext(sourceName string, scanner io.RuneScanner, logGrammar Grammar, verbose bool, next Next) ParserContext {
	context :=  ParserContext{sourceName, 1, 0, 0, scanner, logGrammar, verbose, next}
	context.Verbose("Parsing file [%s] suffix [%s]", sourceName, context.Suffix())
	return context
}

func (context * ParserContext) Open() {
	context.depth += 1
}

func (context * ParserContext) Close() {
	context.depth -= 1
}

func (context * ParserContext) IsInteractive() bool {
	return context.SourceName == STDIN
}

func (context * ParserContext) Suffix() string {
	return path.Ext(context.SourceName)
}

func (context * ParserContext) prompt() {
	if context.IsInteractive() {
		fmt.Print (context.SourceName)
		if context.depth > 0 {
			fmt.Printf (" %d%s", context.depth, PROMPT)
		} else {
			fmt.Print (PROMPT)
		}
	}
}

func (context * ParserContext) ReadRune() (rune, error) {
	ch, _, err := context.scanner.ReadRune()
	switch {
	case err != nil: return ch, err
	case ch == '\n':
		context.line ++
		context.column = 0
		context.Verbose("New line")
		context.prompt()
	default:
		context.column ++
	}
	return ch, nil
}

func (context * ParserContext) UnreadRune() {
	context.scanner.UnreadRune()
	if context.column == 0 {
		context.line --
	} else {
		context.column --
	}
}

func (context * ParserContext) log(format string, level string, args ...interface{}) {
	prefix := fmt.Sprintf("%s at %d, %d depth=%d in '%s': ", level, context.line, context.column, context.depth, context.SourceName)
	suffix := fmt.Sprintf(format, args...)
	log.Print(prefix + suffix)

	tuple := NewTuple()
	tuple.Append(level)
	tuple.Append(context.line)
	tuple.Append(context.column)
	tuple.Append(int64(context.depth))
	tuple.Append(context.SourceName)
	tuple.Append(suffix)
	// TODO context.logGrammar.Print(tuple, func (value string) { fmt.Print(value) })

}

func (context * ParserContext) Error(format string, args ...interface{}) {
	context.log(format, "ERROR", args...)
}

func (context * ParserContext) UnexpectedCloseBracketError(token string) {
	context.Error ("Unexpected close bracket '%s'", token)
}

func (context * ParserContext) UnexpectedEndOfInputErrorBracketError() {
	context.Error ("Unexpected end of input")
}

func (context * ParserContext) Verbose(format string, args ...interface{}) {
	if context.verbose {
		context.log(format, "VERBOSE", args...)
	}
}

func RunParser(args []string, logGrammar Grammar, verbose bool, inputGrammar * Grammar, grammars *Grammars, next Next) {

	// TODO this can be improved
	parse := func (context * ParserContext) {
		suffix := context.Suffix()
		var grammar *Grammar
		if suffix == "" {
			if inputGrammar == nil {
				panic("Input grammar for '" + context.SourceName + "' not given, use -in ...")
			}
			grammar = inputGrammar
		} else {
			grammar = grammars.FindBySuffixOrPanic(suffix)
		}
		context.Verbose("source [%s] suffix [%s]", context.SourceName, (*grammar).FileSuffix ())
		(*grammar).Parse(context)
	}

	if len(args) == 0 {
		reader := bufio.NewReader(os.Stdin)
		context := NewParserContext(STDIN, reader, logGrammar, verbose, next)
		context.prompt()
		parse(&context)
	} else {
		for _, fileName := range args {
			file, err := os.Open(fileName)
			if err != nil {
				log.Fatal(err)
			}
			reader := bufio.NewReader(file)
			context := NewParserContext(fileName, reader, logGrammar, verbose, next)
			parse(&context)
			file.Close()
		}
	}
}


//
//  Set up the translator pipeline.
//
func SimplePipeline (eval bool, queryPattern string, outputGrammar * Grammar) Next {

	prettyPrint := func(tuple interface{}) {
		(*outputGrammar).Print(tuple, func(value string) {
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
