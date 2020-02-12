package visitors

import (
	"fmt"
	"friston/ast"
)

type ASTPrinter struct{}

func (printer ASTPrinter) VisitBinary(b ast.Binary) interface{} {
	fmt.Printf("(")
	b.X.Accept(printer)
	fmt.Printf(" %s ", b.Op.Lexeme)
	b.Y.Accept(printer)
	fmt.Printf(")")
	return nil
}

func (printer ASTPrinter) VisitLogic(l ast.Logic) interface{} {
	l.X.Accept(printer)
	fmt.Printf(" %s ", l.Op.Lexeme)
	l.Y.Accept(printer)
	return nil
}

func (printer ASTPrinter) VisitUnary(u ast.Unary) interface{} {
	fmt.Printf("(")
	fmt.Printf(" %s ", u.Op.Lexeme)
	u.X.Accept(printer)
	fmt.Printf(")")
	return nil
}

func (printer ASTPrinter) VisitGroup(g ast.Group) interface{} {
	fmt.Printf("%sg ", g.Left.Lexeme)
	g.X.Accept(printer)
	fmt.Printf("%s ", g.Right.Lexeme)
	return nil
}

func (printer ASTPrinter) VisitLiteral(l ast.Literal) interface{} {
	fmt.Printf("%v", l.X.Literal)
	return nil
}

func (printer ASTPrinter) VisitVariable(vr ast.Variable) interface{} {
	fmt.Printf("%s", vr.Name.Lexeme)
	return nil
}

func (printer ASTPrinter) VisitAssignment(a ast.Assignment) interface{} {
	fmt.Printf("%s = ", a.Name.Lexeme)
	a.Value.Accept(printer)
	return nil
}

func (printer ASTPrinter) VisitCall(c ast.Call) interface{} {
	c.Callee.Accept(printer)
	fmt.Printf("(")
	for _, arg := range(c.Arguments) {
		arg.Accept(printer)
	}
	fmt.Printf(")\n")

	return nil
}

func (printer ASTPrinter) VisitExprStmt(e ast.ExprStmt) {
	e.Expr.Accept(printer)
	fmt.Printf(";\n")
}

func (printer ASTPrinter) VisitIfStmt(stmt ast.IfStmt) {
	fmt.Printf("if (")
	stmt.Condition.Accept(printer)
	fmt.Printf(") ")

	fmt.Printf("then ")
	stmt.ThenBranch.Accept(printer)

	if stmt.ElseBranch != nil {
		fmt.Printf("else ")
		stmt.ElseBranch.Accept(printer)
	}
}

func (printer ASTPrinter) VisitWhileStmt(stmt ast.WhileStmt) {
	fmt.Printf("while (")
	stmt.Condition.Accept(printer)
	fmt.Printf(") ")
	stmt.LoopBranch.Accept(printer)
}

func (printer ASTPrinter) VisitFuncDecl(f ast.FuncDecl) {
	fmt.Printf("function %s : ", f.Name.Lexeme)
	for _, arg := range(f.ArgumentNames) {
		fmt.Printf(" %s ", arg.Lexeme)
	}
	fmt.Printf(" =\n")
	f.StmtBlock.Accept(printer)
}

func (printer ASTPrinter) VisitVarDecl(d ast.VarDecl) {
	fmt.Printf("let %s = ", d.Name.Lexeme)
	d.Initializer.Accept(printer)
	fmt.Printf(";\n")
}

func (printer ASTPrinter) VisitBlock(b ast.Block) {
	fmt.Printf("{\n")
	for _, s := range(b.Stmts) {
		s.Accept(printer)
	}
	fmt.Printf("}\n")
}