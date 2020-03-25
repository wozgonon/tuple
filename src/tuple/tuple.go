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

/////////////////////////////////////////////////////////////////////////////

// https://en.wikipedia.org/wiki/Arity
type Value interface {
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

func NewTuple(values...Value) Tuple {
	return Tuple{values}
}

