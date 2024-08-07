package ast

import "github.com/jacksonopp/monkey/token"

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}

func (i IntegerLiteral) String() string {
	return i.TokenLiteral()
}

func (i IntegerLiteral) expressionNode() {
}
