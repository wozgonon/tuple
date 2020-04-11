package eval_test

import (
	"testing"
	"tuple/eval"
)

func TestBuildSymbolTable(t *testing.T) {

	symbols := eval.NewSymbolTable(eval.NewErrorIfFunctionNotFound())
	if symbols.Arity() != 0 {
		t.Errorf("Expected empty table got %d", symbols.Arity())

	}
	count := symbols.Arity()
	eval.AddBooleanAndArithmeticFunctions(&symbols)
	eval.AddOperatingSystemFunctions(&symbols)
	if symbols.Arity() <= count {
		t.Errorf("Expected functions to be added to symbol table")
	}


	
}
