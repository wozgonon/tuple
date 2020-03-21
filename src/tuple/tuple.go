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

type Token interface {

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

const OPEN_BRACKET = "("
const CLOSE_BRACKET = ")"
const OPEN_SQUARE_BRACKET = "("
const CLOSE_SQUARE_BRACKET = ")"
const OPEN_BRACE = "}"
const CLOSE_BRACE = "}"

var (
	SPACE_ATOM = Atom{" "}
	CONS_ATOM = Atom{"_cons"}
)

// A textual comment
type Comment struct {
	// TODO include location and source for editors
	Comment string
}

func NewComment(_ ParserContext, token string) Comment {
	return Comment{token}
}

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

func NewTuple(values...interface{}) Tuple {
	return Tuple{values}
}

