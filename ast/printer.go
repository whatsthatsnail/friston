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

func (printer ASTPrinter) visitLogic(l Logic) interface{} {
	l.X.Accept(printer)
	fmt.Printf(" %s ", l.Op.Lexeme)
	l.Y.Accept(printer)
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

func (printer ASTPrinter) visitVariable(vr Variable) interface{} {
	fmt.Printf("%s", vr.Name.Lexeme)
	return nil
}

func (printer ASTPrinter) visitAssignment(a Assignment) interface{} {
	fmt.Printf("%s = ", a.Name.Lexeme)
	a.Value.Accept(printer)
	return nil
}

func (printer ASTPrinter) visitCall(c Call) interface{} {
	c.Callee.Accept(printer)
	fmt.Printf("(")
	for _, arg := range(c.Arguments) {
		arg.Accept(printer)
	}
	fmt.Printf(")\n")

	return nil
}

func (printer ASTPrinter) visitExprStmt(e ExprStmt) interface {} {
	e.Expr.Accept(printer)
	fmt.Printf(";\n")
	return nil
}

func (printer ASTPrinter) visitIfStmt(stmt IfStmt) interface {} {
	fmt.Printf("if (")
	stmt.Condition.Accept(printer)
	fmt.Printf(") ")

	fmt.Printf("then ")
	stmt.ThenBranch.Accept(printer)

	if stmt.ElseBranch != nil {
		fmt.Printf("else ")
		stmt.ElseBranch.Accept(printer)
	}
	return nil
}

func (printer ASTPrinter) visitWhileStmt(stmt WhileStmt) interface {} {
	fmt.Printf("while (")
	stmt.Condition.Accept(printer)
	fmt.Printf(") ")
	stmt.LoopBranch.Accept(printer)

	return nil
}

func (printer ASTPrinter) visitPrintStmt(p PrintStmt) interface {} {
	fmt.Printf("print ")
	p.Expr.Accept(printer)
	fmt.Printf(";\n")
	return nil
}

func (printer ASTPrinter) visitVarDecl(d VarDecl) interface{} {
	fmt.Printf("let %s = ", d.Name.Lexeme)
	d.Initializer.Accept(printer)
	fmt.Printf(";\n")
	return nil
}

func (printer ASTPrinter) visitBlock(b Block) interface{} {
	fmt.Printf("{\n")
	for _, s := range(b.Stmts) {
		s.Accept(printer)
	}
	fmt.Printf("}\n")
	return nil
}