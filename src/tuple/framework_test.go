package tuple_test

import (
	"testing"
	"tuple"
	"strings"
	"math"
)


func TestGrammars(t *testing.T) {

	grammars := tuple.NewGrammars()
	grammars.Add(tuple.NewLispWithInfixGrammar())
	grammars.Add((tuple.NewLispGrammar()))
	grammars.Add((tuple.NewShellGrammar()))
	grammars.Add((tuple.NewInfixExpressionGrammar()))
	grammars.Add((tuple.NewYamlGrammar()))
	grammars.Add((tuple.NewIniGrammar()))
	grammars.Add((tuple.NewPropertyGrammar()))
	grammars.Add((tuple.NewJSONGrammar()))

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
		test(tuple.Atom{"abcde"}, "abcde")
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

