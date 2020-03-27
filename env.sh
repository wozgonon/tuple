#############################################################################
#  To set up the environment for bash and to provide some utilities for bash
#
#   bash$  . env.sh
#
#############################################################################

export GOPATH=`pwd`

alias h=history

#############################################################################
# make test
#############################################################################

mk () {
    make test
}

mck () {
    make clean test
}

alias wexpr=bin/wexpr
alias wozg=bin/wozg
alias wsh=bin/wsh

#############################################################################
# Search for a pattern in the source code
#############################################################################

ff () {
    find . -name '*.go' -exec grep -H "$*" {} \;
}

ffi () {
    find . -name '*.go' -exec grep -Hi "$*" {} \;
}

#############################################################################
# Calculate total lines of code
#############################################################################

loc () {
    wc -l `find . -name '*.go'`
}

# Reports lines of test code and percentage test coverage
loct () {
    wc -l `find . -name '*_test.go'`
    go test tuple -cover
}

cov() {
    go tool cover -html=c.out
}

#############################################################################
#  Run some smoke tests
#############################################################################

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
    run_test tcl tcl tcl
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

