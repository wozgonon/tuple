package main_test


import (
	"testing"
	"strings"
	"bytes"
	"os/exec"
)




func TestExpr(t *testing.T) {

	test := func(expression string, expected string) {
		commands := strings.Split(" ", expression)
		cmd := exec.Command("../../bin/wexpr", commands...)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			t.Errorf("Expected '%s' to succeed, error='%s'", expression, err)
		}
		output := out.String()
		if output !=  expected {
			// TODO t.Errorf("Expected '%s' got '%s'", expected, output)
		}
	}

	test("-1", "-1")
	test("11", "11")
	test("11+2", "22")
	test("1+2*3", "7")
	test("-1+-2*-3", "5")
	test("-1+-2*-3/2", "5")
}
