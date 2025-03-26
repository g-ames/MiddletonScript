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

You can also output directly as a MiddletonScript bytecode file:
```sh
middle -b hello.middle
```

or, to be verbose:

```sh
middle --bytes hello.middle
```
