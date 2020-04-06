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

//var logger = tuple.GetLogger(nil, false)
//var notFound = eval.NewErrorIfFunctionNotFound(logger)
//var symbols = eval.NewHarmlessSymbolTable(notFound)

func TestSafe(t *testing.T) {

	symbols := eval.NewSafeSymbolTable(notFound)
	if symbols.Count() == 0  {
		t.Errorf("Expected functions to be added to symbol table")
	}


	
}


func TestDeclareFunctions(t *testing.T) {

	grammar := parsers.NewShellGrammar()

	var symbols = eval.NewSafeSymbolTable(notFound)  // TODO perhaps another default function would be better
	
	test := func (formula string) {
		val := runner.ParseAndEval(&symbols, grammar, formula)
		if val != tuple.Bool(true) {
			t.Errorf("Expected '%s' to be TRUE", formula)
		}
	}

	test(`
func aa a { a*2 }
aa(2)==4`)

	test("123  == (progn (func a { 123 }) a())")
	test("1234 == (progn (func a b { b }) a(1234))")
	test("15   == (progn (func a b c { b+c }) a(13 2))")
	test("23   == (progn (func a b c d { b+c*d }) a(13 2 5))")

	test("if(true,1,2) == 1")
	test("if(false,1,2) == 2")
	test("if(false,1, cos(PI)) == -1")

	symbols.Add("=", eval.AssignLocal)

	// Test assignment to a variable updates the a local variable or a global variable
	test("progn n=12 { func a b { progn n=b n } } a(234)!=n")
	test("progn n=12 { func a b { progn set(b n) n  } } a(234)==n")
}
