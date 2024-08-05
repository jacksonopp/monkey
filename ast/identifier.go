package ast

import "github.com/jacksonopp/monkey/token"

type Identifier struct {
	Token token.Token // token.IDENT
	Value string      // The value
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
