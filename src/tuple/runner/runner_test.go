package runner_test

import (
	"testing"
	"tuple"
	"tuple/parsers"
	"tuple/runner"
	"math"
)

//var logger = tuple.GetLogger(nil, false)
//var safeEvalContext = runner.NewSafeEvalContext(logger)

func TestEval1(t *testing.T) {
	var grammar = parsers.NewInfixExpressionGrammar()
	val,_ := runner.ParseAndEval(safeEvalContext, grammar, "1+1")
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
		val,_ := runner.ParseAndEval(safeEvalContext, grammar, k)
		if val != tuple.Float64(v) {
			t.Errorf("%s=%d   %s", k, int64(v), val)
		}
	}
}

func TestSimplePipeline(t *testing.T) {
	
	runner.SimplePipeline(safeEvalContext, true, "*", parsers.NewLispGrammar(), func (_ string) {})
	// TODO
}
