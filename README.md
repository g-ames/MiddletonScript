# MiddletonScript

## Overview

MiddletonScript is a fast, multi-paradigm efficient interpreted language.

Here is an example, hello.middle:
```middle
swyk middleton

notes (
    "middleton:fmt"
)

MFunc main() {
    middleton::mout("Hello, MiddletonScript!");
}
```
It compiles to bytecode for efficient operation and can even be used in the browser when compiled to wasm!

## How to use?

MiddletonScript can be run very easily.

```sh
middle hello.middle
```

You can also output directly as a MiddletonScript bytecode file (.middlebytes):
```sh
middle -b hello.middle
```

or, to be verbose:

```sh
middle --bytes hello.middle
```

## How to make notes in MiddletonScript?

To make notes in MiddletonScript, you first need to create a .ltc file. 
The .ltc file is similar to go's 'go.mod' or node.js's 'package.json'.
How do you create an ltc?

```ltc
make hello.middle (are ya winning son?) 
```

This .ltc file, which we will call middle.ltc, is very simple.
For example, 
```sh
middle -b hello.middle
```
becomes
```ltc
middle -b hello.middle (are ya winning son?)
```

## Advanced ltc concepts

The '(are ya winning son?)' option allows you to error handle in a .ltc file.
```ltc
middle hello.middle (are ya winning son?)
```
The 'goget' option allows you to 'goget' the latest version of MiddletonScript notes.
```ltc
goget fmt (are ya winning son?)
```
