package main

import (
	"fmt"
	"bufio";
	"os";
	"io/ioutil";
	"github.com/whatsthatsnail/simple_interpreter/lexer";
	"github.com/whatsthatsnail/simple_interpreter/type_generator";
	"github.com/whatsthatsnail/simple_interpreter/parser";
	"github.com/whatsthatsnail/simple_interpreter/ast";
)

// Gets arguments when using 'go run *.go -- ...'
func main() {
	
	var args []string
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}

	if len(args) >= 1 && args[0] == "repl" {
		repl()
	} else if len(args) >= 2 && args[0] == "file" {
		file(args[1], false)
	} else if len(args) >= 2 && args[0] == "GenASTSource" {
		genASTSource(args[1])
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
	fmt.Printf("Entering REPL:\n>> ")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		
		if line == "exit" {
			os.Exit(0)
		}
		
		lex := lexer.NewLexer(line)
		tokens, errFlag := lex.ScanTokens()
		
		if !errFlag {
			par := parser.NewParser(tokens)
			expr := par.Parse()

			interpreter := ast.Interpreter{}
			out := expr.Accept(interpreter)
			if out != nil { 
				fmt.Printf("%v\n>> ", out)
			}
		}
	}
}

// Reads file into lexer, tokenizes, and prints tokens
func file(path string, quiet bool) {
	dat, err := ioutil.ReadFile(path)
	check(err)

	fmt.Println(path + ":" + "\n")
	if !quiet {
		fmt.Println(string(dat) + "\n")
	}

	lex := lexer.NewLexer(string(dat))
	tokens, errFlag := lex.ScanTokens()

	if !errFlag{
		lexer.PrintTokens(tokens)
		
		par := parser.NewParser(tokens)
		expr := par.Parse()
		
		printer := ast.ASTPrinter{}
		fmt.Printf("\n")
		expr.Accept(printer)
		fmt.Printf("\n\n")

		interpreter := ast.Interpreter{}
		out := expr.Accept(interpreter)
		fmt.Printf("Output: %v\n", out)
	}
}

func genASTSource(path string) {
	dat, err := ioutil.ReadFile(path)
	check(err)

	scanner := lexer.NewLexer(string(dat))
	tokens, errFlag := scanner.ScanTokens()

	if !errFlag{
		type_generator.GenerateNodeTypes(tokens)
	}
}