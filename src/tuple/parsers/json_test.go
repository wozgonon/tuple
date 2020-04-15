package parsers_test

import (
	"testing"
	"tuple"
	"tuple/parsers"
	"reflect"
	"strings"
)

var zero = tuple.Int64(0)
var one = tuple.Int64(1)

func TestEvalJson(t *testing.T) {
	var grammar = NewJSONGrammar()

	if len(grammar.Name()) <= 0 {
		t.Errorf("%s", grammar.Name())
	}
	if grammar.FileSuffix() != ".json" {
		t.Errorf("%s", grammar.FileSuffix())
	}

		
	test := func(formula string, expected tuple.Value) {
		val,_ := ParseAndEval(safeEvalContext, grammar, formula)
		if ! reflect.DeepEqual(val, expected) {
			t.Errorf("%s=%f  expected=%s", formula, val, expected)
		}
	}

	test("true", tuple.Bool(true))
	test("false", tuple.Bool(false))
	test("123", tuple.Int64(123))
	test("0", zero) 
	test("-0", zero)
	// TODO test("-1", tuple.Int64(-1))
	test("123.", tuple.Float64(123))
	// TODO test("-123.", tuple.Float64(-123))
	test(".123", tuple.Float64(.123))
	test("1234.11", tuple.Float64(1234.11))
	test("\"abc\"", tuple.String("abc"))
	test("\"a\\nb\\tc\"", tuple.String("a\nb\tc"))
	test("[]", tuple.NewTuple())
	// TODO...
}

func TestEvalParseJson(t *testing.T) {
	var grammar = NewJSONGrammar()

	test := func(formula string, expected tuple.Value) {
		val, err := parsers.ParseString(logger, grammar, formula)
		if err != nil {
			t.Errorf("Given '%s' expected '%s' got error '%s'", formula, expected, err)
		}
		//t.Logf("Tuple %d %d", val.Arity(), expected.Arity())
		if ! reflect.DeepEqual(val, expected) {
			t.Errorf("%s=%f  expected=%s", formula, val, expected)
		}
	}

	zero := tuple.Int64(0)
	one := tuple.Int64(1)

	t01 := tuple.NewTuple(zero, one)
	t001 := tuple.NewTuple(zero, t01)
	empty := tuple.NewTuple()
	
	test("[]", empty)
	test("[0, 1]", t01)
	test("[0 1]", t01)  // TODO space
	/// TODO   test("[0]", tuple.NewTuple(zero))  FIXME - retain brackets
	test("[0, [0,1]]", t001)

	//a := tuple.String("a")
	//b := tuple.String("b")
	//cons := tuple.CONS_ATOM

	//ac1 := tuple.NewTuple(cons, a, one)
	//bc0 := tuple.NewTuple(cons, b, zero)

	mmap := tuple.NewTagValueMap()
	mmap.Add(Tag{"a"}, empty)

	test("{}", empty)
	test("{\"a\" : [] }", mmap)
	mmap.Add(Tag{"a"}, t01)
	test("{\"a\" : [0, 1] }", mmap)
	mmap.Add(Tag{"a"}, one)
	test("{\"a\" : 1 }", mmap)
	mmap.Add(Tag{"b"}, zero)
	test("{\"a\" : 1, \"b\" : 0 }", mmap)

	// TODO...
}

func TestJsonPrint(t *testing.T) {
	var grammar = NewJSONGrammar()

	test := func(value tuple.Value, expected string) {
		result := ""
		grammar.Print(value, func (value string) {
			result += value
		})
		result = strings.Replace(result, "\n", "", 99999999)
		result = strings.Replace(result, " ", "", 99999999)
		if result != expected {
			t.Errorf("expected '%s' got '%s'", expected, result)
		}
	}

	test(one, "1")  // TODO this is not actually JSON
	test(tuple.NewTuple(), "[]")
	test(tuple.NewTuple(one), "[1]")
	test(tuple.NewTuple(zero, one), "[0,1]")
	// TODO
}
