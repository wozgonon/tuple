## Tuple

This package provides an operator grammer that can be tweaked to 
create simple parsers and Domain Specific Languages (DSL) without having the overhead of writing ones own parser or learning to use a [compiler-compiler](https://en.wikipedia.org/wiki/Compiler-compiler).

The package provides functions for parsing, printing and processing a number of simple Grammars.
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
$ git clone https://github.com/wozgonon/tuple.git
$ cd tuple
$ export GOPATH=`pwd` 
$ export GOBIN=`pwd`/bin
$ go env -w GOBIN=`pwd`/bin
$ make
```

To build and run tests:
```
$ make test
go install src/jml/jml.go
```


# FAQ


## Yet another verion of LISP?

Not quite, it is a set of components that one can put together in various ways one of which is a LISP like language.

LISP syntax is very simple, the syntax S-Expressions, is very easy to parse and the basic interpretter Eval is very easy to write.

LISP has an interesting property called [Homoiconicity](https://en.wikipedia.org/wiki/Homoiconicity).
Whereas a language like C or Java has it's own syntax and might store or communicate data using another syntax such as JSON or XML.
LISP just uses the say format S-Expressions as the language syntax and for data, which makes interesting things easy.

This package provides a simple S-Expression parser and a separate 'Eval' interpretter component that can be put together
to provide a LISP like language.
It is not quite LISP in that LISP is based on CONS cells which are like key/value pairs or head/tail. Lists in LISP
are made up of a series of CONS cells.
This package uses arrays and maps rather than CONS cells, though the syntax looks much the same.

## Support for JSON

JSON syntax is simple: just two operators comma and colon, square brackets for arrays and braces for maps.
So it can easilly be parsed by an operator grammar.
One can interchange between S-Expressions and JSON.

This package also supports an extension of JSON with expressions.


## Why would I want a componentized language?

Gives one scope to experiment, if you have an idea for a particular feature you can try it out to see how well it works.
If you don't want some language feature, perhaps for reasons of security or misuse, you can leave it out.
For instance, if you want to enforce a particular style, you can build it in as a component.
Really the same reasons you would want any software with de-coupled components?


## What lexer support is provided

Currently there is just a hand written lexer that supports C language like operators.
I would like to finish off the regular expression grammar to provide a 'lex' like program then use this as a lexer.


## What syntax is supported.

It provides an extended operator grammar, which is sufficient for S-Expressions, arithmetic, JSON, regular expressions
and DLS/C like grammar.   It provides a fair bit of flexibility without being general purpose.


## Why not just use ANTLR or another grammar generator

With general purpose [compiler-compilers](https://en.wikipedia.org/wiki/Compiler-compiler) like ANTLR or BISON/YACC one can find it hard to resist the temptation to add more and more special syntax
although this does not help you complete the backend.
This package provides and operator grammer that is pretty much finished with some flexibility to tweak it.
So you can focus on the functionality rather than on the syntax.

## A common Abstract Syntax Tree

In practice many computer syntaxes/grammars represent the same things: strings, numbers, records or lists and key/value sets.
We often want to create, query, validate or transform files.
If we can transform these to a common form, an Abstract Syntax Tree (AST),
we can in principle provide just one set of tools for manipulating them rather than having to use bespoke tools for each.

## What is the context/directory

Programs typically have a lot of static data at run time: global variables, meta-data, reflection, environment variables, process information, local filesystem and so on.
Typically one accesses these through a variety of functions and APIs which can be awkward to find in the documentation.

This package presents all the static information as a searchable hierarchy like a directory or registry, with a single 'query' command.
A consequence of this for DevSecOps is that it becomes easy to perform a security review.  The context is essentially the sandbox.


## Is it statically typed

Not currently and as a modular system, adding static types component ought to be fairly straight forward.


## Is it compiled or interpretted

Currently interpretted, I would like to add a resolver to the evaluator and then add one or more code generators, possibly to llvm.


## What's the story

I took some time out between contracts to study for some Certifications when
I cycled two hours to a test centre, I to find it has just closed due to the Coronavirus shutdown.  So at a loss as to what to do next, I decided to improve my 'Go' by knocking up an operator grammar.  From there this grew.
