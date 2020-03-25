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

import "strings"
import "log"
import "reflect"
import 	"path"

const STDIN = "<stdin>"

/////////////////////////////////////////////////////////////////////////////
//  Callback signatures
/////////////////////////////////////////////////////////////////////////////

type StringFunction func(value string)
type Next func(value interface{})
type Logger func(context Context, level string, message string)

/////////////////////////////////////////////////////////////////////////////
//  Lexer
/////////////////////////////////////////////////////////////////////////////

type Lexer interface {
	Printer
	GetNext(context Context, open func(open string), close func(close string), nextAtom func(atom Atom), nextLiteral func (literal interface{})) error
}

/////////////////////////////////////////////////////////////////////////////
//  Context
/////////////////////////////////////////////////////////////////////////////


// The Context interface represents the current state of parsing and translation.
// It can provide: the name of the input and current depth and number of errors
type Context interface {
	SourceName() string
	Line() int64
	Column() int64
	Depth() int
	Open()
	Close()
	ReadRune() (rune, error)
	UnreadRune()
	Log(format string, level string, args ...interface{})
	Errors() int64
}

func Verbose(context Context, format string, args ...interface{}) {
	context.Log(format, "VERBOSE", args...)
}

func Error(context Context, format string, args ...interface{}) {
	context.Log(format, "ERROR", args...)
}

func UnexpectedCloseBracketError(context Context, token string) {
	Error(context,"Unexpected close bracket '%s'", token)
}

func UnexpectedEndOfInputErrorBracketError(context Context) {
	Error(context,"Unexpected end of input")
}

func IsInteractive(context Context) bool {
	return context.SourceName() == STDIN
}

func Suffix(context Context) string {
	return path.Ext(context.SourceName())
}


/////////////////////////////////////////////////////////////////////////////
//  Grammar
/////////////////////////////////////////////////////////////////////////////

// The Grammar interface represents a particular language Grammar or Grammar or File Format.
//
// The print and parse method ought to be inverse functions of each other
// so the output of parse can be passed to print which in principle should be parsable by the parse function.
//
type Grammar interface {
	// A friendly name for the syntax
	Name() string

	// A standard suffix for source files.
	FileSuffix() string
	
	// Parses an input stream of characters into an internal representation (AST)
	// The output ought to be printable by the 'print' method.
	Parse(context Context, next Next) // , next func(tuple Tuple)) (interface{}, error)
	
	// Pretty prints the objects in the given syntax.
	// The output ought to be parsable by the 'parse' method.
	Print(token interface{}, next func(value string))
}

// A set of Grammars
type Grammars struct {
	all map[string]Grammar
}

// Returns a new empty set of syntaxes
func NewGrammars() Grammars{
	return Grammars{make(map[string]Grammar)}
}

func (syntaxes * Grammars) Add(syntax Grammar) {
	suffix := syntax.FileSuffix()
	syntaxes.all[suffix] = syntax
}

func (syntaxes * Grammars) FindBySuffix(suffix string) (Grammar, bool) {
	if ! strings.HasPrefix(suffix, ".") {
		suffix = "." + suffix
	}
	syntax, ok := syntaxes.all[suffix]
	return syntax, ok
}

func (syntaxes * Grammars) FindBySuffixOrPanic(suffix string) Grammar {
	syntax, ok := syntaxes.FindBySuffix(suffix)
	if ! ok {
		panic("Unsupported file suffix: '" + suffix + "'")
	}
	return syntax
}


/////////////////////////////////////////////////////////////////////////////
//  Printer
/////////////////////////////////////////////////////////////////////////////

type Printer interface {

	PrintIndent(depth string, out StringFunction)
	PrintSuffix(depth string, out StringFunction)
	PrintScalarPrefix(depth string, out StringFunction)
	PrintSeparator(depth string, out StringFunction)
	PrintEmptyTuple(depth string, out StringFunction)
	PrintNullaryOperator(depth string, atom Atom, out StringFunction)
	PrintUnaryOperator(depth string, atom Atom, value interface{}, out StringFunction)
	PrintBinaryOperator(depth string, atom Atom, value1 interface{}, value2 interface{}, out StringFunction)
	PrintOpenTuple(depth string, tuple Tuple, out StringFunction) string
	PrintCloseTuple(depth string, tuple Tuple, out StringFunction)
	PrintHeadAtom(value Atom, out StringFunction)
	PrintAtom(depth string, value Atom, out StringFunction)
	PrintInt64(depth string, value int64, out StringFunction)
	PrintFloat64(depth string, value float64, out StringFunction)
	PrintString(depth string, value string, out StringFunction)
	PrintBool(depth string, value bool, out StringFunction)
	PrintComment(depth string, value Comment, out StringFunction)
}

func PrintScalar(printer Printer, depth string, token interface{}, out StringFunction) {
	printer.PrintScalarPrefix(depth, out)
	switch token.(type) {
	case Atom:
		printer.PrintAtom(depth, token.(Atom), out)
	case string:
		printer.PrintString(depth, token.(string), out)
	case bool:
		printer.PrintBool(depth, token.(bool), out)
	case Comment:
		printer.PrintComment(depth, token.(Comment), out)
	case int64:
		printer.PrintInt64(depth, token.(int64), out)
	case float64:
		printer.PrintFloat64(depth, token.(float64), out)
	default:
		log.Printf("ERROR type '%s' not recognised: %s", reflect.TypeOf(token), token);
	}
}

func PrintTuple(printer Printer, depth string, tuple Tuple, out StringFunction) {
	newDepth := printer.PrintOpenTuple(depth, tuple, out)
	printer.PrintSuffix(depth, out)
	ll := len(tuple.List)
	first := false
	if ll > 0 {
		_, first = tuple.List[0].(Atom)
	}
	for k, value := range tuple.List {
		printer.PrintIndent(newDepth, out)
		if first && k == 0 {
			printer.PrintHeadAtom(value.(Atom), out)
		} else {
			PrintExpression1(printer, newDepth, value, out)
		}
		if k < ll-1 {
			printer.PrintSeparator(newDepth, out)
		}
		printer.PrintSuffix(depth, out)
	}
	printer.PrintCloseTuple(depth, tuple, out)
}

func PrintExpression(printer Printer, depth string, token interface{}, out StringFunction) {
	printer.PrintIndent(depth, out)
	PrintExpression1(printer, depth, token, out)
	printer.PrintSuffix(depth, out)
}

func PrintExpression1(printer Printer, depth string, token interface{}, out StringFunction) {

	switch token.(type) {
	case Tuple:
		tuple := token.(Tuple)
		len := len(tuple.List)
		if len == 0 {
			printer.PrintScalarPrefix(depth, out)
			printer.PrintEmptyTuple(depth, out)
		} else {
			head := tuple.List[0]
			atom, ok := head.(Atom)
			//log.Printf("Tuple [%s] %d\n", atom, len)
			if ok {  // TODO and head in a (binary) operator
				switch len {
				case 1:
					printer.PrintNullaryOperator(depth, atom, out)
				case 2:
					printer.PrintUnaryOperator(depth, atom, tuple.List[1], out)
				case 3:
					printer.PrintBinaryOperator(depth, atom, tuple.List[1], tuple.List[2], out)
				default:
					PrintTuple(printer, depth, tuple, out)
				}
			} else {
				PrintTuple(printer, depth, tuple, out)
			}
		}
	default:
		PrintScalar(printer, depth, token, out)
	}
}
