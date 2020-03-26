package tuple_test

import (
	"testing"
	"tuple"
)

func testValueExpression(t *testing.T, grammar tuple.Grammar, formula string, expected tuple.Value) {
	val := tuple.Eval(grammar, formula)
	if val != expected {
		t.Errorf("%s=%f  expected=%s", formula, val, expected)
	}
}

func TestEvalJson(t *testing.T) {
	var grammar = tuple.NewJSONGrammar()
	testValueExpression(t, grammar, "true", tuple.Bool(true))
	testValueExpression(t, grammar, "123", tuple.Int64(123))
	testValueExpression(t, grammar, "1234.11", tuple.Float64(1234.11))
	//testValueExpression(t, grammar, "[]", tuple.NewTuple())
	//testFloatExpression(t, grammar, "", )

	// TODO...
}

