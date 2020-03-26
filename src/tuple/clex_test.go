package tuple_test

import (
	"testing"
	"tuple"
	"bufio"
	"strings"
)

func testGetNext(t *testing.T, expression string, expected string) {

	reader := bufio.NewReader(strings.NewReader(expression))
	context := tuple.NewRunnerContext("<eval>", reader, tuple.GetLogger(nil), false)

	result := "..."
	style := tuple.LispWithInfixStyle
	err := style.GetNext(context,
		func(open string) {
			result = open
		},
		func(close string) {
			result = close
		},
		func(atom tuple.Atom) {
			result = atom.Name
		},
		func (literal tuple.Value) {
			// TODO result = literal
		})

	if err != nil {
		t.Errorf("Error: %s", err)
	} else if result != expression {
		t.Errorf("%s==%s", expected, result)
	}
}


func TestLex1(t *testing.T) {
	//testGetNext(t, "1", "1")
	//testGetNext(t, "1.1", "1.1")
	testGetNext(t, "abc123", "abc123")
	testGetNext(t, "+", "+")
	testGetNext(t, ">=", ">=")
	testGetNext(t, "(", "(")
	//testGetNext(t, "[", "]")
	//testGetNext(t, "{", "}")
}
