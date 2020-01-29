package parser

import (
	"github.com/whatsthatsnail/simple_interpreter/lexer";
	"github.com/whatsthatsnail/simple_interpreter/ast";
	"github.com/whatsthatsnail/simple_interpreter/errors"
)

type parser struct {
	tokens []lexer.Token
	current int
}

// Parser constructor, initializes default vaules
func NewParser(tokens []lexer.Token) parser {
	p := parser{}
	p.tokens = tokens
	p.current = 0

	return p
}

// TODO: integrate error handling
func (p *parser) Parse() ast.Expression {
	return p.expresssion()
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
		p.current =+ 1
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
func (p *parser) expresssion() ast.Expression {
	return p.equality()
}

func (p *parser) equality() ast.Expression {
	expr := p.comparison()

	for p.match([]lexer.TokenType{lexer.BANG_EQUAL, lexer.EQUAL}) {
		operator := p.previous()
		right := p.comparison()
		expr = ast.Equality{right, operator, expr}
	}

	return expr
}

func (p *parser) comparison() ast.Expression {
	expr := p.addition()

	for p.match([]lexer.TokenType{lexer.LESS, lexer.LESS_EQUAL, lexer.GREATER, lexer.GREATER_EQUAL}) {
		operator := p.previous()
		right := p.addition()
		expr = ast.Equality{right, operator, expr}
	}

	return expr
}

func (p *parser) addition() ast.Expression {
	expr := p.multiplication()

	for p.match([]lexer.TokenType{lexer.PLUS, lexer.MINUS}) {
		operator := p.previous()
		right := p.multiplication()
		expr = ast.Equality{right, operator, expr}
	}

	return expr
}

func (p *parser) multiplication() ast.Expression {
	expr := p.unary()

	for p.match([]lexer.TokenType{lexer.STAR, lexer.SLASH}) {
		operator := p.previous()
		right := p.unary()
		expr = ast.Equality{right, operator, expr}
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
		return ast.Literal{p.advance()}
	} else if p.match([]lexer.TokenType{lexer.NUMBER, lexer.STRING}) {
		return ast.Literal{p.advance()}
	} else if p.match([]lexer.TokenType{lexer.LEFT_PAREN}) {
		left := p.previous()
		expr := p.expresssion()
		p.consume(lexer.RIGHT_PAREN, "Expect ')' after expression.")
		right := p.previous()
		return ast.Group{left, expr, right}
	} else {
		p.parseError(p.peek(), "Expect expression.")
		return nil
	}
}

// Error handling:

func (p *parser) consume(tType lexer.TokenType, message string) {
	if p.check(tType) {
		p.advance()
	} else {
		p.parseError(p.peek(), message)
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