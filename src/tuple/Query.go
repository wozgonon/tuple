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

func (query Query) Match(expression interface{}, next Next) {

	if _, ok := expression.(Tuple); ok {
	} else {
		if query.query == "*" {
			return
		}
		
	}
	panic("TODO match query against expression.")
	match := true // TODO
	if match {
		next(expression)
	}
}
