package vm

import (
	"fmt"

	"github.com/icholy/monkey/code"
	"github.com/icholy/monkey/compiler"
	"github.com/icholy/monkey/object"
)

const StackSize = 2048

type VM struct {
	constants    []object.Object
	instructions code.Instructions

	// always points to the next value, top of the stack is sp-1
	sp    int
	stack []object.Object
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		constants:    bytecode.Constants,
		instructions: bytecode.Instructions,
		sp:           0,
		stack:        make([]object.Object, StackSize),
	}
}

func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

func (vm *VM) LastPopped() object.Object {
	return vm.stack[vm.sp]
}

func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])
		switch op {
		case code.OpConstant:
			index := code.ReadUint16(vm.instructions[ip+1:])
			if err := vm.push(vm.constants[index]); err != nil {
				return err
			}
			ip += 2
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
			right := vm.pop()
			left := vm.pop()
			if err := vm.binaryOp(op, left, right); err != nil {
				return err
			}
		case code.OpPop:
			vm.pop()
		default:
			return fmt.Errorf("unexpected opcode: %d", op)
		}
	}
	return nil
}

func (vm *VM) binaryOp(op code.Opcode, left, right object.Object) error {
	if left.Type() == object.INTEGER && right.Type() == object.INTEGER {
		return vm.binaryIntegerOp(op, left.(*object.Integer), right.(*object.Integer))
	}
	return fmt.Errorf("unsuported types for binary operator: %s, %s", left.Type(), right.Type())
}

func (vm *VM) binaryIntegerOp(op code.Opcode, left, right *object.Integer) error {
	var result int64
	switch op {
	case code.OpAdd:
		result = left.Value + right.Value
	case code.OpSub:
		result = left.Value - right.Value
	case code.OpMul:
		result = left.Value * right.Value
	case code.OpDiv:
		result = left.Value / right.Value
	default:
		return fmt.Errorf("unknown integer operator: %d", op)
	}
	return vm.push(&object.Integer{Value: result})
}

func (vm *VM) pop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	v := vm.StackTop()
	vm.sp--
	return v
}

func (vm *VM) push(v object.Object) error {
	if vm.sp > StackSize {
		return fmt.Errorf("stack overflow")
	}
	vm.stack[vm.sp] = v
	vm.sp++
	return nil
}
