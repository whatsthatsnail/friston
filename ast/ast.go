package ast

import "friston/lexer"

// Visitor interface
type Visitor interface {
	VisitBinary(b Binary) interface{}
	VisitLogic(l Logic) interface{}
	VisitUnary(u Unary) interface{}
	VisitGroup(g Group) interface{}
	VisitLiteral(l Literal) interface{}
	VisitVariable(vr Variable) interface{}
	VisitAssignment(a Assignment) interface{}
	VisitCall(c Call) interface{}

	VisitExprStmt(e ExprStmt)
	VisitIfStmt(stmt IfStmt)
	VisitWhileStmt(stmt WhileStmt)
	VisitFuncDecl(f FuncDecl)
	VisitVarDecl(d VarDecl)
	VisitBlock(b Block)
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
	return v.VisitBinary(b)
}

type Logic struct {
	X Expression
	Op lexer.Token
	Y Expression
}

func (l Logic) Accept(v Visitor) interface{} {
	return v.VisitLogic(l)
}

type Unary struct {
	Op lexer.Token
	X Expression
}

func (u Unary) Accept(v Visitor) interface{} {
	 return v.VisitUnary(u)
}

type Group struct {
	Left lexer.Token
	X Expression
	Right lexer.Token
}

func (g Group) Accept(v Visitor) interface{} {
	 return v.VisitGroup(g)
}

type Literal struct {
	X lexer.Token
}

func (l Literal) Accept(v Visitor) interface{} {
	 return v.VisitLiteral(l)
}

type Variable struct {
	Name lexer.Token
}

func (vr Variable) Accept(v Visitor) interface{} {
	return v.VisitVariable(vr)
}

type Assignment struct {
	Name lexer.Token
	Value Expression
}

func (a Assignment) Accept(v Visitor) interface{} {
	return v.VisitAssignment(a)
}

type Call struct {
	Callee Expression
	Paren lexer.Token
	Arguments []Expression
}

func (c Call) Accept(v Visitor) interface{} {
	return v.VisitCall(c)
}

//Statement types:

type Statement interface {
	Accept(v Visitor) interface{}
}

type ExprStmt struct {
	Expr Expression
}

func (e ExprStmt) Accept(v Visitor) interface{} {
	v.VisitExprStmt(e)
	return nil
}

type IfStmt struct {
	Condition Expression
	ThenBranch Statement
	ElseBranch Statement
}

func (i IfStmt) Accept(v Visitor) interface{} {
	v.VisitIfStmt(i)
	return nil
}

type WhileStmt struct {
	Condition Expression
	LoopBranch Statement
}

func (w WhileStmt) Accept(v Visitor) interface{} {
	v.VisitWhileStmt(w)
	return nil
}

type FuncDecl struct {
	Name lexer.Token
	ArgumentNames []lexer.Token
	StmtBlock Statement
}

func (f FuncDecl) Accept(v Visitor) interface{} {
	v.VisitFuncDecl(f)
	return nil
}

type VarDecl struct {
	Name lexer.Token
	Initializer Expression
}

func (d VarDecl) Accept(v Visitor) interface{} {
	v.VisitVarDecl(d)
	return nil
}

type Block struct {
	Stmts []Statement
}

func (b Block) Accept(v Visitor) interface{} {
	v.VisitBlock(b)
	return nil
}