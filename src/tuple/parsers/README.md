This package provides an number of parsers for the kind of grammars commonly used in computing.

It is not intended as a general purpose language translation package like Bison or ANTLR,
the idea is to avoid you having to write any grammar, parsers or lexer since those provided should be sufficient for many/most purposes.
All you have to do is provide just the backend functionality that you need.

The package does not interpret the input at all, it just generates a parse tree or Abstract System Tree (AST) which can be processed by another package

## Can read, write and convert between:

* Arithmetic expression grammar similar to those used used by most programming languages, [expr](https://en.wikipedia.org/wiki/Expr) and EXCEL.
* [JSON](https://en.wikipedia.org/wiki/JSON), which is a subset of Javascript.
* [Lisp](https://en.wikipedia.org/wiki/Lisp_(programming_language)) like grammar prefix notation.
* [Lisp](https://en.wikipedia.org/wiki/Lisp_(programming_language)) like grammar with [infix notation](https://en.wikipedia.org/wiki/Infix_notation)
* JSON extended with variables and expressions (or Javascript without loops, objects and functions)
* A [shell](https://en.wikipedia.org/wiki/Unix_shell) like grammar similar to that used by command line interpreters and [TCL](https://en.wikipedia.org/wiki/Tcl)

## Write only

Can write:
* YAML
* Properties
* INI file


# Read, write and processing

## Processing:

* Translation between grammars: for example prefix to and from infix, LISP to JSON, etc.
* A command line [REPL](ehttps://en.wikipedia.org/wiki/Read%E2%80%93eval%E2%80%93print_loop) interface
* Eval - an [evaluator](https://en.wikipedia.org/wiki/Eval) or interpetter
* Query - a utility to filter and select from the Abstract Syntax Tree, similar idea to [jq](https://stedolan.github.io/jq/)


