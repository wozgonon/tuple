package parsers_test

import (
	"testing"
	"tuple"
	"math"
)

func testIntExpression(t *testing.T, grammar tuple.Grammar, formula string, expected int64) {
	//t.Logf("TRY: %s==%f", formula, expected)
	val,_ := ParseAndEval(safeEvalContext, grammar, formula)
	intExpected := tuple.Int64(expected)
	intVal, ok := val.(tuple.Int64)
	if ok {
		if intVal != intExpected {
			t.Errorf("ERROR: %s=%d  val=%d %d", formula, expected, intVal, intExpected)
		}
	} else {
		t.Errorf("ERROR: Expected int %s=%d  expected=%d got %s", formula, expected, intExpected, val)
	}
}

func testFloatExpression(t *testing.T, grammar tuple.Grammar, formula string, expected float64) {
	//t.Logf("TRY: %s==%f", formula, expected)
	val,err := ParseAndEval(safeEvalContext, grammar, formula)
	if err != nil {
		t.Errorf("Given '%s' expected got error %s", formula, err)
	}
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

func TestLispCons(t *testing.T) {

	test := func(formula string) {
		var grammar = NewLispGrammar()
		c,err := ParseAndEval(safeEvalContext, grammar, formula)
		if err != nil {
			t.Errorf("Given '%s' expected got error %s", formula, err)
		}
		if b, ok := c.(tuple.Bool); !ok || ! bool(b) {
			t.Errorf("Given '%s' expected true got %s", formula, c)
		}
	}
	//test("a.b")
	// test("((cons a b) ())")  // TODO investigate
	
	test("(ismap (a.b))")
	test("(ismap (a.b c.d))")
	test("(ismap (a.b c.d e.f))")
}

func TestEvalLisp1(t *testing.T) {
	var grammar = NewLispGrammar()
	testFloatExpression(t, grammar, "(+ 1 1)", 2)
}

func TestEvalLispToInt64(t *testing.T) {
	var grammar = NewLispGrammar()
	tests := map[string]float64{
		//"(-1.)" : -1.,
		"(+ 1 1)" : 2,
		"(1.)" : 1,
	//	"(-1.)" : -1,
	//	"(((-(1)))" : -1,
		"(+ 1 (* 2 3))" : 7,
		"(* (+ 1 2) 3)" : 9,
		"(/ (* (* (+ 1 2) 3) 3) 9)" : 3,
		"(- (* (** 1 7) 2))" : -2,
		"(+ (* (** (-1) 7) 2) 3)" : 1,
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
	var grammar = NewLispWithInfixGrammar()
	testFloatExpression(t, grammar, "(1+1)", 2)
}

func TestEvalLispWithInfixGrammarToInt64(t *testing.T) {
	var grammar = NewLispWithInfixGrammar()
	tests := map[string]float64{
		"(1+1)" : 2,
		"(1.)" : 1,
		"(-1.)" : -1,
		"((((-1.))))" : -1,
		"(1+2*3)" : 7,
		"((1+2)*3)" : 9,
		"(((1+2)*3))" : 9,
		"(((1+2)*3*3/9))" : 3,
		"(-1**7*2)" : -2,
		"(-1**7*2+3)" : 1,
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
		val,_ := ParseAndEval(safeEvalContext, grammar, formula)
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
	test("(-3 == -(-(-1)+2))")
}

