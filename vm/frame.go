package vm

import (
	"github.com/icholy/monkey/code"
	"github.com/icholy/monkey/object"
)

type Frame struct {
	instructions code.Instructions
	cl           *object.Closure
	ip           int
	bp           int
}

func NewFrame(cl *object.Closure, bp int) *Frame {
	return &Frame{
		instructions: cl.Fn.Instructions,
		cl:           cl,
		ip:           -1,
		bp:           bp,
	}
}

func (f *Frame) next() bool {
	f.ip++
	return f.ip < len(f.instructions)
}

func (f *Frame) Opcode() code.Opcode {
	return code.Opcode(f.instructions[f.ip])
}

func (f *Frame) JumpTo(ip int) {
	f.ip = ip - 1
}

func (f *Frame) ReadUint8() int {
	x := code.ReadUint8(f.instructions[f.ip+1:])
	f.ip++
	return int(x)
}

func (f *Frame) ReadUint16() int {
	x := code.ReadUint16(f.instructions[f.ip+1:])
	f.ip += 2
	return int(x)
}
