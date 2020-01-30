package ast

import "github.com/whatsthatsnail/simple_interpreter/lexer"

// Visitor interface (all other visitors must implement this)
type Visitor interface {
	visitBinary(b Binary)
	visitUnary(u Unary)
	visitGroup(g Group)
	visitLiteral(l Literal)
}

// Node types:
type Expression interface {
	Accept(v Visitor)
}

type Binary struct {
	X Expression
	Op lexer.Token
	Y Expression
}

func (b Binary) Accept(v Visitor) {
	v.visitBinary(b)
}

type Unary struct {
	Op lexer.Token
	X Expression
}

func (u Unary) Accept(v Visitor) {
	 v.visitUnary(u)
}

type Group struct {
	Left lexer.Token
	X Expression
	Right lexer.Token
}

func (g Group) Accept(v Visitor) {
	 v.visitGroup(g)
}

type Literal struct {
	X lexer.Token
}

func (l Literal) Accept(v Visitor) {
	 v.visitLiteral(l)
}
