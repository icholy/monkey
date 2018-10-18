package code

import (
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
	}

	for _, tt := range tests {
		actual := Make(tt.op, tt.operands...)
		assert.DeepEqual(t, tt.expected, actual)
	}
}
