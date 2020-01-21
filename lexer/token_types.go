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

	EOF
)

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
		return "EOF"
	}

	return "Invalid TokenType"
}