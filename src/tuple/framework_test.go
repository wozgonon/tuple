package tuple_test

import (
	"testing"
	"tuple"
)


func TestGrammars(t *testing.T) {

	grammars := tuple.NewGrammars()
	grammars.Add(tuple.NewLispWithInfixGrammar())
	grammars.Add((tuple.NewLispGrammar()))

	count := 0
	grammars.Forall(func (grammar tuple.Grammar) {
		count += 1
	})
	if count != 2 {
		t.Errorf("Expected %d got %d", 2, count)
	}
}

