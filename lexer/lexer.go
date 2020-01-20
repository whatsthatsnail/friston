package lexer

import "fmt"

type Token struct {
	t_type TokenType
	lexeme string
	literal int
	line int
}

// Prints tokens in a readable manner
func PrintTokens(tokens []Token) {
	for _, tok := range(tokens) {
		s := fmt.Sprintf("{%s, %s, %d, %d}", tok.t_type.typeString(), tok.lexeme, tok.literal, tok.line)
		fmt.Println(s)
	}
}

type lexer struct {
	start int
	current int
	line int
	tokens []Token
	source string
}

// Lexer constructor, initializes default values
func NewLexer(code string) lexer {
	l := lexer{}
	l.start = 0
	l.current = 0
	l.line = 1
	l.source = code

	return l
}

// Checks if current has reaced the end of the source
func (l *lexer) isAtEnd() bool {
	return l.current >= len(l.source)
}

// Consumes the current character and returns it
func (l *lexer) advance() rune {
	l.current++
	return rune(l.source[l.current - 1])
}

// Advances current and creates the next token
func (l *lexer) scanToken() Token {
	char := l.advance()

	// Creates single-character tokens
	switch char {
	case '+':
		return Token{PLUS, string(char), 1, l.line}
	case '-':
		return Token{MINUS, string(char), 1, l.line}
	case '(':
		return Token{LEFT_PAREN, string(char), 1, l.line}
	case ')':
		return Token{RIGHT_PAREN, string(char), 1, l.line}
	case '{':
		return Token{LEFT_BRACE, string(char), 1, l.line}
	case '}':
		return Token{RIGHT_BRACE, string(char), 1, l.line}
	case ',':
		return Token{COMMA, string(char), 1, l.line}
	case ';':
		return Token{SEMICOLON, string(char), 1, l.line}
	case '.':
		return Token{DOT, string(char), 1, l.line}
	case '*':
		return Token{STAR, string(char), 1, l.line}
	}

	return Token{}
}

// Loops over all characters in souce, creating tokens as it goes, places EOF token at end of source
func (l *lexer) ScanTokens() []Token {

	for !l.isAtEnd() {
		l.start = l.current
		l.tokens = append(l.tokens, l.scanToken())
	}

	l.tokens = append(l.tokens, Token{EOF, "", 0, l.line})
	return l.tokens
}