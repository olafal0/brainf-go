package main

import (
	"fmt"
	"os"
)

func compilego(command []byte, filename string) {
	outf, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outf.Close()

	outf.WriteString("package main\nfunc main () {\n")
	outf.WriteString("mem := make([]int,30000)\nmemptr := 0\nreader := bufio.NewReader(os.Stdin)\n")
	for i := 0; i < len(command); i++ {
		switch command[i] {
		case '<':
			outf.WriteString("memptr--\n")
		case '>':
			outf.WriteString("memptr++\n")
		case '+':
			outf.WriteString("mem[memptr]++\n")
		case '-':
			outf.WriteString("mem[memptr]--\n")
		case '.':
			outf.WriteString("fmt.Printf(\"%c\",mem[memptr])\n")
		case ',':
			outf.WriteString(`input, err := reader.ReadByte()
			if err != nil {
				mem[memptr] = 0
			} else {
				mem[memptr] = int(input)
			}
			`)
		case '[':
			outf.WriteString("for mem[memptr] != 0 {\n")
		case ']':
			outf.WriteString("}\n")
		}
	}
	outf.WriteString("\n}")
}

func compilecc(command []byte, filename string) {
	outf, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outf.Close()

	outf.WriteString("#include <stdio.h>\n#include <stdlib.h>\nint main () {\n")
	outf.WriteString("int *mem = (int*)malloc(sizeof(int)*30000);\nint memptr = 0;\n")
	for i := 0; i < len(command); i++ {
		switch command[i] {
		case '<':
			outf.WriteString("memptr--;\n")
		case '>':
			outf.WriteString("memptr++;\n")
		case '+':
			outf.WriteString("mem[memptr]++;\n")
		case '-':
			outf.WriteString("mem[memptr]--;\n")
		case '.':
			outf.WriteString("putchar(mem[memptr]);\n")
		case ',':
			outf.WriteString("mem[memptr] = getchar();\n")
		case '[':
			outf.WriteString("while (mem[memptr]) {\n")
		case ']':
			outf.WriteString("}\n")
		}
	}
	outf.WriteString("\nfree(mem);\nreturn 0;\n}")
}

func compileccfast(command []byte, filename string) {
	outf, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outf.Close()

	// create state to collapse multiple add or shifts to a single instruction
	lShift := 0
	rShift := 0
	madd := 0
	msub := 0

	outf.WriteString("#include <stdio.h>\n#include <stdlib.h>\nint main () {\n")
	outf.WriteString("int *mem = (int*)malloc(sizeof(int)*30000);\nint memptr = 0;\n")
	for i := 0; i < len(command); i++ {
		switch command[i] {
		case '<':
			if (rShift | madd | msub) != 0 {
				// switch states
				outf.WriteString(getStateChangeString(lShift, rShift, madd, msub))
				rShift = 0
				madd = 0
				msub = 0
			}
			lShift++
		case '>':
			if (lShift | madd | msub) != 0 {
				// switch states
				outf.WriteString(getStateChangeString(lShift, rShift, madd, msub))
				lShift = 0
				madd = 0
				msub = 0
			}
			rShift++
		case '+':
			if (lShift | rShift | msub) != 0 {
				// switch states
				outf.WriteString(getStateChangeString(lShift, rShift, madd, msub))
				lShift = 0
				rShift = 0
				msub = 0
			}
			madd++
		case '-':
			if (lShift | rShift | madd) != 0 {
				// switch states
				outf.WriteString(getStateChangeString(lShift, rShift, madd, msub))
				lShift = 0
				rShift = 0
				madd = 0
			}
			msub++
		case '.':
			if (lShift | rShift | madd | msub) != 0 {
				// switch states
				outf.WriteString(getStateChangeString(lShift, rShift, madd, msub))
				lShift = 0
				rShift = 0
				madd = 0
				msub = 0
			}
			outf.WriteString("putchar(mem[memptr]);\n")
		case ',':
			if (lShift | rShift | madd | msub) != 0 {
				// switch states
				outf.WriteString(getStateChangeString(lShift, rShift, madd, msub))
				lShift = 0
				rShift = 0
				madd = 0
				msub = 0
			}
			outf.WriteString("mem[memptr] = getchar();\n")
		case '[':
			if (lShift | rShift | madd | msub) != 0 {
				// switch states
				outf.WriteString(getStateChangeString(lShift, rShift, madd, msub))
				lShift = 0
				rShift = 0
				madd = 0
				msub = 0
			}
			outf.WriteString("while (mem[memptr]) {\n")
		case ']':
			if (lShift | rShift | madd | msub) != 0 {
				// switch states
				outf.WriteString(getStateChangeString(lShift, rShift, madd, msub))
				lShift = 0
				rShift = 0
				madd = 0
				msub = 0
			}
			outf.WriteString("}\n")
		}
	}
	outf.WriteString("\nfree(mem);\nreturn 0;\n}")
}

func getStateChangeString(lshif, rshif, madd, msub int) string {
	if rshif != 0 {
		return fmt.Sprintf("memptr += %d;\n", rshif)
	}
	if lshif != 0 {
		return fmt.Sprintf("memptr -= %d;\n", lshif)
	}
	if madd != 0 {
		return fmt.Sprintf("mem[memptr] += %d;\n", madd)
	}
	if msub != 0 {
		return fmt.Sprintf("mem[memptr] -= %d;\n", msub)
	}
	fmt.Println("Error! No state found")
	return ""
}
