#############################################################################
#  To set up the environment for bash and to provide some utilities for bash
#
#   bash$  . env.sh
#
#############################################################################

export GOPATH=`pwd`

ff () {
    # Search all the source files
    find . -name *.go -exec grep -Hi $* {} \;
}

run_test() {
    local suffix=$1
    local suffix2=$2
    local file="src/lisp/testdata/test.${suffix2}"
    cat ${file}
    echo "----"
    bin/lisp --in .${suffix} --out .${suffix} ${file}
}

rl () {
    run_test l l
}

rtcl () {
    run_test tcl fl.tcl
}

rtup () {
    run_test tuple tuple
}

