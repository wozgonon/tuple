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

