package lexer

import (
	"fmt"
)

// Simple helper functions to avoid importing a whole module for a one-liner
func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

// Return underscore as alpha to allow '_' in idenifiers and keywords
func isAlpha(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r == '_')
}

// Note: number literals are still stored as a string, cast as int before use.
// TODO: Develop a more elegant solution for the above issue.
type Token struct {
	tType TokenType
	lexeme string
	literal string
	line int
}

// Prints tokens in a readable manner as {Token_Type, lexeme, (literal), line}
func PrintTokens(tokens []Token) {
	for _, tok := range(tokens) {
		s := ""
		if tok.literal == "" {
			s = fmt.Sprintf("{%s, '%s', %d}", tok.tType.typeString(), tok.lexeme, tok.line)
		} else {
			s = fmt.Sprintf("{%s, '%s', %s, %d}", tok.tType.typeString(), tok.lexeme, tok.literal, tok.line)
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

// Print any line dependant message (error, warning, etc.)
func (l *lexer) report(line int, message string) {
	fmt.Printf("[Line %d] %s\n", line, message)
}

// Checks if current position has reaced the end of the source
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

// Returns the next character without conuming it, as long as there is anothe character peek
func (l *lexer) peekNext() rune {
	if !l.isAtEnd() && !((l.current + 1) > len(l.source) - 1) {
		return rune(l.source[l.current + 1])
	} else {
		return '\n'
	}
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
func (l *lexer) addToken(tType TokenType, literal string) {
	l.tokens = append(l.tokens, Token{tType, l.source[l.start:l.current], literal, l.line})
}

// Consumes a string literal, including new lines, and creates a STRING token
func (l *lexer) getString() {
	for l.peek() != '"' && !l.isAtEnd() {
		// Handle multi-line comments
		if l.peek() == '\n' {
			l.line++
			l.advance()
		}
		// Advance through all characters until reaching the terminating "
		l.advance()
	}

	// If we haven't reached the end of l.source, but find terminating "
	if l.peek() == '"' {
		// Consume the " and store the STRING token
		l.advance()
		literal := l.source[l.start + 1 : l.current - 1]
		l.addToken(STRING, literal)
	// If there's no terminating ", throw an error
	} else {
		l.throwError("Unterminated string")
	}
}

// Consumes number literal and creates a NUMBER token
func (l *lexer) getNumber() {
	// Advance until end of number (while peek == a digit or a dot, as long as there's another number after the dot)
	for (isDigit(l.peek()) || (l.peek() == '.' && isDigit(l.peekNext()))) && !l.isAtEnd() {
		l.advance()
	}

	// Store a NUMBER token
	number := l.source[l.start:l.current]
	l.addToken(NUMBER, number)
}

func (l *lexer) getWord() {
	// Advance to en of word
	for (isAlpha(l.peek()) && !l.isAtEnd()) {
		l.advance()
	}

	// Store word lexeme
	word := l.source[l.start:l.current]

	switch keywords[word] {
	case AND:
		l.addToken(AND, word)
	case CLASS:
		l.addToken(CLASS, word)
	case ELSE:
		l.addToken(ELSE, word)
	case FALSE:
		l.addToken(FALSE, word)
	case FOR:
		l.addToken(FOR, word)
	case FUNC:
		l.addToken(FUNC, word)
	case IF:
		l.addToken(IF, word)
	case OR:
		l.addToken(OR, word)
	case THIS:
		l.addToken(THIS, word)
	case TRUE:
		l.addToken(TRUE, word)
	case VAR:
		l.addToken(VAR, word)
	case WHILE:
		l.addToken(WHILE, word)
	default:
		l.addToken(IDENTIFIER, word)
	}
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

	// String literals
	case '"':
		l.getString()

	// Check for literals without an identifying characer (numbers and words)
	// Throw an error for unidentified characters
	default:
		if isDigit(char) {
			l.getNumber()
		} else if isAlpha(char){
			l.getWord()
		} else {
			l.throwError(fmt.Sprintf("Invalid character '%c'", char))
		}
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