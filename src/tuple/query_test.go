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
			count += 1
		})

		if count != expectedCount {
			t.Errorf("Expected %d match(es) got %d input=%s query=%s", expectedCount, count, input, queryString)
		}
	}

	a := (tuple.Atom{"a"})
	b := (tuple.Atom{"b"})
	c := (tuple.Atom{"c"})

	empty := tuple.NewTuple()
	abc := tuple.NewTuple(a, b, c)
	bc := tuple.NewTuple(b, c)
	aaa := tuple.NewTuple(a, bc, abc, bc, abc, abc)

	test(abc, 1, "*")
	test(abc, 0, "q")
	test(abc, 1, "a")
	test(aaa, 2, "*.b")
	test(aaa, 2, "a.b")
	test(aaa, 4, "a.a")
	test(aaa, 6, "a.*")
	test(aaa, 6, "*.*")
	test(aaa, 0, "*.z")
	test(aaa, 0, "a.z")
	test(a, 1, "a")
	test(a, 1, "a")
	test(b, 0, "a")
	test(empty, 0, "a")
	test(empty, 1, "*")
	test(empty, 1, "")

	test(tuple.NewTuple(tuple.String("a")), 1, "a")
	test(tuple.NewTuple(tuple.String("a")), 1, "*")
}

