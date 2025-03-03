package parsers_test

import (
	"testing"
	"tuple"
	"tuple/parsers"
	"bufio"
	"strings"
	"strconv"
	"fmt"
)

const NO_RESULT = "..."

func testGetNext(t *testing.T, logger tuple.LocationLogger, expression string, expected string) {

	reader := bufio.NewReader(strings.NewReader(expression))
	context := parsers.NewParserContext("<eval>", reader, logger)

	result := NO_RESULT
	style := parsers.LispStyle()
	err := style.GetNext(&context,
		func() {},
		func(open string) {
			result = open
		},
		func(close string) {
			result = close
		},
		func(tag tuple.Tag) {
			result = tag.Name
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
		t.Errorf("expected '%s' got '%s' given expression='%s'", expected, result, expression)
	}
	if context.Errors() > 0 {
		t.Errorf("Expected no errors: %d", context.Errors())
	}
}

func TestLex1(t *testing.T) {

	testLex1(t, tuple.NewDefaultLocationLogger())
	testLex1(t, tuple.NewGrammarLogger(parsers.NewLispWithInfixGrammar()))
}

func testLex1(t *testing.T, logger tuple.LocationLogger) {

	testGetNext(t, logger, "1", "1")
	testGetNext(t, logger, "-", "-")
	testGetNext(t, logger, ".1", "0.1")
	//testGetNext(t, logger, "1.", "-1.")
	// TODO testGetNext(t, logger, "-.1", ".1")
	testGetNext(t, logger, "abc123", "abc123")
	testGetNext(t, logger, "+", "+")
	testGetNext(t, logger, ">=", ">=")
	testGetNext(t, logger, "(", "(")
	//testGetNext(t, logger, "[", "]")
	//testGetNext(t, logger, "{", "}")

	testGetNext(t, logger, ";", NO_RESULT)  // Comments are currently ignored
	testGetNext(t, logger, ";comment", NO_RESULT)
}

func TestCLanguageOperators(t *testing.T) {

	logger := tuple.NewGrammarLogger(parsers.NewLispGrammar())
	operators := parsers.NewOperators(parsers.LispStyle())
	parsers.AddStandardCOperators(&operators)
	operators.Forall(func(operator string) {
		if operator != " " && operator != ";" && operator != "," {  // TODO handle space
			testGetNext(t, logger, operator, operator)
		}
	})
}
