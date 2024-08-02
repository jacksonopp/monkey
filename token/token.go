package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// IDENTIFIERS + LITERALS
	IDENT = "IDENT"
	INT   = "INT"

	// OPERATORS
	ASSIGN = "="
	PLUS   = "+"

	//	DELIMITERS
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"

	// KEYWORDS
	FUNCTION = "FUNCTION"
	LET      = "LET"
)

var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

// LookupIdent get's the TokenType based on the Token's Literal value
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
