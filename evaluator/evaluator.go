package evaluator

import (
	"strings"

	"github.com/icholy/monkey/ast"
	"github.com/icholy/monkey/object"
)

var (
	TRUE  = &object.Boolean{true}
	FALSE = &object.Boolean{false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node, env *object.Env) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.CallExpression:
		fn := Eval(node.Function, env)
		if isError(fn) {
			return fn
		}
		var params []object.Object
		for _, arg := range node.Arguments {
			val := Eval(arg, env)
			if isError(val) {
				return val
			}
			params = append(params, val)
		}
		return applyFunction(fn, params)
	case *ast.FunctionLiteral:
		return &object.Function{
			Parameters: node.Parameters,
			Body:       node.Body,
			Env:        env,
		}
	case *ast.Identifier:
		return evalIdent(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.BooleanExpression:
		return boolToObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return NULL
	case *ast.ArrayLiteral:
		return evalArray(node, env)
	case *ast.HashLiteral:
		return evalHash(node, env)
	case *ast.IndexExpression:
		left := Eval(node.Value, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndex(left, index)
	default:
		return object.Errorf("invalid node: %#v", node)
	}
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	if builtin, ok := fn.(*object.Builtin); ok {
		return builtin.Fn(args...)
	}
	function, ok := fn.(*object.Function)
	if !ok {
		return object.Errorf("not a function: %s", fn.Type())
	}
	env := object.NewEnv(function.Env)
	for i, param := range function.Parameters {
		env.Set(param.Value, args[i])
	}
	val := Eval(function.Body, env)
	if ret, ok := val.(*object.ReturnValue); ok {
		return ret.Value
	}
	return val
}

func evalIndex(left, index object.Object) object.Object {
	if hash, ok := left.(*object.Hash); ok {
		if value, ok := hash.Get(index); ok {
			return value
		}
		return NULL
	}
	arr, ok := left.(*object.Array)
	if !ok {
		return object.Errorf("cannot index into %s", left.Type())
	}
	idx, ok := index.(*object.Integer)
	if !ok {
		return object.Errorf("index must be an integer %s", index.Type())
	}
	if idx.Value < 0 || idx.Value >= int64(len(arr.Elements)) {
		return object.Errorf("index out of range %d", idx.Value)
	}
	return arr.Elements[idx.Value]
}

func evalHash(h *ast.HashLiteral, env *object.Env) object.Object {
	hash := object.NewHash()
	for _, p := range h.Pairs {
		key := Eval(p.Key, env)
		if isError(key) {
			return key
		}
		value := Eval(p.Value, env)
		if isError(value) {
			return value
		}
		hash.Set(key, value)
	}
	return hash
}

func evalArray(a *ast.ArrayLiteral, env *object.Env) object.Object {
	var elements []object.Object
	for _, e := range a.Elements {
		val := Eval(e, env)
		if isError(val) {
			return val
		}
		elements = append(elements, val)
	}
	return &object.Array{Elements: elements}
}

func evalIdent(i *ast.Identifier, env *object.Env) object.Object {
	if val, ok := builtins[i.Value]; ok {
		return val
	}
	if val, ok := env.Get(i.Value); ok {
		return val
	}
	return object.Errorf("identifier not found: %s", i.Value)
}

func evalProgram(p *ast.Program, env *object.Env) object.Object {
	var last object.Object
	for _, stmt := range p.Statements {
		last = Eval(stmt, env)
		if isError(last) {
			return last
		}
		if ret, ok := last.(*object.ReturnValue); ok {
			return ret.Value
		}
	}
	return last
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Env) object.Object {
	var last object.Object
	for _, stmt := range block.Statements {
		last = Eval(stmt, env)
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
	if left.Type() != right.Type() {
		return object.Errorf("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	}
	switch left.Type() {
	case object.INTEGER:
		return evalIntegerInfixExpression(operator, left.(*object.Integer), right.(*object.Integer))
	case object.STRING:
		return evalStringInfixExpression(operator, left.(*object.String), right.(*object.String))
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

func evalStringInfixExpression(operator string, left, right *object.String) object.Object {
	switch operator {
	case "+":
		return &object.String{Value: left.Value + right.Value}
	case "==":
		return &object.Boolean{Value: left.Value == right.Value}
	case "!=":
		return &object.Boolean{Value: left.Value != right.Value}
	case ">":
		return &object.Boolean{Value: strings.Compare(left.Value, right.Value) > 0}
	case "<":
		return &object.Boolean{Value: strings.Compare(left.Value, right.Value) < 0}
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

func evalIfExpression(i *ast.IfExpression, env *object.Env) object.Object {
	cond := Eval(i.Condition, env)
	if isError(cond) {
		return cond
	}
	if isTruthy(cond) {
		return Eval(i.Concequence, env)
	}
	if i.Alternative != nil {
		return Eval(i.Alternative, env)
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
