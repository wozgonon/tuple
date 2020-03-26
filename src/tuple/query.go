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
//import "fmt"

//  The Query type is used for filtering the AST produced by the parser.
type Query struct {
	query string
	components []string
	depth int
}

func NewQuery(query string) Query {
	components := strings.Split(query, ".")
	return Query{query, components, 0}
}

func (query Query) match(depth int, name string) bool {
	//TODO Handles cons cells
	ll := len(query.components)
	//fmt.Printf("name=%s ll=%d depth=%d\n", name, ll, depth)
	if depth <= ll {
		component := query.components[depth-1]
		//fmt.Printf("component=%s", component)
		if name == component || component == "*" {
			if depth == ll-1 {
				return true
			}
			
		}
		return false
	}
	return true
}

func (query Query) filter(depth int, token Value, next Next) {
	///fmt.Printf("match=%s\n", token)
	if tuple, ok := token.(Tuple); ok {

		/*if len(tuple.List) == 0 {
			return
		}
		head := tuple.List[0]
		//fmt.Printf("head=%s\n", head)
		atom, ok := head.(Atom)
		var name string
		if ok {
			name = atom.Name
			if query.match(depth, name) {
				next(token)
			}
		} else {
			// TODO test string
			return // TODO
		}
               */
		for _, token := range tuple.List {
			query.filter(depth+1, token, next)
		}

	} else if atom, ok := token.(Atom); ok {
		if query.match(depth+1, atom.Name) {
			next(token)
		}
	}
}


func (query Query) Match(expression Value, next Next) {

	query.filter(-1, expression, next)
}
