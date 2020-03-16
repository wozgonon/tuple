#############################################################################
#   Makefile to build components
#
#   make or make all - to build everything
#   make clean       - to remove all the build files
#   make test        - to run tests
#
#############################################################################

all: bin/lisp pkg/linux_amd64/tuple.a bin/jml # forth json

bin/lisp: lisp
bin/jml: jml
pkg/linux_amd64/tuple.a: tuple

lisp: tuple src/lisp/lisp.go
	go install $@

jml: tuple src/jml/jml.go
	go install $@

tuple: src/tuple/tuple.go
	go install $@

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

test: test_basic test_arithmetic test_tcl # test_infix

TDIR=src/lisp/testdata/
T1DIR=target/test/1/
T2DIR=target/test/2/

DIFF=" -y --suppress-common-lines "

test_dirs: 
	mkdir -p ${T1DIR}/${TDIR} ${T2DIR}/${TDIR}

test_basic: ${TDIR}test.l test_dirs all
	bin/lisp $<  > ${T1DIR}$<
	@bin/lisp ${T1DIR}$<  > ${T2DIR}$<
	@diff -y --suppress-common-lines ${T1DIR}$< ${T2DIR}$<
	@diff -y --suppress-common-lines ${T2DIR}$< $<.golden

test_arithmetic: ${TDIR}arithmetic.l test_dirs all
	bin/lisp --eval $<  > ${T1DIR}$<
	@bin/lisp --eval ${T1DIR}$<  > ${T2DIR}$<
	@diff -y --suppress-common-lines ${T1DIR}$< ${T2DIR}$<

test_tcl: ${TDIR}test.fl.tcl test_dirs all
	bin/lisp --in .tcl --out .tcl $<  > ${T1DIR}$<
	@bin/lisp --in .tcl --out .tcl ${T1DIR}$<  > ${T2DIR}$<
	@diff -y --suppress-common-lines ${T1DIR}$< ${T2DIR}$<

test_tuple: ${TDIR}test.tuple test_dirs all
	bin/lisp --out .tuple $<  > ${T1DIR}$<
	@bin/lisp --out .tuple ${T1DIR}$<  > ${T2DIR}$<
	@diff -y --suppress-common-lines ${T1DIR}$< ${T2DIR}$<

test_infix: ${TDIR}infix.l test_dirs  all
	bin/lisp $<  > ${T1DIR}$<
	@bin/lisp ${T1DIR}$<  > ${T2DIR}$<
	@diff -y --suppress-common-lines ${T1DIR}$< ${T2DIR}$<

smoke: test test_dirs 
	bin/lisp --out tcl ${TDIR}test.fl.tcl

#############################################################################
#  Clean up
#############################################################################

clean:
	rm -rf bin target # pkg
