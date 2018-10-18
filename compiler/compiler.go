package compiler

import (
	"github.com/icholy/monkey/ast"
	"github.com/icholy/monkey/code"
	"github.com/icholy/monkey/object"
)

type Compiler struct {
	instructions code.Instructions
	constants    []object.Object
}

func New() *Compiler {
	return &Compiler{}
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
		return c.Compile(node.Expression)
	case *ast.InfixExpression:
		if err := c.Compile(node.Left); err != nil {
			return err
		}
		if err := c.Compile(node.Right); err != nil {
			return err
		}
	case *ast.IntegerLiteral:
		v := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(v))
	}
	return nil
}

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	return c.addInstructions(ins)
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
