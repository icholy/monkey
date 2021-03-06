package code

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

type Opcode byte

const (
	OpConstant Opcode = iota
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpPop
	OpTrue
	OpFalse
	OpEqual
	OpNotEqual
	OpGreaterThan
	OpMinus
	OpBang
	OpJump
	OpJumpNotTruthy
	OpNull
	OpSetGlobal
	OpGetGlobal
	OpGetBuiltin
	OpSetLocal
	OpGetLocal
	OpGetFree
	OpArray
	OpHash
	OpIndex
	OpCall
	OpReturn
	OpClosure
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

func (d Definition) String() string {
	if len(d.OperandWidths) == 0 {
		return d.Name
	}
	var widths []string
	for _, w := range d.OperandWidths {
		widths = append(widths, strconv.Itoa(w))
	}
	return fmt.Sprintf("%s(%s)", d.Name, strings.Join(widths, ", "))
}

var definitions = map[Opcode]*Definition{
	OpConstant:      {"OpConstant", []int{2}},
	OpAdd:           {"OpAdd", []int{}},
	OpSub:           {"OpSub", []int{}},
	OpMul:           {"OpMul", []int{}},
	OpDiv:           {"OpDiv", []int{}},
	OpPop:           {"OpPop", []int{}},
	OpTrue:          {"OpTrue", []int{}},
	OpFalse:         {"OpFalse", []int{}},
	OpNull:          {"OpNull", []int{}},
	OpEqual:         {"OpEqual", []int{}},
	OpNotEqual:      {"OpNotEqual", []int{}},
	OpGreaterThan:   {"OpGreaterThan", []int{}},
	OpMinus:         {"OpMinus", []int{}},
	OpBang:          {"OpBang", []int{}},
	OpJump:          {"OpJump", []int{2}},
	OpJumpNotTruthy: {"OpJumpNotTruthy", []int{2}},
	OpSetGlobal:     {"OpSetGlobal", []int{2}},
	OpGetGlobal:     {"OpGetGlobal", []int{2}},
	OpSetLocal:      {"OpSetLocal", []int{1}},
	OpGetLocal:      {"OpGetLocal", []int{1}},
	OpGetBuiltin:    {"OpGetBuiltin", []int{1}},
	OpArray:         {"OpArray", []int{2}},
	OpHash:          {"OpHash", []int{2}},
	OpIndex:         {"OpIndex", []int{}},
	OpCall:          {"OpCall", []int{1}},
	OpReturn:        {"OpReturn", []int{}},
	OpClosure:       {"OpClosure", []int{2, 1}},
	OpGetFree:       {"OpGetFree", []int{1}},
}

type Instructions []byte

func (ins Instructions) String() string {
	var lines []string
	for i := 0; i < len(ins); i++ {
		var b strings.Builder
		def, err := Lookup(ins[i])
		if err != nil {
			lines = append(lines, fmt.Sprintf("ERROR: %s", err))
			continue
		}
		fmt.Fprintf(&b, "%04d %s", i, def.Name)
		operands, n := ReadOperands(def, ins[i+1:])
		for _, o := range operands {
			fmt.Fprintf(&b, " %d", o)
		}
		i += n
		lines = append(lines, b.String())
	}
	return strings.Join(lines, "\n")
}

func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0
	for i, width := range def.OperandWidths {
		if len(ins) < width {
			return nil, 0
		}
		switch width {
		case 1:
			operands[i] = int(ReadUint8(ins[offset:]))
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))
		}
		offset += width
	}
	return operands, offset
}

func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}

func ReadUint8(ins Instructions) uint8 {
	return uint8(ins[0])
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
		case 1:
			instruction[offset] = uint8(o)
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
