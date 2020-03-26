package tuple_test

import (
	"testing"
	"tuple"
)


func TestQuery(t *testing.T) {

	test := func(input tuple.Value, expectedCount int, queryString string) {
		count := 0
		query := tuple.NewQuery(queryString)
		query.Match(input, func (value tuple.Value) {

			if atom, ok := value.(tuple.Atom); ok {
				if atom.Name == queryString {
					count += 1
				} else {
					t.Errorf("Expected (%s) got: %s input=%s ", queryString, value, input)
				}
			} else {
				tupleValue := value.(tuple.Tuple)
				if tupleValue.Arity() != 1 {
					head := tupleValue.List[0]
					if headAtom, ok := head.(tuple.Atom); ok && headAtom.Name == queryString {
						count += 1
						//t.Logf("Match")
					} else {
						t.Errorf("Expected (%s) got: %s input=%s head=%s ", queryString, value, input, head)
					}
				}
			}
		})

		if count != expectedCount {
			t.Errorf("Expected %d match(es) got %d input=%s query=%s", expectedCount, count, input, queryString)
		}
	}
	//a := tuple.NewTuple(tuple.Atom{"a"})
	//b := tuple.NewTuple(tuple.Atom{"b"})
	//c := tuple.NewTuple(tuple.Atom{"c"})

	a := (tuple.Atom{"a"})
	b := (tuple.Atom{"b"})
	c := (tuple.Atom{"c"})

	abc := tuple.NewTuple(a, b, c)
	//aabbcca:= tuple.NewTuple(a, a, b, b, c, c, a)

	//test(abc, 3, "*")
	test(abc, 0, "q")
	//test(abc, 1, "a")
	//test(abc, 1, "b")
	//test(aabbcca, 3, "a")
}

