package lexer

import (
	"github.com/jacksonopp/monkey/token"
	"testing"
)

type testToken struct {
	expectedType    token.TokenType
	expectedLiteral string
}

func TestNextToken(t *testing.T) {

	t.Run("basic input", func(t *testing.T) {
		input := `=+(){},;`

		tests := []testToken{
			{token.ASSIGN, "="},
			{token.PLUS, "+"},
			{token.LPAREN, "("},
			{token.RPAREN, ")"},
			{token.LBRACE, "{"},
			{token.RBRACE, "}"},
			{token.COMMA, ","},
			{token.SEMICOLON, ";"},
		}

		l := New(input)

		for i, tt := range tests {
			tok := l.NextToken()
			assertTokenIsExpected(t, tok, tt, i)
		}
	})

	t.Run("realistic monkey code", func(t *testing.T) {
		input := `let five = 5;

let ten = 10;

let add = fn(x, y) {
  x + y;
}

let result = add(five, ten);
`

		tests := []testToken{
			{token.LET, "let"},
			{token.IDENT, "five"},
			{token.ASSIGN, "="},
			{token.INT, "5"},
			{token.SEMICOLON, ";"},

			{token.LET, "let"},
			{token.IDENT, "ten"},
			{token.ASSIGN, "="},
			{token.INT, "10"},
			{token.SEMICOLON, ";"},

			{token.LET, "let"},
			{token.IDENT, "add"},
			{token.ASSIGN, "="},
			{token.FUNCTION, "fn"},
			{token.LPAREN, "("},
			{token.IDENT, "x"},
			{token.COMMA, ","},
			{token.IDENT, "y"},
			{token.RPAREN, ")"},
			{token.LBRACE, "{"},
			{token.IDENT, "x"},
			{token.PLUS, "+"},
			{token.IDENT, "y"},
			{token.SEMICOLON, ";"},
			{token.RBRACE, "}"},

			{token.LET, "let"},
			{token.IDENT, "result"},
			{token.ASSIGN, "="},
			{token.IDENT, "add"},
			{token.LPAREN, "("},
			{token.IDENT, "five"},
			{token.COMMA, ","},
			{token.IDENT, "ten"},
			{token.RPAREN, ")"},
			{token.SEMICOLON, ";"},
		}

		l := New(input)

		for i, tt := range tests {
			tok := l.NextToken()
			assertTokenIsExpected(t, tok, tt, i)
		}
	})
}

func assertTokenIsExpected(t *testing.T, tok token.Token, tt testToken, i int) {
	if tok.Type != tt.expectedType {
		t.Fatalf("test[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
	}
	if tok.Literal != tt.expectedLiteral {
		t.Fatalf("test[%d] - tokenliteral wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
	}
}
