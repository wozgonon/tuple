A simple web server that accepts requests.

Each request has the Content-type header set to one of the grammars provided by the 'tuple' framework.
The body of the request is parsed using the given grammar and the results of 'eval' are returned to the client.


# Example

## Server

```
$ whd
``

## Client

The body is just a simple sum: 1+2+3

```
$ curl localhost:8888/do -XPOST -d"1+2+3" -HContent-type:expr
{[{+} {[{+} %!s(tuple.Int64=1) %!s(tuple.Int64=2)]} %!s(tuple.Int64=3)]}
```
