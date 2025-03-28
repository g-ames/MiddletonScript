package main

import (
	"fmt"
	"log"
	"slices"
	"os"
)

func repl() {
	fmt.Println("NOTE: the MiddletonScript REPL is experimental!")
	running := true

	var input string
	middleton := MiddletonInterpreter{}

	for running {
		fmt.Print(">> ")
		fmt.Scan(&input)
		middlebytes := middleton.ToBytecode(input)
		fmt.Println(middlebytes)
	}
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("MiddletonScript ready! Interpreting a 'Hello World' example.")
		middleton := MiddletonInterpreter{}
		middlebytes := middleton.ToBytecode(`
		swyk middleton

		notes (
			"middleton:fmt"
		)

		MFunc main() {
    			middleton::mout("Hello, MiddletonScript!");
		}
		`)
		fmt.Println(middlebytes)
		return
	}

	var flags []rune

	// Flags pass
	for _, arg := range args {
		if arg == "--repl" {
			repl()
			return
		}
	
		if arg[0] == '-' {
			flags = append(flags, rune(arg[1]))
			continue
		}
	}

	if slices.Contains(flags, 'd') {
		dirargs := []string{}
		
		for _, arg := range args {
			if arg[0] == '-' { continue }
			
			entries, err := os.ReadDir(arg)
			if err != nil {
				log.Fatal(err)
			}
			for _, entry := range entries {
				if entry.IsDir() { continue }
				dirargs = append(dirargs, arg + "/" + entry.Name())
			}
		}
		
		args = dirargs
	}

	for _, arg := range args {
		// Ignore flags as we have already parsed them
		if arg[0] == '-' { continue }

		fmt.Println("#", arg)
	
		data, err := os.ReadFile(arg)

		if err != nil {
			fmt.Printf("Error reading file with name '%s'!\n", arg)
			continue
		}

		middleton := MiddletonInterpreter{}
		middlebytes := middleton.ToBytecode(string(data))

		if slices.Contains(flags, 'b') {
			file, err := os.Create(arg + ".middlebytes")
			if err != nil {
				log.Fatal(fmt.Errorf("Error opening MiddletonScript bytes file! %s", err))
			}

			_, err = file.Write(middlebytes)
			if err != nil {
				log.Fatal(fmt.Errorf("Error writing to MiddletonScript bytes file! %s", err))
			}

			file.Close()
		}
	}
}
