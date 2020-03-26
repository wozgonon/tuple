package tuple_test

import (
	"testing"
	"tuple"
)


func TestEvalJson(t *testing.T) {
	var grammar = tuple.NewJSONGrammar()

	test := func(formula string, expected tuple.Value) {
		val := tuple.Eval(grammar, formula)
		if val != expected {
			t.Errorf("%s=%f  expected=%s", formula, val, expected)
		}
	}

	test("true", tuple.Bool(true))
	test("false", tuple.Bool(false))
	test("123", tuple.Int64(123))
	test("0", tuple.Int64(0))
	test("-0", tuple.Int64(0))
	test("-1", tuple.Int64(-1))
	test("123.", tuple.Float64(123))
	test("-123.", tuple.Float64(-123))
	test(".123", tuple.Float64(.123))
	test("1234.11", tuple.Float64(1234.11))
	test("\"abc\"", tuple.String("abc"))
	test("\"a\\nb\\tc\"", tuple.String("a\nb\tc"))
	//test("[]", tuple.NewTuple())
	//testFloatExpression(t, grammar, "", )

	// TODO...
}

