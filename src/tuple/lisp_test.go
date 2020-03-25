package tuple_test

import (
	"testing"
	"tuple"
	"math"
)


func TestEvalLisp1(t *testing.T) {
	var grammar = tuple.NewLispGrammar()
	k := "(+ 1 1)"
	val := tuple.Eval(grammar, k)
	if val != tuple.Float64(2) {
		t.Errorf("%s=%s", k, val)
	}
}

func TestEvalLispToInt64(t *testing.T) {
	var grammar = tuple.NewLispGrammar()
	tests := map[string]float64{
		"(+ 1 1)" : 2,
		"(1.)" : 1,
		"(-1)" : -1,
		"(((-(1)))" : -1,
		"(+ 1 (* 2 3))" : 7,
		"(* (+ 1 2) 3)" : 9,
		"(/ (* (* (+ 1 2) 3) 3) 9)" : 3,
		"(- (* (^ 1 7) 2))" : -2,
		"(+ (* (^ (-1) 7) 2) 3)" : 1,
		"(* 0 (7))" : 0,
		"(/ 8 4)" : 2,
		"(cos PI)" : -1,
		"(acos (cos PI)))" : math.Pi,
	}
	for k, v := range tests {
		val := tuple.Eval(grammar, k)
		if val != tuple.Float64(v) {
			t.Errorf("%s=%d   %s", k, int64(v), val)
		}
	}
}
