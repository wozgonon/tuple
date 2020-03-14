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

test: all
	bin/lisp src/lisp/test.l

clean:
	rm -rf bin pkg
