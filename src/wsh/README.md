A simple command shell based on the 'tuple' framework.

Expressions work as normal:

```
$ 1+2
3
>$ cos(PI)
-1
```

Any unrecognised function is treated as a call to 'exec' which runs an executable process

For example to list files in the current directory
```
$ pwd
/
$ ls
/bin
...
```

```
$ exec "ls"
...
```
