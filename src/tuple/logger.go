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

import "log"
import "fmt"

/////////////////////////////////////////////////////////////////////////////
// Logger
/////////////////////////////////////////////////////////////////////////////

type Logger interface {
	Log (level string, format string, args ...interface{})
}

func Verbose(logger Logger, format string, args ...interface{}) {
	logger.Log("VERBOSE", format, args...)
}

func Trace(logger Logger, format string, args ...interface{}) {
	logger.Log("TRACE", format, args...)
}

func Error(logger Logger, format string, args ...interface{}) {
	AssertNotNil(logger)
	logger.Log("ERROR", format, args...)
}

/////////////////////////////////////////////////////////////////////////////
// Location
/////////////////////////////////////////////////////////////////////////////

type LocationLogger func (location Location, level string, message string)

// The Context interface represents the current state of parsing and translation.
// It can provide: the name of the input and current depth and number of errors
type Location struct {
	sourceName string
	line int64
	column int64
	depth int
}

func NewLocation(sourceName string, line int64, column int64, depth int) Location {
	return Location{sourceName, 1, 0, 0}
}

func (location Location) SourceName() string {
	return location.sourceName
}
func (location Location) Line() int64 {
	return location.line
}
func (location Location) Column() int64 {
	return location.column
}
func (location Location) Depth() int {
	return location.depth
}
func (location * Location) IncrLine() {
	location.line += 1
	location.column = 0
}
func (location * Location) IncrColumn() {
	location.column += 1
}
func (location * Location) IncrDepth() {
	location.depth += 1
}
func (location * Location) DecrDepth() {
	if location.depth > 0 {
		location.depth -= 1
	}
}

/////////////////////////////////////////////////////////////////////////////

func GetLogger(logGrammar Grammar, verbose bool) LocationLogger {

	var logger LocationLogger
	if logGrammar == nil {
		logger = func (context Location, level string, message string) {
			prefix := fmt.Sprintf("%s at %d, %d depth=%d in '%s': %s", level, context.Line(), context.Column(), context.Depth(), context.SourceName(), message)
			log.Print(prefix)
		}
	} else {
		logger = func(context Location, level string, message string) {
			record := NewTuple()
			record.Append(String(level))
			record.Append(Int64(context.Line()))
			record.Append(Int64(context.Column()))
			record.Append(Int64(context.Depth()))
			record.Append(String(context.SourceName()))
			record.Append(String(message))
			logGrammar.Print(record, func (value string) { fmt.Print(value) })
		}
	}

	if verbose {
		return logger
	} else {
		return func (context Location, level string, message string) {
			if level != "VERBOSE" && level != "TRACE" {
				logger(context, level, message)
			}
		}
	}
}
