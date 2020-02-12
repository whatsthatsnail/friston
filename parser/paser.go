package parser

import (
	"friston/ast"
	"friston/errors"
	"friston/lexer"
)

type parser struct {
	tokens []lexer.Token
	current int
	statements []ast.Statement
	errFlag bool
}

// Parser constructor, initializes default vaules
func NewParser(tokens []lexer.Token) parser {
	p := parser{}
	p.tokens = tokens
	p.current = 0

	return p
}

func (p *parser) Parse() ([]ast.Statement, bool) {
	for !p.isAtEnd() {
		p.statements = append(p.statements, p.statement())
	}

	return p.statements, p.errFlag
}

// Helper methods:
// Return current token without consuming it (advancing).
func (p *parser) peek() lexer.Token {
	return p.tokens[p.current]
}

// Check if the current position is the last token (an EOF token).
func (p *parser) isAtEnd() bool {
	return p.peek().TType == lexer.EOF
}

// Return the token directly before the current position.
func (p *parser) previous() lexer.Token {
	return p.tokens[p.current - 1]
}

// Advance the current position and return the current token.
func (p *parser) advance() lexer.Token {
	if !p.isAtEnd() {
		p.current++
		return p.previous()
	} else {
		return p.tokens[p.current]
	}
}

// Compare the type of the current token to a given TokenType.
func (p *parser) check(tType lexer.TokenType) bool {
	if !p.isAtEnd() {
		return p.peek().TType == tType
	} else {
		return false
	}
}

// Compre the curent token type against a list of given TokenTypes (and advance).
func (p *parser) match(tTypes []lexer.TokenType) bool {
	for _, tType := range(tTypes) {
		if p.check(tType) {
			p.advance()
			return true
		}
	}

	return false
}

// Node creator methods:
// Each method calls a method of higher precedence if possible,
// so high precedence expressions are evaluated first.
func (p *parser) expression() ast.Expression {
	return p.assignment()
}

func (p *parser) assignment() ast.Expression {
	expr := p.or()

	// IDENTIFIER "=" assignment
	if p.match([]lexer.TokenType{lexer.EQUAL}) {
		equals := p.previous()
		value := p.assignment()

		vr, ok := expr.(ast.Variable)
		if ok {
			name := vr.Name
			return ast.Assignment{name, value}
		}

		errors.ThrowError(equals.Line, "Invalid assignment target.")
	}

	// Implement increment (++) and decrement (--) as sugar, ex: translate a++ to a = a + 1
	if p.match([]lexer.TokenType{lexer.PLUS_PLUS}) {
		operator := lexer.Token{lexer.PLUS, "+", nil, p.previous().Line}
		right := ast.Literal{lexer.Token{lexer.NUMBER, "1", 1.0, p.previous().Line}}
		binary := ast.Binary{expr, operator, right}

		vr, ok := expr.(ast.Variable)
		if ok {
			name := vr.Name
			return ast.Assignment{name, binary}
		}

		errors.ThrowError(p.previous().Line, "Invalid increment target.")
	}

	// Decrement
	if p.match([]lexer.TokenType{lexer.MINUS_MINUS}) {
		operator := lexer.Token{lexer.MINUS, "-", nil, p.previous().Line}
		right := ast.Literal{lexer.Token{lexer.NUMBER, "1", 1.0, p.previous().Line}}
		binary := ast.Binary{expr, operator, right}

		vr, ok := expr.(ast.Variable)
		if ok {
			name := vr.Name
			return ast.Assignment{name, binary}
		}

		errors.ThrowError(p.previous().Line, "Invalid decrement target.")
	}

	return expr
}

func (p *parser) or() ast.Expression {
	expr := p.and()

	for p.match([]lexer.TokenType{lexer.OR}) {
		operator := p.previous()
		value := p.and()
		expr = ast.Logic{expr, operator, value}
	}

	return expr
}

func (p *parser) and() ast.Expression {
	expr := p.equality()

	for p.match([]lexer.TokenType{lexer.AND}) {
		operator := p.previous()
		value := p.equality()
		expr = ast.Logic{expr, operator, value}
	}

	return expr
}

func (p *parser) equality() ast.Expression {
	expr := p.comparison()

	for p.match([]lexer.TokenType{lexer.BANG_EQUAL, lexer.EQUAL_EQUAL}) {
		operator := p.previous()
		right := p.comparison()
		expr = ast.Binary{expr, operator, right}
	}

	return expr
}

func (p *parser) comparison() ast.Expression {
	expr := p.addition()

	for p.match([]lexer.TokenType{lexer.LESS, lexer.LESS_EQUAL, lexer.GREATER, lexer.GREATER_EQUAL}) {
		operator := p.previous()
		right := p.addition()
		expr = ast.Binary{expr, operator, right}
	}

	return expr
}

func (p *parser) addition() ast.Expression {
	expr := p.multiplication()

	for p.match([]lexer.TokenType{lexer.PLUS, lexer.MINUS}) {
		operator := p.previous()
		right := p.multiplication()
		expr = ast.Binary{expr, operator, right}
	}

	return expr
}

func (p *parser) multiplication() ast.Expression {
	expr := p.unary()

	for p.match([]lexer.TokenType{lexer.STAR, lexer.SLASH}) {
		operator := p.previous()
		right := p.unary()
		expr = ast.Binary{expr, operator, right}
	}

	return expr
}

func (p *parser) unary() ast.Expression {
	if p.match([]lexer.TokenType{lexer.BANG, lexer.MINUS}) {
		operator := p.previous()
		right := p.unary()
		return ast.Unary{operator, right}
	}

	return p.call()
}

func (p *parser) call() ast.Expression {
	expr := p.primary()
	
	var arguments []ast.Expression
	if p.match([]lexer.TokenType{lexer.LEFT_PAREN}) {
		paren := p.previous()

		for !p.check(lexer.RIGHT_PAREN) {
			arg := p.expression()
			arguments = append(arguments, arg)
			if p.peek().TType != lexer.RIGHT_PAREN {
				p.consume(lexer.COMMA, "Arguments must be separated by ','.")
			}
		}

		p.consume(lexer.RIGHT_PAREN, "Arguments must end with ')'.")
		return ast.Call{expr, paren, arguments}
	}

	return expr
}

func (p *parser) primary() ast.Expression {
	if p.match([]lexer.TokenType{lexer.TRUE, lexer.FALSE, lexer.NIL}) {
		return ast.Literal{p.previous()}
	} else if p.match([]lexer.TokenType{lexer.NUMBER, lexer.STRING}) {
		return ast.Literal{p.previous()}
	} else if p.match([]lexer.TokenType{lexer.LEFT_PAREN}) {
		left := p.previous()
		expr := p.expression()
		p.consume(lexer.RIGHT_PAREN, "Expect ')' after expression.")
		right := p.previous()
		return ast.Group{left, expr, right}
	} else if p.match([]lexer.TokenType{lexer.IDENTIFIER}) {
		return ast.Variable{p.previous()}
	} else {
		p.parseError(p.peek(), "Expect expression.")
		return nil
	}
}

// Statement/Declaration creator methods:

func (p *parser) declaration() ast.Statement {
	if p.match([]lexer.TokenType{lexer.LET}) {
		return p.varDecl()
	} else if p.match([]lexer.TokenType{lexer.FUNCTION}) {
		return p.funcDecl()
	} else {
		return p.varDecl()
	}
}

func (p *parser) funcDecl() ast.Statement {
	var name lexer.Token
	p.consume(lexer.IDENTIFIER, "Expect function name.")
	name = p.previous()

	p.consume(lexer.COLON, "Expect ':' in function declaration.")

	var parameters []lexer.Token
	for !p.check(lexer.EQUAL) {
		param := p.advance()
		parameters = append(parameters, param)
		if p.peek().TType != lexer.EQUAL {
			p.consume(lexer.COMMA, "Parameters must be separated by ','.")
		}
	}

	p.consume(lexer.EQUAL, "Arguments must end with '='.")
	p.consume(lexer.INDENT, "Function blocks must begin with an indent.")

	block := p.block()

	return ast.FuncDecl{name, parameters, block}
}

func (p *parser) varDecl() ast.Statement {
	var name lexer.Token
	p.consume(lexer.IDENTIFIER, "Expect variable name.")
	name = p.previous()

	var initializer ast.Expression
	if p.match([]lexer.TokenType{lexer.EQUAL}) {
		initializer = p.expression()
	}

	p.consumeMatch([]lexer.TokenType{lexer.NEWLINE, lexer.SEMICOLON}, "Expect ';' or new line after variable declaration.")
	return ast.VarDecl{name, initializer}
}

func (p *parser) statement() ast.Statement {
	switch p.peek().TType {
	case lexer.IF:
		p.advance()
		return p.ifStmt()
	case lexer.WHILE:
		p.advance()
		return p.whileStmt()
	case lexer.FOR:
		p.advance()
		return p.forStmt()
	case lexer.FUNCTION:
		p.advance()
		return p.funcDecl()
	case lexer.LET:
		p.advance()
		return p.varDecl()
	case lexer.INDENT:
		p.advance()
		return p.block()
	case lexer.RETURN:
		p.advance()
		return p.returnStmt()
	}

	return p.exprStmt()
}

func (p *parser) exprStmt() ast.Statement {
	expr := p.expression()

	// Check for non control flow expressions
	//if !p.check(lexer.THEN) {
	p.consumeMatch([]lexer.TokenType{lexer.NEWLINE, lexer.SEMICOLON}, "Expect ';' or new line after expression.")
	//}

	return ast.ExprStmt{expr}
}

func (p *parser) ifStmt() ast.Statement {
	condition := p.expression()
	p.consume(lexer.THEN, "Expect 'then' after if condition.")

	thenBranch := p.statement()
	var elseBranch ast.Statement = nil

	if p.match([]lexer.TokenType{lexer.ELSE}) {
		p.consume(lexer.THEN, "Expect 'then' after else statement.")
		elseBranch = p.statement()
	}

	return ast.IfStmt{condition, thenBranch, elseBranch}
}

func (p *parser) whileStmt() ast.Statement {
	condition := p.expression()
	p.consume(lexer.THEN, "Expect 'then' after while condition.")

	loopBranch := p.statement()

	return ast.WhileStmt{condition, loopBranch}
}

// For loops are syntactic sugar, they are expressed as while loops.
func (p *parser) forStmt() ast.Statement {
	declaration := p.declaration()

	condition := p.equality()
	p.consume(lexer.SEMICOLON, "Expect ';' after condition statement.")

	increment := p.assignment()
	p.consume(lexer.THEN, "Expect 'then' after increment statement.")

	loopBranch := p.statement()

	block := ast.Block{[]ast.Statement{loopBranch, increment}}
	forLoop := []ast.Statement{declaration, ast.WhileStmt{condition, block}}

	return ast.Block{forLoop}
}

func (p *parser) returnStmt() ast.Statement {
	keyword := p.previous()
	var expr ast.Expression = nil

	if !p.check(lexer.NEWLINE) {
		expr = p.expression()
	}

	p.consume(lexer.NEWLINE, "Return statement must end in a new line.")

	return ast.ReturnStmt{keyword, expr}
}

func (p *parser) block() ast.Block {
	var stmts []ast.Statement
	for !p.check(lexer.DEDENT) && !p.isAtEnd() {
		stmts = append(stmts, p.statement())
	}

	if !p.isAtEnd() {
		p.consume(lexer.DEDENT, "Expect dedent after block statement.")
	}

	return ast.Block{stmts}
}

// Error handling:

func (p *parser) consume(tType lexer.TokenType, message string) {
	if p.check(tType) {
		p.advance()
	} else {
		p.parseError(p.peek(), message)
	}
}

func (p *parser) consumeMatch(types []lexer.TokenType, message string) {
	if p.match(types) {
		return
	} else {
		p.parseError(p.peek(), message)
	}
}

func (p *parser) parseError(token lexer.Token, message string) {
	p.errFlag = true
	errors.ThrowError(token.Line, message)
	p.synchronize()
}

// TODO: Synchronize to previous statement when a parseError is called.
func (p *parser) synchronize() {
	p.advance()

	for !p.isAtEnd() && !(p.peek().TType > lexer.IDENTIFIER && p.peek().TType < lexer.EOF) {
		p.advance()
	}
}