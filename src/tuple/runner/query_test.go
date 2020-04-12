package runner_test

import (
	"testing"
	"tuple"
	"tuple/runner"
	"tuple/parsers"
)

type SimpleLogger struct {}
func (logger * SimpleLogger) Log (level string, format string, args ...interface{}) {}


func TestQuery(t *testing.T) {

	test := func(input tuple.Value, expectedCount int, queryString string) {
		count := 0
		query := runner.NewQuery(queryString)
		logger := SimpleLogger{}
		query.Match(&logger, input, func (value tuple.Value) error {
			count += 1
			return nil
		})

		if count != expectedCount {
			t.Errorf("Expected %d match(es) got %d input=%s query=%s", expectedCount, count, input, queryString)
		}
	}

	a := (tuple.Tag{"a"})
	b := (tuple.Tag{"b"})
	c := (tuple.Tag{"c"})

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


func TestExprQuery(t *testing.T) {
	test := func (t *testing.T, formula string) {

		runner.AddSafeQueryFunctions(safeEvalContext)
		var grammar = parsers.NewInfixExpressionGrammar()
		val,err := runner.ParseAndEval(safeEvalContext, grammar, formula)
		if err != nil || ! bool(val.(tuple.Bool)) {
			t.Errorf("Expected success but got '%s' formula='%s' err=%s", val, formula, err)
		}
	}
	test(t, "eq (list {c:1 d:\"w\"}) (query \"os.b\" { os: { a:1 b: {c:1 d:\"w\"}}})")
	test(t, "eq (list \"w\") (query \"os.b.d\" { os: { a:1 b: {c:1 d:\"w\"}}})")
	test(t, "eq (list 1) (query \"os.*.c\" { os: { a:1 b: {c:1 d:\"w\"}}})")
	test(t, "eq (list 1) (query \"*.b.c\" { os: { a:1 b: {c:1 d:\"w\"}}})")

	// query "a.b.c" ("a" ("b" ("c" 1 2 3) ("c" 4 5 6)) ("d" (8 9 0)))  
}
