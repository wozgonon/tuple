package runner_test

import (
	"testing"
	"tuple"
	"tuple/parsers"
	"tuple/runner"
	"math"
	"strings"
)

var logger = tuple.GetLogger(nil, false)
var safeEvalContext = runner.NewSafeEvalContext(logger)

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

func testFiles(t *testing.T) {

	expected := 5

	count := 0
	files := []string{"../wozg/testdata/test.l"}
	grammars := runner.NewGrammars(parsers.NewLispGrammar())
	errors := runner.RunFiles(&grammars, logger, files, func (next tuple.Value) error {
		count += 1
		return nil
	})

	if errors > 0 {
		t.Errorf("Found unexpected errors: %d", errors)
	}
	if count != expected {
		t.Errorf("Expected %d got %d", expected, count)
	}
}

func TestSimplePipeline(t *testing.T) {
	
	runner.SimplePipeline(safeEvalContext, true, "*", parsers.NewLispGrammar(), func (_ string) {})
	// TODO
}


func TestGrammars(t *testing.T) {

	grammars := runner.NewGrammars(parsers.NewLispGrammar())
	runner.AddAllKnownGrammars(&grammars)

	count := 0
	grammars.Forall(func (grammar tuple.Grammar) {
		count += 1

		// TODO check uniqueness
		name := grammar.Name()
		suffix := grammar.FileSuffix()
		if name!= "" && suffix != "" && ! strings.HasPrefix(suffix, ".") {
			t.Errorf("Expected name and suffix, got '%s', '%s'", name, suffix)
		}

		if g, ok := grammars.FindBySuffix(suffix); ! ok || g.Name() != name {
			t.Errorf("Expected find by suffix '%s' to return grammar: '%s'", suffix, name)
		}

		suffixWithoutDot := strings.Replace(suffix, ".", "", 999)
		if g, ok := grammars.FindBySuffix(suffixWithoutDot); ! ok || g.Name() != name {
			t.Errorf("Expected find by suffix '%s' to return grammar: '%s'", suffixWithoutDot, name)
		}

		test := func (value tuple.Value, expected string) {
			printed := ""
			grammar.Print(value, func (value string) {
				printed += value
			})
			if ! strings.Contains(printed, expected) {
				t.Errorf("Expected '%s' in output", expected)
			}
		}
		test(tuple.Tag{"abcde"}, "abcde")
		test(tuple.Float64(-1.123), "-1.123")
		test(tuple.Float64(math.NaN()), "NaN")
		test(tuple.Float64(math.Inf(1)), "Inf")
		test(tuple.Int64(123), "123")
		test(tuple.String("abc"), "abc")
		test(tuple.NewTuple(tuple.Int64(-1234)), "-1234")
		test(tuple.Bool(false), "false")  //  'false' might not be valid for all grammars
	})
	if count != 8 {
		t.Errorf("Expected %d got %d", 2, count)
	}
}
