package vm

import (
	"fmt"

	"github.com/icholy/monkey/code"
	"github.com/icholy/monkey/compiler"
	"github.com/icholy/monkey/object"
)

const (
	StackSize   = 2048
	GlobalsSize = 65536
	MaxFrames   = 1024
)

var (
	True  = object.New(true)
	False = object.New(false)
	Null  = object.New(nil)
)

type VM struct {
	constants []object.Object

	// always points to the next value, top of the stack is sp-1
	sp      int
	stack   []object.Object
	globals []object.Object

	frames   []*Frame
	frameIdx int
}

func New(bytecode *compiler.Bytecode) *VM {

	frames := make([]*Frame, MaxFrames)
	frames[0] = NewFrame(bytecode.Instructions, 0)

	return &VM{
		constants: bytecode.Constants,
		sp:        0,
		stack:     make([]object.Object, StackSize),
		globals:   make([]object.Object, GlobalsSize),
		frames:    frames,
	}
}

func NewWithGlobals(bytecode *compiler.Bytecode, globals []object.Object) *VM {
	vm := New(bytecode)
	vm.globals = globals
	return vm
}

func (vm *VM) frame() *Frame {
	return vm.frames[vm.frameIdx]
}

func (vm *VM) pushFrame(f *Frame) {
	vm.frameIdx++
	vm.frames[vm.frameIdx] = f
}

func (vm *VM) popFrame() *Frame {
	f := vm.frames[vm.frameIdx]
	vm.frameIdx--
	return f
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

	frame := vm.frame()

	for frame.next() {

		op := frame.Opcode()

		switch op {
		case code.OpConstant:
			index := frame.ReadUint16()
			if err := vm.push(vm.constants[index]); err != nil {
				return err
			}
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
			right := vm.pop()
			left := vm.pop()
			if err := vm.binaryOp(op, left, right); err != nil {
				return err
			}
		case code.OpEqual, code.OpNotEqual, code.OpGreaterThan:
			right := vm.pop()
			left := vm.pop()
			if err := vm.compareOp(op, left, right); err != nil {
				return err
			}
		case code.OpTrue:
			if err := vm.push(True); err != nil {
				return err
			}
		case code.OpFalse:
			if err := vm.push(False); err != nil {
				return err
			}
		case code.OpNull:
			if err := vm.push(Null); err != nil {
				return err
			}
		case code.OpArray:
			n := frame.ReadUint16()
			elements := make([]object.Object, n)
			for i := 0; i < n; i++ {
				elements[n-1-i] = vm.pop()
			}
			v := &object.Array{Elements: elements}
			if err := vm.push(v); err != nil {
				return err
			}
		case code.OpHash:
			n := frame.ReadUint16()
			h := object.NewHash()
			for i := 0; i < n; i++ {
				value := vm.pop()
				key := vm.pop()
				h.Set(key, value)
			}
			if err := vm.push(h); err != nil {
				return err
			}
		case code.OpIndex:
			if err := vm.indexOp(); err != nil {
				return err
			}
		case code.OpPop:
			vm.pop()
		case code.OpMinus:
			if err := vm.minusOp(); err != nil {
				return err
			}
		case code.OpBang:
			if err := vm.bangOp(); err != nil {
				return err
			}
		case code.OpJump:
			pos := frame.ReadUint16()
			frame.JumpTo(pos)
		case code.OpJumpNotTruthy:
			pos := frame.ReadUint16()
			condition := vm.pop()
			if !isTruthy(condition) {
				frame.JumpTo(pos)
			}
		case code.OpSetGlobal:
			index := frame.ReadUint16()
			vm.globals[index] = vm.pop()
		case code.OpGetGlobal:
			index := frame.ReadUint16()
			if err := vm.push(vm.globals[index]); err != nil {
				return err
			}
		case code.OpSetLocal:
			index := frame.ReadUint8()
			vm.stack[frame.bp+index] = vm.pop()
		case code.OpGetLocal:
			index := frame.ReadUint8()
			vm.push(vm.stack[frame.bp+index])
		case code.OpGetBuiltin:
			index := frame.ReadUint8()
			vm.push(object.Builtins[index])
		case code.OpCall:
			nArgs := frame.ReadUint8() // num args
			fn := vm.stack[vm.sp-1-nArgs]

			switch fn := fn.(type) {
			case *object.Builtin:
				args := make([]object.Object, nArgs)
				for i := 0; i < nArgs; i++ {
					args[i] = vm.pop()
				}
				ret, err := fn.Fn(args...)
				if err != nil {
					return err
				}
				vm.push(ret)
			case *object.CompiledFunction:
				if nArgs != fn.NumParameters {
					return fmt.Errorf("wrong number of arguments: want %d, got %d", fn.NumParameters, nArgs)
				}
				frame = NewFrame(fn.Instructions, vm.sp-nArgs)
				vm.pushFrame(frame)
				vm.sp = frame.bp + fn.NumLocals
			default:
				return fmt.Errorf("calling non-function")
			}
		case code.OpReturn:
			retVal := vm.pop()
			vm.sp = vm.popFrame().bp - 1
			frame = vm.frame()
			if err := vm.push(retVal); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unexpected opcode: %d", op)
		}
	}
	return nil
}

func (vm *VM) DumpStack() {
	for i := 0; i < vm.sp; i++ {
		v := vm.stack[i]
		fmt.Println(v.Type(), v.Inspect(0))
	}
}

func isTruthy(v object.Object) bool {
	switch v := v.(type) {
	case *object.Boolean:
		return v.Value
	case *object.Null:
		return false
	default:
		return true
	}
}

func boolObject(v bool) object.Object {
	if v {
		return True
	}
	return False
}

func (vm *VM) indexOp() error {
	index := vm.pop()
	value := vm.pop()

	switch value := value.(type) {
	case *object.Array:
		i, ok := index.(*object.Integer)
		if !ok {
			return fmt.Errorf("cannot index into array with: %s", index.Type())
		}
		el, err := value.At(int(i.Value))
		if err != nil {
			return err
		}
		return vm.push(el)
	case *object.Hash:
		el, ok := value.Get(index)
		if !ok {
			return vm.push(Null)
		}
		return vm.push(el)
	default:
		return fmt.Errorf("cannot index into: %s", value.Type())
	}
}

func (vm *VM) minusOp() error {
	right := vm.pop()
	value, ok := right.(*object.Integer)
	if !ok {
		return fmt.Errorf("cannot use minus on type: %s", right.Type())
	}
	return vm.push(&object.Integer{Value: -value.Value})
}

func (vm *VM) bangOp() error {
	switch vm.pop() {
	case True:
		return vm.push(False)
	case False:
		return vm.push(True)
	case Null:
		return vm.push(True)
	default:
		return vm.push(False)
	}
}

func (vm *VM) compareOp(op code.Opcode, left, right object.Object) error {
	if left.Type() == object.INTEGER && left.Type() == object.INTEGER {
		return vm.compareIntegerOp(op, left.(*object.Integer), right.(*object.Integer))
	}
	switch op {
	case code.OpEqual:
		return vm.push(boolObject(left == right))
	case code.OpNotEqual:
		return vm.push(boolObject(left != right))
	default:
		return fmt.Errorf("unknown operator: %d (%s, %s)", op, left.Type(), right.Type())
	}
}

func (vm *VM) compareIntegerOp(op code.Opcode, left, right *object.Integer) error {
	switch op {
	case code.OpEqual:
		return vm.push(boolObject(left.Value == right.Value))
	case code.OpNotEqual:
		return vm.push(boolObject(left.Value != right.Value))
	case code.OpGreaterThan:
		return vm.push(boolObject(left.Value > right.Value))
	default:
		return fmt.Errorf("unknown operator: %d", op)
	}
}

func (vm *VM) binaryOp(op code.Opcode, left, right object.Object) error {
	if left.Type() == object.INTEGER && right.Type() == object.INTEGER {
		return vm.binaryIntegerOp(op, left.(*object.Integer), right.(*object.Integer))
	}
	if left.Type() == object.STRING && right.Type() == object.STRING {
		return vm.binaryStringOp(op, left.(*object.String), right.(*object.String))
	}
	return fmt.Errorf("unsuported types for binary operator: %s, %s", left.Type(), right.Type())
}

func (vm *VM) binaryStringOp(op code.Opcode, left, right *object.String) error {
	var result string
	switch op {
	case code.OpAdd:
		result = left.Value + right.Value
	default:
		return fmt.Errorf("unknown string operator: %d", op)
	}
	return vm.push(&object.String{Value: result})
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

func (vm *VM) peek() object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

func (vm *VM) pop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	vm.sp--
	return vm.stack[vm.sp]
}

func (vm *VM) push(v object.Object) error {
	if vm.sp > StackSize {
		return fmt.Errorf("stack overflow")
	}
	vm.stack[vm.sp] = v
	vm.sp++
	return nil
}
