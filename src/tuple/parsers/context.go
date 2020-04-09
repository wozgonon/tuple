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

type Location = tuple.Location
type LocationLogger = tuple.LocationLogger

/////////////////////////////////////////////////////////////////////////////
// For running the language translations
/////////////////////////////////////////////////////////////////////////////


type ParserContext struct {
	location Location
	errors int64
	scanner io.RuneScanner
	logger LocationLogger
	eolCallback func(context Context)
}

func NewParserContext(sourceName string, scanner io.RuneScanner, logger LocationLogger) ParserContext {
	return NewParserContext2(sourceName, scanner, logger, func(context Context) {})
}

func NewParserContext2(sourceName string, scanner io.RuneScanner, logger LocationLogger, eol func(context Context)) ParserContext {
	initialLocation := tuple.NewLocation(sourceName, 1, 0, 0)
	context :=  ParserContext{initialLocation, 0, scanner, logger, eol}
	tuple.Verbose(&context,"Parsing file [%s] suffix [%s]", sourceName, tuple.Suffix(&context))
	return context
}

func (context * ParserContext) Location() Location {
	return context.location
}

func (context * ParserContext) Errors() int64 {
	return context.errors
}

func (context * ParserContext) Open() {
	tuple.Verbose(context, "*OPEN")
	context.location.IncrDepth()
}

func (context * ParserContext) Close() {
	context.location.DecrDepth()
	tuple.Verbose(context, "*CLOSE")
}

func (context * ParserContext) EOL() {
	context.eolCallback(context)
}

func (context * ParserContext) ReadRune() (rune, error) {
	ch, _, err := context.scanner.ReadRune()
	switch {
	case err != nil: return ch, err
	case ch == '\n':
		context.location.IncrLine()
		tuple.Verbose(context,"New line")
	default:
		context.location.IncrColumn()
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
	case "ERROR": context.errors += 1
	default:
	}
	suffix := fmt.Sprintf(format, args...)
	context.logger(context.location, level, suffix)
}
