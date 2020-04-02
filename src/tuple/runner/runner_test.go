package runner_test

import (
	"testing"
	"tuple"
	"tuple/parsers"
	"tuple/runner"
	"math"
)


func TestEval1(t *testing.T) {
	var grammar = parsers.NewInfixExpressionGrammar()
	val := runner.ParseAndEval(grammar, symbols, "1+1")
	if val != tuple.Float64(2) {
		t.Errorf("1+1=%s", val)
	}
}

func TestEvalToInt64(t *testing.T) {
	var grammar = parsers.NewInfixExpressionGrammar()
	tests := map[string]float64{
		"1+1" : 2,
		"1." : 1,
		//"-1" : -1,
		//"(((-1)))" : -1,
		"1+2*3" : 7,
		"(1+2)*3" : 9,
		"((1+2)*3)" : 9,
		"((1+2)*3*3/9)" : 3,
		"-1**7*2" : -2,
		"-1**7*2+3" : 1,
		"0*(7)" : 0,
		"8/4" : 2,
		"cos(PI)" : -1,
		"acos(cos(PI))" : math.Pi,
	}
	for k, v := range tests {
		val := runner.ParseAndEval(grammar, symbols, k)
		if val != tuple.Float64(v) {
			t.Errorf("%s=%d   %s", k, int64(v), val)
		}
	}
}

func testFiles(t *testing.T) {

	grammars := tuple.NewGrammars()
	grammars.Add((tuple.NewLispGrammar()))
	expected := 5

	count := 0
	files := []string{"../wozg/testdata/test.l"}
	errors := tuple.RunFiles(files, tuple.GetLogger(nil), false, tuple.NewLispGrammar(), &grammars, func (next tuple.Value) { count += 1})

	if errors > 0 {
		t.Errorf("Found unexpected errors: %d", errors)
	}
	if count != expected {
		t.Errorf("Expected %d got %d", expected, count)
	}
}

func TestSimplePipeline(t *testing.T) {

	table := tuple.NewSafeSymbolTable(&tuple.ErrorIfFunctionNotFound{})
	tuple.SimplePipeline(&table, "*", tuple.NewLispGrammar(), func (_ string) {})
	// TODO
}
