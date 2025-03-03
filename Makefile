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

all: bin/wozg bin/whd bin/wsh bin/wozg bin/wexpr pkg/linux_amd64/tuple.a

bin/wsh: wsh
bin/whd: whd
bin/wexpr: wexpr
bin/wozg: wozg
pkg/linux_amd64/tuple.a: tuple

##FLAGS=-ldflags "-X main.Version=$(VERSION)"

wexpr: tuple src/wexpr/wexpr.go ${VERSION_FILE}
	go install $@  ${FLAGS}

wozg: tuple src/wozg/wozg.go ${VERSION_FILE}
	go install $@ ${FLAGS}

wsh: tuple src/wsh/wsh.go ${VERSION_FILE}
	go install $@ ${FLAGS}

whd: tuple src/whd/whd.go ${VERSION_FILE}
	go install $@ ${FLAGS}

tuple: src/tuple/*.go
	go install $@

#############################################################################
# Automatically create a file with version, date and commit infomration
#############################################################################

version_file: ${VERSION_FILE}

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

test: version go_test  test_arithmetic test_expr test_wsh   test_lisp  test_yaml test_wexpr test_infix test_query examples

# test_json


examples: test_dirs bin/wsh examples/availability.wsh
	bin/wsh examples/availability.wsh > ${T1DIR}/availability.txt  && true
	bin/wsh examples/fibonnacci.wsh > ${T1DIR}/fibonnacci.txt  && true

TDIR=src/wozg/testdata/
T1DIR=target/test/1/
T2DIR=target/test/2/
T3DIR=target/test/3/

DIFF=" -y --suppress-common-lines "

version: wozg
	bin/wozg --version

go_test: all
	go test tuple wexpr tuple/eval tuple/parsers tuple/runner -coverprofile=c.out
	echo RUN: go tool cover -html=c.out

test_dirs: 
	mkdir -p ${T1DIR}${TDIR} ${T2DIR}${TDIR} ${T3DIR}${TDIR}

test_lisp: ${TDIR}test.l test_dirs all
	bin/wozg $<  > ${T1DIR}$<
	@bin/wozg ${T1DIR}$<  > ${T2DIR}$<
	@bin/wozg ${T2DIR}$<  > ${T3DIR}$<
	diff -y --suppress-common-lines ${T2DIR}$< ${T3DIR}$<
	diff -y --suppress-common-lines ${T2DIR}$< $<.golden

test_arithmetic: ${TDIR}arithmetic.l test_dirs all
	bin/wozg --eval $<  > ${T1DIR}$<
	@bin/wozg --eval ${T1DIR}$<  > ${T2DIR}$<
	@diff -y --suppress-common-lines ${T1DIR}$< ${T2DIR}$<

test_wsh: ${TDIR}test.wsh test_dirs all
	bin/wozg --in .wsh --out .wsh $<  > ${T1DIR}$<
	@bin/wozg --in .wsh --out .wsh ${T1DIR}$<  > ${T2DIR}$<
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

# This no longer works due to unordered maps
test_json: ${TDIR}test.json all
	@bin/wozg -out json ${TDIR}test.json > ${T1DIR}test.json
	@bin/wozg -out json ${T1DIR}test.json > ${T2DIR}test.json
	##@bin/wozg -out json ${T2DIR}test.json > ${T3DIR}test.json
	##diff -y --suppress-common-lines ${T1DIR}test.json  ${T2DIR}test.json
	diff -y --suppress-common-lines ${T1DIR}test.json ${TDIR}test.json.golden

test_wexpr: bin/wexpr all target
	./src/wexpr/testdata/test.sh > target/test_wexpr.out

smoke: test test_dirs 
	bin/wozg --out sh ${TDIR}test.wsh

#############################################################################
#  Clean up
#############################################################################

clean:
	rm -rf bin target ${VERSION_FILE} # pkg
