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

	star := Tag{"*"}
	plus := Tag{"+"}
	question := Tag{"?"}
	or := Tag{"|"}
	hythen := Tag{"-"}
	a := tuple.String("a")
	b := tuple.String("b")
	c := tuple.String("c")
	z := tuple.String("z")
	abc := tuple.NewTuple(a, b, c)
	haz := tuple.NewTuple(hythen, a, z)
	bstar := tuple.NewTuple(star, b)
	aquestion := tuple.NewTuple(question, a)

	test := func (expression string, expected Value) {
		reader := bufio.NewReader(strings.NewReader(expression))
		context := parsers.NewParserContext("<eval>", reader, runner.GetLogger(nil, false))
		value := parsers.ParseRegexp(&context)
		if ! reflect.DeepEqual(expected, value) {
			t.Errorf("Expected '%s' got '%s'", expected, value)
		}
	}

	test("[a-z]", haz)
	test("[a-z]|abc", NewTuple(or, haz, abc))
	test("a[a-z]|abc", NewTuple(or, NewTuple(a, haz), abc))
	test("abc", abc)
	test("a|c", NewTuple(or, a, c)) 
	test("abc|c", NewTuple(or, abc, c)) 
	test("abc|abc", NewTuple(or, abc, abc)) 
	test("abc|[a-z]|z", NewTuple(or, NewTuple(or, abc, haz), z))  // TODO should collapse into a single Tuple for 
	test("a?", aquestion) 
	test("b*", bstar) 
	test("c+", NewTuple(plus, c)) 
	test("b*|a?", NewTuple(or, bstar, aquestion)) 
	test("a?|b*", NewTuple(or, aquestion, bstar))
	// TODO
	
}
