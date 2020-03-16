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

test: test_basic test_arithmetic

test_basic: src/lisp/test.l all
	@mkdir -p /tmp/1/src/lisp /tmp/2/src/lisp
	@bin/lisp $<  > /tmp/1/$<
	@bin/lisp /tmp/1/$<  > /tmp/2/$<
	@diff /tmp/1/$< /tmp/2/$<

test_arithmetic: src/lisp/arithmetic.l all
	@mkdir -p /tmp/1/src/lisp /tmp/2/src/lisp
	@bin/lisp --eval $<  > /tmp/1/$<
	@bin/lisp --eval /tmp/1/$<  > /tmp/2/$<
	@diff /tmp/1/$< /tmp/2/$<

test_tcl: src/lisp/test.tcl all
	bin/lisp --tcl $<

smoke: test
	bin/lisp --tuple src/lisp/test.l

#############################################################################
#  Clean up
#############################################################################

clean:
	rm -rf bin pkg
