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
		"--1." : 1,
		"-- 1." : 1,
		"--- 1." : -1,
		"(((-1.)))" : -1,
		"1+2*3" : 7,
		"(1+2)*3" : 9,
		"((1+2)*3)" : 9,
		"((1+2)*3*3/9)" : 3,
		"-1^7*2" : -2,
		"-1^7*2+3" : 1,
		"0*(7)" : 0,
		"8/4" : 2,
		"-(-(-1)+2)" : 1,
		"cos(PI)" : -1,
		"acos(cos(PI))" : math.Pi,
		"asin(sin(0))" : 0,
		"atan(tan(0))" : 0,
		"exp(log(1))" : 1,
	}
	for k, v := range tests {
		//t.Logf("*** %s %f", k, v)
		testFloatExpression(t, grammar, k, v)
	}

}

func TestExpr(t *testing.T) {
	testFloatExpressionFailParse(t, grammar, "-")

	test := func(formula string, expected tuple.Value) {
		val := tuple.Eval(grammar, formula)
		if val != expected {
			t.Errorf("%s=%f  expected=%s", formula, val, expected)
		}
	}

	test("true", tuple.Bool(true))
	test("false", tuple.Bool(false))

	test("0", tuple.Int64(0))
	test("1", tuple.Int64(1))
	test("-1", tuple.Int64(-1))

	test("-1.", tuple.Float64(-1.))
	test(".0", tuple.Float64(.0))
	test(".1", tuple.Float64(.1))

	test("round(1234.1234)", tuple.Float64(1234))
	test("atan2(0,0)", tuple.Float64(0))

	test("\"abc\"", tuple.String("abc"))
}



func TestExprEquals(t *testing.T) {

	test := func (formula string) {
		val := tuple.Eval(grammar, formula)
		if val != tuple.Bool(true) {
			t.Errorf("Expected '%s' to be TRUE", formula)
		}
	}

	test("11 == 11")
	test("7 == 1+2*3")
	test("5 == 1*2+3")
	test("120 == 1*2*3*4*5")
	test("6 == (1)+((2))+(((3)))")
	test("22 == ((22))")
	test("22 == ((((22))))")
	test("3 == (1+2)")
	test("9 == (1+2)*3")
	test("10 == 1+2+3+4")
	test("10 == 1+(2+3)+4")
	test("10 == (1+2+3+4)")
	test("10 == ((1+((2)+3))+(4))")
	test("-123 == -123")
	test("-123 == -(123)")
	test("3 == (0- - 3)")
	test("-3 == -(0- - 3)")
	test("-2 == -(1--1)")
	test("-3 == -(0- - - - 3)")
	test("-3 == -(0--3)")
	test("1 == (cos 0)")
	test("-1 == (cos PI)")
	test("3.141592653589793 == (acos (cos (PI)))")
	test("(acos (cos(PI)))==PI")
	test("true == ((acos (cos(PI)))==PI)")

	// TODO test("+1 == 1")
	// TODO test("++1 == 1")
	// TODO test("+-+1 == -1")


	test("3 != -(1+2)")
	test("1 == -(-(-1)+2)")

	test("1 >= 1")
	test("1 >= 0")
	test("1 >  0")
	test("1 <  2")
	test("1 <= 1")
	test("0 <= 1")
	test("1 != -1")
	test("1 == 1")

	test("true")
	test("! false")

	test("(!1) == 0")
	test("(!0) == 1")
	// TODO test("!1 == 0")  - TODO priority of unary operators
	// TODO test("!0 == 1")

	test("len(\"abcde\")==5")
	test("1+len(\"abcde\")==6")
}

