package main

import (
	"fmt"
	"bufio";
	"os";
	"io/ioutil";
	"friston/lexer";
	"friston/type_generator";
	"friston/parser";
	"friston/ast";
)

// Gets arguments when using 'go run *.go -- ...'
func main() {

	var args []string
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}

	if len(args) >= 1 && args[0] == "repl" {
		repl()
	} else if len(args) >= 3 && args[0] == "file" && args[2] == "-v" {
		file(args[1], false)
	} else if len(args) >= 2 && args[0] == "file" {
		file(args[1], true)
	} else if len(args) >= 2 && args[0] == "GenASTSource" {
		genASTSource(args[1])
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
	fmt.Printf("Entering REPL:\n>>> ")

	scanner := bufio.NewScanner(os.Stdin)

	interpreter := ast.NewInterpreter(true)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "exit" {
			os.Exit(0)
		}

		lex := lexer.NewLexer(line, true)
		tokens, lexErr := lex.ScanTokens()

		if !lexErr {
			par := parser.NewParser(tokens)
			stmts, parErr := par.Parse()

			if !parErr {
				interpreter.Interpret(stmts)
				fmt.Printf(">>> ")
			}
		}
	}
}

// Reads file into lexer, tokenizes, and prints tokens
func file(path string, quiet bool) {
	dat, err := ioutil.ReadFile(path)
	check(err)

	if !quiet {
		fmt.Println(path + ":" + "\n")
		fmt.Println(string(dat) + "\n")
	}

	lex := lexer.NewLexer(string(dat), false)
	tokens, lexErr := lex.ScanTokens()

	if !lexErr{
		if !quiet {
			lexer.PrintTokens(tokens)
		}

		par := parser.NewParser(tokens)
		stmts, parErr := par.Parse()

		if !quiet && !parErr {
			printer := ast.ASTPrinter{}
			fmt.Printf("\n")
			for _, s := range(stmts) {
				s.Accept(printer)
			}
			fmt.Printf("\n")
		}

		if !parErr {
			interpreter := ast.NewInterpreter(false)
			interpreter.Interpret(stmts)
		}
	}
}

func genASTSource(path string) {
	dat, err := ioutil.ReadFile(path)
	check(err)

	scanner := lexer.NewLexer(string(dat), false)
	tokens, errFlag := scanner.ScanTokens()

	if !errFlag{
		type_generator.GenerateNodeTypes(tokens)
	}
}