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

test: all
	bin/lisp src/lisp/test.l  > /tmp/test.out.l
	bin/lisp /tmp/test.out.l  > /tmp/test.out2.l
	wc -l /tmp/test.out.l /tmp/test.out2.l
	diff -w /tmp/test.out.l /tmp/test.out2.l  # Ignore whitespace
	diff -y /tmp/test.out.l /tmp/test.out2.l

test_tcl: all
	bin/lisp --tcl src/lisp/test.tcl

smoke: test
	bin/lisp --tuple src/lisp/test.l

#############################################################################
#  Clean up
#############################################################################

clean:
	rm -rf bin pkg
