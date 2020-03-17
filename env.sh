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
    local outSuffix=$2
    local suffix2=$3
    local file="src/lisp/testdata/test.${suffix2}"
    echo "-- In:"
    cat ${file}
    echo "-- Out:"
    bin/lisp --in .${suffix} --out .${outSuffix} ${file}
}

rl () {
    run_test l l l
}

rtcl () {
    run_test tcl tcl fl.tcl
}

rtup () {
    run_test tuple tuple tuple
}

ry () {
    run_test l yaml l
}

