class MiddletonInterpreter {
    constructor() {
        this.middletonianState = {}; 
    }

    toBytecode(source) {
        let tokens = lex(source);
        let parsed = parse(tokens);
        console.log(tokens, parsed);
    }
}

let middleton = new MiddletonInterpreter();
middleton.toBytecode(`package middleton

import (
    "middleton:fmt"
)

MFunc main() {
    middleton::mout("Hello, MiddletonScript!");
}`);