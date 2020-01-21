package main

import (
	"fmt";
	"github.com/whatsthatsnail/simple_interpreter/lexer"
)

func main() {
	//scanner := lexer.NewLexer("+-(){},.;*===!!=<<=>>=")
	scanner := lexer.NewLexer(fmt.Sprintf("// this is a comment \n {()} // grouping \n !*/+-=<> <= == // operators \n \"test 2\""))
	test, err := scanner.ScanTokens()
	if !err || err {
		lexer.PrintTokens(test)
	}
}