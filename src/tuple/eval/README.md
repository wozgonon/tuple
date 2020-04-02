This package provides a simple (Eval function)[https://en.wikipedia.org/wiki/Eval] or the backend of an interpretter.

# De-coupled from parser

The idea is that you do no have to use this implementation but can use write your own interpretter or
compiler for use with one of the grammars provided in the parser package
or write your own parser and use it with this package.

# Sandboxing

The package supports (Sandboxing)[https://en.wikipedia.org/wiki/Sandbox_(computer_security)].

## Provided functions

The package provides a number of functions by default

* harmless functions - just do basic arithmetic and compare, run in determinate time and do not allocate any resources.
* safe functions - allocate memory or potentially recurse and so could potentially though perhaps unlikely cause a process memory.
* less safe functions - access external resources such as the operating environment which could potentially be used for nefarious purposes.

Risks:
* One could include 'harmless functions' in a configuration file sent between client and server or between uServices without any risks at all.
* One could include 'safe functions' in a file sent between processes with the risk that an operation might cause an memory outage so only behind a firewall.
* Include 'less safe functions' in a file sent between processes might potentially be used to access remote resources, so best server (or client) side only.

y
##  DSL and DevSecOps

In terms of security one can see this as DENY by default, one can only use a function if explicitly ALLOWED.
In Groovy and I believe in Scala a DSL has access to the entire environment, one can disable features at the JVM level
but this is ALLOW by default and less safe than DENY by default.
