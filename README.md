## The Faceless Lang

We have to work with few standard and bespoke syntaxes for representing data, code and markup.
We often want to create, query, validate or transform these files.
The data typically comes down to just strings, numbers, records or lists and key/value sets.

Often we need to customize them, sometimes a bit or sometimes  to create entire Domain Specific Languages (DSL)s.
The problem is we have a disperate set of tools for working with them,
general purpose scripts, APIs and tools for particular formats.

This tool provides a solution to all these issues by.


convert, query or process them and customize these.

## Macros


The faceless language is a little programming system that one might like to use
for embedded languages, domain specific languages, expressions data and configuration,
configuration.   Querying and converting JSON and YAML.

The name faceless lang comes from the Game of Thrones faceless men
with with various 


* Conversion
* COde generators
* DSL
* Pretty printing
* Query
* Sccripts

S-Expression

## The Computer Science

A [Homoiconic](https://en.wikipedia.org/wiki/Homoiconicity) language treats "code as data".
All code can be accessed and treated as if it is data,

Observing that JSON (and YAML) is now 2020 used as the common configuration languages and consists of nothing but
nested scalars, arrays (fixed length list) and maps.   Yet LISP from the 1950's supported the same concepts: nested scalars, lists and maps.

Of course one thinks of LISP as an executable programming language and JSON just as a file format LISP always was a file format
with an eval function.  One could add an eval function to JSON (which one can with faceless).

## Conversions

To tabular format for excel ...


# A history of Little languages parsers

Parsers for little languages, Written in golang:

* Reverse Polish Notation as used by Forth and Postscript
* Expression similar to that provided by EXCEL
* Lisp - a very regular syntax that works equally well for data and code.
* TCL - for code and data
* JSON - for data
* jml - for data, code and markup
* prolog - code, data and expressions
* occam/python/yaml - indented

My experience is that it is very easy to write a Lisp parser and very hard to write anything more complex because there are so many special cases.
It would be nice not to have to put quotes around strings, in particular working interactice on the command line (CLI) with a shell such as bash or
dos command line.  It is very nice not to have to include any boiler plate when just entering commands but this quickly becomes very awkward to add any more complex syntax.

```
cd cmd/jml
go run main.go in.jml
```


```
export GOBIN=`pwd`/bin
go env -w GOBIN=`pwd`/bin
go install src/jml/jml.go


go install src/tuple/tuple.go
go install src/lisp/lisp.go

```

## Personal preferences

Syntax preferences can be very subjective.
if you have ever seen programmers get into religious battles over where to put spaces and brackets.

Place a githook to format into project standard format on commit and into personal preference.


TODO Comments
TODO Errors in any format - eat own shit
TODO github git remote - so in code commit and github
