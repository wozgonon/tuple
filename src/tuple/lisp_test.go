package tuple_test

import (
	"testing"
	"tuple"
	"math"
)

func testFloatExpression(t *testing.T, grammar tuple.Grammar, formula string, expected float64) {
	val := tuple.Eval(grammar, formula)
	floatExpected := tuple.Float64(expected)
	floatVal, ok := val.(tuple.Float64)
	if ok {
		if floatVal != floatExpected {
			t.Errorf("%s=%f  val=%f %f", formula, expected, floatVal, floatExpected)
		}
	} else {
		t.Errorf("Expected float %s=%f  expected=%f got %s", formula, expected, floatExpected, val)
	}
}


func TestEvalLisp1(t *testing.T) {
	var grammar = tuple.NewLispGrammar()
	testFloatExpression(t, grammar, "(+ 1 1)", 2)
}

func TestEvalLispToInt64(t *testing.T) {
	var grammar = tuple.NewLispGrammar()
	tests := map[string]float64{
		//"(-1.)" : -1.,
		"(+ 1 1)" : 2,
		"(1.)" : 1,
	//	"(-1.)" : -1,
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
		testFloatExpression(t, grammar, k, v)
	}
}

func TestEvalLispWithInfixGrammar1(t *testing.T) {
	var grammar = tuple.NewLispWithInfixGrammar()
	testFloatExpression(t, grammar, "(1+1)", 2)
}

func TestEvalLispWithInfixGrammarToInt64(t *testing.T) {
	var grammar = tuple.NewLispWithInfixGrammar()
	tests := map[string]float64{
		"(1+1)" : 2,
		"(1.)" : 1,
//		"(-1)" : -1,
		"((((-1.))))" : -1,
		"(1+2*3)" : 7,
		"((1+2)*3)" : 9,
		"(((1+2)*3))" : 9,
		"(((1+2)*3*3/9))" : 3,
		"(-1^7*2)" : -2,
		"(-1^7*2+3)" : 1,
		"(0*(7))" : 0,
		"(8/4)" : 2,
		"(cos PI)" : -1,
		"(acos (cos PI))" : math.Pi,
	}
	for k, v := range tests {
		testFloatExpression(t, grammar, k, v)
	}
}

