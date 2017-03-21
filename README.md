# brainf-go
A simple brainf*** interpreter, written in Go.

# Usage
Use `go build` to build or `go run bf.go` to simply run it. If the first command-line argument is a filename,
that file will be interpreted as brainf*** code. Otherwise, an interpreter state will be created, and you can write
brainf*** directly (as long as there are no mismatched square brackets on a single line).
