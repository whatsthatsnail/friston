package main

import (
	"fmt";
	"os";
	"io/ioutil";
	"github.com/whatsthatsnail/simple_interpreter/lexer"
)

// Gets arguments when using 'go run *.go -- ...'
func main() {
	args := os.Args[2:]

	if len(args) >= 1 && args[0] == "repl" {
		repl()
	} else if len(args) >= 2 && args[0] == "file" {
		file(args[1], false)
	} else if len(args) >= 3 && args[0] == "file" && args[2] == "-q" {
		file(args[1], true)
	} else {
		repl()
	}
}

// Helper function to check for errors when reading files
func check(err error) {
	if err != nil {
		panic(err)
	}
}

// TODO: Implement a REPL
func repl() {
	fmt.Println("Entering REPL:")
}

// Reads file into lexer, tokenizes, and prints tokens
func file(path string, quiet bool) {
	dat, err := ioutil.ReadFile(path)
	check(err)

	fmt.Println(path + ":" + "\n")
	if !quiet {
		fmt.Println(string(dat) + "\n")
	}

	scanner := lexer.NewLexer(string(dat))
	tokens, errflag := scanner.ScanTokens()

	if !errflag{
		lexer.PrintTokens(tokens)
	}
}