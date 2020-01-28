package parser

import (
	"github.com/whatsthatsnail/simple_interpreter/lexer";
	"github.com/whatsthatsnail/simple_interpreter/ast"
)

type Parser struct {
	tokens []lexer.Token
	current int
}

// Parser constructor, initializes default vaules
func NewParser(tokens []lexer.Token) Parser {
	p := Parser{}
	p.tokens = tokens
	p.current = 0

	return p
}

// Helper methods:

// Return current token without consuming it (advancing).
func (p* Parser) peek() lexer.Token {
	return p.tokens[p.current]
}

// Check if the current position is the last token (an EOF token).
func (p* Parser) isAtEnd() bool {
	return p.peek().GetType() == lexer.EOF
}

// Return te token directly before the current position.
func (p* Parser) previous() lexer.Token {
	return p.tokens[p.current - 1]
}

// Advance the current position and return the current token.
func (p* Parser) advance() lexer.Token {
	if !p.isAtEnd() {
		p.current =+ 1
		return p.previous()
	} else {
		return p.tokens[p.current]
	}
}

// Compare the type of the current token to a given TokenType.
func (p* Parser) check(tType lexer.TokenType) bool {
	if !p.isAtEnd() {
		return p.peek().GetType() == tType 
	} else {
		return false
	}
}

// Compre the curent token type against a list of given TokenTypes (and advance).
func (p* Parser) match(tTypes []lexer.TokenType) bool {
	for _, tType := range(tTypes) {
		if p.check(tType) {
			p.advance()
			return true
		}
	}

	return false
}

// Node creator methods:
func (p* Parser) expresssion() ast.Expression {
	return p.equality()
}

func (p* Parser) equality() ast.Expression {
	expr := p.comparison()

	for p.match([]lexer.TokenType{lexer.BANG_EQUAL, lexer.EQUAL}) {
		operator := p.previous()
		right := p.comparison()
		expr = ast.Equality{right, operator, expr}
	}

	return expr
}

func (p* Parser) comparison() ast.Expression {
	expr := p.addition()

	for p.match([]lexer.TokenType{lexer.LESS, lexer.LESS_EQUAL, lexer.GREATER, lexer.GREATER_EQUAL}) {
		operator := p.previous()
		right := p.addition()
		expr = ast.Equality{right, operator, expr}
	}

	return expr
}

func (p* Parser) addition() ast.Expression {
	expr := p.multiplication()

	for p.match([]lexer.TokenType{lexer.PLUS, lexer.MINUS}) {
		operator := p.previous()
		right := p.multiplication()
		expr = ast.Equality{right, operator, expr}
	}

	return expr
}

func (p* Parser) multiplication() ast.Expression {
	expr := p.unary()

	for p.match([]lexer.TokenType{lexer.STAR, lexer.SLASH}) {
		operator := p.previous()
		right := p.unary()
		expr = ast.Equality{right, operator, expr}
	}

	return expr
}

func (p* Parser) unary() ast.Expression {
	if p.match([]lexer.TokenType{lexer.BANG, lexer.MINUS}) {
		operator := p.previous()
		right := p.unary()
		return ast.Unary{operator, right}
	}

	return p.primary()
}

// TODO: Implement primary method, needs Literal type, check for grouping, etc
func (p* Parser) primary() ast.Primary {

}