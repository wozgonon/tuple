package parsers_test

import (
	"testing"
	"tuple"
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
		context := parsers.NewParserContext("<eval>", reader, logger)
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
	

	test_regexp := func (regexp string, input string, expected bool) {
		reader := bufio.NewReader(strings.NewReader(regexp))
		context := parsers.NewParserContext("<eval>", reader, logger)
		regexpTree := parsers.ParseRegexp(&context)

		reader = bufio.NewReader(strings.NewReader(input))
		err := parsers.MatchRegexp(reader, regexpTree)
		match := err != nil
		if match == expected {
			t.Errorf("Expected match of input '%s' against regxep '%s' to be %t err=%s", input, regexp, expected, err)
		}
	}
	test_match := func (regexp string, input string) { test_regexp(regexp, input, true) }
	test_not_match := func (regexp string, input string) { test_regexp(regexp, input, false) }
	test_match("a", "a")
	test_not_match("b", "a")
	test_match(".", "a")
	test_match("\\a", "a")
	test_match("\\\\", "\\")
	test_match(".\\\\.", "a\\b")
	test_match("aa", "aa")
	test_match(".a", "aa")
	test_match("..", "aa")
	test_match("a?", "a")
	test_match("a+", "a")
	test_match("a*", "a")
	test_match(".*", "a")
	test_match("ab*", "abbb")
	test_match("a-e", "a")
	test_match("a-e", "c")
	test_match("a-e", "e")
	test_not_match("a-e", "A")
	test_not_match("a-e", "f")
	test_not_match("bb*", "abbb")
	// TODO

}

