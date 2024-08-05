package ast

import "github.com/jacksonopp/monkey/token"

// LetStatement
// ex: `let x = 3;`
type LetStatement struct {
	Token token.Token // token.LET
	Name  *Identifier // Identifier of the binding (x)
	Value Expression  // Expression that produces the value (3)
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
