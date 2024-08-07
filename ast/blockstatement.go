package ast

import (
	"bytes"
	"github.com/jacksonopp/monkey/token"
)

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (b BlockStatement) TokenLiteral() string {
	return b.Token.Literal
}

func (b BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range b.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

func (b BlockStatement) statementNode() {

}
