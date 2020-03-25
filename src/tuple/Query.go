package tuple

import "strings"

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
