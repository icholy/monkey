package compiler

import (
	"testing"

	"github.com/icholy/monkey/object"

	"gotest.tools/assert"

	"github.com/icholy/monkey/code"
	"github.com/icholy/monkey/parser"
)

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
					code.Make(code.OpConstant, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpPop),
				),
				Constants: []object.Object{
					object.New(1),
					object.New(2),
				},
			},
		},
		{
			input: "12; 43",
			expected: &Bytecode{
				Instructions: code.Concat(
					code.Make(code.OpConstant, 0),
					code.Make(code.OpPop),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpPop),
				),
				Constants: []object.Object{
					object.New(12),
					object.New(43),
				},
			},
		},
		{
			input: "1 - 2",
			expected: &Bytecode{
				Instructions: code.Concat(
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpSub),
					code.Make(code.OpPop),
				),
				Constants: []object.Object{
					object.New(1),
					object.New(2),
				},
			},
		},
		{
			input: "1 * 2",
			expected: &Bytecode{
				Instructions: code.Concat(
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpMul),
					code.Make(code.OpPop),
				),
				Constants: []object.Object{
					object.New(1),
					object.New(2),
				},
			},
		},
		{
			input: "1 / 2",
			expected: &Bytecode{
				Instructions: code.Concat(
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpDiv),
					code.Make(code.OpPop),
				),
				Constants: []object.Object{
					object.New(1),
					object.New(2),
				},
			},
		},
		{
			input: "true; false",
			expected: &Bytecode{
				Instructions: code.Concat(
					code.Make(code.OpTrue),
					code.Make(code.OpPop),
					code.Make(code.OpFalse),
					code.Make(code.OpPop),
				),
			},
		},
		{
			input: "1 > 2",
			expected: &Bytecode{
				Instructions: code.Concat(
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpGreaterThan),
					code.Make(code.OpPop),
				),
				Constants: []object.Object{
					object.New(1),
					object.New(2),
				},
			},
		},
		{
			input: "true == false",
			expected: &Bytecode{
				Instructions: code.Concat(
					code.Make(code.OpTrue),
					code.Make(code.OpFalse),
					code.Make(code.OpEqual),
					code.Make(code.OpPop),
				),
			},
		},
		{
			input: "1 < 2",
			expected: &Bytecode{
				Instructions: code.Concat(
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpGreaterThan),
					code.Make(code.OpPop),
				),
				Constants: []object.Object{
					object.New(2),
					object.New(1),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program, err := parser.Parse(tt.input)
			assert.NilError(t, err)
			actual, err := Compile(program)
			assert.NilError(t, err)
			assert.DeepEqual(t, tt.expected.Constants, actual.Constants)
			assert.Equal(t, tt.expected.Instructions.String(), actual.Instructions.String())
		})
	}
}
