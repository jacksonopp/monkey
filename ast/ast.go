package ast

// Node is the root interface that all AST nodes implement
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement is a unit of execution
//
// ex. let x = 3;
// ex. let foo = fn(a, b) {a + b}
type Statement interface {
	Node // all Statement's implement Node
	statementNode()
}

// Expression produces a value when executed
//
// ex. 3 + 4
type Expression interface {
	Node
	expressionNode()
}

type TypeChecker interface {
	IsType() bool
}
