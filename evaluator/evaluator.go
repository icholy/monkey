package evaluator

import (
	"github.com/icholy/monkey/ast"
	"github.com/icholy/monkey/object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.BlockStatement:
		return evalStatements(node.Statements)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.BooleanExpression:
		return &object.Boolean{Value: node.Value}
	default:
		return &object.Null{}
	}
}

func evalStatements(stmts []ast.Statement) object.Object {
	var last object.Object
	for _, stmt := range stmts {
		last = Eval(stmt)
	}
	return last
}
