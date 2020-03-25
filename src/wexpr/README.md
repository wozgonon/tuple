A simple shell utility for evaluating [expressions](https://en.wikipedia.org/wiki/Expression_(mathematics))
similar to the UNIX [expr](https://en.wikipedia.org/wiki/Expr) utility.

For example:

```
$ wexpr 1+1
2
```

```
$ wexpr "acos(cos(PI))"
3.141592653589793
```


```
$ wexpr -- "-1-2"
-3
```

The '-ast' option can be used to print the [AST](https://en.wikipedia.org/wiki/Abstract_syntax_tree)
before evaluation.

```
$ wexpr -ast -- "-1+2*3^4"
(
  (
    -
    1
  )
  +
  (
    2
    *
    (
      3
      ^
      4
    )
  )
)
```



