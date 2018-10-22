package compiler

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	is "gotest.tools/assert/cmp"

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
		{
			input: "-1",
			expected: &Bytecode{
				Instructions: code.Concat(
					code.Make(code.OpConstant, 0),
					code.Make(code.OpMinus),
					code.Make(code.OpPop),
				),
				Constants: []object.Object{
					object.New(1),
				},
			},
		},
		{
			input: "!true",
			expected: &Bytecode{
				Instructions: code.Concat(
					code.Make(code.OpTrue),
					code.Make(code.OpBang),
					code.Make(code.OpPop),
				),
			},
		},
		{
			input: "if true { 10 }; 3333",
			expected: &Bytecode{
				Instructions: code.Concat(
					code.Make(code.OpTrue),              // 0000
					code.Make(code.OpJumpNotTruthy, 10), // 0001
					code.Make(code.OpConstant, 0),       // 0004
					code.Make(code.OpJump, 11),          // 0007
					code.Make(code.OpNull),              // 0010
					code.Make(code.OpPop),               // 0011
					code.Make(code.OpConstant, 1),       // 0012
					code.Make(code.OpPop),               // 0015
				),
				Constants: []object.Object{
					object.New(10),
					object.New(3333),
				},
			},
		},
		{
			input: "if true { 10 } else { 20 }; 3333;",
			expected: &Bytecode{
				Instructions: code.Concat(
					code.Make(code.OpTrue),              // 0000
					code.Make(code.OpJumpNotTruthy, 10), // 0001
					code.Make(code.OpConstant, 0),       // 0004
					code.Make(code.OpJump, 13),          // 0007
					code.Make(code.OpConstant, 1),       // 0010
					code.Make(code.OpPop),               // 0013
					code.Make(code.OpConstant, 2),       // 0014
					code.Make(code.OpPop),               // 0017
				),
				Constants: []object.Object{
					object.New(10),
					object.New(20),
					object.New(3333),
				},
			},
		},
		{
			input: "let x = 33; let y = x; y;",
			expected: &Bytecode{
				Instructions: code.Concat(
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpSetGlobal, 1),
					code.Make(code.OpGetGlobal, 1),
					code.Make(code.OpPop),
				),
				Constants: []object.Object{
					object.New(33),
				},
			},
		},
		{
			input: `"one" + "two"`,
			expected: &Bytecode{
				Instructions: code.Concat(
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpPop),
				),
				Constants: []object.Object{
					object.New("one"),
					object.New("two"),
				},
			},
		},
		{
			input: "[1, 2, 3]",
			expected: &Bytecode{
				Instructions: code.Concat(
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpConstant, 2),
					code.Make(code.OpArray, 3),
					code.Make(code.OpPop),
				),
				Constants: []object.Object{
					object.New(1),
					object.New(2),
					object.New(3),
				},
			},
		},
		{
			input: "{ 1: 1, 2: 2, 3: 3}",
			expected: &Bytecode{
				Instructions: code.Concat(
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpConstant, 2),
					code.Make(code.OpConstant, 3),
					code.Make(code.OpConstant, 4),
					code.Make(code.OpConstant, 5),
					code.Make(code.OpHash, 3),
					code.Make(code.OpPop),
				),
				Constants: []object.Object{
					object.New(1),
					object.New(1),
					object.New(2),
					object.New(2),
					object.New(3),
					object.New(3),
				},
			},
		},
		{
			input: "[1, 2, 3][1]",
			expected: &Bytecode{
				Instructions: code.Concat(
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpConstant, 2),
					code.Make(code.OpArray, 3),
					code.Make(code.OpConstant, 3),
					code.Make(code.OpIndex),
					code.Make(code.OpPop),
				),
				Constants: []object.Object{
					object.New(1),
					object.New(2),
					object.New(3),
					object.New(1),
				},
			},
		},
		{
			input: "fn() { return 5 + 10 }",
			expected: &Bytecode{
				Instructions: code.Concat(
					code.Make(code.OpConstant, 2),
					code.Make(code.OpPop),
				),
				Constants: []object.Object{
					object.New(5),
					object.New(10),
					&object.CompiledFunction{
						Instructions: code.Concat(
							code.Make(code.OpConstant, 0),
							code.Make(code.OpConstant, 1),
							code.Make(code.OpAdd),
							code.Make(code.OpReturn),
						),
					},
				},
			},
		},
		{
			input: "fn() { 5 + 10 }",
			expected: &Bytecode{
				Instructions: code.Concat(
					code.Make(code.OpConstant, 2),
					code.Make(code.OpPop),
				),
				Constants: []object.Object{
					object.New(5),
					object.New(10),
					&object.CompiledFunction{
						Instructions: code.Concat(
							code.Make(code.OpConstant, 0),
							code.Make(code.OpConstant, 1),
							code.Make(code.OpAdd),
							code.Make(code.OpReturn),
						),
					},
				},
			},
		},
		{
			input: "fn() {}",
			expected: &Bytecode{
				Instructions: code.Concat(
					code.Make(code.OpConstant, 0),
					code.Make(code.OpPop),
				),
				Constants: []object.Object{
					&object.CompiledFunction{
						Instructions: code.Concat(
							code.Make(code.OpNull),
							code.Make(code.OpReturn),
						),
					},
				},
			},
		},
		{
			input: "fn() { return; }",
			expected: &Bytecode{
				Instructions: code.Concat(
					code.Make(code.OpConstant, 0),
					code.Make(code.OpPop),
				),
				Constants: []object.Object{
					&object.CompiledFunction{
						Instructions: code.Concat(
							code.Make(code.OpNull),
							code.Make(code.OpReturn),
						),
					},
				},
			},
		},
		{
			input: "fn() { return 24; }()",
			expected: &Bytecode{
				Instructions: code.Concat(
					code.Make(code.OpConstant, 1),
					code.Make(code.OpCall, 0),
					code.Make(code.OpPop),
				),
				Constants: []object.Object{
					object.New(24),
					&object.CompiledFunction{
						Instructions: code.Concat(
							code.Make(code.OpConstant, 0),
							code.Make(code.OpReturn),
						),
					},
				},
			},
		},
		{
			input: "let f = fn() { 24 }; f();",
			expected: &Bytecode{
				Instructions: code.Concat(
					code.Make(code.OpConstant, 1),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpCall, 0),
					code.Make(code.OpPop),
				),
				Constants: []object.Object{
					object.New(24),
					&object.CompiledFunction{
						Instructions: code.Concat(
							code.Make(code.OpConstant, 0),
							code.Make(code.OpReturn),
						),
					},
				},
			},
		},
		{
			input: "let one = fn() { 1 }; one() + one()",
			expected: &Bytecode{
				Instructions: code.Concat(
					code.Make(code.OpConstant, 1),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpCall, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpCall, 0),
					code.Make(code.OpAdd),
					code.Make(code.OpPop),
				),
				Constants: []object.Object{
					object.New(1),
					&object.CompiledFunction{
						Instructions: code.Concat(
							code.Make(code.OpConstant, 0),
							code.Make(code.OpReturn),
						),
					},
				},
			},
		},
		{
			input: "let num = 55; fn() { num }",
			expected: &Bytecode{
				Instructions: code.Concat(
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpPop),
				),
				Constants: []object.Object{
					object.New(55),
					&object.CompiledFunction{
						Instructions: code.Concat(
							code.Make(code.OpGetGlobal, 0),
							code.Make(code.OpReturn),
						),
					},
				},
			},
		},
		{
			input: "fn() { let num = 55; num }",
			expected: &Bytecode{
				Instructions: code.Concat(
					code.Make(code.OpConstant, 1),
					code.Make(code.OpPop),
				),
				Constants: []object.Object{
					object.New(55),
					&object.CompiledFunction{
						NumLocals: 1,
						Instructions: code.Concat(
							code.Make(code.OpConstant, 0),
							code.Make(code.OpSetLocal, 0),
							code.Make(code.OpGetLocal, 0),
							code.Make(code.OpReturn),
						),
					},
				},
			},
		},
		{
			input: "fn() { let a = 55; let b = 77; return a + b }",
			expected: &Bytecode{
				Instructions: code.Concat(
					code.Make(code.OpConstant, 2),
					code.Make(code.OpPop),
				),
				Constants: []object.Object{
					object.New(55),
					object.New(77),
					&object.CompiledFunction{
						NumLocals: 2,
						Instructions: code.Concat(
							code.Make(code.OpConstant, 0),
							code.Make(code.OpSetLocal, 0),
							code.Make(code.OpConstant, 1),
							code.Make(code.OpSetLocal, 1),
							code.Make(code.OpGetLocal, 0),
							code.Make(code.OpGetLocal, 1),
							code.Make(code.OpAdd),
							code.Make(code.OpReturn),
						),
					},
				},
			},
		},
		{
			input: "fn(){}(1, 2)",
			expected: &Bytecode{
				Instructions: code.Concat(
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpConstant, 2),
					code.Make(code.OpCall, 2),
					code.Make(code.OpPop),
				),
				Constants: []object.Object{
					&object.CompiledFunction{
						Instructions: code.Concat(
							code.Make(code.OpNull),
							code.Make(code.OpReturn),
						),
					},
					object.New(1),
					object.New(2),
				},
			},
		},
		{
			input: "fn(a) {}",
			expected: &Bytecode{
				Instructions: code.Concat(
					code.Make(code.OpConstant, 0),
					code.Make(code.OpPop),
				),
				Constants: []object.Object{
					&object.CompiledFunction{
						NumLocals:     1,
						NumParameters: 1,
						Instructions: code.Concat(
							code.Make(code.OpNull),
							code.Make(code.OpReturn),
						),
					},
				},
			},
		},
		{
			input: "let oneArg = fn(a, b) { a + b }; oneArg(24)",
			expected: &Bytecode{
				Instructions: code.Concat(
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpCall, 1),
					code.Make(code.OpPop),
				),
				Constants: []object.Object{
					&object.CompiledFunction{
						NumLocals:     2,
						NumParameters: 2,
						Instructions: code.Concat(
							code.Make(code.OpGetLocal, 0),
							code.Make(code.OpGetLocal, 1),
							code.Make(code.OpAdd),
							code.Make(code.OpReturn),
						),
					},
					object.New(24),
				},
			},
		},
		{
			input: "len([]); append([], 1)",
			expected: &Bytecode{
				Instructions: code.Concat(
					code.Make(code.OpGetBuiltin, object.FindBuiltin("len")),
					code.Make(code.OpArray, 0),
					code.Make(code.OpCall, 1),
					code.Make(code.OpPop),
					code.Make(code.OpGetBuiltin, object.FindBuiltin("append")),
					code.Make(code.OpArray, 0),
					code.Make(code.OpConstant, 0),
					code.Make(code.OpCall, 2),
					code.Make(code.OpPop),
				),
				Constants: []object.Object{
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
			assert.DeepEqual(t, tt.expected.Constants, actual.Constants, cmp.Transformer("Instructions", code.Instructions.String))
			assert.DeepEqual(t, tt.expected.Instructions, actual.Instructions, cmp.Transformer("Instructions", code.Instructions.String))
			assert.DeepEqual(t, tt.expected.Constants, actual.Constants)
			assert.DeepEqual(t, tt.expected.Instructions, actual.Instructions)
		})
	}
}
func TestScopes(t *testing.T) {
	compiler := New()
	global := compiler.symbols

	compiler.emit(code.OpMul)
	assert.Assert(t, is.Len(compiler.scopes, 1))
	assert.Assert(t, compiler.scope().prev.Is(code.OpMul))

	compiler.enterScope()
	assert.Assert(t, is.Len(compiler.scopes, 2))
	assert.Equal(t, global, compiler.symbols.Outer)

	compiler.emit(code.OpSub)
	assert.Assert(t, is.Len(compiler.scope().instructions, 1))
	assert.Assert(t, compiler.scope().prev.Is(code.OpSub))

	compiler.leaveScope()
	assert.Assert(t, is.Len(compiler.scopes, 1))
	assert.Equal(t, global, compiler.symbols)
	assert.Assert(t, is.Nil(compiler.symbols.Outer))

	compiler.emit(code.OpAdd)
	assert.Assert(t, is.Len(compiler.scope().instructions, 2))
	assert.Assert(t, compiler.scope().prev.Is(code.OpAdd))
	assert.Assert(t, compiler.scope().prevprev.Is(code.OpMul))
}
