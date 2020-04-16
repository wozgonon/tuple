package runner_test

import (
	"testing"
	"tuple"
	"tuple/parsers"
	"tuple/runner"
	"tuple/eval"
	"math"
	"strings"
	"reflect"
)

var NewTuple = tuple.NewTuple
type Int64 = tuple.Int64
type Tag = tuple.Tag
var logger = tuple.NewDefaultLocationLogger()
var safeEvalContext = runner.NewSafeEvalContext(logger)


func testFiles(t *testing.T) {

	expected := 5

	count := 0
	files := []string{"../wozg/testdata/test.l"}
	grammars := runner.NewGrammars(parsers.NewLispGrammar())
	errors := grammars.RunFiles(logger, files, func (next tuple.Value) error {
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

func TestGrammars(t *testing.T) {

	defaultGrammar := parsers.NewLispGrammar()
	grammars := runner.NewGrammars(defaultGrammar)
	grammars.AddAllKnownGrammars()

	if grammars.Default().Name() != defaultGrammar.Name() {
		t.Errorf("Expected '%s'", defaultGrammar)
	}
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
		test(NewTuple(tuple.Int64(-1234)), "-1234")
		test(tuple.Bool(false), "false")  //  'false' might not be valid for all grammars
	})
	if count != 8 {
		t.Errorf("Expected %d got %d", 2, count)
	}
}

func TestGrammarsFunctions(t *testing.T) {

	defaultGrammar := parsers.NewLispGrammar()
	grammars := runner.NewGrammars(defaultGrammar)
	grammars.AddAllKnownGrammars()

	ifNotFound := eval.NewErrorIfFunctionNotFound()
	evalContext := eval.NewRunner(ifNotFound, logger)
	eval.AddSafeFunctions(&evalContext)
	grammars.AddSafeGrammarFunctions(&evalContext)
	runner.AddTranslatedSafeFunctions(&evalContext)

	val,_ := runner.ParseAndEval(&evalContext, grammars.Default(), "(grammars)")
	if val.Arity() != grammars.Arity() {
		t.Errorf("Expected %d, got %d", val.Arity(), grammars.Arity())
	}
	//val,_ := runner.ParseAndEval(&evalContext, grammars.Default(), "(grammars)")
	//if ! reflect.DeepEqual(val, grammars) {
	//	t.Errorf("Expected grammars, got %s", reflect.TypeOf(val)
	//}
	val,_ = runner.ParseAndEval(&evalContext, grammars.Default(), "ctx")
	if ! reflect.DeepEqual(val, evalContext.GlobalScope().Root()) {
		t.Errorf("Expected global, got %s", val)
	}
	val,_ = runner.ParseAndEval(&evalContext, grammars.Default(), "second ( 11 22 33 )")
	if val != Int64(22) {
		t.Errorf("Expected 22, got %s", val)
	}
	val,_ = runner.ParseAndEval(&evalContext, grammars.Default(), "expr(\"1+2\")")
	if val == tuple.Int64(3) {
		t.Errorf("Expected 3, got %s", val)
	}
	val,_ = runner.ParseAndEval(&evalContext, grammars.Default(), "(expr2 \"l\" \"(+ 1 2)\")")
	if val == tuple.Int64(3) {
		t.Errorf("Expected 3, got %s", val)
	}
	/*
	p12 := NewTuple(Tag{"+"}, Int64(1), Int64(2))
	val,_ = runner.ParseAndEval(&evalContext, grammars.Default(), "ast2(\"l\" \"(+ 1 2)\")")
	if ! reflect.DeepEqual(val, p12) {
		t.Errorf("Expected '%s', got %s type=%s", p12, val, reflect.TypeOf(val))
	}*/


}

	

func TestRunFilesFunctions(t *testing.T) {
	defaultGrammar := parsers.NewLispGrammar()
	grammars := runner.NewGrammars(defaultGrammar)
	grammars.AddAllKnownGrammars()
	files := []string{"../../wozg/testdata/test.l", "../../wozg/testdata/test.json", "../../wozg/testdata/test.infix"}
	errors := grammars.RunFiles(logger, files, func(value tuple.Value) error { return nil })
	if errors != 0 {
		t.Errorf("Expected no errors got %d", errors)
	}


	files = []string{"../../wozg/testdata/test.unknown"}
	errors = grammars.RunFiles(logger, files, func(value tuple.Value) error { return nil })
	if errors != 1 {
		t.Errorf("Expected one errors got %d", errors)
	}

	files = []string{"../../wozg/testdata/test.unknown"}
	_, err := grammars.RunFile(logger, "unknown_file_skjjs.l", func(value tuple.Value) error { return nil })
	if err == nil {
		t.Errorf("Expected errors got %d", errors)
	}

}
