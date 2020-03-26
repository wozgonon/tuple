The tuple package provides functions for parsing, printing and processing a number of simple Grammars.
The parsers produce a common [parse tree](https://en.wikipedia.org/wiki/Parse_tree)
and [AST](https://en.wikipedia.org/wiki/Abstract_syntax_tree).


# Read, write and processing

Can read, write and convert between:

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

## Processing:

* Translation between grammars: for example prefix to and from infix, LISP to JSON, etc.
* A command line [REPL](ehttps://en.wikipedia.org/wiki/Read%E2%80%93eval%E2%80%93print_loop) interface
* Eval - an [evaluator](https://en.wikipedia.org/wiki/Eval) or interpetter
* Query - a utility to filter and select from the Abstract Syntax Tree, similar idea to [jq](https://stedolan.github.io/jq/)

# API:

* Supports UNICODE/UTF-8
* Supports reading a writing from a text string
* Supports reading a writing from files
* Errors and warnings logs can optionally be formatted as any of the supported grammars.

# Embdedded processing

To embed within your own programs add your own functions the the evaluator.

* Supports expressions parsing and evaluation: e.g. 1+2, cos(PI) etc
* Embedded scripting

```
... give example
```

# Other User Cases

## Pretty printing

The 'Printer' interface is provided for pretty printing.
This does a depth first traversal of the AST and writes to a text stream.

One can provide one's own pretty printer producing HTML or Markup for example.

## Code generator

A simple example of when you might want to write a code generator is when
you want to automatically make sure your  OpenAPI objects defined in YAML
to have the same fields and types as those in your SQL database schema.
Specify the types and fields in a DSL and then write generators
YAML and SQL DML.  In this case the generators would work much like
pretty printers, doing a depth first traversal of the AST and writing to a file.

## [DSL](https://en.wikipedia.org/wiki/Domain-specific_language) and min-languages

The package can be used for writing 'mini-languages',
it parses and evaluates expressions while you implement just the functionality you need.

* Import the tuple package
* Add your own functions to the evaluator

You can use any or all of the provided grammars, these are interchangable.
The LISP like grammars provide regular syntax that works equally well for data and code,
the 'Shell' grammar is similar gradle, bash or TCL.

If you need to tweak the grammars please do so and if you really need something complex then use
a general purpose language recognition tool like [Yacc/Bison](https://en.wikipedia.org/wiki/GNU_Bison) or [ANTLR](https://en.wikipedia.org/wiki/ANTLR).

However,the grammars provided are sufficient:

* One will not get more functionality from a complex grammar.
* JSON (and YAML) are the mordern configuration languages and consists of nested scalars, arrays (fixed length list) and maps.
  Yet LISP from the 1950's supported the same concepts: nested scalars, lists and maps.
* You can write complete programming languages with trivial syntax for instance Forth and Postscript use reverse polish notation and LISP syntax uses S-Expressions.



The basic grammars provided are variations on S-Expressions and operator precedence grammars.
A basic lexer is recognises C-Language like tokens

### Security

A DSL written in Groovy or Scala typically has access to functionality you don't necessarilly need.
With 'tuple' you specify exactly what functions are available so there is no risk of any backdoors.

## Teaching

The package might be useful to those interested in understanding or in teaching others how compilers and interpreters work.
It provides implementations of each of the stages (lex, parse, resolve, eval, generate)  which you can tweak or completely replace with your own.


# Further work:

* Read and write YAML
* TODO Write XML
* TODO Read XML
* TODO Extend query
* TODO Schema extraction
* TODO Read from relational database
* TODO read from tabular formats TSV and CSV
* TODO write TSV,CSV, SQL Insert statements
* TODO validation
* TODO Querying and converting JSON and YAML.
* TODO Macro language
* TODO Markup grammar
* TODO Can pretty print HTML etc
* TODO Comments
* TODO generalize eval
* TODO eval to llvm and or JVM
* TODO SQL select
* TODO Recognise Reverse Polish Notation, aka postfix (postscript, forth)
* TODO Recognise indented syntax (occam/python/yaml)
* TODO basic prolog like eval (operator grammar syntax postfix operator)
* TODO cloudformation like expressions in strings in json:   abc "${,,,}"
* TODO resolver