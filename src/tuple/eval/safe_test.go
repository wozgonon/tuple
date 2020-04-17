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
	"tuple/runner"
	"tuple/eval"
	"tuple/parsers"
)

var safeEvalContext = runner.NewSafeEvalContext(logger)


func TestDeclareFunctions(t *testing.T) {

	grammar := parsers.NewShellGrammar()
	
	test := func (formula string) {
		val,_ := runner.ParseAndEval(safeEvalContext, grammar, formula)
		if val != tuple.Bool(true) {
			t.Errorf("Expected '%s' to be TRUE", formula)
		}
	}

	test("eq (join \"-\" (1 2 3)) \"1-2-3\"")
	test("eq (concat 1 2 3) \"123\"")
	test("eq (concat true 2.2 3 \"x\") \"true2.23x\"")

	test("eq (values (a:2 b:3)) (2 3)")
	test("eq (arity (keys (a:2 b:3))) 2")
	test("eq (arity (keys (a:2 b:3))) 2")

	test("eq (list 1 2 3) (list 1 2 3)")
	test("eq (list) (list)")
	test("eq (list 1 2 3) (1 2 3)")
	// TODO test("eq (list) ()")
	
	test(`
func aa a { a*2 }
aa(2)==4`)

	test("123  == (progn (func a { 123 }) a())")
	test("1234 == (progn (func a b { b }) a(1234))")
	test("15   == (progn (func a b c { b+c }) a(13 2))")
	test("23   == (progn (func a b c d { b+c*d }) a(13 2 5))")

	test("eq (while false 1) ()")
	// eq (for v (list true) { v }) (list true)
	test("if(true,1,2) == 1")
	test("if(false,1,2) == 2")
	test("if(false,1, cos(PI)) == -1")

	safeEvalContext.Add("=", eval.AssignLocal)

	// Test assignment to a variable updates the a local variable or a global variable
	test("progn n=12 { func a b { progn n=b n } } a(234)!=n")
	test("progn n=12 { func a b { progn set(b n) n  } } a(234)==n")
	test("progn abcda=2  2==abcda")

	test(`progn m=(a:1 b:2) (ismap m)`)
	test(`eq "1-2-3" (join "-" (for v (1 2 3) { v }))`)
	test(`eq "1 2" (join " " (for v ( a:1 b:2 ) { v }))`)
	test(`progn a=(forkv k v (a:1 b:2) { concat k v }) (eq "b2" (nth 1 a))`)
}
