/*
    This file is part of WOZG.

    WOZG is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    WOZG is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with WOZG.  If not, see <https://www.gnu.org/licenses/>.
*/

package eval_test

import (
	"testing"
	"tuple"
	"tuple/eval"
	"tuple/parsers"
	"tuple/runner"
	"math"
	"math/rand"
)

var NewTuple = tuple.NewTuple
type Bool = tuple.Bool
type Int64 = tuple.Int64
type Float64 = tuple.Float64
type Tag = tuple.Tag

// A random float64 for testing.
// Using random rather than fixed values increases 'statistical sample size' and reduces 'statistical bias'.
func randomFloat64() float64 {
	return rand.Float64()  // TODO perhaps use a normal distribution to be more representative
}

var logger = tuple.GetLogger(nil, false)
var symbols = runner.NewHarmlessEvalContext(logger)

func TestHarmless(t *testing.T) {

	test := func (expression tuple.Tuple) {
		value, _ := eval.Eval(symbols, expression)
		if ! bool(value.(Bool)) {
			t.Errorf("Expected '%s' to be true", expression)
		}
	}

	ONE := Int64(1)
	TWO := Int64(2)
	THREE := Int64(3)
	A1 := NewTuple(Tag{"++"}, ONE)
	A12 := NewTuple(Tag{"+"}, ONE, TWO)
	M23 := NewTuple(Tag{"*"}, TWO, THREE)
	test(NewTuple(Tag{"=="}, THREE, A12))
	test(NewTuple(Tag{"=="}, A12, A12))
	test(NewTuple(Tag{"=="}, M23, M23))
	test(NewTuple(Tag{"=="}, TWO, A1))
}

func TestBinaryFloat64(t *testing.T) {

	test := func (arg string, aa float64, bb float64, expected float64) {
		op := Tag{arg}
		a1 := Float64(aa)
		b1 := Float64(bb)
		lhs := Float64(expected)
		rhs := NewTuple(op, a1, b1)
		expression := NewTuple(Tag{"=="}, lhs, rhs)
		value, err := eval.Eval(symbols, expression)
		if err != nil {
			t.Errorf("Expected '%s' to be true, got error: %s", expression, err)
		} else if ! bool(value.(Bool)) {
			t.Errorf("Expected '%s' to be true", expression)
		}
	}

	r1 := randomFloat64()
	r2 := randomFloat64()
	test("*", r1, r2, r1*r2)
	test("+", r1, r2, r1+r2)
	test("-", r1, r2, r1-r2)
	test("/", r1, r2, r1/r2)
	test("atan2", r1, r2, math.Atan2(r1,r2))
}

func TestBinaryFloat64Bool(t *testing.T) {

	test := func (arg string, aa float64, bb float64) {
		op := Tag{arg}
		a1 := Float64(aa)
		b1 := Float64(bb)
		expression := NewTuple(op, a1, b1)
		value, _ := eval.Eval(symbols, expression)
		if ! bool(value.(Bool)) {
			t.Errorf("Expected '%s' to be true", expression)
		}
	}

	r1 := randomFloat64()
	r2 := randomFloat64() + r1
	test("==", r1, r1)
	test("!=", r1, r2)
	test("<", r1, r2)
	test("<=", r1, r1)
	test("<=", r1, r2)
	test(">", r2, r1)
	test(">=", r2, r2)
	test(">=", r2, r1)
}

func TestUnaryFloat64(t *testing.T) {

	test := func (arg string, aa float64, expected float64) {
		op := Tag{arg}
		a1 := Float64(aa)
		lhs := Float64(expected)
		rhs := NewTuple(op, a1)
		expression := NewTuple(Tag{"=="}, lhs, rhs)
		value, _ := eval.Eval(symbols, expression)
		if ! bool(value.(Bool)) {
			t.Errorf("Expected '%s' to be true", expression)

		}
	}

	r1 := randomFloat64()
	test("-", r1, -r1)
	test("+", r1, +r1)
	test("sqrt", r1, math.Sqrt(r1))
	test("exp", r1, math.Exp(r1))
	test("log", r1, math.Log(r1))
	test("sin", r1, math.Sin(r1))
	test("cos", r1, math.Cos(r1))
	test("tan", r1, math.Tan(r1))
	test("asin", r1, math.Asin(r1))
	test("acos", r1, math.Acos(r1))
	test("atan", r1, math.Atan(r1))
	test("round", r1, math.Round(r1))
}


func TestTestFunctions(t *testing.T) {

	grammar := parsers.NewShellGrammar()
	
	test := func (formula string) {
		val,_ := runner.ParseAndEval(safeEvalContext, grammar, formula)
		if val != tuple.Bool(true) {
			t.Errorf("Expected '%s' to be TRUE", formula)
		}
	}

	test("istuple (1 2 3)")
	test("! (istuple \"a\")")
	test("ismap (1:2 3:4)")
	test("! (ismap \"a\")")
	test("eq (nth 1 (11 22 33)) 22")
}
