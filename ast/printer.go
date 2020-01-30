package ast

import "fmt"

type ASTPrinter struct{}

func (printer ASTPrinter) visitBinary(b Binary) interface{} {
	fmt.Printf("(")
	b.X.Accept(printer)
	fmt.Printf(" %s ", b.Op.Lexeme)
	b.Y.Accept(printer)
	fmt.Printf(")")
	return nil
}

func (printer ASTPrinter) visitUnary(u Unary) interface{} {
	fmt.Printf("(")
	fmt.Printf(" %s ", u.Op.Lexeme)
	u.X.Accept(printer)
	fmt.Printf(")")
	return nil
}

func (printer ASTPrinter) visitGroup(g Group) interface{} {
	fmt.Printf("%sg ", g.Left.Lexeme)
	g.X.Accept(printer)
	fmt.Printf("%s ", g.Right.Lexeme)
	return nil
}

func (printer ASTPrinter) visitLiteral(l Literal) interface{} {
	fmt.Printf("%v", l.X.Literal)
	return nil
}
