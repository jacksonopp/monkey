package lexer

import (
	"github.com/jacksonopp/monkey/token"
)

type Lexer struct {
	input        string
	position     int  // the index current position being read
	readPosition int  // the index of the next position to read
	ch           byte // the value of the current position being read
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	default:
		return l.handleIdentifier()
	}

	l.readChar()

	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) handleIdentifier() token.Token {
	var tok token.Token

	if isLetter(l.ch) {
		tok.Literal = l.readIdentifier()
		tok.Type = token.LookupIdent(tok.Literal)
		return tok
	}

	if isDigit(l.ch) {
		tok.Type = token.INT
		tok.Literal = l.readNumber()

		return tok
	}

	return newToken(token.ILLEGAL, l.ch)
}

func (l *Lexer) readNumber() string {
	pos := l.position

	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[pos:l.position]
}

// readIdentifier will parse an entire identifier
func (l *Lexer) readIdentifier() string {
	pos := l.position

	// advance till you hit a non-letter character
	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[pos:l.position]
}

// readChar gives us the next character and advances the position of the input string
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1
}
