package interpreter

import (
	"friston/ast"
)

type Function interface {
	Call(i Interpreter, args []interface{}) interface{}
	Arity() int
}

type UserFunc struct {
	Args int
	ArgumentNames []string
	StmtBlock ast.Statement
}

func (u UserFunc) Call(i Interpreter, args []interface{}) interface{} {
	for n, arg := range args {
		i.environment.Declare(u.ArgumentNames[n], arg)
	}  

	return i.evaluate(u.StmtBlock)
}

func (u UserFunc) Arity() int { return u.Args }