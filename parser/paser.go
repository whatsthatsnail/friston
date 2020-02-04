package parser

import (
	"github.com/whatsthatsnail/simple_interpreter/lexer"
	"github.com/whatsthatsnail/simple_interpreter/ast";
	"github.com/whatsthatsnail/simple_interpreter/errors"
)

type parser struct {
	tokens []lexer.Token
	current int
	statements []ast.Statement
}

// Parser constructor, initializes default vaules
func NewParser(tokens []lexer.Token) parser {
	p := parser{}
	p.tokens = tokens
	p.current = 0

	return p
}

func (p *parser) Parse() []ast.Statement {
	for !p.isAtEnd() {
		p.statements = append(p.statements, p.statement())
	}

	return p.statements
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
	expr := p.equality()

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

	return p.primary()
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
	if p.match([]lexer.TokenType{lexer.VAR}) {
		return p.varDecl()
	} else {
		return p.statement()
	}
}

func (p *parser) varDecl() ast.Statement {
	var name lexer.Token
	if p.consume(lexer.IDENTIFIER, "Expect variable name.") {
		name = p.previous()
	}

	var initializer ast.Expression
	if p.match([]lexer.TokenType{lexer.EQUAL}) {
		initializer = p.expression()
	}

	p.consume(lexer.SEMICOLON, "Expect ';' after variable declaration.")
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
	case lexer.PRINT:
		p.advance()
		return p.printStmt()
	case lexer.VAR:
		p.advance()
		return p.varDecl()
	case lexer.LEFT_BRACE:
		p.advance()
		return p.block()
	}

	return p.exprStmt()
}

func (p *parser) exprStmt() ast.Statement {
	expr := p.expression()
	p.consume(lexer.SEMICOLON, "Expect ';' after expression.")
	return ast.ExprStmt{expr}
}

func (p *parser) ifStmt() ast.Statement {
	p.consume(lexer.LEFT_PAREN, "Expect '(' before if condition.")
	condition := p.expression()
	p.consume(lexer.RIGHT_PAREN, "Expect ')' after if condition.")
	
	thenBranch := p.statement()
	var elseBranch ast.Statement = nil

	if p.match([]lexer.TokenType{lexer.ELSE}) {
		elseBranch = p.statement() 
	}

	return ast.IfStmt{condition, thenBranch, elseBranch}
}

func (p *parser) whileStmt() ast.Statement {
	p.consume(lexer.LEFT_PAREN, "Expect '(' before while condition.")
	condition := p.expression()
	p.consume(lexer.RIGHT_PAREN, "Expect ')' after while condition.")

	loopBranch := p.statement()

	return ast.WhileStmt{condition, loopBranch}
}

// For loops are syntactic sugar, they are expressed as while loops.
func (p *parser) forStmt() ast.Statement {
	p.consume(lexer.LEFT_PAREN, "Expect '(' before initialization.")
	declaration := p.declaration()
	
	condition := p.equality()
	p.consume(lexer.SEMICOLON, "Expect ';' after condition statement.")
	
	increment := p.assignment()
	p.consume(lexer.RIGHT_PAREN, "Expect ')' after increment statement.")
	
	loopBranch := p.statement()

	block := ast.Block{[]ast.Statement{loopBranch, increment}}
	forLoop := []ast.Statement{declaration, ast.WhileStmt{condition, block}}

	return ast.Block{forLoop}
}

func (p *parser) printStmt() ast.Statement {
	expr := p.expression()
	p.consume(lexer.SEMICOLON, "Expect ';' after value.")
	return ast.PrintStmt{expr}
}

func (p *parser) block() ast.Statement {
	var stmts []ast.Statement
	for !p.check(lexer.RIGHT_BRACE) && !p.isAtEnd() {
		stmts = append(stmts, p.statement())
	}

	if p.consume(lexer.RIGHT_BRACE, "Expect '}' after block.") {
		return ast.Block{stmts}
	} else {
		// TODO: Proper parser errors!
		return nil
	}
}

// Error handling:

func (p *parser) consume(tType lexer.TokenType, message string) bool {
	if p.check(tType) {
		p.advance()
		return true
	} else {
		p.parseError(p.peek(), message)
		return false
	}
}

func (p *parser) parseError(token lexer.Token, message string) {
	errors.ThrowError(token.Line, message)
}

// TODO: Sychronize to previous statement when a parseError is called.
func (p *parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().TType == lexer.SEMICOLON {
			return
		}
	}
}