package lexer

import "fmt"

type Token struct {
	t_type TokenType
	lexeme string
	literal string
	line int
}

// Prints tokens in a readable manner
func PrintTokens(tokens []Token) {
	for _, tok := range(tokens) {
		s := ""
		if tok.literal == "" {
			s = fmt.Sprintf("{%s, %s, %d}", tok.t_type.typeString(), tok.lexeme, tok.line)
		} else {
			s = fmt.Sprintf("{%s, %s, %s, %d}", tok.t_type.typeString(), tok.lexeme, tok.literal, tok.line)
		}
		fmt.Println(s)
	}
}

type lexer struct {
	start int
	current int
	line int
	tokens []Token
	source string
	hadError bool
}

// Lexer constructor, initializes default values
func NewLexer(code string) lexer {
	l := lexer{}
	l.start = 0
	l.current = 0
	l.line = 1
	l.source = code
	l.hadError = false

	return l
}

// Prints error message and sets error flag
func (l *lexer) throwError(message string) {
	l.report(l.line, "Error: " + message)
	l.hadError = true
}

// Useable to print any line dependant message (error, warning, etc.)
func (l *lexer) report(line int, message string) {
	fmt.Printf("[Line %d] %s\n", line, message)
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

// Peeks at next character without consuming it
func (l *lexer) peek() rune {
	if !l.isAtEnd() {
		return rune(l.source[l.current])
	}
	return ' '
}

// If peek() == x -> a and advance, else b
func (l *lexer) match(x rune, a TokenType, b TokenType) TokenType {
	if l.peek() == x {
		l.advance()
		return a
	} else {
		return b
	}
}

// Returns a new Token instance using input type and literal, and infered lexeme and line
func (l *lexer) addToken(t_type TokenType, literal string) Token {
	return Token{t_type, l.source[l.start:l.current], literal, l.line}
}

// Advances current and creates the next token
func (l *lexer) scanToken() Token {
	char := l.advance()

	switch char {

	// Creates single-character tokens
	case '+':
		return l.addToken(PLUS, "")
	case '-':
		return l.addToken(MINUS, "")
	case '(':
		return l.addToken(LEFT_PAREN, "")
	case ')':
		return l.addToken(RIGHT_PAREN, "")
	case '{':
		return l.addToken(LEFT_BRACE, "")
	case '}':
		return l.addToken(RIGHT_BRACE, "")
	case ',':
		return l.addToken(COMMA, "")
	case ';':
		return l.addToken(SEMICOLON, "")
	case '.':
		return l.addToken(DOT, "")
	case '*':
		return l.addToken(STAR, "")

	// Create one/two-character tokens
	// '=' or '=='
	case '=':
		return l.addToken(l.match('=', EQUAL_EQUAL, EQUAL), "")
	// '!' or '!='
	case '!':
		return l.addToken(l.match('=', BANG_EQUAL, BANG), "")
	// '<' or '<='
	case '<':
		return l.addToken(l.match('=', LESS_EQUAL, LESS), "")
	// '>' or '>='
	case '>':
		return l.addToken(l.match('=', GREATER_EQUAL, GREATER), "")
	}

	l.throwError(fmt.Sprintf("Invalid character '%c'", char))
	return Token{}
}

// Loops over all characters in souce, creating tokens as it goes, places EOF token at end of source
func (l *lexer) ScanTokens() ([]Token, bool) {

	for !l.isAtEnd() {
		l.start = l.current
		l.tokens = append(l.tokens, l.scanToken())
	}

	l.tokens = append(l.tokens, Token{EOF, "EOF", "", l.line})
	return l.tokens, l.hadError
}