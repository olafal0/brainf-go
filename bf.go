package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	var err error
	var command []byte
	reader := bufio.NewReader(os.Stdin)
	memsize := 30000

	// construct a state for the interpreter
	bfinterpreter := NewInterpreter(memsize)

	if len(os.Args) > 2 {
		command, err = ioutil.ReadFile(os.Args[2])
		if err != nil {
			fmt.Printf("File %s not found!\n", os.Args[2])
			return
		}
		if os.Args[1] == "make" {
			compilego(command, "compiled/compiled.bf.go")
		} else if os.Args[1] == "cc" {
			compilecc(command, "compiled/compiled.bf.c")
		} else if os.Args[1] == "fastcc" {
			compileccfast(command, "compiled/compiled.bf.c")
		} else if os.Args[1] == "run" {
			bfinterpreter.run(command, reader)
		}
	} else {
		// run as live interpreter
		fmt.Print("Running in REPL mode (persistent state)\n")
		for {
			fmt.Print("# ")
			command, err = reader.ReadBytes('\n')
			if err != nil {
				return
			}
			bfinterpreter.run(command, reader)
		}
	}
}
