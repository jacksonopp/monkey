package evaluator

import (
	"github.com/jacksonopp/monkey/ast"
	"github.com/jacksonopp/monkey/object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	//statements
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	//	expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	}
	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)
	}

	return result
}
