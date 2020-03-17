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
    local inSuffix=$1
    local outSuffix=$2
    local suffix2=$3
    local file="src/wozg/testdata/test.${suffix2}"
    echo "-- In: ${inSuffix}"
    cat ${file}
    echo "-- Out: ${outSuffix}"
    bin/wozg --in .${inSuffix} --out .${outSuffix} ${file}
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

ri () {
    make
    run_test l ini l
}

rp () {
    make
    run_test l properties l
}

