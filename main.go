//go:build !(js || wasm)

package main

import (
	"fmt"
	"log"
	"slices"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func __update__() error {
	fmt.Println("Running on", runtime.GOOS)
	if !(runtime.GOOS == "linux" || runtime.GOOS == "darwin") {
		return fmt.Errorf("this function only works on Linux")
	}

	updateDir := ".middleupdate"

	if err := os.MkdirAll(updateDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", updateDir, err)
	}

	repoURL := "https://github.com/g-ames/MiddletonScript.git"
	cloneCmd := exec.Command("git", "clone", "--depth", "1", repoURL, filepath.Join(updateDir, "MiddletonScript"))
	cloneCmd.Dir = updateDir // Set the working directory to the .middleupdate folder

	if err := cloneCmd.Run(); err != nil {
		return fmt.Errorf("failed to clone the repo: %v", err)
	}

	buildScriptPath := filepath.Join(updateDir, "MiddletonScript", "build.sh")
	if _, err := os.Stat(buildScriptPath); os.IsNotExist(err) {
		return fmt.Errorf("build.sh not found in the repository")
	}

	chmodCmd := exec.Command("chmod", "+x", buildScriptPath)
	if err := chmodCmd.Run(); err != nil {
		return fmt.Errorf("failed to make build.sh executable: %v", err)
	}

	buildCmd := exec.Command("bash", buildScriptPath)
	buildCmd.Dir = filepath.Join(updateDir, "MiddletonScript") // Set the working directory to the repo folder

	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("failed to run build.sh: %v", err)
	}

	fmt.Println("Build completed successfully.")
	return nil
}

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

		if arg == "--update" {
			err := __update__()
			if err != nil {
				fmt.Println("Update unsuccessful.")
				log.Fatal(err)
			}
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
