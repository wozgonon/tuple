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
import "path"
import "math"
//import "fmt"
//import "strconv"

/////////////////////////////////////////////////////////////////////////////
//  Callback signatures
/////////////////////////////////////////////////////////////////////////////

type StringFunction func(value string)
type Next func(value Value)
type Logger func(context Context, level string, message string)

var CONS_ATOM = Atom{"cons"}
var NAN Float64 = Float64(math.NaN())
var EMPTY Tuple = NewTuple()

/////////////////////////////////////////////////////////////////////////////
//  Lexer and Values
/////////////////////////////////////////////////////////////////////////////

type Lexer interface {
	Printer
	GetNext(context Context, eol func(), open func(open string), close func(close string), nextAtom func(atom Atom), nextLiteral func (literal Value)) error
}

// The Value interface must be implected by any of the small number of types that can be produced by the lexer
type Value interface {
	// See: https://en.wikipedia.org/wiki/Arity
	Arity() int
}

type String string
type Float64 float64
type Int64 int64
type Bool bool

// An Atom - a name for something, an identifier or operator
// TODO include location and source for editors
type Atom struct {
	Name string
}

func (atom Atom) Arity() int { return 0 }
func (value String) Arity() int { return 0 }
func (comment Comment) Arity() int { return 0 }
func (value Float64) Arity() int { return 0 }
func (value Int64) Arity() int { return 0 }
func (value Bool) Arity() int { return 0 }

// A textual comment
type Comment struct {
	// TODO include location and source for editors
	Comment string
}

func NewComment(_ Context, token string) Comment {
	return Comment{token}
}

func (tuple Tuple) Arity() int { return len(tuple.List) }

/////////////////////////////////////////////////////////////////////////////
// Tuple
/////////////////////////////////////////////////////////////////////////////

type Tuple struct {
	// TODO include location and source for editors
	List []Value
}

func (tuple *Tuple) Append(token Value) {
	tuple.List = append(tuple.List, token)
}

func (tuple *Tuple) Length() int {
	return len(tuple.List)
}

func (tuple *Tuple) IsCons() bool {
	if tuple.Length() > 0 {
		head := tuple.List[0]
		atom, ok := head.(Atom)
		if ok && atom == CONS_ATOM {
			return true
		}
	}
	return false
}

// TODO this may not make sense cons is embedded in another tuple
func (tuple *Tuple) IsConsInTuple() bool {
	cons := false
	if tuple.Length() > 0 {
		t, ok := tuple.List[0].(Tuple)
		cons = ok && t.IsCons()
	}
	return cons
}

func NewTuple(values...Value) Tuple {
	return Tuple{values}
}

/////////////////////////////////////////////////////////////////////////////
//  Context
/////////////////////////////////////////////////////////////////////////////

// The Context interface represents the current state of parsing and translation.
// It can provide: the name of the input and current depth and number of errors
type LocationContext interface {
	SourceName() string
	Line() int64
	Column() int64
	Depth() int
	Log(level string, format string, args ...interface{})
}

// The Context interface represents the current state of parsing and translation.
// It can provide: the name of the input and current depth and number of errors
type Context interface {
	// TODO LocationContext  and change name to ParseContext
	SourceName() string
	Line() int64
	Column() int64
	Depth() int
	Open()
	Close()
	EOL()
	ReadRune() (rune, error)
	LookAhead() rune
	Log(level string, format string, args ...interface{})
	Errors() int64
}

func Verbose(context Context, format string, args ...interface{}) {
	context.Log("VERBOSE", format, args...)
}

func Error(context Context, format string, args ...interface{}) {
	context.Log("ERROR", format, args...)
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
	Parse(context Context, next Next) // , next func(tuple Tuple)) (Value, error)
	
	// Pretty prints the objects in the given syntax.
	// The output ought to be parsable by the 'parse' method.
	Print(token Value, next func(value string))
}

// A set of Grammars
type Grammars struct {
	all map[string]Grammar
}

// Returns a new empty set of syntaxes
func NewGrammars() Grammars{
	return Grammars{make(map[string]Grammar)}
}

func (syntaxes * Grammars) Forall(next func(grammar Grammar)) {
	for _, grammar := range syntaxes.all {
		next (grammar)
	}
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
	PrintUnaryOperator(depth string, atom Atom, value Value, out StringFunction)
	PrintBinaryOperator(depth string, atom Atom, value1 Value, value2 Value, out StringFunction)
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

func PrintScalar(printer Printer, depth string, token Value, out StringFunction) {
	printer.PrintScalarPrefix(depth, out)
	switch token.(type) {
	case Atom:
		printer.PrintAtom(depth, token.(Atom), out)
	case String:
		printer.PrintString(depth, string(token.(String)), out)
	case Bool:
		printer.PrintBool(depth, bool(token.(Bool)), out)
	case Comment:
		printer.PrintComment(depth, token.(Comment), out)
	case Int64:
		printer.PrintInt64(depth, int64(token.(Int64)), out)
	case Float64:
		printer.PrintFloat64(depth, float64(token.(Float64)), out)
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

func PrintExpression(printer Printer, depth string, token Value, out StringFunction) {
	printer.PrintIndent(depth, out)
	PrintExpression1(printer, depth, token, out)
	printer.PrintSuffix(depth, out)
}

func PrintExpression1(printer Printer, depth string, token Value, out StringFunction) {

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
