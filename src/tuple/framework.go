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
import "reflect"
import "path"
import "math"

/////////////////////////////////////////////////////////////////////////////
//  Callback signatures
/////////////////////////////////////////////////////////////////////////////

type StringFunction func(value string)
type Next func(value Value)

var CONS_ATOM = Tag{"cons"}
var NAN Float64 = Float64(math.NaN())
var EMPTY Tuple = NewTuple()

/////////////////////////////////////////////////////////////////////////////
//  Lexer and Values
/////////////////////////////////////////////////////////////////////////////

type Lexer interface {
	Printer
	GetNext(context Context, eol func(), open func(open string), close func(close string), nextTag func(tag Tag), nextLiteral func (literal Value)) error
}

// The Value interface must be implected by any of the small number of types that can be produced by the lexer
type Value interface {
	// See: https://en.wikipedia.org/wiki/Arity
	Arity() int
	Get(index int) Value
//	GetKeyValue(index int) (Value, Value)  // Key should probably be restricted to Scalars/Atoms
}

type Scalar interface {
	Value
	ToString() String
}

func IsAtom(value Value) bool {
	return value.Arity() == 0
}

func Forall(value Value, next func(value Value) error) error {
	ll := value.Arity()
	for k := 0; k < ll; k+=1 {
		err := next(value.Get(k))
		if err != nil {
			return err
		}
	}
	return nil
}

func ForallN(value Value, next func(index int, value Value) error) error {
	ll := value.Arity()
	for k := 0; k < ll; k+=1 {
		err := next(k, value.Get(k))
		if err != nil {
			return err
		}
	}
	return nil
}

func  IsCons(value Value) bool {
	if value.Arity() > 0 {
		head := value.Get(0)
		tag, ok := head.(Tag)
		if ok && tag == CONS_ATOM {
			return true
		}
	}
	return false
}

// TODO this may not make sense cons is embedded in another tuple
func IsConsInTuple(value Value) bool {
	return value.Arity() > 0 && IsCons(value.Get(0))
}

type String string
type Float64 float64
type Int64 int64
type Bool bool

// An Tag - a name for something, an identifier or operator
// TODO include location and source for editors
type Tag struct {
	Name string
}

func (tag Tag) Arity() int { return 0 }
func (value String) Arity() int { return 0 }
func (comment Comment) Arity() int { return 0 }
func (value Float64) Arity() int { return 0 }
func (value Int64) Arity() int { return 0 }
func (value Bool) Arity() int { return 0 }
func (tag Tag) Get(_ int) Value { return EMPTY }
func (value String) Get(_ int) Value { return EMPTY }
func (comment Comment) Get(_ int) Value { return EMPTY }
func (value Float64) Get(_ int) Value { return EMPTY }
func (value Int64) Get(_ int) Value { return EMPTY }
func (value Bool) Get(_ int) Value { return EMPTY }

// A textual comment
type Comment struct {
	// TODO include location and source for editors
	Comment string
}

func NewComment(_ Context, token string) Comment {
	return Comment{token}
}

/////////////////////////////////////////////////////////////////////////////
// Tuple
/////////////////////////////////////////////////////////////////////////////

type Tuple struct {
	// TODO include location and source for editors
	List []Value  // TODO change to elements
}

func NewTuple(values... Value) Tuple {
	return Tuple{values}
}

func (tuple Tuple) Arity() int { return len(tuple.List) }
func (tuple Tuple) Get(index int) Value { return tuple.List[index] }


func (tuple *Tuple) Append(token Value) {
	tuple.List = append(tuple.List, token)
}

/////////////////////////////////////////////////////////////////////////////
//  Context
/////////////////////////////////////////////////////////////////////////////

// The Context interface represents the current state of parsing and translation.
// It can provide: the name of the input and current depth and number of errors
// TODO  change name to ParseContext
type Context interface {
	Location() Location
	
	Open()
	Close()
	EOL()
	ReadRune() (rune, error)
	LookAhead() rune
	Log(level string, format string, args ...interface{})
	Errors() int64
}

func Suffix(context Context) string {
	return path.Ext(context.Location().SourceName())
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

/////////////////////////////////////////////////////////////////////////////
//  Printer
/////////////////////////////////////////////////////////////////////////////

type Printer interface {

	PrintIndent(depth string, out StringFunction)
	PrintSuffix(depth string, out StringFunction)
	PrintScalarPrefix(depth string, out StringFunction)
	PrintSeparator(depth string, out StringFunction)
	PrintEmptyTuple(depth string, out StringFunction)
	PrintNullaryOperator(depth string, tag Tag, out StringFunction)
	PrintUnaryOperator(depth string, tag Tag, value Value, out StringFunction)
	PrintBinaryOperator(depth string, tag Tag, value1 Value, value2 Value, out StringFunction)
	PrintOpenTuple(depth string, tuple Value, out StringFunction) string
	PrintCloseTuple(depth string, tuple Value, out StringFunction)
	PrintHeadTag(value Tag, out StringFunction)
	PrintTag(depth string, value Tag, out StringFunction)
	PrintInt64(depth string, value int64, out StringFunction)
	PrintFloat64(depth string, value float64, out StringFunction)
	PrintString(depth string, value string, out StringFunction)
	PrintBool(depth string, value bool, out StringFunction)
	PrintComment(depth string, value Comment, out StringFunction)
}

func PrintScalar(printer Printer, depth string, token Value, out StringFunction) {
	printer.PrintScalarPrefix(depth, out)
	switch token.(type) {
	case Tag:
		printer.PrintTag(depth, token.(Tag), out)
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
	case Tuple:
		if token.Arity() == 0 {
			printer.PrintEmptyTuple(depth, out)
		}
		log.Printf("ERROR unexpected tuple '%s", token);  // TODO return error or prevent from ever happening
	default:
		log.Printf("ERROR type '%s' not recognised: %s", reflect.TypeOf(token), token);  // TODO return error or prevent from ever happening
	}
}

func PrintTuple(printer Printer, depth string, tuple Value, out StringFunction) {
	newDepth := printer.PrintOpenTuple(depth, tuple, out)
	printer.PrintSuffix(depth, out)
	ll := tuple.Arity()
	first := false
	if ll > 0 {
		_, first = tuple.Get(0).(Tag)
	}
	ForallN(tuple, func (k int, value Value) error {
		printer.PrintIndent(newDepth, out)
		if first && k == 0 {
			printer.PrintHeadTag(value.(Tag), out)
		} else {
			PrintExpression1(printer, newDepth, value, out)
		}
		if k < ll-1 {
			printer.PrintSeparator(newDepth, out)
		}
		printer.PrintSuffix(depth, out)
		return nil
	})
	printer.PrintCloseTuple(depth, tuple, out)
}

func PrintExpression(printer Printer, depth string, token Value, out StringFunction) {
	printer.PrintIndent(depth, out)
	PrintExpression1(printer, depth, token, out)
	printer.PrintSuffix(depth, out)
}

func PrintExpression1(printer Printer, depth string, token Value, out StringFunction) {

	if IsAtom(token) {
		PrintScalar(printer, depth, token, out)
	} else {
		tuple := token
		len := tuple.Arity()
		if len == 0 {
			printer.PrintScalarPrefix(depth, out)
			printer.PrintEmptyTuple(depth, out)
		} else {
			head := tuple.Get(0)
			tag, ok := head.(Tag)
			//log.Printf("Tuple [%s] %d\n", tag, len)
			if ok {  // TODO and head in a (binary) operator
				switch len {
				case 1:
					printer.PrintNullaryOperator(depth, tag, out)
				case 2:
					printer.PrintUnaryOperator(depth, tag, tuple.Get(1), out)
				case 3:
					printer.PrintBinaryOperator(depth, tag, tuple.Get(1), tuple.Get(2), out)
				default:
					PrintTuple(printer, depth, tuple, out)
				}
			} else {
				PrintTuple(printer, depth, tuple, out)
			}
		}
	}
}
