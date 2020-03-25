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

func (query Query) match(depth int, token Value, next Next) {
	if tuple, ok := token.(Tuple); ok {

		if len(tuple.List) == 0 {
			return
		}
		head := tuple.List[0]
		atom, ok := head.(Atom)
		var name string
		if ok {
			name = atom.Name
		} else {
			return // TODO
		}
		//TODO Handles cons cells
		ll := len(query.components)
		if query.depth < ll {
			component := query.components[depth]
			if name == component || component == "*" {
				if depth == ll-1 {
					next(token)
					return
				}
				for _, token := range tuple.List {
					query.match(depth+1, token, next)
				}
				
			}
		}
	}
}


func (query Query) Match(expression Value, next Next) {

	query.match(0, expression, next)
}
