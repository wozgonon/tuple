package tuple_test

import (
	"testing"
	"tuple"
	"bufio"
	"strings"
	"strconv"
	"fmt"
)

const NO_RESULT = "..."

func testGetNext(t *testing.T, expression string, expected string) {

	reader := bufio.NewReader(strings.NewReader(expression))
	context := tuple.NewRunnerContext("<eval>", reader, tuple.GetLogger(nil), false)

	result := NO_RESULT
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
			switch literal.(type) {
			case tuple.Int64: result = strconv.FormatInt(int64(literal.(tuple.Int64)), 10)
			case tuple.Float64: result = fmt.Sprint(float64(literal.(tuple.Float64)))
			case tuple.String: result = string(literal.(tuple.String))
			case tuple.Bool: result = strconv.FormatBool(bool(literal.(tuple.Bool)))
			}
			// TODO result = literal
		})

	if err != nil {
		t.Errorf("Error: %s", err)
	} else if expected != result {
		t.Errorf("%s==%s   exp=%s", expected, result, expression)
	}
	if context.Errors() > 0 {
		t.Errorf("Expected no errors: %d", context.Errors())
	}
}

func TestLex1(t *testing.T) {
	testGetNext(t, "1", "1")
	testGetNext(t, "-1", "-1")
	testGetNext(t, ".1", "0.1")
	//testGetNext(t, "-1.", "-1.")
	// TODO testGetNext(t, "-.1", ".1")
	testGetNext(t, "abc123", "abc123")
	testGetNext(t, "+", "+")
	testGetNext(t, ">=", ">=")
	testGetNext(t, "(", "(")
	//testGetNext(t, "[", "]")
	//testGetNext(t, "{", "}")

	testGetNext(t, ";", NO_RESULT)  // Comments are currently ignored
	testGetNext(t, ";comment", NO_RESULT)
}

func TestCLanguageOperators(t *testing.T) {

	operators := tuple.NewOperators(tuple.LispWithInfixStyle)
	tuple.AddStandardCOperators(&operators)
	operators.Forall(func(operator string) {
		if operator != " " {  // TODO handle space
			testGetNext(t, operator, operator)
		}
	})
}
