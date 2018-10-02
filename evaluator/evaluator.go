package evaluator

import (
	"github.com/icholy/monkey/ast"
	"github.com/icholy/monkey/object"
)

var (
	TRUE  = &object.Boolean{true}
	FALSE = &object.Boolean{false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.BlockStatement:
		return evalBlockStatement(node)
	case *ast.IfExpression:
		return evalIfExpression(node)
	case *ast.ReturnStatement:
		obj := Eval(node.ReturnValue)
		if isError(obj) {
			return obj
		}
		return &object.ReturnValue{Value: obj}
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.BooleanExpression:
		return boolToObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		if isError(left) {
			return left
		}
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	default:
		return NULL
	}
}

func evalProgram(p *ast.Program) object.Object {
	var last object.Object
	for _, stmt := range p.Statements {
		last = Eval(stmt)
		if isError(last) {
			return last
		}
		if ret, ok := last.(*object.ReturnValue); ok {
			return ret.Value
		}
	}
	return last
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var last object.Object
	for _, stmt := range block.Statements {
		last = Eval(stmt)
		if isError(last) {
			return last
		}
		if last.Type() == object.RETURN {
			return last
		}
	}
	return last
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	if left.Type() == object.INTEGER && right.Type() == object.INTEGER {
		return evalIntegerInfixExpression(operator, left.(*object.Integer), right.(*object.Integer))
	}
	if left.Type() != right.Type() {
		return object.Errorf("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	}
	switch operator {
	case "==":
		return boolToObject(left == right)
	case "!=":
		return boolToObject(left != right)
	default:
		return object.Errorf("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right *object.Integer) object.Object {
	switch operator {
	case "+":
		return &object.Integer{Value: left.Value + right.Value}
	case "-":
		return &object.Integer{Value: left.Value - right.Value}
	case "*":
		return &object.Integer{Value: left.Value * right.Value}
	case "/":
		return &object.Integer{Value: left.Value / right.Value}
	case "<":
		return boolToObject(left.Value < right.Value)
	case ">":
		return boolToObject(left.Value > right.Value)
	case "==":
		return boolToObject(left.Value == right.Value)
	case "!=":
		return boolToObject(left.Value != right.Value)
	default:
		return object.Errorf("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func isError(obj object.Object) bool {
	_, ok := obj.(*object.Error)
	return ok
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case FALSE, NULL:
		return false
	default:
		return true
	}
}

func evalIfExpression(i *ast.IfExpression) object.Object {
	cond := Eval(i.Condition)
	if isError(cond) {
		return cond
	}
	if isTruthy(cond) {
		return Eval(i.Concequence)
	}
	if i.Alternative != nil {
		return Eval(i.Alternative)
	}
	return NULL
}

func evalPlusOperator(left, right object.Object) object.Object {
	leftVal, ok := left.(*object.Integer)
	if !ok {
		return NULL
	}
	rightVal, ok := right.(*object.Integer)
	if !ok {
		return NULL
	}
	return &object.Integer{
		Value: leftVal.Value + rightVal.Value,
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return object.Errorf("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if obj, ok := right.(*object.Integer); ok {
		return &object.Integer{Value: -obj.Value}
	}
	return object.Errorf("unknown operator: -%s", right.Type())
}

func boolToObject(b bool) object.Object {
	if b {
		return TRUE
	}
	return FALSE
}
