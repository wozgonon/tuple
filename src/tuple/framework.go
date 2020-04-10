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

import "path"
import "math"
import "strconv"
import "fmt"

/////////////////////////////////////////////////////////////////////////////
//  Values, maps and callbacks
/////////////////////////////////////////////////////////////////////////////

// The Value interface:
// * must be implemented by any of the small number of types that can be produced by the lexer.
// * is used in the parsers to represent the Parse Tree or AST
// * is used in the backend Eval function at runtime
// TODO include location and source for editors
type Value interface {
	// See: https://en.wikipedia.org/wiki/Arity
	Arity() int
	ForallValues(next func(value Value) error) error
}

type Array interface {
	Value
	Get(index int) Value
}

type Map interface {
	Value
	ForallKeyValue(next KeyValueFunction)
}

type KeyValueFunction = func(key Tag, value Value)
type StringFunction = func(value string)
type Next = func(value Value) error

/////////////////////////////////////////////////////////////////////////////
//  Lexer
//
//  The Value interface must be implemented by any of the small number of types
//  that can be produced by the lexer.
/////////////////////////////////////////////////////////////////////////////

type Lexer interface {
	Printer
	GetNext(context Context, eol func(), open func(open string), close func(close string), nextTag func(tag Tag), nextLiteral func (literal Value)) error
}

type String string
type Float64 float64
type Int64 int64
type Bool bool

// A Tag - a name for something, an identifier or operator
type Tag struct {
	Name string
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
	Parse(context Context, next Next) error // , next func(tuple Tuple)) (Value, error)
	
	// Pretty prints the objects in the given syntax.
	// The output ought to be parsable by the 'parse' method.
	Print(token Value, next func(value string))
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
//  
/////////////////////////////////////////////////////////////////////////////

// LISP cons operator (https://en.wikipedia.org/wiki/Cons)
var CONS_ATOM = Tag{"_cons"}
var NAN Float64 = Float64(math.NaN())
var EMPTY Tuple = NewTuple()
const DOUBLE_QUOTE = "\""

// TODO an atom is pretty subjective, should be grammar specific
func IsAtom(value Value) bool {
	return value.Arity() == 0
}

func ForallInArray(value Array, next func(value Value) error) error {  // Make this a method on the interface
	ll := value.Arity()
	for k := 0; k < ll; k+=1 {
		err := next(value.Get(k))
			if err != nil {
				return err
			}
	}
	return nil
}

func ForallKeyValuesInArray(value Array, next func(index int, value Value) error) error {
	ll := value.Arity()
	for k := 0; k < ll; k+=1 {
		err := next(k, value.Get(k))
		if err != nil {
			return err
		}
	}
	return nil
}

func  Head(value Value) (Tag, bool) {
	if value.Arity() > 0 {
		if array, ok := value.(Array); ok {
			first := array.Get(0)
			head, ok := first.(Tag)
			return head, ok
		}
	}
	return Tag{""}, false
}

/////////////////////////////////////////////////////////////////////////////
// 
/////////////////////////////////////////////////////////////////////////////

func (value Tag) Arity() int { return 0 }
func (value String) Arity() int { return 0 }
func (value Float64) Arity() int { return 0 }
func (value Int64) Arity() int { return 0 }
func (value Bool) Arity() int { return 0 }

func (value Tag) ForallValues(next func(value Value) error) error { return ForallInArray(value, next) }
func (value String) ForallValues(next func(value Value) error) error { return ForallInArray(value, next) }
func (value Float64) ForallValues(next func(value Value) error) error { return ForallInArray(value, next) }
func (value Int64) ForallValues(next func(value Value) error) error { return ForallInArray(value, next) }
func (value Bool) ForallValues(next func(value Value) error) error { return ForallInArray(value, next) }
func (value Tuple) ForallValues(next func(value Value) error) error { return ForallInArray(value, next) }

func (value Tag) Get(index int) Value {
	if index == 0 {
		return String(value.Name)
	}
	return EMPTY
}
func (value String) Get(index int) Value {
	if index >=0 && index < len(string(value)) {
		return String(value[index])
	}
	return EMPTY
}
func (value Float64) Get(index int) Value { return Int64(int64(value)) }
func (value Int64) Get(index int) Value { return Bool(NthBitOfInt(int64(value), index)) }
func (value Bool) Get(_ int) Value { return value }  // TODO should this return EMPTY or just itself??

func (value Tag) GetKeyValue(index int) (Tag, Value) { return IntToTag(index), value.Get(index) }
func (value String) GetKeyValue(index int) (Tag, Value) { return IntToTag(index), value.Get(index) }
func (value Float64) GetKeyValue(index int) (Tag, Value) { return IntToTag(index), value.Get(index) }
func (value Int64) GetKeyValue(index int) (Tag, Value) { return IntToTag(index), value.Get(index) }
func (value Bool) GetKeyValue(index int) (Tag, Value) { return IntToTag(index), value.Get(index) }

////////////////////////////////////////////////////////////////////////////
// Tuple
/////////////////////////////////////////////////////////////////////////////

type Tuple struct {
	List []Value  // TODO change to elements
}

func NewTuple(values... Value) Tuple {
	return Tuple{values}
}

func (tuple Tuple) Arity() int { return len(tuple.List) }
func (tuple Tuple) Get(index int) Value {
	if index >= 0 && index < len(tuple.List) {
		return tuple.List[index]
	}
	return NAN  // TODO perhaps this should be an error or EMPTY
}
func (tuple Tuple) GetKeyValue(index int) (Tag,Value) {
	return IntToTag(index), tuple.Get(index)
}

func (tuple *Tuple) Append(token Value) {
	if token == nil {
		panic("Unexpected nil value")
	}
	tuple.List = append(tuple.List, token)
}

func (tuple *Tuple) Set(index int, token Value) {
	if token == nil {
		panic("Unexpected nil value")
	}
	tuple.List[index] = token
}

/////////////////////////////////////////////////////////////////////////////

type TagValueMap struct {
	elements map[Tag]Value
}

func NewTagValueMap() TagValueMap {
	return TagValueMap{make(map[Tag]Value)}
}

func (mapp * TagValueMap) Add(key Tag, value Value) { mapp.elements[key] = value }

func (mapp TagValueMap) Arity() int { return len(mapp.elements) }

func (mapp TagValueMap) ForallKeyValue(next KeyValueFunction) {
	for k, v := range mapp.elements {
		next(k, v)
	}
}

func (mapp TagValueMap) ForallValues(next func(value Value) error) error {
	for _, v := range mapp.elements {
		err := next(v)
		if err != nil {
			return nil
		}
	}
	return nil
}

/////////////////////////////////////////////////////////////////////////////

// A Finite Stream it implements the Value interface and so can be used in place of a map or an array collection.
// It is an important abstraction for efficiency, as it avoids having to copy data into arrays or maps.
// It involves no overhead of allocating memory and passes on the first element to next immediately so there is
// no delay allocating and copying the whole collection.
// A finite stream is 'safe' because it allocates no memory and runs in finite time.
// It is not necessarilly 'harmless' in in thate finite time is variable and might be large.
type FiniteStream struct {
	value Value
	next func(v Value, next func (v Value) error) error
}

func NewFiniteStream(value Value, next func(v Value, next func(v Value) error) error) FiniteStream {
	return FiniteStream{value, next}
}
func (stream FiniteStream) Arity() int { return stream.value.Arity() }
func (stream FiniteStream) ForallValues(next func(value Value) error) error {
	return stream.value.ForallValues(func (v Value) error {
		return stream.next(v, next)
	})}

/////////////////////////////////////////////////////////////////////////////
// Basic conversion functions between scalar types
/////////////////////////////////////////////////////////////////////////////

func IntToString(value int64) string {
	return strconv.FormatInt(value, 10)
}

func FloatToString(value float64) string {
	if math.IsInf(value, 64) {
		return "Inf" // style.INF)  // Do not print +Inf
	} else {
		return fmt.Sprint(value)
	}
}

func Float64ToString(value Float64) string {
	return FloatToString(float64(value))
}

func Int64ToString(value Int64) string {
	return IntToString(int64(value))
}

func IntToTag(value int) Tag {
	return Tag{IntToString(int64(value))}
}

func BoolToFloat(value Bool) float64 {
	if bool(value) {
		return 1.
	}
	return 0.0
}

func BoolToInt(value Bool) int64 {
	if bool(value) {
		return 1
	}
	return 0
}

func Quote(value string, out func(value string)) {
	out(DOUBLE_QUOTE)
	out(value)   // TODO Escape
	out(DOUBLE_QUOTE)
}

func BoolToString(value bool) string {
	return fmt.Sprintf("%t", value)
}

func NthBitOfInt(value int64, index int) bool {
	bit := uint64(value) & (1<<uint(index))
	return bit != 0
}
