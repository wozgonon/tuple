package eval_test

import (
	"testing"
	"tuple/eval"
)

func TestBuildSymbolTable(t *testing.T) {

	symbols := eval.NewSymbolTable(&eval.ErrorIfFunctionNotFound{})
	if symbols.Count() != 0 {
		t.Errorf("Expected empty table got %d", symbols.Count())

	}
	count := symbols.Count()
	eval.AddBooleanAndArithmeticFunctions(&symbols)
	eval.AddOperatingSystemFunctions(&symbols)
	if symbols.Count() <= count {
		t.Errorf("Expected functions to be added to symbol table")
	}


	
}
