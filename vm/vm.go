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
		case code.OpAdd:
			right := vm.pop()
			left := vm.pop()

			leftVal := left.(*object.Integer).Value
			rightVal := right.(*object.Integer).Value
			vm.push(&object.Integer{Value: leftVal + rightVal})
		}
	}
	return nil
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
