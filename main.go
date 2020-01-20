package main

import (
	//"fmt";
	"github.com/whatsthatsnail/simple_interpreter/lexer"
)

func main() {
	scanner := lexer.NewLexer("+-(){},.;*===!!=<<=>>=")
	test := scanner.ScanTokens()
	lexer.PrintTokens(test)
}