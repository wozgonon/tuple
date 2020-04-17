package tuple_test

import (
	"testing"
	"tuple"
//	"strings"
	"math"
//	"tuple/parsers"
)


var NthBitOfInt = tuple.NthBitOfInt

func TestStructs(t *testing.T) {

	var array tuple.Array = tuple.NewTuple()
	if _, ok := array.(tuple.Array); ! ok {
		t.Errorf("Expected array got %s", array)
	}

	var mapp tuple.Map = tuple.NewTagValueMap()
	if _, ok := mapp.(tuple.Map); ! ok {
		t.Errorf("Expected map got %s", mapp)
	}

}

func TestNthBitOfInt(t *testing.T) {

	if NthBitOfInt(0, 0) {
		t.Errorf("ERROR: expected false")
	}
	if NthBitOfInt(0, 1) {
		t.Errorf("ERROR: expected false")
	}
	if NthBitOfInt(0, 2) {
		t.Errorf("ERROR: expected false")
	}


	if ! NthBitOfInt(1, 0) {
		t.Errorf("ERROR: expected true")
	}
	if NthBitOfInt(1, 1) {
		t.Errorf("ERROR: expected false")
	}
	if NthBitOfInt(1, 2) {
		t.Errorf("ERROR: expected false")
	}

	if NthBitOfInt(2, 0) {
		t.Errorf("ERROR: expected false")
	}
	if ! NthBitOfInt(2, 1) {
		t.Errorf("ERROR: expected true")
	}
	if NthBitOfInt(2, 2) {
		t.Errorf("ERROR: expected false")
	}
}

func TestConversions(t *testing.T) {
	if tuple.IntToString(1234) != "1234" {
		t.Errorf("Expected 1234")
	}
	if tuple.FloatToString(123.456) != "123.456" {
		t.Errorf("Expected")
	}
	if tuple.FloatToString(math.Inf(1)) != "Inf" {
		t.Errorf("Expected")
	}
	if tuple.Float64ToString(tuple.Float64(145.002)) != "145.002" {
		t.Errorf("Expected")
	}
	if tuple.Int64ToString(tuple.Int64(1234)) != "1234" {
		t.Errorf("Expected")
	}
	expected := tuple.Tag{"12345"}
	if tuple.IntToTag(12345) != expected {
		t.Errorf("Expected")
	}
	if tuple.BoolToFloat(true) != 1. {
		t.Errorf("Expected")
	}
	if tuple.BoolToInt(false) != 0 {
		t.Errorf("Expected")
	}
	if tuple.BoolToFloat(false) != 0. {
		t.Errorf("Expected")
	}
	if tuple.BoolToInt(true) != 1 {
		t.Errorf("Expected")
	}
	if tuple.BoolToString(true) != "true" {
		t.Errorf("Expected")
	}

	if tuple.DoubleQuotedString("abc") != "\"abc\"" {
		t.Errorf("Expected")
	}
	if tuple.DoubleQuotedString("a\nb\tc") != "\"a\\nb\\tc\"" {
		t.Errorf("Expected")
	}
}
