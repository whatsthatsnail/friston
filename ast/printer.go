package ast

import "fmt"

type ASTPrinter struct{}

func (printer ASTPrinter) visitBinary(b Binary) {
	fmt.Printf("(")
	b.X.Accept(printer)
	fmt.Printf(" %s ", b.Op.Lexeme)
	b.Y.Accept(printer)
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
