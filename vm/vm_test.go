package vm

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"gotest.tools/assert"

	"github.com/icholy/monkey/compiler"
	"github.com/icholy/monkey/object"
	"github.com/icholy/monkey/parser"
)

func TestRun(t *testing.T) {
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
		{"true", object.New(true)},
		{"false", object.New(false)},
		{"1 > 2", object.New(false)},
		{"1 > 2", object.New(false)},
		{"1 < 1", object.New(false)},
		{"1 != 2", object.New(true)},
		{"false != true", object.New(true)},
		{"true == true", object.New(true)},
		{"(1 > 2) == true", object.New(false)},
		{"(1 > 2) == false", object.New(true)},
		{"!true", object.New(false)},
		{"!false", object.New(true)},
		{"-1", object.New(-1)},
		{"if true { 10 }", object.New(10)},
		{"if true { 10 } else { 20 }", object.New(10)},
		{"if false { 10 } else { 20 }", object.New(20)},
		{"if (1) { 10 }", object.New(10)},
		{"if 1 > 2 { 10 } else { 20 }", object.New(20)},
		{"if 1 < 2 { 10 } else { 20 }", object.New(10)},
		{"if false { 10 }", object.New(nil)},
		{"if 1 > 2 { 10 }", object.New(nil)},
		{"!(if false { 4; })", object.New(true)},
		{"if null { 10 } else { 20 }", object.New(20)},
		{"let one = 1; one", object.New(1)},
		{"let one = 1; let two = 2; one + two", object.New(3)},
		{`"hello" + " " + "world"`, object.New("hello world")},
		{"[1, 2, 3]", object.New([]interface{}{1, 2, 3})},
		{`[1 + 2, "test", true == false]`, object.New([]interface{}{3, "test", false})},
		{"{1: 1, 2: 2, 3:3 }", object.New(map[interface{}]interface{}{1: 1, 2: 2, 3: 3})},
		{"[1, 2, 3][1]", object.New(2)},
		{"fn() { 15 }()", object.New(15)},
		{"let one = fn() { 1 }; one() + one()", object.New(2)},
		{"let x = fn() { return 1; return 2; }; x()", object.New(1)},
		{"let x = fn() { return; }; x()", object.New(nil)},
		{"let x = fn() {}; x()", object.New(nil)},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program, err := parser.Parse(tt.input)
			assert.NilError(t, err)
			bytecode, err := compiler.Compile(program)
			assert.NilError(t, err)
			vm := New(bytecode)
			assert.NilError(t, vm.Run())
			assert.DeepEqual(t, vm.LastPopped(), tt.expected, cmp.AllowUnexported(object.Hash{}))
		})
	}
}
