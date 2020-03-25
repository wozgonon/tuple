package tuple_test

import (
	"testing"
	"tuple"
	"math"
)

var grammar = tuple.NewInfixExpressionGrammar()

func TestEval1(t *testing.T) {
	val := tuple.Eval(grammar, "1+1")
	if val != tuple.Float64(2) {
		t.Errorf("1+1=%s", val)
	}
}

func TestEvalToInt64(t *testing.T) {
	tests := map[string]float64{
		"1+1" : 2,
		"1." : 1,
		"-1" : -1,
		"(((-1)))" : -1,
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
		val := tuple.Eval(grammar, k)
		if val != tuple.Float64(v) {
			t.Errorf("%s=%d   %s", k, int64(v), val)
		}
	}
}
