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
	"tuple/runner"
	//"tuple/parsers"
	"tuple/eval"
	"os"
)

func TestLessSafe(t *testing.T) {

	var safeEvalContext = runner.NewSafeEvalContext(logger)
	eval.AddLessSafeFunctions(safeEvalContext, symbols.GlobalScope())
	// TODO

/*		test := func (formula string) {
		val,_ := runner.ParseAndEval(safeEvalContext, grammar, formula)
		if val != tuple.Bool(true) {
			t.Errorf("Expected '%s' to be TRUE", formula)
		}
	}*/

	//test("eq (join \"-\" (1 2 3)) \"1-2-3\"")

}

func TestLessSafeOs(t *testing.T) {

	os1 := eval.Os{}
	if os1.Get(0) != Int64(os.Getpid()) {
		t.Errorf("Expected")
	}
	//GetKeyValue
}
