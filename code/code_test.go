package code

import (
	"strings"
	"testing"

	"gotest.tools/assert"
)

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expected Instructions
	}{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
		{OpAdd, []int{}, []byte{byte(OpAdd)}},
		{OpGetLocal, []int{255}, []byte{byte(OpGetLocal), 255}},
	}

	for _, tt := range tests {
		actual := Make(tt.op, tt.operands...)
		assert.DeepEqual(t, tt.expected, actual)
	}
}

func TestInstructionsString(t *testing.T) {
	instructions := Concat(
		Make(OpConstant, 1),
		Make(OpConstant, 2),
		Make(OpConstant, 65535),
		Make(OpAdd),
	)

	expected := strings.Join([]string{
		"0000 OpConstant 1",
		"0003 OpConstant 2",
		"0006 OpConstant 65535",
		"0009 OpAdd",
	}, "\n")

	assert.Equal(t, instructions.String(), expected)
}

func TestReadOperands(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		nBytes   int
	}{
		{OpConstant, []int{65535}, 2},
		{OpAdd, []int{}, 0},
	}

	for _, tt := range tests {
		ins := Make(tt.op, tt.operands...)
		def, err := Lookup(byte(tt.op))
		assert.NilError(t, err)
		operands, n := ReadOperands(def, ins[1:])
		assert.Equal(t, n, tt.nBytes)
		assert.DeepEqual(t, operands, tt.operands)
	}
}
