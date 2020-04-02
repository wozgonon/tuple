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
)

var NewTuple = tuple.NewTuple
type Bool = tuple.Bool
type Int64 = tuple.Int64
type Float64 = tuple.Float64
type Atom = tuple.Atom

func TestHarmless(t *testing.T) {

	symbols := eval.NewHarmlessSymbolTable(&eval.ErrorIfFunctionNotFound{})
	if symbols.Count() == 0  {
		t.Errorf("Expected functions to be added to symbol table")
	}

	test := func (expression tuple.Tuple) {
		if ! bool((eval.Eval(&symbols, expression)).(Bool)) {
			t.Errorf("Expected '%s' to be true", expression)

		}
	}
	
	ONE := Int64(1)
	TWO := Int64(2)
	THREE := Int64(3)
	A1 := NewTuple(Atom{"++"}, ONE)
	A12 := NewTuple(Atom{"+"}, ONE, TWO)
	M23 := NewTuple(Atom{"*"}, TWO, THREE)
	test(NewTuple(Atom{"=="}, THREE, A12))
	test(NewTuple(Atom{"=="}, A12, A12))
	test(NewTuple(Atom{"=="}, M23, M23))
	test(NewTuple(Atom{"=="}, TWO, A1))
}
