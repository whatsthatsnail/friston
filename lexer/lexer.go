package lexer

import (
	"fmt";
	"strconv";
	"github.com/whatsthatsnail/simple_interpreter/errors"
)

// Simple helper functions to avoid importing a whole module for a one-liner
func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

// Return underscore as alpha to allow '_' in idenifiers and keywords
func isAlpha(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r == '_')
}

// Literals stores as empty interface, use type assertions when parsing
type Token struct {
	TType TokenType
	Lexeme string
	Literal interface{}
	Line int
}

// Print an instance of a token.
func (tok Token) Print() {
	if tok.Literal == nil {
		fmt.Printf("{%s, '%s', %d}\n", tok.TType.typeString(), tok.Lexeme, tok.Line)
	} else {
		fmt.Printf("{%s, '%s', %v, %d}\n", tok.TType.typeString(), tok.Lexeme, tok.Literal, tok.Line)
	}
}

// Prints a list of tokens in a readable manner as {Token_Type, lexeme, (literal), line}
func PrintTokens(tokens []Token) {
	for _, tok := range(tokens) {
		tok.Print()
	}
}

type lexer struct {
	start int
	current int
	line int
	tokens []Token
	source string
	hadError bool
	repl bool
	depth int
}

// Lexer constructor, initializes default values
func NewLexer(code string, replFlag bool) lexer {
	l := lexer{}
	l.start = 0
	l.current = 0
	l.line = 1
	l.source = code
	l.hadError = false
	l.repl = replFlag
	l.depth = 0

	return l
}

// Error handling:
func (l *lexer) throwError(message string) {
	errors.ThrowError(l.line, message)
	l.hadError = true
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

// Returns the next character without consuming it, as long as there is another character peek
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
func (l *lexer) addToken(tType TokenType, literal interface{}) {
	if l.repl == true {
		l.tokens = append(l.tokens, Token{tType, l.source[l.start:l.current], literal, 0})
	} else {
		l.tokens = append(l.tokens, Token{tType, l.source[l.start:l.current], literal, l.line})
	}
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

	floatFlag := false

	// Advance until end of number (while peek == a digit or a dot, as long as there's another number after the dot)
	for (isDigit(l.peek()) || (l.peek() == '.' && isDigit(l.peekNext()))) && !l.isAtEnd() {
		if l.peek() == '.' {floatFlag = true}
		l.advance()
	}

	// Store a NUMBER token
	number := l.source[l.start:l.current]
	if floatFlag {
		num, _ := strconv.ParseFloat(number, 10)
		l.addToken(NUMBER, num)
	} else {
		num, _ := strconv.ParseFloat(number, 10)
		l.addToken(NUMBER, num)
	}
}

func (l *lexer) getWord() {
	// Advance to en of word
	for (isAlpha(l.peek()) && !l.isAtEnd()) {
		l.advance()
	}

	// Store word lexeme
	word := l.source[l.start:l.current]

	if keywords[word] == NIL {
		l.addToken(NIL, nil)
	} else if keywords[word] == 0 {
		l.addToken(IDENTIFIER, word)
	} else {
		l.addToken(keywords[word], word)
	}
}

func (l *lexer) countIndent() int {
	count := 0
	for !l.isAtEnd() && l.peek() == ' ' {
		count++
		l.advance()
	}

	if count % 4 == 0 {
		count = count / 4
		return count
	} else {
		l.throwError("Indents must be four spaces")
	}

	return count
}

// Add INDENT, DEDENT, and adjust l.depth accordin to space count.
func (l *lexer) getDent() {
	count := l.countIndent()
	difference := count - l.depth
	
	if difference > 0 {
		for i := 0; i < difference; i++ {
			l.tokens = append(l.tokens, Token{INDENT, "", nil, l.line})
		}
	} else if difference < 0 {
		for i := 0; i < -difference; i++ {
			l.tokens = append(l.tokens, Token{DEDENT, "", nil, l.line})
		}
	}

	l.depth = count
}

func (l *lexer) getNewline() {
	previousToken := l.tokens[len(l.tokens) - 1]
	
	// Only append NEWLINE if the previous character is not a newline or keyword
	_, ok := keywords[previousToken.Lexeme]
	if !ok && previousToken.TType != NEWLINE && previousToken.TType != SEMICOLON {
		l.tokens = append(l.tokens, Token{NEWLINE, "", nil, l.line})
	}
}

// Advances current and adds the current token
func (l *lexer) scanToken() {
	char := l.advance()

	switch char {

	// Creates single-character tokens
	case '+':
		l.addToken(PLUS, nil)
	case '-':
		l.addToken(MINUS, nil)
	case '(':
		l.addToken(LEFT_PAREN, nil)
	case ')':
		l.addToken(RIGHT_PAREN, nil)
	case '{':
		l.addToken(LEFT_BRACE, nil)
	case '}':
		l.addToken(RIGHT_BRACE, nil)
	case ',':
		l.addToken(COMMA, nil)
	case ';':
		l.addToken(SEMICOLON, nil)
	case ':':
		l.addToken(COLON, nil)
	case '.':
		l.addToken(DOT, nil)
	case '*':
		l.addToken(STAR, nil)

	// Create one or two-character tokens
	case '=':
		l.addToken(l.match('=', EQUAL_EQUAL, EQUAL), nil)
	case '!':
		l.addToken(l.match('=', BANG_EQUAL, BANG), nil)
	case '<':
		l.addToken(l.match('=', LESS_EQUAL, LESS), nil)
	case '>':
		l.addToken(l.match('=', GREATER_EQUAL, GREATER), nil)

	// Differentiate betweek SLASH and a comment (which ignores the rest of the line)
	case '/':
		if l.peek() == '/' {
			for l.peek() != '\n' && !l.isAtEnd() {
				l.advance()
			}

			if l.peek() == '\n' {
				l.advance()
			}

			break
		} else {
			l.addToken(SLASH, nil)
		}

	// Whitespace and meanginless characters
	case '\n':
		l.getNewline()
		l.line++
		if l.peek() != '\n' {
			l.getDent()
		}
	case '~':
		// Skip a newline if it's preceded by a '~' to allow a statement to continue to a new line of text.
		if l.peek() == '\n' {
			l.advance()
		}
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
	l.getDent()

	for !l.isAtEnd() {
		l.start = l.current
		l.scanToken()
	}

	l.getNewline()
	l.tokens = append(l.tokens, Token{EOF, "EOF", nil, l.line})
	return l.tokens, l.hadError
}