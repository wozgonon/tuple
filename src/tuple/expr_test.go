package tuple_test

import (
	"testing"
	"tuple"
	"math"
)

var grammar = tuple.NewInfixExpressionGrammar()

func TestExpr1(t *testing.T) {
	testFloatExpression(t, grammar, "1+1", 2)
}

func TestExprToInt64(t *testing.T) {
	tests := map[string]float64{
		"1+1" : 2,
		"1." : 1,
		"-1." : -1,
		"(((-1.)))" : -1,
		"1+2*3" : 7,
		"(1+2)*3" : 9,
		"((1+2)*3)" : 9,
		"((1+2)*3*3/9)" : 3,
		"-1^7*2" : -2,
		"-1^7*2+3" : 1,
		"0*(7)" : 0,
		"8/4" : 2,
		"cos(PI)" : -1,
		"acos(cos(PI))" : math.Pi,
	}
	for k, v := range tests {
		testFloatExpression(t, grammar, k, v)
	}
}

func testFloatExpressionFailParse(t *testing.T, grammar tuple.Grammar, formula string) {
	val := tuple.Eval(grammar, formula)
	f,ok := val.(tuple.Float64);
	if !ok || (ok && ! math.IsNaN(float64(f))) {
		t.Errorf("%s != %s", grammar, val)
	}
}


func TestExpr(t *testing.T) {
	testFloatExpressionFailParse(t, grammar, "-")
}
