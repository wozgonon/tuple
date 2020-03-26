package tuple_test

import (
	"testing"
	"tuple"
	"math"
)

func testFloatExpression(t *testing.T, grammar tuple.Grammar, formula string, expected float64) {
	//t.Logf("TRY: %s==%f", formula, expected)
	val := tuple.Eval(grammar, formula)
	floatExpected := tuple.Float64(expected)
	floatVal, ok := val.(tuple.Float64)
	if ok {
		if floatVal != floatExpected {
			t.Errorf("ERROR: %s=%f  val=%f %f", formula, expected, floatVal, floatExpected)
		}
	} else {
		t.Errorf("ERROR: Expected float %s=%f  expected=%f got %s", formula, expected, floatExpected, val)
	}
}

func testFloatExpressionFailParse(t *testing.T, grammar tuple.Grammar, formula string) {
       val := tuple.Eval(grammar, formula)
       f,ok := val.(tuple.Float64);
       if !ok || (ok && ! math.IsNaN(float64(f))) {
               t.Errorf("%s != %s", grammar, val)
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
	//	"(((-(1)))" : -1,
		"(+ 1 (* 2 3))" : 7,
		"(* (+ 1 2) 3)" : 9,
		"(/ (* (* (+ 1 2) 3) 3) 9)" : 3,
		"(- (* (^ 1 7) 2))" : -2,
		"(+ (* (^ (-1) 7) 2) 3)" : 1,
		"(* 0 (7))" : 0,
		"(/ 8 4)" : 2,
		"(cos PI)" : -1,
		"(acos (cos PI))" : math.Pi,
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
		"(-1.)" : -1,
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
		"(0- - -3)" : -3,
	}
	for k, v := range tests {
		testFloatExpression(t, grammar, k, v)
	}
}

func TestLispInfixEquals(t *testing.T) {

	test := func (formula string) {
		val := tuple.Eval(grammar, formula)
		if val != tuple.Bool(true) {
			t.Errorf("Expected '%s' to be TRUE", formula)
		}
	}

	test("(11 == 11)")
	test("(7 == 1+2*3)")
	test("(5 == 1*2+3)")
	test("(120 == 1*2*3*4*5)")
	test("(6 == (1)+((2))+(((3))))")
	test("(22 == ((22)))")
	test("(22 == ((((22)))))")
	test("(3 == (1+2))")
	test("(9 == (1+2)*3)")
	test("(10 == 1+2+3+4)")
	test("(10 == 1+(2+3)+4)")
	test("(10 == (1+2+3+4))")
	test("(10 == ((1+((2)+3))+(4)))")
	test("(-123 == -123)")
	test("(-123 == -(123))")
	test("(3 == (0- - 3))")
	test("(-3 == -(0- - 3))")
	test("(-2 == -(1--1))")
	test("(-3 == -(0- - - - 3))")
	test("(-3 == -(0--3))")
	test("(1 == (cos 0))")
	test("(-1 == (cos PI))")
	test("(3.141592653589793 == (acos (cos PI)))")
	test("((acos (cos PI ))==PI)")
	test("(true == ((acos (cos PI))==PI))")


	test("(3 != -(1+2))")
	test("(1 == -(-(-1)+2))")
}

