The tuple package provides functions for parsing, printing and processing a number of simple Grammars.
The parsers produce a common [parse tree](https://en.wikipedia.org/wiki/Parse_tree)
and [AST](https://en.wikipedia.org/wiki/Abstract_syntax_tree).

Can read and write:

* [Lisp](https://en.wikipedia.org/wiki/Lisp_(programming_language)) like grammar prefix notation
* [Lisp](https://en.wikipedia.org/wiki/Lisp_(programming_language)) like grammar with [infix notation](https://en.wikipedia.org/wiki/Infix_notation)
* Arithmetic expression grammar like that used used by most programming languages, [expr](https://en.wikipedia.org/wiki/Expr) and EXCEL
* [JSON](https://en.wikipedia.org/wiki/JSON), which is a subset of Javascript.
* JSON extended with variables and expressions (or Javascript without loops, objects and functions)
* A [shell](https://en.wikipedia.org/wiki/Unix_shell) like grammar similar to that used by command line interpreters and [TCL](https://en.wikipedia.org/wiki/Tcl)

Can write:
* YAML
* Properties
* INI file

Processing:
* Translation between grammars: for example prefix to and from infix, LISP to JSON, etc.
* A command line [REPL](ehttps://en.wikipedia.org/wiki/Read%E2%80%93eval%E2%80%93print_loop) interface
* Eval - an [evaluator](https://en.wikipedia.org/wiki/Eval) or interpetter
* Query - a utility to filter and select from the Abstract Syntax Tree, similar idea to [jq](https://stedolan.github.io/jq/)

API:
* Supports UNICODE/UTF-8
* Supports reading a writing from a text string
* Supports reading a writing from files
* Errors and warnings logs can optionally be formatted as any of the supported grammars.

Further work:
* Read and write YAML
* TODO Write XML
* TODO Read XML
* TODO read from tabular formats TSV and CSV
* TODO Extend query
* TODO Schema extraction
* TODO Read from relational database
* TODO validation
