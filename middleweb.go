//go:build js && wasm

package main

import (
	"fmt"
	"C"
	"syscall/js"
)

//export MiddletonScriptToBytecode
func MiddletonScriptToBytecode(source string) {
	mt := MiddletonInterpreter{}
	return mt.ToBytecode(source)
}

// Nothing to see here.
func main() {
	fmt.Println(".")
}
