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

package runner_test

import (
	"testing"
	"tuple"
	"tuple/runner"
	"tuple/parsers"
)

func TestDeclareFunctions(t *testing.T) {

	var safeEvalContext = runner.NewSafeEvalContext(logger)
	grammar := parsers.NewShellGrammar()
	runner.AddTranslatedSafeFunctions(safeEvalContext)
	
	test := func (formula string) {
		val,err := runner.ParseAndEval(safeEvalContext, grammar, formula)
		if val != tuple.Bool(true) {
			t.Errorf("Expected '%s' to be TRUE, val=%s err=%s", formula, val, err)
		}
	}

	test("eq 1 (first (1 2 3))")
	test("eq 2 (second (1 2 3))")
	test("eq 3 (third (1 2 3))")
	test(`eq (print (1 2 3 (a:1 b:(abc 2 3 4 (1 () 2 3))))) "(1 2 3 (a:1 b:(abc 2 3 4 (1 () 2 3))))"`)
	//...
}
