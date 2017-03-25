package main

import (
	"bufio"
	"io"
	"os"
)

// Interpreter represents a brainf interpreter state
type Interpreter struct {
	mem          [30000]byte
	memptr       int
	bracketLevel int
	idx          int
	reader       io.Reader
	writer       io.Writer
	buf          []byte
}

// NewInterpreter constructs a new interpreter state
func NewInterpreter(memorySize int) *Interpreter {
	i := Interpreter{
		buf: make([]byte, 1),
	}
	return &i
}

func (bf *Interpreter) run(command []byte, reader *bufio.Reader) {
	bf.reader = reader
	bf.writer = io.Writer(os.Stdout)
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
			if bf.memptr >= 30000 {
				bf.memptr = 29999
			}
		case '+':
			bf.mem[bf.memptr]++
		case '-':
			bf.mem[bf.memptr]--
		case '.':
			bf.buf[0] = bf.mem[bf.memptr]
			bf.writer.Write(bf.buf)
		case ',':
			input, err := reader.ReadByte()
			if err != nil {
				bf.mem[bf.memptr] = 0
			} else {
				bf.mem[bf.memptr] = input
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
