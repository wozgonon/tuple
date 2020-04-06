package parsers_test

import (
	"testing"
	"tuple"
	"tuple/eval"
	"math"
)


var grammar = NewInfixExpressionGrammar()
var logger = tuple.GetLogger(nil, false)
var symbols = NewSafeSymbolTable(eval.NewErrorIfFunctionNotFound(logger))  // TODO perhaps another default function would be better


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
		"+1." : 1,
		"(((-1.)))" : -1,
		"1+2*3" : 7,
		"(1+2)*3" : 9,
		"((1+2)*3)" : 9,
		"((1+2)*3*3/9)" : 3,
		"-1**7*2" : -2,
		"-1**7*2+3" : 1,
		"0*(7)" : 0,
		"8/4" : 2,
		"-(-(-1)+2)" : -3,
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
		val := ParseAndEval(&symbols, grammar, formula)
		if val != expected {
			t.Errorf("%s=%f  expected=%s", formula, val, expected)
		}
	}

	test("true", tuple.Bool(true))
	test("false", tuple.Bool(false))

	test("0", tuple.Int64(0))
	test("1", tuple.Int64(1))
	test("-1", tuple.Float64(-1))

	test("-1.", tuple.Float64(-1.))
	test(".0", tuple.Float64(.0))
	test(".1", tuple.Float64(.1))

	test("round(1234.1234)", tuple.Float64(1234))
	test("atan2(0,0)", tuple.Float64(0))

	test("\"abc\"", tuple.String("abc"))
}


func TestArithmeticAndLogic(t *testing.T) {

	testArithmeticAndLogic(t, NewInfixExpressionGrammar())
	testArithmeticAndLogic(t, NewShellGrammar())
}


func testArithmeticAndLogic(t *testing.T, grammar tuple.Grammar) {

	test := func (formula string) {
		val := ParseAndEval(&symbols, grammar, formula)
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

	test("log(E) == 1")
	test("PHI == PHI")
	test("PI != PHI")
	test("PI == PI")

	test("3 != -(1+2)")
	test("-3 == -(-(-1)+2)")

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

	test("true && ! false")
	test("true || false")

	test("1<3 && 3>2")

	test("(!1) == 0")
	test("(!0) == 1")
	// TODO test("!1 == 0")  - TODO priority of unary operators
	// TODO test("!0 == 1")

	test("len(\"abcde\")==5")
	test("1+len(\"abcde\")==6")
	test("len(upper(\"abc\"))==3")
	test("eq(upper(\"abc\"),\"ABC\")")
	test("eq(upper(lower(\"aBc\")),\"ABC\")")

	test("eq(\"ab\",\"ab\")")
	test("eq(eq(\"ad\",\"ab\"),false)")
	test("eq(concat(\"ab\",\"cde\"),\"abcde\")")
	test("eq(join(\"-\",(\"abc\",\"def\",\"ghi\")),\"abc-def-ghi\")")
}

func TestNewLinesInBraces(t *testing.T) {

	grammar := NewShellGrammar()

	//var symbols = NewSafeSymbolTable(&ErrorIfFunctionNotFound{})  // TODO perhaps another default function would be better
	
	test := func (formula string) {
		val := ParseAndEval(&symbols, grammar, formula)
		if val != tuple.Bool(true) {
			t.Errorf("Expected '%s' to be TRUE", formula)
		}
	}

	test("eqt { 1 2 ; 3 4 ; } { 1 2 ; 3 4 }")

	test(`eqt { 1 2 ; 3 4 ; 5 6 } { 1 2
 3 4
 5 6
}`)

	test(`eqt { 1 2 ; 3 4 ; 5 6 } {
 1 2
 3 4
 5 6
}`)

	test(`eqt { 1 2 ; 3 4 ; 5 6 } {
 1 2
 3 4
 5 6}`)

}

func TestExprDeclareFunctions(t *testing.T) {

	grammar := NewShellGrammar()

	//var symbols = NewSafeSymbolTable(&ErrorIfFunctionNotFound{})  // TODO perhaps another default function would be better
	
	test := func (formula string) {
		val := ParseAndEval(&symbols, grammar, formula)
		if val != tuple.Bool(true) {
			t.Errorf("Expected '%s' to be TRUE", formula)
		}
	}

	test("if(true,1,2) == 1")
	test("if(false,1,2) == 2")
	test("if(false,1, cos(PI)) == -1")


	test("nth(0  ( 1 2 3 )) == 1")
	test("nth(1  ( 1 2 3 )) == 2")
	test("nth(2  ( 1 2 3 )) == 3")
	test("nth(-1 ( 1 2 3 )) != 1")
	test("nth(3  ( 1 2 3 )) != 4")

	test("-1 == progn(1+2 3+4 cos(PI)")
	// TODO uses assign test("6==progn (m=3) (s=2) (m*s)")
	
	//test("for a (1 2) { for b (4 5) { a+b }} == ((5 6) (6 7))")
	test(`
func aa a { a*2 }
aa(2)==4`)

	test(`
func aa a { a*2 }
func bb b { aa(b) }
bb(2)==4`)

	test(`
func aa a { a*2 }
func bb b { aa(b*3) }
bb(2)==12`)

	test(`
func aa a { a*2 }
func bb b { aa(b*3) }
bb(2+1)==18`)

	test(`
func aa a { a*2 }
func bb b { aa(b*3)*aa(1-b) }
bb(3)==18*-4`)

}
