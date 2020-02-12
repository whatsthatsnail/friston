package interpreter

import (
	"fmt"
	"friston/ast"
	"friston/environment"
	"friston/lexer"
)

type Function interface {
	Call(i Interpreter, args []interface{}) interface{}
	Arity() int
}

type UserFunc struct {
	Identifier lexer.Token
	Parameters []string
	Block ast.Block
}

func (u UserFunc) Call(i Interpreter, args []interface{}) interface{} {
	// Each function call gets its own environment to allow recursion.
	i.environment = environment.NewEnclosed(i.environment)

	for n, arg := range args {
		i.environment.Declare(u.Parameters[n], arg)
	}

	value := i.executeBlock(u.Block)
	fmt.Printf("Call returning: %v\n", value)
	return value
}

func (u UserFunc) Arity() int { return len(u.Parameters) }

// String representation to allow code to print UserFunction types.
func (u UserFunc) String() string {
	return "<fn " + u.Identifier.Lexeme + ">"
}