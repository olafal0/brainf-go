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

	if len(os.Args) > 1 {
		command, err = ioutil.ReadFile(os.Args[1])
		if err != nil {
			fmt.Printf("File %s not found!\n", os.Args[1])
			return
		}
		bfinterpreter.run(command, reader)
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

// Interpreter represents a brainf interpreter state
type Interpreter struct {
	mem          []int
	memptr       int
	memsize      int
	bracketLevel int
	idx          int
}

// NewInterpreter constructs a new interpreter state
func NewInterpreter(memorySize int) *Interpreter {
	i := Interpreter{
		mem:     make([]int, memorySize),
		memsize: memorySize,
	}
	return &i
}

func (bf *Interpreter) run(command []byte, reader *bufio.Reader) {
	var bracketLevel, i int
	for i < len(command) {
		switch command[i] {
		case '<':
			bf.memptr--
			if bf.memptr < 0 {
				bf.memptr = 0
			}
		case '>':
			bf.memptr++
			if bf.memptr >= bf.memsize {
				bf.memptr = bf.memsize - 1
			}
		case '+':
			bf.mem[bf.memptr]++
		case '-':
			bf.mem[bf.memptr]--
		case '.':
			fmt.Printf("%c", bf.mem[bf.memptr])
		case ',':
			input, err := reader.ReadByte()
			if err != nil {
				bf.mem[bf.memptr] = 0
			} else {
				bf.mem[bf.memptr] = int(input)
			}
		case '[':
			if bf.mem[bf.memptr] == 0 {
				bracketLevel = 1
				for bracketLevel != 0 {
					i++
					if command[i] == '[' {
						bracketLevel++
					} else if command[i] == ']' {
						bracketLevel--
					}
				}
			}
			bracketLevel++
		case ']':
			if bf.mem[bf.memptr] != 0 {
				bracketLevel = -1
				for bracketLevel != 0 {
					i--
					if command[i] == '[' {
						bracketLevel++
					} else if command[i] == ']' {
						bracketLevel--
					}
				}
			}
		}
		i++
	}
}
