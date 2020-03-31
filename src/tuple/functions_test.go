package tuple_test

import (
	"testing"
	"tuple"
)

func TestBuildSymbolTable(t *testing.T) {

	symbols := tuple.NewSymbolTable(tuple.ErrorIfFunctionNotFound)
	//if symbols.FunctionNotFound != tuple.ErrorIfFunctionNotFound {
	//	t.Errorf("Expected error function")
	//}
	if symbols.Count() != 0 {
		t.Errorf("Expected empty table got %d", symbols.Count())

	}
	count := symbols.Count()
	tuple.AddBooleanAndArithmeticFunctions(symbols)
	tuple.AddDeclareFunctions(symbols)
	tuple.AddOperatinSystemFunctions(symbols)
	if symbols.Count() <= count {
		t.Errorf("Expected functions to be added to symbol table")
	}

}
