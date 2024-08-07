package ast

import (
	"bytes"
	"github.com/jacksonopp/monkey/token"
)

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (i IfExpression) TokenLiteral() string {
	return i.Token.Literal
}

func (i IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(i.Condition.String())
	out.WriteString(" ")
	out.WriteString(i.Consequence.String())
	out.WriteString(" ")

	if i.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(i.Alternative.String())
	}

	return out.String()
}

func (i IfExpression) expressionNode() {

}
