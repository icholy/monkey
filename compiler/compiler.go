package compiler

import (
	"fmt"

	"github.com/icholy/monkey/ast"
	"github.com/icholy/monkey/code"
	"github.com/icholy/monkey/object"
)

type Instruction struct {
	Opcode   code.Opcode
	Position int
}

func (i Instruction) Is(op code.Opcode) bool {
	return i.Opcode == op
}

type Compiler struct {
	constants []object.Object
	symbols   *SymbolTable

	scopes []*Scope
}

func New() *Compiler {
	return &Compiler{
		symbols: NewSymbolTable(nil),
		scopes: []*Scope{
			&Scope{},
		},
	}
}

type Scope struct {
	instructions code.Instructions
	prev         Instruction
	prevprev     Instruction
}

func (s *Scope) undo() {
	s.instructions = s.instructions[:s.prev.Position]
	s.prev = s.prevprev
}

func (s *Scope) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := len(s.instructions)
	s.instructions = append(s.instructions, ins...)
	s.prevprev = s.prev
	s.prev = Instruction{
		Opcode:   op,
		Position: pos,
	}
	return pos
}

func (s *Scope) rewrite(pos int, op code.Opcode, operands ...int) {
	ins := s.instructions
	for i, x := range code.Make(op, operands...) {
		ins[pos+i] = x
	}
}

func NewWithState(symbols *SymbolTable, constants []object.Object) *Compiler {
	c := New()
	c.symbols = symbols
	c.constants = c.constants
	return c
}

func Compile(node ast.Node) (*Bytecode, error) {
	c := New()
	if err := c.Compile(node); err != nil {
		return nil, err
	}
	return c.Bytecode(), nil
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			if err := c.Compile(s); err != nil {
				return err
			}
		}
	case *ast.ExpressionStatement:
		if err := c.Compile(node.Expression); err != nil {
			return err
		}
		c.emit(code.OpPop)
	case *ast.BlockStatement:
		for _, s := range node.Statements {
			if err := c.Compile(s); err != nil {
				return err
			}
		}
	case *ast.InfixExpression:
		if node.Operator == "<" {
			if err := c.Compile(node.Right); err != nil {
				return err
			}
			if err := c.Compile(node.Left); err != nil {
				return err
			}
			c.emit(code.OpGreaterThan)
			return nil
		}
		if err := c.Compile(node.Left); err != nil {
			return err
		}
		if err := c.Compile(node.Right); err != nil {
			return err
		}
		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		case ">":
			c.emit(code.OpGreaterThan)
		case "==":
			c.emit(code.OpEqual)
		case "!=":
			c.emit(code.OpNotEqual)
		default:
			return fmt.Errorf("unknown operator: %s", node.Operator)
		}
	case *ast.IfExpression:
		scope := c.scope()
		if err := c.Compile(node.Condition); err != nil {
			return err
		}
		jumpNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999)
		if err := c.Compile(node.Concequence); err != nil {
			return err
		}
		if scope.prev.Is(code.OpPop) {
			scope.undo()
		}

		jumpPos := c.emit(code.OpJump, 9999)
		c.rewrite(jumpNotTruthyPos, code.OpJumpNotTruthy, len(c.instructions()))

		if node.Alternative != nil {
			if err := c.Compile(node.Alternative); err != nil {
				return err
			}
			if scope.prev.Is(code.OpPop) {
				scope.undo()
			}
		} else {
			c.emit(code.OpNull)
		}

		c.rewrite(jumpPos, code.OpJump, len(c.instructions()))

		return nil
	case *ast.PrefixExpression:
		if err := c.Compile(node.Right); err != nil {
			return err
		}
		switch node.Operator {
		case "-":
			c.emit(code.OpMinus)
		case "!":
			c.emit(code.OpBang)
		default:
			return fmt.Errorf("unknown operator: %s", node.Operator)
		}
	case *ast.IntegerLiteral:
		v := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(v))
	case *ast.StringLiteral:
		v := &object.String{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(v))
	case *ast.BooleanExpression:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}
	case *ast.ArrayLiteral:
		for _, e := range node.Elements {
			if err := c.Compile(e); err != nil {
				return err
			}
		}
		c.emit(code.OpArray, len(node.Elements))
	case *ast.HashLiteral:
		for _, p := range node.Pairs {
			if err := c.Compile(p.Key); err != nil {
				return err
			}
			if err := c.Compile(p.Value); err != nil {
				return err
			}
		}
		c.emit(code.OpHash, len(node.Pairs))
	case *ast.LetStatement:
		if err := c.Compile(node.Value); err != nil {
			return err
		}
		symbol := c.symbols.Define(node.Name.Value)
		if symbol.Scope == GlobalScope {
			c.emit(code.OpSetGlobal, symbol.Index)
		} else {
			c.emit(code.OpSetLocal, symbol.Index)
		}
	case *ast.Identifier:
		symbol, ok := c.symbols.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("invalid identifier: %s", node.Value)
		}
		if symbol.Scope == GlobalScope {
			c.emit(code.OpGetGlobal, symbol.Index)
		} else {
			c.emit(code.OpGetLocal, symbol.Index)
		}
	case *ast.IndexExpression:
		if err := c.Compile(node.Value); err != nil {
			return err
		}
		if err := c.Compile(node.Index); err != nil {
			return err
		}
		c.emit(code.OpIndex)
	case *ast.FunctionLiteral:
		c.enterScope()

		// make the parameters locals
		for _, p := range node.Parameters {
			c.symbols.Define(p.Name.Value)
		}

		if err := c.Compile(node.Body); err != nil {
			return err
		}

		// handle implicit return
		scope := c.scope()
		if scope.prev.Is(code.OpPop) {
			scope.undo()
			scope.emit(code.OpReturn)
		}

		// handle empty function
		if !scope.prev.Is(code.OpReturn) {
			scope.emit(code.OpNull)
			scope.emit(code.OpReturn)
		}

		nLocals := c.symbols.Count
		compiledFn := &object.CompiledFunction{
			NumParameters: len(node.Parameters),
			NumLocals:     nLocals,
			Instructions:  c.leaveScope(),
		}
		c.emit(code.OpConstant, c.addConstant(compiledFn))
	case *ast.ReturnStatement:
		if node.ReturnValue != nil {
			if err := c.Compile(node.ReturnValue); err != nil {
				return err
			}
		} else {
			c.emit(code.OpNull)
		}
		c.emit(code.OpReturn)
	case *ast.CallExpression:
		if err := c.Compile(node.Function); err != nil {
			return err
		}
		for _, arg := range node.Arguments {
			if err := c.Compile(arg); err != nil {
				return err
			}
		}
		c.emit(code.OpCall, len(node.Arguments))
	case *ast.NullExpression:
		c.emit(code.OpNull)
	}
	return nil
}

func (c *Compiler) instructions() code.Instructions {
	return c.scope().instructions
}

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	return c.scope().emit(op, operands...)
}

func (c *Compiler) rewrite(pos int, op code.Opcode, operands ...int) {
	c.scope().rewrite(pos, op, operands...)
}

func (c *Compiler) addInstructions(ins code.Instructions) int {
	scope := c.scope()
	pos := len(scope.instructions)
	scope.instructions = append(scope.instructions, ins...)
	return pos
}

func (c *Compiler) addConstant(v object.Object) int {
	c.constants = append(c.constants, v)
	return len(c.constants) - 1
}

func (c *Compiler) scope() *Scope {
	return c.scopes[len(c.scopes)-1]
}

func (c *Compiler) enterScope() {
	c.symbols = NewSymbolTable(c.symbols)
	c.scopes = append(c.scopes, &Scope{})
}

func (c *Compiler) leaveScope() code.Instructions {
	ins := c.instructions()
	c.scopes = c.scopes[:len(c.scopes)-1]
	c.symbols = c.symbols.Outer
	return ins
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions(),
		Constants:    c.constants,
	}
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}
