#############################################################################
#   Makefile to build components
#
#   make or make all - to build everything
#   make clean       - to remove all the build files
#   make test        - to run tests
#
#############################################################################

VERSION_FILE=src/wozg/version.go
VERSION="0.1"

all: bin/wozg bin/wozg bin/wexpr pkg/linux_amd64/tuple.a


bin/wexpr: wexpr
bin/wozg: wozg
pkg/linux_amd64/tuple.a: tuple

wexpr: tuple src/wexpr/wexpr.go ${VERSION_FILE}
	go install $@

wozg: tuple src/wozg/wozg.go ${VERSION_FILE}
	go install $@

tuple: src/tuple/*.go
	go install $@

#############################################################################
# Automatically create a file with version, date and commit infomration
#############################################################################

${VERSION_FILE}:
	@echo "// Auto generated"  > $@
	@echo "package main" >> $@
	@echo "const BUILT = \"`date -u -Iseconds`\"" >> $@
	@echo "const COMMIT = \"`git rev-parse  HEAD`\""  >> $@
	@echo "const VERSION = \"${VERSION}\""  >> $@


#############################################################################
#  Run tests
# - this software really lends itself to system testing
# - since one can run treat the output as input to a second execution.
# - then compare that output from the second execution is the same as the first.
# - provided the initial input is rich, such a system test will provide a very high
# - test coverage and very little effort.
# - This is a lot to be said for Quality over Quantity when it comes to testing
#   - a small test like this that do a lot of coverage
#   - rather than large numbers of tests with little coverage.
#############################################################################

test: version go_test test_basic test_arithmetic test_expr test_tcl test_yaml test_json test_wexpr test_infix test_query 

TDIR=src/wozg/testdata/
T1DIR=target/test/1/
T2DIR=target/test/2/

DIFF=" -y --suppress-common-lines "

version: wozg
	bin/wozg --version

go_test: all
	go test tuple -coverprofile=c.out
	echo RUN: go tool cover -html=c.out

test_dirs: 
	mkdir -p ${T1DIR}${TDIR} ${T2DIR}${TDIR}

test_basic: ${TDIR}test.l test_dirs all
	bin/wozg $<  > ${T1DIR}$<
	@bin/wozg ${T1DIR}$<  > ${T2DIR}$<
	diff -y --suppress-common-lines ${T1DIR}$< ${T2DIR}$<
	diff -y --suppress-common-lines ${T2DIR}$< $<.golden

test_arithmetic: ${TDIR}arithmetic.l test_dirs all
	bin/wozg --eval $<  > ${T1DIR}$<
	@bin/wozg --eval ${T1DIR}$<  > ${T2DIR}$<
	@diff -y --suppress-common-lines ${T1DIR}$< ${T2DIR}$<

test_tcl: ${TDIR}test.tcl test_dirs all
	bin/wozg --in .tcl --out .tcl $<  > ${T1DIR}$<
	@bin/wozg --in .tcl --out .tcl ${T1DIR}$<  > ${T2DIR}$<
	@diff -y --suppress-common-lines ${T1DIR}$< ${T2DIR}$<

test_expr: ${TDIR}test.expr test_dirs all
	bin/wozg --out .expr $<  > ${T1DIR}$<
	@bin/wozg --out .expr ${T1DIR}$<  > ${T2DIR}$<
	@diff -y --suppress-common-lines ${T1DIR}$< ${T2DIR}$<

test_yaml: ${TDIR}test.l  ${TDIR}test.yaml.golden test_dirs all
	bin/wozg --out .yaml $<  > ${T1DIR}test.yaml
	diff -y --suppress-common-lines ${T1DIR}test.yaml ${TDIR}test.yaml.golden

test_infix: ${TDIR}test.infix test_dirs  all
	bin/wozg $<  > ${T1DIR}test.infix.l
	@bin/wozg ${T1DIR}test.infix.l  > ${T2DIR}test.infix.l
	diff -y --suppress-common-lines ${T1DIR}test.infix.l ${T2DIR}test.infix.l
	diff -y --suppress-common-lines ${T1DIR}test.infix.l ${TDIR}test.infix.l.golden

	@bin/wozg -out .infix ${TDIR}test.infix  > ${T1DIR}test.infix
	@bin/wozg -out .infix ${T1DIR}test.infix  > ${T2DIR}test.infix
	diff -y --suppress-common-lines ${T1DIR}test.infix ${T2DIR}test.infix

	@bin/wozg -out l ${TDIR}test.infix  > ${T1DIR}test.infix.l
	diff -y --suppress-common-lines ${T1DIR}test.infix.l ${T2DIR}test.infix.l
	@bin/wozg -out infix ${T1DIR}test.infix.l  > ${T1DIR}test.infix
	diff -y --suppress-common-lines ${T1DIR}test.infix ${T2DIR}test.infix

test_query: all
	bin/wozg --query a.*.c ${TDIR}nested.l > ${T1DIR}nested.l
	diff -y --suppress-common-lines ${T1DIR}nested.l ${TDIR}nested.l.golden

test_json: ${TDIR}test.json all
	@bin/wozg -out json ${TDIR}test.json > ${T1DIR}test.json
	@bin/wozg -out json ${T1DIR}test.json > ${T2DIR}test.json
	diff -y --suppress-common-lines ${T1DIR}test.json  ${T2DIR}test.json
	diff -y --suppress-common-lines ${T1DIR}test.json ${TDIR}test.json.golden

test_wexpr: bin/wexpr all
	test 11 = `bin/wexpr "11"`
	test () = `bin/wexpr "()"`
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
	test 1 = `bin/wexpr -- "-(-(-1)+2)"`
	test 3 = `bin/wexpr -- "(0- - 3)"`
	test x-3 = x`bin/wexpr -- "-(0- - 3)"`
	test x-2 = x`bin/wexpr -- "-(1--1)"`
	test x-3 = x`bin/wexpr -- "-(0- - - - 3)"`
	test x-3 = x`bin/wexpr -- "-(0--3)"`
	test 1 = `bin/wexpr -- "cos(0)"`
	test -1 = `bin/wexpr -- "cos(PI)"`
	test 3.141592653589793 = `bin/wexpr -- "acos(cos(PI))"`
	test true = `bin/wexpr -- "(acos(cos(PI)))==PI"`
	@bin/wexpr  +     || true
	@bin/wexpr  "(+"  || true
	@bin/wexpr  "+("  || true
	@bin/wexpr  "("  || true
	@bin/wexpr  ")"  || true
	test "0 == `bin/wexpr "atan2(0 1)"`"
	test "0.7853981633974483 == `bin/wexpr "atan2(1 1)"`"




smoke: test test_dirs 
	bin/wozg --out tcl ${TDIR}test.tcl

#############################################################################
#  Clean up
#############################################################################

clean:
	rm -rf bin target ${VERSION_FILE} # pkg
