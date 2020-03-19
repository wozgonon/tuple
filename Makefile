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

all: bin/wozg pkg/linux_amd64/tuple.a

bin/wozg: wozg
pkg/linux_amd64/tuple.a: tuple

wozg: tuple src/wozg/wozg.go ${VERSION_FILE}
	go install $@

tuple: src/tuple/tuple.go
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

test: version test_basic test_arithmetic test_tcl test_yaml test_query test_json # test_infix

TDIR=src/wozg/testdata/
T1DIR=target/test/1/
T2DIR=target/test/2/

DIFF=" -y --suppress-common-lines "

version: wozg
	bin/wozg --version

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

test_tuple: ${TDIR}test.tuple test_dirs all
	bin/wozg --out .tuple $<  > ${T1DIR}$<
	@bin/wozg --out .tuple ${T1DIR}$<  > ${T2DIR}$<
	@diff -y --suppress-common-lines ${T1DIR}$< ${T2DIR}$<

test_yaml: ${TDIR}test.l  ${TDIR}test.yaml.golden test_dirs all
	bin/wozg --out .yaml $<  > ${T1DIR}test.yaml
	diff -y --suppress-common-lines ${T1DIR}test.yaml ${TDIR}test.yaml.golden

test_infix: ${TDIR}infix.l test_dirs  all
	bin/wozg $<  > ${T1DIR}$<
	@bin/wozg ${T1DIR}$<  > ${T2DIR}$<
	@diff -y --suppress-common-lines ${T1DIR}$< ${T2DIR}$<

test_query:
	bin/wozg --query a.*.c ${TDIR}nested.l > ${T1DIR}nested.l
	diff -y --suppress-common-lines ${T1DIR}nested.l ${TDIR}nested.l.golden

test_json: ${TDIR}test.json all
	@bin/wozg -out json ${TDIR}test.json > ${T1DIR}test.json
	@bin/wozg -out json ${T1DIR}test.json > ${T2DIR}test.json
	diff -y --suppress-common-lines ${T1DIR}test.json  ${T2DIR}test.json
	diff -y --suppress-common-lines ${T1DIR}test.json ${TDIR}test.json.golden

smoke: test test_dirs 
	bin/wozg --out tcl ${TDIR}test.tcl

#############################################################################
#  Clean up
#############################################################################

clean:
	rm -rf bin target ${VERSION_FILE} # pkg
