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

type String struct {
	value string
}

type Number struct {
	float bool
	value string
}

// An Atom - a name for something, an identifier or operator
type Atom struct {
	// TODO include location and source for editors
	Name string
}

func (atom Atom) Arity() int { return 0 }
func (number Number) Arity() int { return 0 }
func (number String) Arity() int { return 0 }
func (comment Comment) Arity() int { return 0 }

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
	List []interface{}
}

func (tuple *Tuple) Append(token interface{}) {
	tuple.List = append(tuple.List, token)
}

func (tuple *Tuple) Length() int {
	return len(tuple.List)
}

//func NewTuple() Tuple {
//	return Tuple{make([]interface{}, 0)}
//}

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

func NewTuple(values...interface{}) Tuple {
	return Tuple{values}
}

