package parsers_test

import (
	"testing"
	"tuple"
	//"math"
	"tuple/runner"
	"tuple/parsers"
	"bufio"
	"strings"
	"reflect"
)

func TestRegexp(t *testing.T) {

	or := Atom{"|"}
	a := tuple.String("a")
	b := tuple.String("b")
	c := tuple.String("c")
	abc := tuple.NewTuple(a, b, c)

	test := func (expression string, expected Value) {
		reader := bufio.NewReader(strings.NewReader(expression))
		context := runner.NewRunnerContext("<eval>", reader, runner.GetLogger(nil), false)
		value := parsers.ParseRegexp(&context)
		if ! reflect.DeepEqual(expected, value) {
			t.Errorf("Expected '%s' got '%s'", expected, value)
		}
	}

	test("abc", abc)
	test("a|c", tuple.NewTuple(or, a, c)) 
	// TODO
	
}
