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
		{"let x = fn() { return }; x()", object.New(nil)},
		{"let x = fn() {}; x()", object.New(nil)},
		{"let one = fn() { let one = 1; one }; one()", object.New(1)},
		{"let x = fn() { let one = 1; let two = 2; one + two }; x()", object.New(3)},
		{"let x = fn() { let one = 1; let two = 2; one + two }; let y = fn() { let three = 3; let four = 4; three + four }; x() + y()", object.New(10)},
		{"let seed = 50; let minusOne = fn() { let x = 1; seed - x }()", object.New(49)},
		{"let x = fn() { let x = 1; x }; let y = fn() { let x = 1; x }; x() + y()", object.New(2)},
		{"fn(x) { x }(1)", object.New(1)},
		{"let twice = fn(f, x) { f(f(x)) }; let double = fn(x) { x * 2 }; twice(double, 1)", object.New(4)},
		{"let x = fn(x, y) { let x = 10; x + y }; x(1, 1)", object.New(11)},
		{"len([])", object.New(0)},
		{"append([], 1)", object.New([]interface{}{1})},
		{"len([]); 1", object.New(1)},
		{`len("hello world")`, object.New(11)},
		{"last([1, 2, 3])", object.New(3)},
		{"let x = len([1, 2, 3]); let y = len([1, 2, 3]); y + x", object.New(6)},
		{"let make = fn(a) { fn() {a} }; make(1)()", object.New(1)},
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
