package ast

import (
	"github.com/whatsthatsnail/simple_interpreter/lexer";
	"fmt"
)

// Visitors:

type Visitor interface {
	visitLiteral(l Literal)
	visitUnary(u Unary)
	visitBinary(b Binary)
	visitGrouping(g Grouping)
}

// Print an AST tree in a reasonable manner
type ASTprinter struct {}

func (a ASTprinter) visitLiteral(l Literal) {
	fmt.Print(l.x.GetLiteral())
}

func (a ASTprinter) visitUnary(u Unary) {
	fmt.Print(" (")
	fmt.Print(u.op.GetLexeme() + " ")
	u.x.accept(a)
	fmt.Print(")")
}

func (a ASTprinter) visitBinary(b Binary) {
	fmt.Print(" (" + b.op.GetLexeme())
	b.x.accept(a)
	b.y.accept(a)
	fmt.Print(")")
}

func (a ASTprinter) visitGrouping(g Grouping) {
	fmt.Print(" (group ")
	g.x.accept(a)
	fmt.Print(")")
}

// Nodes:

type Expression interface {
	accept(v Visitor)
}

type Literal struct {
	x lexer.Token
}

func (l Literal) accept(v Visitor) {
	v.visitLiteral(l)
}

type Unary struct {
	op lexer.Token
	x Expression
}

func (u Unary) accept(v Visitor) {
	v.visitUnary(u)
}

type Binary struct {
	x Expression
	op lexer.Token
	y Expression
}

func (b Binary) accept(v Visitor) {
	v.visitBinary(b)
}

type Grouping struct {
	left lexer.Token
	x Expression
	right lexer.Token
}

func (g Grouping) accept(v Visitor) {
	v.visitGrouping(g)
}
