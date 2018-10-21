package vm

import (
	"github.com/icholy/monkey/code"
)

type Frame struct {
	instructions code.Instructions
	ip           int
}

func NewFrame(ins code.Instructions) *Frame {
	return &Frame{instructions: ins, ip: -1}
}

func (f *Frame) next() bool {
	f.ip++
	return f.ip < len(f.instructions)
}

func (f *Frame) Opcode() code.Opcode {
	return code.Opcode(f.instructions[f.ip])
}

func (f *Frame) ReadOperand() int {
	x := code.ReadUint16(f.instructions[f.ip+1:])
	f.ip += 2
	return int(x)
}
