package tuple_test

import (
	"testing"
	"tuple"
)

func TestEvalJson(t *testing.T) {
	var grammar = tuple.NewJSONGrammar()
	testFloatExpression(t, grammar, "123.", 123)
	...
}

