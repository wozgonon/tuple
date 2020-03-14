## The Faceless Lang

The faceless language is a little programming system that one might like to use
for embedded languages, domain specific languages, expressions data and configuration,
configuration.   Querying and converting JSON and YAML.

The name faceless lang comes from the Game of Thrones faceless men
with with various 

## Commonalities

Observing that JSON (and YAML) is now 2020 used as the common configuration languages and consists of nothing but
scalars, arrays (fixed length list) and maps.   Yet LISP from the 1950's supported the same concepts: scalars, lists and maps.

Of course one thinks of LISP as an executable programming language and JSON just as a file format LISP always was a file format
with an eval function.  One could add an eval function to JSON (which one can with faceless).

## Conversions

To tabular format for excel ...


# Little languages parsers

Parsers for little languages, Written in golang:

* Reverse Polish Notation as used by Forth and Postscript
* Expression similar to that provided by EXCEL
* Lisp - a very regular syntax that works equally well for data and code.
* TCL - for code and data
* JSON - for data
* jml - for data, code and markup
* prolog - code, data and expressions
* occam/python/yaml - indented

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

