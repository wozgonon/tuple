package tuple_test

import (
	"testing"
	"tuple"
	"strings"
)


func TestGrammars(t *testing.T) {

	grammars := tuple.NewGrammars()
	grammars.Add(tuple.NewLispWithInfixGrammar())
	grammars.Add((tuple.NewLispGrammar()))
	grammars.Add((tuple.NewTclGrammar()))
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

		smoke := 123
		printed := ""
		grammar.Print(tuple.Int64(smoke), func (value string) {
			printed += value
		})

		if ! strings.Contains(printed, "123") {
			t.Errorf("Expected '%d' in output", smoke)
		}
	})
	if count != 8 {
		t.Errorf("Expected %d got %d", 2, count)
	}
}

