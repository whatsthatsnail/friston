package ast

import "github.com/whatsthatsnail/simple_interpreter/lexer"

// Visitor interface (all other visitors must implement this)
type Visitor interface {
	visitEquality(e Equality)
	visitComparison(c Comparison)
	visitAddition(a Addition)
	visitMultiplication(m Multiplication)
	visitUnary(u Unary)
	visitGroup(g Group)
	visitLiteral(l Literal)
}

// Node types:
type Expression interface {
	Accept(v Visitor)
}

type Equality struct {
	X Expression
	Op lexer.Token
	Y Expression
}

func (e Equality) Accept(v Visitor) {
	 v.visitEquality(e)
}

type Comparison struct {
	X Expression
	Op lexer.Token
	Y Expression
}

func (c Comparison) Accept(v Visitor) {
	 v.visitComparison(c)
}

type Addition struct {
	X Expression
	Op lexer.Token
	Y Expression
}

func (a Addition) Accept(v Visitor) {
	 v.visitAddition(a)
}

type Multiplication struct {
	X Expression
	Op lexer.Token
	Y Expression
}

func (m Multiplication) Accept(v Visitor) {
	 v.visitMultiplication(m)
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
