package main

import (
	//"fmt";
	"github.com/whatsthatsnail/simple_interpreter/lexer"
)

func main() {
	scanner := lexer.NewLexer("+-(){},.;*===!!=<<=>>=")
	test, err := scanner.ScanTokens()
	if !err {
		lexer.PrintTokens(test)
	}
}