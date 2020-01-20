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

// TODO: Report column of error, print line of error, etc.
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
	} else {
		return '\n'
	}
	return '\n'
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

// Adds a new Token instance to l.tokens using input type and literal, and infered lexeme and line
func (l *lexer) addToken(t_type TokenType, literal string) {
	l.tokens = append(l.tokens, Token{t_type, l.source[l.start:l.current], literal, l.line})
}

// Advances current and adds the current token
func (l *lexer) scanToken() {
	char := l.advance()

	switch char {

	// Creates single-character tokens
	case '+':
		l.addToken(PLUS, "")
	case '-':
		l.addToken(MINUS, "")
	case '(':
		l.addToken(LEFT_PAREN, "")
	case ')':
		l.addToken(RIGHT_PAREN, "")
	case '{':
		l.addToken(LEFT_BRACE, "")
	case '}':
		l.addToken(RIGHT_BRACE, "")
	case ',':
		l.addToken(COMMA, "")
	case ';':
		l.addToken(SEMICOLON, "")
	case '.':
		l.addToken(DOT, "")
	case '*':
		l.addToken(STAR, "")

	// Create one or two-character tokens
	case '=':
		l.addToken(l.match('=', EQUAL_EQUAL, EQUAL), "")
	case '!':
		l.addToken(l.match('=', BANG_EQUAL, BANG), "")
	case '<':
		l.addToken(l.match('=', LESS_EQUAL, LESS), "")
	case '>':
		l.addToken(l.match('=', GREATER_EQUAL, GREATER), "")

	// Differentiate betweek SLASH and a comment (which ignores the rest of the line)
	case '/':
		if l.peek() == '/' {
			for l.peek() != '\n' && !l.isAtEnd() {
				l.advance()
			}
			break
		} else {
			l.addToken(SLASH, "")
		}

	// Whitespace and meanginless characters
	case '\n':
		l.line++
	case ' ':
	case '\r':
	case '\t':

	// Throw an error for unidentified characters
	default:
		l.throwError(fmt.Sprintf("Invalid character '%c'", char))
	}
}

// Loops over all characters in souce, creating tokens as it goes, places EOF token at end of source
func (l *lexer) ScanTokens() ([]Token, bool) {

	for !l.isAtEnd() {
		l.start = l.current
		l.scanToken()
	}

	l.tokens = append(l.tokens, Token{EOF, "EOF", "", l.line})
	return l.tokens, l.hadError
}