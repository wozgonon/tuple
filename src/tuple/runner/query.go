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
import "tuple/eval"
import "reflect"

type Tuple = tuple.Tuple
type Map = tuple.Map
type Tag = tuple.Tag
type Logger = tuple.Logger
var Verbose = tuple.Verbose
var IntToString = tuple.IntToString

func AddSafeQueryFunctions(table eval.LocalScope) {

	table.Add("query", func (context eval.EvalContext, path string, what Value) (Value, error) {
		evaluated, err := eval.Eval(context, what)
		if err != nil {
			return nil, err
		}
		result := tuple.NewTuple()
		query := NewQuery(path)
		query.Match(context, evaluated, func (value Value) error {
			result.Append(value)
			return nil
		})
		return result, nil
	})
}

/////////////////////////////////////////////////////////////////////////////

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

func (query Query) Match(logger Logger, expression Value, next Next) {

	Verbose(logger, "QUERY typeofexpression=%s components=%", expression, query.components)
	query.filter(logger, 0, expression, next)
}

func (query Query) filter(logger Logger, depth int, value Value, next Next) {

	ll := len(query.components)
	matchedAll := depth >= ll
	if matchedAll {
		next(value)
		return
	}

	if tuple, ok := value.(Tuple); ok {
		Verbose(logger, "QUERY depth=%d, tuple arity=%d", depth, value.Arity())
		if tuple.Arity() == 0 {
			if query.matchLeaf(logger, depth, "") {
				next(value)
			}
			return
		}
		head := tuple.Get(0)
		Verbose(logger, "  head=%s", head)
		if tag, ok := head.(Tag); ok {
			name := tag.Name
			if query.matchComponent(logger, depth, name) {
				if depth+1 == ll {
					next(value)
					return
				}
				for _, value := range tuple.List {
					query.filter(logger, depth+1, value, next)
				}
			}
			return
		} else if str, ok := head.(String); ok{
			if query.matchComponent(logger, depth, string(str)) {
				if depth+1 == ll {
					next(value)
					return
				}
				for _, value := range tuple.List {
					query.filter(logger, depth+1, value, next)
				}
			}
			return
		}
		for k, value := range tuple.List {
			if query.matchComponent(logger, depth, IntToString(int64(k))) {
				if depth+1 == ll {
					query.filter(logger, depth+1, value, next)
				}
			}
		}
		
	} else if mapp, ok := value.(Map); ok {
		Verbose(logger, "QUERY depth=%d, map arity=%d", depth, mapp.Arity())
		mapp.ForallKeyValue(func (key Tag, value Value) {
			Verbose(logger, "   QUERY depth=%d, match key=%s", depth, key.Name)
			if query.matchComponent(logger, depth, key.Name) {
				query.filter(logger, depth+1, value, next)
			}
			return
		})
	} else if tag, ok := value.(Tag); ok {
		Verbose(logger, "QUERY depth=%d, tag=%d", depth, tag)
		if query.matchLeaf(logger, depth, tag.Name) {
			next(value)
		}
	} else {
		Verbose(logger, "QUERY depth=%d, Ignore value type='%s'", depth, reflect.TypeOf(value))
	}
}

func (query Query) matchLeaf(logger Logger, depth int, name string) bool {
	ll := len(query.components)
	return depth == ll-1 && query.matchComponent(logger, depth, name)
}

func (query Query) matchComponent(logger Logger, depth int, name string) bool {
	ll := len(query.components)
	if depth < ll {
		component := query.components[depth]
		Verbose(logger, "   QUERY TRY MATCH depth=%d, ll=%d name=%s component=%s", depth, ll, name, component)
		if name == component || component == "*" {
			Verbose(logger, "   QUERY MATCH depth=%d, ll=%d name=%s component=%s", depth, ll, name, component)
			return true
			
		}
	}
	Verbose(logger, "   QUERY NO MATCH depth=%d, ll=%d name=%s", depth, ll, name)
	return false
}

