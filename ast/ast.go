package ast

import "github.com/whatsthatsnail/simple_interpreter/lexer"

// Visitor interface (all other visitors must implement this)
type Visitor interface {
	visitBinary(b Binary) interface{}
	visitUnary(u Unary) interface{}
	visitGroup(g Group) interface{}
	visitLiteral(l Literal) interface{}
}

// Node types:
type Expression interface {
	Accept(v Visitor) interface{}
}

type Binary struct {
	X Expression
	Op lexer.Token
	Y Expression
}

func (b Binary) Accept(v Visitor) interface{} {
	return v.visitBinary(b)
}

type Unary struct {
	Op lexer.Token
	X Expression
}

func (u Unary) Accept(v Visitor) interface{} {
	 return v.visitUnary(u)
}

type Group struct {
	Left lexer.Token
	X Expression
	Right lexer.Token
}

func (g Group) Accept(v Visitor) interface{} {
	 return v.visitGroup(g)
}

type Literal struct {
	X lexer.Token
}

func (l Literal) Accept(v Visitor) interface{} {
	 return v.visitLiteral(l)
}
