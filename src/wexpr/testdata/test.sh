#!/bin/bash
# -ev

test 11 = `bin/wexpr "11"`
test "()" = `bin/wexpr "()"`
test 7 = `bin/wexpr "1+2*3"`
test 5 = `bin/wexpr "1*2+3"`
test 120 = `bin/wexpr 1*2*3*4*5`
test 6 = `bin/wexpr "(1)+((2))+(((3)))"`
test 22 = `bin/wexpr "((22))"`
test 22 = `bin/wexpr "((((22))))"`
test 3 = `bin/wexpr "(1+2)"`
test 9 = `bin/wexpr "(1+2)*3"`
test 10 = `bin/wexpr "1+2+3+4"`
test 10 = `bin/wexpr "1+(2+3)+4"`
test 10 = `bin/wexpr "(1+2+3+4)"`
test 10 = `bin/wexpr "((1+((2)+3))+(4))"`
test x-123 = x`bin/wexpr -- "-123"`
test x-123 = x`bin/wexpr -- "-(123)"`
test -3 = `bin/wexpr -- "-(1+2)"`
test -3 = `bin/wexpr -- "-(-(-1)+2)"`
test 3 = `bin/wexpr -- "(0- - 3)"`
test x-3 = x`bin/wexpr -- "-(0- - 3)"`
test x-2 = x`bin/wexpr -- "-(1--1)"`
test x-3 = x`bin/wexpr -- "-(0- - - - 3)"`
test x-3 = x`bin/wexpr -- "-(0--3)"`
test 1 = `bin/wexpr -- "cos(0)"`
test -1 = `bin/wexpr -- "cos(PI)"`
test 3.141592653589793 = `bin/wexpr -- "acos(cos(PI))"`
test true = `bin/wexpr -- "(acos(cos(PI)))==PI"`
test "0 == `bin/wexpr "atan2(0 1)"`"
test "0.7853981633974483 == `bin/wexpr "atan2(1 1)"`"

# Test for syntax errors

expect_fail () {
    bin/wexpr "$1" 2> /dev/null  || true
}

expect_fail  "()-"
expect_fail  "()/"
expect_fail  "(*)"
expect_fail  "+"
expect_fail  "+/"
expect_fail  "(+"
expect_fail  "+("
expect_fail  "("
expect_fail  ")"
