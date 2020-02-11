package lexer

// Declare constant names for TokenTypes, use iota to increment values
type TokenType int

const (
	// Single-character
	LEFT_PAREN TokenType = iota
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
	PLUS
	PLUS_PLUS
	MINUS
	MINUS_MINUS
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
	THEN
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
	"then" : THEN,
	"this" : THIS,
	"true" : TRUE,
	"let" : LET,
	"while" : WHILE,
}


// Ugly, ugly method to retun string type names from TokenType constants (I hate it)
func (t TokenType) typeString() string {
	switch t {
	case 0:
		return "LEFT_PAREN"
	case 1:
		return "RIGHT_PAREN"
	case 2:
		return "LEFT_BRACE"
	case 3:
		return "RIGHT_BRACE"
	case 4:
		return "COMMA"
	case 5:
		return "DOT"
	case 6:
		return "SEMICOLON"
	case 7:
		return "COLON"
	case 8:
		return "STAR"
	case 9:
		return "SLASH"
	case 10:
		return "PLUS"
	case 11:
		return "PLUS_PLUS"
	case 12:
		return "MINUS"
	case 13:
		return "MINUS_MINUS"
	case 14:
		return "EQUAL"
	case 15:
		return "EQUAL_EQUAL"
	case 16:
		return "BANG"
	case 17:
		return "BANG_EQUAL"
	case 18:
		return "LESS"
	case 19:
		return "LESS_EQUAL"
	case 20:
		return "GREATER"
	case 21:
		return "GREATER_EQUAL"
	case 22:
		return "NUMBER"
	case 23:
		return "STRING"
	case 24:
		return "IDENTIFIER"
	case 25:
		return "AND"
	case 26:
		return "CLASS"
	case 27:
		return "ELSE"
	case 28:
		return "FALSE"
	case 29:
		return "FOR"
	case 30:
		return "FUNC"
	case 31:
		return "IF"
	case 32:
		return "NIL"
	case 33:
		return "OR"
	case 34:
		return "THEN"
	case 35:
		return "THIS"
	case 36:
		return "TRUE"
	case 37:
		return "LET"
	case 38:
		return "WHILE"
	case 39:
		return "INDENT"
	case 40:
		return "DEDENT"
	case 41:
		return "NEWLINE"
	case 42:
		return "EOF"
	}

	return "Invalid TokenType"
}