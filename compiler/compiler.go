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
	instructions code.Instructions
	constants    []object.Object
	symbols      *SymbolTable

	prev     Instruction
	prevprev Instruction
}

func New() *Compiler {
	return &Compiler{
		symbols: NewSymbolTable(),
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
		if err := c.Compile(node.Condition); err != nil {
			return err
		}
		jumpNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999)
		if err := c.Compile(node.Concequence); err != nil {
			return err
		}
		if c.prev.Is(code.OpPop) {
			c.undo()
		}

		jumpPos := c.emit(code.OpJump, 9999)
		c.rewrite(jumpNotTruthyPos, code.OpJumpNotTruthy, len(c.instructions))

		if node.Alternative != nil {
			if err := c.Compile(node.Alternative); err != nil {
				return err
			}
			if c.prev.Is(code.OpPop) {
				c.undo()
			}
		} else {
			c.emit(code.OpNull)
		}

		c.rewrite(jumpPos, code.OpJump, len(c.instructions))

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
	case *ast.LetStatement:
		if err := c.Compile(node.Value); err != nil {
			return err
		}
		symbol := c.symbols.Define(node.Name.Value)
		c.emit(code.OpSetGlobal, symbol.Index)
	case *ast.Identifier:
		symbol, ok := c.symbols.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("invalid identifier: %s", node.Value)
		}
		c.emit(code.OpGetGlobal, symbol.Index)
	case *ast.NullExpression:
		c.emit(code.OpNull)
	}
	return nil
}

func (c *Compiler) undo() {
	c.instructions = c.instructions[:c.prev.Position]
	c.prev = c.prevprev
}

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstructions(ins)
	c.prevprev = c.prev
	c.prev = Instruction{
		Opcode:   op,
		Position: pos,
	}
	return pos
}

func (c *Compiler) rewrite(pos int, op code.Opcode, operands ...int) {
	ins := code.Make(op, operands...)
	for i := range ins {
		c.instructions[pos+i] = ins[i]
	}
}

func (c *Compiler) addInstructions(ins code.Instructions) int {
	pos := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return pos
}

func (c *Compiler) addConstant(v object.Object) int {
	c.constants = append(c.constants, v)
	return len(c.constants) - 1
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}
