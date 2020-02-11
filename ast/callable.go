package ast

type Function interface {
	Call(i Interpreter, args []interface{}) interface{}
	Arity() int
}