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
package parsers

import "tuple"
import 	"io"
import 	"fmt"

type Logger = tuple.Logger

/////////////////////////////////////////////////////////////////////////////
// For running the language translations
/////////////////////////////////////////////////////////////////////////////

type ParserContext struct {
	sourceName string
	line int64
	column int64
	depth int
	errors int64
	scanner io.RuneScanner
	logger Logger
	verbose bool
	eol func(context Context)
}

func NewParserContext(sourceName string, scanner io.RuneScanner, logger Logger, verbose bool) ParserContext {
	return NewParserContext2(sourceName, scanner, logger, verbose, func(context Context) {})
}

func NewParserContext2(sourceName string, scanner io.RuneScanner, logger Logger, verbose bool, eol func(context Context)) ParserContext {
	context :=  ParserContext{sourceName, 1, 0, 0, 0, scanner, logger, verbose, eol}
	tuple.Verbose(&context,"Parsing file [%s] suffix [%s]", sourceName, tuple.Suffix(&context))
	return context
}

func (context * ParserContext) Line() int64 {
	return context.line
}
func (context * ParserContext) Column() int64 {
	return context.column
}
func (context * ParserContext) Depth() int {
	return context.depth
}

func (context * ParserContext) Errors() int64 {
	return context.errors
}

func (context * ParserContext) SourceName() string {
	return context.sourceName
}

func (context * ParserContext) Open() {
	tuple.Verbose(context, "*OPEN")
	context.depth += 1
}

func (context * ParserContext) Close() {
	if context.depth > 0 {
		context.depth -= 1
	}
	tuple.Verbose(context, "*CLOSE")
}

func (context * ParserContext) EOL() {
	context.eol(context)
}

func (context * ParserContext) ReadRune() (rune, error) {
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

func (context * ParserContext) LookAhead() rune {
	ch, _, err := context.scanner.ReadRune()
	if err != nil {
		// TODO Is this okay to just return false rather than an error
		context.scanner.UnreadRune()
		return ' '
	}
	context.scanner.UnreadRune()
	return ch

}

func (context * ParserContext) Log(level string, format string, args ...interface{}) {

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
