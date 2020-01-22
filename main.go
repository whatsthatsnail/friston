package main

import (
	"fmt";
	"github.com/whatsthatsnail/simple_interpreter/lexer"
)

func main() {
	//scanner := lexer.NewLexer("+-(){},.;*===!!=<<=>>=")
	scanner := lexer.NewLexer(fmt.Sprintf("// this is a comment \n {()} // grouping \n !*/+-=<> <= == // operators \n \"test 2\" 5.5 876. \n and class else false for func if or this true var while x y test"))
	test, err := scanner.ScanTokens()
	if !err || err {
		lexer.PrintTokens(test)
	}
}

