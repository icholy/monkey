package evaluator

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/icholy/monkey/ast"
	"github.com/icholy/monkey/object"
	"github.com/icholy/monkey/parser"
)

var (
	TRUE  = &object.Boolean{true}
	FALSE = &object.Boolean{false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node, env *object.Env) (object.Object, error) {
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
		fn, err := Eval(node.Function, env)
		if err != nil {
			return nil, err
		}
		var params []object.Object
		for _, arg := range node.Arguments {
			val, err := Eval(arg, env)
			if err != nil {
				return nil, err
			}
			params = append(params, val)
		}
		return applyFunction(fn, params)
	case *ast.FunctionStatement:
		env.Set(node.Name.Value, &object.Function{
			Parameters: node.Parameters,
			Body:       node.Body,
			Env:        env,
		})
		return NULL, nil
	case *ast.FunctionLiteral:
		return &object.Function{
			Parameters: node.Parameters,
			Body:       node.Body,
			Env:        env,
		}, nil
	case *ast.ImportStatement:
		return evalImport(node, env)
	case *ast.Identifier:
		return evalIdent(node, env)
	case *ast.WhileStatement:
		return evalWhile(node, env)
	case *ast.ReturnStatement:
		if node.ReturnValue == nil {
			return &object.ReturnValue{Value: NULL}, nil
		}
		val, err := Eval(node.ReturnValue, env)
		if err != nil {
			return nil, err
		}
		return &object.ReturnValue{Value: val}, nil
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}, nil
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}, nil
	case *ast.BooleanExpression:
		return boolToObject(node.Value), nil
	case *ast.PrefixExpression:
		right, err := Eval(node.Right, env)
		if err != nil {
			return nil, err
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left, err := Eval(node.Left, env)
		if err != nil {
			return nil, err
		}
		right, err := Eval(node.Right, env)
		if err != nil {
			return nil, err
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.LetStatement:
		val, err := Eval(node.Value, env)
		if err != nil {
			return nil, err
		}
		if node.Type != nil {
			if err := typeCheck(node.Type, val); err != nil {
				return nil, err
			}
			typ, _ := object.LookupType(node.Type.Value)
			val = &object.TypedObject{
				ObjectType: typ,
				Object:     val,
			}
		}
		env.Set(node.Name.Value, val)
		return NULL, nil
	case *ast.PropertyExpression:
		left, err := Eval(node.Value, env)
		if err != nil {
			return nil, err
		}
		return evalProperty(left, node.Name, env)
	case *ast.ArrayLiteral:
		return evalArray(node, env)
	case *ast.HashLiteral:
		return evalHash(node, env)
	case *ast.AssignmentExpression:
		value, err := Eval(node.Value, env)
		if err != nil {
			return nil, err
		}
		return evalAssign(node.Left, value, env)
	case *ast.IndexExpression:
		left, err := Eval(node.Value, env)
		if err != nil {
			return nil, err
		}
		index, err := Eval(node.Index, env)
		if err != nil {
			return nil, err
		}
		return evalIndex(left, index)
	default:
		return nil, fmt.Errorf("invalid node: %#v", node)
	}
}

func typeCheck(typeName *ast.Identifier, val object.Object) error {
	typ, ok := object.LookupType(typeName.Value)
	if !ok {
		return fmt.Errorf("invalid type name: %s", typeName)
	}
	if val.Type() != typ {
		return fmt.Errorf("wrong type: expected %s, got %s", typ, val.Type())
	}
	return nil
}

func applyFunction(fn object.Object, args []object.Object) (object.Object, error) {
	if builtin, ok := fn.(*object.Builtin); ok {
		return builtin.Fn(args...)
	}
	function, ok := fn.(*object.Function)
	if !ok {
		return nil, fmt.Errorf("not a function: %s", fn.Type())
	}
	if len(function.Parameters) != len(args) {
		return nil, fmt.Errorf("invalid number of function parameters")
	}
	env := object.NewEnv(function.Env)
	for i, param := range function.Parameters {
		if param.Type != nil {
			if err := typeCheck(param.Type, args[i]); err != nil {
				return nil, err
			}
		}
		env.Set(param.Name.Value, args[i])
	}
	val, err := Eval(function.Body, env)
	if err != nil {
		return nil, err
	}
	return object.UnwrapReturn(val), nil
}

func evalWhile(w *ast.WhileStatement, env *object.Env) (object.Object, error) {
	for {
		ok, err := Eval(w.Condition, env)
		if err != nil {
			return nil, err
		}
		if !isTruthy(ok) {
			break
		}
		if _, err := Eval(w.Body, env); err != nil {
			return nil, err
		}
	}
	return NULL, nil
}

func evalAssign(left ast.Expression, val object.Object, env *object.Env) (object.Object, error) {
	switch node := left.(type) {
	case *ast.Identifier:
		if err := env.Update(node.Value, val); err != nil {
			return nil, err
		}
		return NULL, nil
	case *ast.PropertyExpression:
		dest, err := Eval(node.Value, env)
		if err != nil {
			return nil, err
		}
		hash, ok := dest.(*object.Hash)
		if !ok {
			return nil, fmt.Errorf("cannot access property on %s", dest.Type())
		}
		hash.Set(&object.String{Value: node.Name.Value}, val)
		return NULL, nil
	case *ast.IndexExpression:
		dest, err := Eval(node.Value, env)
		if err != nil {
			return nil, err
		}
		index, err := Eval(node.Index, env)
		if err != nil {
			return nil, err
		}
		return evalAssignIndex(dest, index, val, env)
	default:
		return nil, fmt.Errorf("invalid assignment target")
	}
}

func evalAssignIndex(dest, index, val object.Object, env *object.Env) (object.Object, error) {
	switch obj := dest.(type) {
	case *object.Array:
		idx, ok := index.(*object.Integer)
		if !ok {
			return nil, fmt.Errorf("index must be an integer %s", index.Type())
		}
		obj.SetAt(int(idx.Value), val)
	case *object.Hash:
		obj.Set(index, val)
	default:
		return nil, fmt.Errorf("cannot index into %s", dest.Type())
	}
	return NULL, nil
}

func evalImport(i *ast.ImportStatement, env *object.Env) (object.Object, error) {
	if i.Program == nil {
		data, err := ioutil.ReadFile(i.Value)
		if err != nil {
			return nil, fmt.Errorf("import: %s", err)
		}
		p, err := parser.Parse(string(data))
		if err != nil {
			return nil, fmt.Errorf("import: %s", err)
		}
		i.Program = p
	}
	return Eval(i.Program, env)
}

func evalProperty(left object.Object, name *ast.Identifier, env *object.Env) (object.Object, error) {
	hash, ok := left.(*object.Hash)
	if !ok {
		return nil, fmt.Errorf("cannot access '%s' of %s", name, left.Type())
	}
	val, ok := hash.Get(&object.String{Value: name.Value})
	if !ok {
		return nil, fmt.Errorf("property not found: %s", name)
	}
	return val, nil
}

func evalIndex(left, index object.Object) (object.Object, error) {
	switch obj := left.(type) {
	case *object.Hash:
		if value, ok := obj.Get(index); ok {
			return value, nil
		}
		return NULL, nil
	case *object.Array:
		idx, ok := index.(*object.Integer)
		if !ok {
			return nil, fmt.Errorf("index must be an integer %s", index.Type())
		}
		return obj.At(int(idx.Value))
	case *object.String:
		idx, ok := index.(*object.Integer)
		if !ok {
			return nil, fmt.Errorf("index must be an integer %s", index.Type())
		}
		return obj.At(int(idx.Value))
	default:
		return nil, fmt.Errorf("cannot index into %s", left.Type())
	}
}

func evalHash(h *ast.HashLiteral, env *object.Env) (object.Object, error) {
	hash := object.NewHash()
	for _, p := range h.Pairs {
		key, err := Eval(p.Key, env)
		if err != nil {
			return nil, err
		}
		value, err := Eval(p.Value, env)
		if err != nil {
			return nil, err
		}
		hash.Set(key, value)
	}
	return hash, nil
}

func evalArray(a *ast.ArrayLiteral, env *object.Env) (object.Object, error) {
	var elements []object.Object
	for _, e := range a.Elements {
		val, err := Eval(e, env)
		if err != nil {
			return nil, err
		}
		elements = append(elements, val)
	}
	return &object.Array{Elements: elements}, nil
}

func evalIdent(i *ast.Identifier, env *object.Env) (object.Object, error) {
	if val, ok := builtins[i.Value]; ok {
		return val, nil
	}
	if val, ok := env.Get(i.Value); ok {
		return val, nil
	}
	return nil, fmt.Errorf("identifier not found: %s", i.Value)
}

func evalProgram(p *ast.Program, env *object.Env) (object.Object, error) {
	var (
		last object.Object = NULL
		err  error
	)
	for _, stmt := range p.Statements {
		last, err = Eval(stmt, env)
		if err != nil {
			return nil, err
		}
		if last.Type() == object.RETURN {
			return object.UnwrapReturn(last), nil
		}
	}
	return last, nil
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Env) (object.Object, error) {
	var (
		last object.Object = NULL
		err  error
	)
	for _, stmt := range block.Statements {
		last, err = Eval(stmt, env)
		if err != nil {
			return nil, err
		}
		if last.Type() == object.RETURN {
			return last, nil
		}
	}
	return last, nil
}

func evalInfixExpression(operator string, left, right object.Object) (object.Object, error) {
	if left.Type() != right.Type() {
		return nil, fmt.Errorf("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	}
	switch left.Type() {
	case object.INTEGER:
		return evalIntegerInfixExpression(operator, left.(*object.Integer), right.(*object.Integer))
	case object.STRING:
		return evalStringInfixExpression(operator, left.(*object.String), right.(*object.String))
	}

	switch operator {
	case "==":
		return boolToObject(left == right), nil
	case "!=":
		return boolToObject(left != right), nil
	default:
		return nil, fmt.Errorf("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right *object.String) (object.Object, error) {
	switch operator {
	case "+":
		return &object.String{Value: left.Value + right.Value}, nil
	case "==":
		return &object.Boolean{Value: left.Value == right.Value}, nil
	case "!=":
		return &object.Boolean{Value: left.Value != right.Value}, nil
	case ">":
		return &object.Boolean{Value: strings.Compare(left.Value, right.Value) > 0}, nil
	case "<":
		return &object.Boolean{Value: strings.Compare(left.Value, right.Value) < 0}, nil
	case "<=":
		return &object.Boolean{Value: strings.Compare(left.Value, right.Value) <= 0}, nil
	case ">=":
		return &object.Boolean{Value: strings.Compare(left.Value, right.Value) >= 0}, nil
	default:
		return nil, fmt.Errorf("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right *object.Integer) (object.Object, error) {
	switch operator {
	case "+":
		return &object.Integer{Value: left.Value + right.Value}, nil
	case "-":
		return &object.Integer{Value: left.Value - right.Value}, nil
	case "*":
		return &object.Integer{Value: left.Value * right.Value}, nil
	case "/":
		return &object.Integer{Value: left.Value / right.Value}, nil
	case "<":
		return boolToObject(left.Value < right.Value), nil
	case ">":
		return boolToObject(left.Value > right.Value), nil
	case ">=":
		return boolToObject(left.Value >= right.Value), nil
	case "<=":
		return boolToObject(left.Value <= right.Value), nil
	case "==":
		return boolToObject(left.Value == right.Value), nil
	case "!=":
		return boolToObject(left.Value != right.Value), nil
	default:
		return nil, fmt.Errorf("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case FALSE, NULL:
		return false
	default:
		return true
	}
}

func evalIfExpression(i *ast.IfExpression, env *object.Env) (object.Object, error) {
	cond, err := Eval(i.Condition, env)
	if err != nil {
		return nil, err
	}
	if isTruthy(cond) {
		return Eval(i.Concequence, env)
	}
	if i.Alternative != nil {
		return Eval(i.Alternative, env)
	}
	return NULL, nil
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

func evalPrefixExpression(operator string, right object.Object) (object.Object, error) {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return nil, fmt.Errorf("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) (object.Object, error) {
	switch right {
	case TRUE:
		return FALSE, nil
	case FALSE:
		return TRUE, nil
	case NULL:
		return TRUE, nil
	default:
		return FALSE, nil
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) (object.Object, error) {
	if obj, ok := right.(*object.Integer); ok {
		return &object.Integer{Value: -obj.Value}, nil
	}
	return nil, fmt.Errorf("unknown operator: -%s", right.Type())
}

func boolToObject(b bool) object.Object {
	if b {
		return TRUE
	}
	return FALSE
}
