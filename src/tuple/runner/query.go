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
package runner

import "strings"
import "tuple"

type Tuple = tuple.Tuple
type Tag = tuple.Tag

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
	if depth < ll {
		component := query.components[depth]
		if name == component || component == "*" {
			if depth == ll-1 {
				//fmt.Printf("Match component=%s name=%s depth=%d\n", component, name, depth)
				return true
			}
			
		}
	}
	return false
}

func (query Query) filter(depth int, token Value, next Next) {
	if tuple, ok := token.(Tuple); ok {

		if tuple.Arity() == 0 {
			if query.match(depth, "") {
				next(token)
			}
			return
		}
		head := tuple.Get(0)
		if tag, ok := head.(Tag); ok {
			name := tag.Name
			if query.match(depth, name) {
				next(token)
			}
		} else if str, ok := head.(String); ok{
			if query.match(depth, string(str)) {
				next(token)
			}
		}
		for _, token := range tuple.List {
			query.filter(depth+1, token, next)
		}

	} else if tag, ok := token.(Tag); ok {
		if query.match(depth, tag.Name) {
			next(token)
		}
	}
}


func (query Query) Match(expression Value, next Next) {

	query.filter(0, expression, next)
}
