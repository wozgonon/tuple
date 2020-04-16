package parsers_test

import (
//	"testing"
	"tuple"
	"tuple/parsers"
	"tuple/runner"
	"tuple/eval"
//	"bufio"
//	"strings"
//	"strconv"
//	"fmt"
)

type Grammar = tuple.Grammar
type Context = tuple.Context
type Tag = tuple.Tag
type Value = tuple.Value
type StringFunction = tuple.StringFunction
type String = tuple.String
type Tuple = tuple.Tuple
type Next = tuple.Next
type Lexer = tuple.Lexer
type Float64 = tuple.Float64
type Int64 = tuple.Int64

var CONS_ATOM = tuple.CONS_ATOM
var NewTuple = tuple.NewTuple
var Error = tuple.Error
var Verbose = tuple.Verbose

var NewJSONGrammar = parsers.NewJSONGrammar
var NewInfixExpressionGrammar = parsers.NewInfixExpressionGrammar
var NewShellGrammar = parsers.NewShellGrammar
var ParseAndEval = runner.ParseAndEval
var NewLispGrammar = parsers.NewLispGrammar
var NewLispWithInfixGrammar = parsers.NewLispWithInfixGrammar

type ErrorIfFunctionNotFound = eval.ErrorIfFunctionNotFound
