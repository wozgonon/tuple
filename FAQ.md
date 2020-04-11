# FAQ

At the time of writing this is really anticipated questions.


## Why the name 'wozg'

Everything else I could think of seems to be taken, please suggest a better name.


## Is 'wozg' yet another verion of LISP?

No and sort of.   No it is a set of components that one can put together in various ways
one of which is a LISP like language.

LISP is actually very simple, the syntax S-Expressions, is very easy to parse and
the basic interpretter Eval is very easy to write.

What is interesting that LISP that a property called (Homoiconicity)[https://en.wikipedia.org/wiki/Homoiconicity].
Whereas a language like C or Java has it's own syntax and might store or communicate data using another syntax such as JSON or XML.
LISP just uses the say format S-Expressions as the language syntax and for data,
which make it easy to do some interesting things.

'Wozg' provides a simple S-Expression parser and a separate 'Eval' interpretter component that can be put together
to provide a LISP like language.
It is not quite LISP in that LISP is based on CONS cells which are like key/value pairs or head/tail. Lists in LISP
are made up of a series of CONS cells.
The basic wozg uses arrays and maps rather than CONS c ells, though at a syntax level it looks much the same.


## Why would I want a componentized language?

This gives you scope to experiment, if you have an idea for a particular feature you can
try it out to see how well it works.

If you don't want some language feature, perhaps for reasons of security or misuse, you can leave it out.
For instance, if you want to enforce a particular style, you can build this in as a component.

Why would you want any software with de-coupled components?

Programming languages


## What's with JSON

JSON syntax is simple: just two operators comma and colon, square brackets for arrays and braces for maps.
So it can easilly be parsed by an operator grammar.
One can interchange between S-Expressions and JSON.

It also supports an extension of JSON with expressions.


## What lexer support is provided

Currently there is just a hand written lexer that supports C language like operators.
I would like to finish off the regular expression grammar to provide a 'lex' like program
then use this as a lexer.


## What syntax is supported.

It provides an extended operator grammar, which is sufficient for S-Expressions, arithmetic, JSON, regular expressions
and DLS/C like grammar.   It provides a fair bit of flexibility without being general purpose.


## Why not just use ANTLR

Is did, and it's a wonderful tool.  It's just that I found it hard to resist the temptation to add more and more special syntax because it is easy to do
so and this becomes difficult to finish.   It doesn't really provide any help with the backend, the transition from ANTLR 3 to 4 was painful and I was tied into Java.
With 'wozg' the bones of the operator grammar is pretty much finished.

With 'wozg' the grammar is done - or the bones of the implementation at least - you can parameterize it but largely one can just focus on the functionality rather than syntax.


## What is the context/directory

Programs typically have a lot of static data at run time: global variables, meta-data, reflection, environment variables, process information, local filesystem and so on.
Typically one acceses these through a variety of different functions and APIs and one typically have to manually search the API documentation of Google to find out the details.

In 'wozg' all of this static information is presented as a searchable hierarchy like a directory or registry.
So one only has to know one 'query' command.

A consequence of this for DevSecOps is that it becomes easy to perform a security review.  The context is essentially the sandbox.


## Is it statically typed

Not currently and as a modular system, adding static types component ought to be fairly straight forward.


## Is it compiled or interpretted

Currently interpretted, I would like to add a resolved to the evaluator and then add one or more code generators, possibly to llvm.



## What's the story

I had taken some time out between contracts to study for some Certifications,
I cycled two hours to a test centre only to find it has just closed due to the Coronavirus shutdown.
So as a loss two what to do next, I decided to learn some 'Go' by knocking up an operator grammar.
From there this grew.
