The 'wogz' utility that presents the functionality of the 'tuple' package on the command line.

It can be used to parse and process a number of common grammars.


...

## Recipes

###  Convert from prefix to infix notation

To convert between LISP in prefix notation (+ 1 2) to LISP with infix notation (1 + 2), use:

```
$ wogz -out infix prefix.l > prefix.infix
$ wogz -out l prefix.infix > prefix.l
```

### Query

```
$wogz -query a
...
```

### List supported Grammars

```
$ wozg -list-grammars
```
