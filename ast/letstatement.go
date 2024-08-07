package ast

import (
	"bytes"
	"github.com/jacksonopp/monkey/token"
)

// LetStatement
// ex: `let x = 3;`
type LetStatement struct {
	Token token.Token // token.LET
	Name  *Identifier // Identifier of the binding (x)
	Value Expression  // Expression that produces the value (3)
}

func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")
	return out.String()
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
