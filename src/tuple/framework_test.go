package tuple_test

import (
	"testing"
	"tuple"
//	"strings"
//	"math"
//	"tuple/parsers"
)


var NthBitOfInt = tuple.NthBitOfInt

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




