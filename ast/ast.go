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

	VisitExprStmt(e ExprStmt) interface{}
	VisitIfStmt(stmt IfStmt) interface{}
	VisitWhileStmt(stmt WhileStmt) interface{}
	VisitFuncDecl(f FuncDecl) interface{}
	VisitVarDecl(d VarDecl) interface{}
	VisitReturn(d ReturnStmt) interface{}
	VisitBlock(b Block) interface{}
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
	return v.VisitExprStmt(e)
}

type IfStmt struct {
	Condition Expression
	ThenBranch Statement
	ElseBranch Statement
}

func (i IfStmt) Accept(v Visitor) interface{} {
	return v.VisitIfStmt(i)
}

type WhileStmt struct {
	Condition Expression
	LoopBranch Statement
}

func (w WhileStmt) Accept(v Visitor) interface{} {
	return v.VisitWhileStmt(w)
}

type FuncDecl struct {
	Name lexer.Token
	Parameters []lexer.Token
	Block Block
}

func (f FuncDecl) Accept(v Visitor) interface{} {
	return v.VisitFuncDecl(f)
}

type VarDecl struct {
	Name lexer.Token
	Initializer Expression
}

func (d VarDecl) Accept(v Visitor) interface{} {
	return v.VisitVarDecl(d)
}

type ReturnStmt struct {
	Keyword lexer.Token
	Value Expression
}

func (r ReturnStmt) Accept(v Visitor) interface{} {
	return v.VisitReturn(r)
}

type Block struct {
	Stmts []Statement
}

func (b Block) Accept(v Visitor) interface{} {
	return v.VisitBlock(b)
}