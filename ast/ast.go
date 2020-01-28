package ast

import "github.com/whatsthatsnail/simple_interpreter/lexer"

// Visitor interface (all other visitors must implement this)
type Visitor interface {
	visitEquality(e Equality)
	visitComparison(c Comparison)
	visitComparison(c Comparison)
	visitAddition(a Addition)
	visitMultiplication(m Multiplication)
	visitUnary(u Unary)
	visitPrimary(p Primary)
}

// Node types:
type Expression interface {
	accept(v Visitor)
}

type Equality struct {
	x Expression
	op Token
	y Expression
}

func (e Equality) accept(v Visitor) {
	 v.visitEquality(e)
}

type Comparison struct {
	x Expression
	op Token
	y Expression
}

func (c Comparison) accept(v Visitor) {
	 v.visitComparison(c)
}

type Comparison struct {
	x Expression
	op Token
	y Expression
}

func (c Comparison) accept(v Visitor) {
	 v.visitComparison(c)
}

type Addition struct {
	x Expression
	op Token
	y Expression
}

func (a Addition) accept(v Visitor) {
	 v.visitAddition(a)
}

type Multiplication struct {
	x Expression
	op Token
	y Expression
}

func (m Multiplication) accept(v Visitor) {
	 v.visitMultiplication(m)
}

type Unary struct {
	x Expression
	op Token
}

func (u Unary) accept(v Visitor) {
	 v.visitUnary(u)
}

type Primary struct {
	x Token
}

func (p Primary) accept(v Visitor) {
	 v.visitPrimary(p)
}
