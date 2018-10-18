package vm

import (
	"testing"

	"gotest.tools/assert"

	"github.com/icholy/monkey/compiler"
	"github.com/icholy/monkey/object"
	"github.com/icholy/monkey/parser"
)

func TestIntegerArithmetic(t *testing.T) {
	tests := []struct {
		input    string
		expected object.Object
	}{
		{"1", object.New(1)},
		{"2", object.New(2)},
		{"1 + 2", object.New(3)},
		{"2 - 2", object.New(0)},
		{"10 * 10", object.New(100)},
		{"1 / 1", object.New(1)},
		{"1 + 4 * 2", object.New(9)},
		{"10 + 10 / 5", object.New(12)},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program, err := parser.Parse(tt.input)
			assert.NilError(t, err)
			bytecode, err := compiler.Compile(program)
			assert.NilError(t, err)
			vm := New(bytecode)
			assert.NilError(t, vm.Run())
			assert.DeepEqual(t, vm.LastPopped(), tt.expected)
		})
	}
}
