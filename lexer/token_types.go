package lexer

// Declare constant names for TokenTypes, use iota to increment values
type TokenType int

const (
	// Single-character
	PLUS TokenType = iota
	MINUS
	LEFT_PAREN
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	SEMICOLON
	COLON
	STAR
	SLASH

	// One or two characters
	EQUAL
	EQUAL_EQUAL
	BANG
	BANG_EQUAL
	LESS
	LESS_EQUAL
	GREATER
	GREATER_EQUAL

	// Literals
	NUMBER
	STRING
	IDENTIFIER

	// Reserved keywords
	AND
	CLASS
	ELSE
	FALSE
	FOR
	FUNC
	IF
	NIL
	OR
	PRINT
	THIS
	TRUE
	LET
	WHILE

	INDENT
	DEDENT
	NEWLINE
	EOF
)

var keywords = map[string]TokenType{
	"and" : AND,
	"class" : CLASS,
	"else" : ELSE,
	"false" : FALSE,
	"for" : FOR,
	"func" : FUNC,
	"if" : IF,
	"nil" : NIL,
	"or" : OR,
	"print" : PRINT,
	"this" : THIS,
	"true" : TRUE,
	"let" : LET,
	"while" : WHILE,
}


// Ugly, ugly method to retun string type names from TokenType constants (I hate it)
func (t TokenType) typeString() string {
	switch t {
	case 0:
		return "PLUS"
	case 1:
		return "MINUS"
	case 2:
		return "LEFT_PAREN"
	case 3:
		return "RIGHT_PAREN"
	case 4:
		return "LEFT_BRACE"
	case 5:
		return "RIGHT_BRACE"
	case 6:
		return "COMMA"
	case 7:
		return "DOT"
	case 8:
		return "SEMICOLON"
	case 9:
		return "COLON"
	case 10:
		return "STAR"
	case 11:
		return "SLASH"
	case 12:
		return "EQUAL"
	case 13:
		return "EQUAL_EQUAL"
	case 14:
		return "BANG"
	case 15:
		return "BANG_EQUAL"
	case 16:
		return "LESS"
	case 17:
		return "LESS_EQUAL"
	case 18:
		return "GREATER"
	case 19:
		return "GREATER_EQUAL"
	case 20:
		return "NUMBER"
	case 21:
		return "STRING"
	case 22:
		return "IDENTIFIER"
	case 23:
		return "AND"
	case 24:
		return "CLASS"
	case 25:
		return "ELSE"
	case 26:
		return "FALSE"
	case 27:
		return "FOR"
	case 28:
		return "FUNC"
	case 29:
		return "IF"
	case 30:
		return "NIL"
	case 31:
		return "OR"
	case 32:
		return "PRINT"
	case 33:
		return "THIS"
	case 34:
		return "TRUE"
	case 35:
		return "LET"
	case 36:
		return "WHILE"
	case 37:
		return "INDENT"
	case 38:
		return "DEDENT"
	case 39:
		return "NEWLINE"
	case 40:
		return "EOF"
	}

	return "Invalid TokenType"
}