package type_generator

import (
	"fmt";
	"strings";
	"friston/lexer"
)

// An ugly, ugly, function to create a complete ast.go source file from an input of tokens.
// TODO: Hangle package reference in field types ex:"lexer.Token"
func GenerateNodeTypes(tokens []lexer.Token) {
	
	// Print package and imports
	fmt.Printf("package ast\n\nimport \"github.com/whatsthatsnail/simple_interpreter/lexer\"\n\n")

	// Create Visitor Interface:
	fmt.Println("// Visitor interface (all other visitors must implement this)")
	fmt.Println("type Visitor interface {")
	for i, token := range(tokens) {
		if token.TType != lexer.EOF && tokens[i + 1].TType == lexer.COLON {
			s := token.Lexeme
			fmt.Printf("\tvisit%s(%s %s)\n", s, strings.ToLower(string(s[0])), s)
		}
	}
	fmt.Printf("}\n\n")

	// Expression
	fmt.Println("// Node types:")
	fmt.Printf("type Expression interface {\n\taccept(v Visitor)\n}\n\n")

	// Create custom node types (which implement the expression interface)
	var currentNode string
	for i, token := range(tokens) {
		var nextType lexer.TokenType

		if token.TType != lexer.EOF {
			nextType = tokens[i + 1].TType 
		} else {
			fmt.Printf("}\n\n")
			fmt.Printf("func (%s %s) accept(v Visitor) {\n\t v.visit%s(%s)\n}\n", strings.ToLower(string(currentNode[0])), currentNode, currentNode, strings.ToLower(string(currentNode[0])))
		}
		
		if nextType == lexer.COLON{
			if i != 0 {
				fmt.Printf("}\n\n")
				fmt.Printf("func (%s %s) accept(v Visitor) {\n\t v.visit%s(%s)\n}\n\n", strings.ToLower(string(currentNode[0])), currentNode, currentNode, strings.ToLower(string(currentNode[0])))
			}

			currentNode = token.Lexeme
			fmt.Printf("type %s struct {\n", currentNode)
		} else if nextType == lexer.EQUAL {
			fmt.Printf("\t%s ", token.Lexeme)
		} else if nextType == lexer.SEMICOLON {
			fmt.Printf("%s\n", token.Lexeme)
		}
	}
}
