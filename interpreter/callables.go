package interpreter

import (
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
	Block      ast.Block
	Closure    environment.Environment
}

func (u UserFunc) Call(i Interpreter, args []interface{}) interface{} {
	// Call a function within it's eclosed environment, making an environment chain all the way up to globals through nested functions.
	i.environment = environment.NewEnclosed(u.Closure)

	for n, arg := range args {
		i.environment.Declare(u.Parameters[n], arg)
	}

	value := i.executeBlock(u.Block)
	return value
}

func (u UserFunc) Arity() int { return len(u.Parameters) }

// String representation to allow code to print UserFunction types.
func (u UserFunc) String() string {
	return "<fn " + u.Identifier.Lexeme + ">"
}
