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
	OR
	THIS
	TRUE
	VAR
	WHILE

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
	"or" : OR,
	"this" : THIS,
	"true" : TRUE,
	"var" : VAR,
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
		return "STAR"
	case 10:
		return "SLASH"
	case 11:
		return "EQUAL"
	case 12:
		return "EQUAL_EQUAL"
	case 13:
		return "BANG"
	case 14:
		return "BANG_EQUAL"
	case 15:
		return "LESS"
	case 16:
		return "LESS_EQUAL"
	case 17:
		return "GREATER"
	case 18:
		return "GREATER_EQUAL"
	case 19:
		return "NUMBER"
	case 20:
		return "STRING"
	case 21:
		return "IDENTIFIER"
	case 22:
		return "AND"
	case 23:
		return "CLASS"
	case 24:
		return "ELSE"
	case 25:
		return "FALSE"
	case 26:
		return "FOR"
	case 27:
		return "FUNC"
	case 28:
		return "IF"
	case 29:
		return "OR"
	case 30:
		return "THIS"
	case 31:
		return "TRUE"
	case 32:
		return "VAR"
	case 33:
		return "WHILE"
	case 34:
		return "EOF"
	}

	return "Invalid TokenType"
}