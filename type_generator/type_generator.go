package type_generator

import (
	"fmt";
	"strings"
)

func GenASTTypes(types []string) {
	fmt.Printf("package ast\n\n")

	// Create Visitor Interface:
	fmt.Println("// Visitor interface (all other visitors must implement this)")
	fmt.Println("type Visitor interface {")
	for _, s := range(types) {
		fmt.Printf("\tvisit%s(%s %s)\n", s, strings.ToLower(string(s[0])), s)
	}
	fmt.Printf("}\n\n")

	// Create node types:

	// Expression
	fmt.Println("// Node types:")
	fmt.Printf("type Expression interface {\n\taccept(v Visitor)\n}\n\n")

	// User specified nodes:
	for _, s := range(types) {
		// Struct:
		fmt.Printf("type %s struct {\n\t// TODO: define fields\n}\n\n", s)
		// accept(v Visitor):
		fmt.Printf("func (%s %s) accept(v Visitor) {\n\t v.visit%s(%s)\n}\n\n", strings.ToLower(string(s[0])), s, s, strings.ToLower(string(s[0])))
	}
}