package ast

import "fmt"

type ASTPrinter struct{}

func (printer ASTPrinter) visitEquality(e Equality) {
	fmt.Printf("(")
	e.X.Accept(printer)
	fmt.Printf(" %s ", e.Op.Lexeme)
	e.Y.Accept(printer)
	fmt.Printf(")")
}

func (printer ASTPrinter) visitComparison(c Comparison) {
	fmt.Printf("(")
	c.X.Accept(printer)
	fmt.Printf(" %s ", c.Op.Lexeme)
	c.Y.Accept(printer)
	fmt.Printf(")")
}

func (printer ASTPrinter) visitAddition(a Addition) {
	fmt.Printf("(")
	a.X.Accept(printer)
	fmt.Printf(" %s ", a.Op.Lexeme)
	a.Y.Accept(printer)
	fmt.Printf(")")
}

func (printer ASTPrinter) visitMultiplication(m Multiplication) {
	fmt.Printf("(")
	m.X.Accept(printer)
	fmt.Printf(" %s ", m.Op.Lexeme)
	m.Y.Accept(printer)
	fmt.Printf(")")
}

func (printer ASTPrinter) visitUnary(u Unary) {
	fmt.Printf("(")
	fmt.Printf(" %s ", u.Op.Lexeme)
	u.X.Accept(printer)
	fmt.Printf(")")
}

func (printer ASTPrinter) visitGroup(g Group) {
	fmt.Printf("%sg ", g.Left.Lexeme)
	g.X.Accept(printer)
	fmt.Printf("%s ", g.Right.Lexeme)
}

func (printer ASTPrinter) visitLiteral(l Literal) {
	fmt.Printf("%v", l.X.Literal)
}
