## Wogz/Tuple

We have to work with few standard and bespoke syntaxes for representing data, code and markup.
We often want to create, query, validate or transform these files.
The data typically comes down to just strings, numbers, records or lists and key/value sets.

Often we need to customize them, sometimes a bit or sometimes  to create entire Domain Specific Languages (DSL)s.
The problem is we have a disperate set of tools for working with them,
general purpose scripts, APIs and tools for particular formats.

The tuple package provides functions for parsing, printing and processing a number of simple Grammars.
The parsers produce a common [parse tree](https://en.wikipedia.org/wiki/Parse_tree)
and [AST](https://en.wikipedia.org/wiki/Abstract_syntax_tree).

For:
* Conversion
* Query
* Embdedded expression processing and embedded scripting
* Code generators and pretty printing
* The package can be used for writing 'mini-languages' and [DSL](https://en.wikipedia.org/wiki/Domain-specific_language)


To build:
```
$ export GOBIN=`pwd`/bin
$ go env -w GOBIN=`pwd`/bin
$ make
```

To build and run tests:
```
$ make test
go install src/jml/jml.go
```
