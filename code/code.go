package code

import (
	"encoding/binary"
	"fmt"
)

type Instructions []byte

type Opcode byte

const (
	OpConstant Opcode = iota
)

type Definition struct {
	Name          string
	OperandWidths []int
}

func (d Definition) Width() int {
	width := 1
	for _, w := range d.OperandWidths {
		width += w
	}
	return width
}

var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", []int{2}},
}

func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return def, nil
}

func Make(op Opcode, operands ...int) Instructions {
	def, ok := definitions[op]
	if !ok {
		return nil
	}
	if len(operands) != len(def.OperandWidths) {
		panic("number of operands doesn't match opcode definition")
	}

	instruction := make([]byte, def.Width())
	instruction[0] = byte(op)

	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}
		offset += width
	}

	return instruction
}

func Concat(ii ...Instructions) Instructions {
	var concatted Instructions
	for _, i := range ii {
		concatted = append(concatted, i...)
	}
	return concatted
}
