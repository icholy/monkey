package compiler

import (
	"testing"

	"github.com/icholy/monkey/object"

	"gotest.tools/assert"

	"github.com/icholy/monkey/code"
	"github.com/icholy/monkey/parser"
)

type compilerTestCase struct {
	input    string
	expected *Bytecode
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []struct {
		input    string
		expected *Bytecode
	}{
		{
			input: "1 + 2",
			expected: &Bytecode{
				Instructions: code.Concat(
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 2),
				),
				Constants: []object.Object{
					object.New(1),
					object.New(2),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program, err := parser.Parse(tt.input)
			assert.NilError(t, err)
			compiler := New()
			assert.NilError(t, compiler.Compile(program))
			assert.DeepEqual(t, tt.expected, compiler.Bytecode())
		})
	}
}
