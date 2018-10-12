package evaluator

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/icholy/monkey/object"
	"github.com/icholy/monkey/parser"
)

func TestEvaluator(t *testing.T) {
	t.Run("integer", func(t *testing.T) {
		RequireEqualEval(t, "5", &object.Integer{5})
		RequireEqualEval(t, "-10", &object.Integer{-10})
		RequireEqualEval(t, "1 + 1", &object.Integer{2})
		RequireEqualEval(t, "(2 * 2) + 1", &object.Integer{5})
	})

	t.Run("boolean expressions", func(t *testing.T) {
		RequireEqualEval(t, "true", TRUE)
		RequireEqualEval(t, "false", FALSE)
		RequireEqualEval(t, "!true", FALSE)
		RequireEqualEval(t, "!!true", TRUE)
		RequireEqualEval(t, "!false", TRUE)
		RequireEqualEval(t, "!!false", FALSE)
		RequireEqualEval(t, "1 < 2", TRUE)
		RequireEqualEval(t, "1 <= 2", TRUE)
		RequireEqualEval(t, "null == null", TRUE)
		RequireEqualEval(t, "null != null", FALSE)
	})

	t.Run("if expressions", func(t *testing.T) {
		RequireEqualEval(t, "if (2 > 1) { 123 } else { 4 }", &object.Integer{123})
		RequireEqualEval(t, "if (false) { 123 } else { 4 }", &object.Integer{4})
		RequireEqualEval(t, "if (false) { 123 }", NULL)
	})

	t.Run("returns", func(t *testing.T) {
		RequireEqualEval(t, "return 10;", &object.Integer{10})
		RequireEqualEval(t, "return 1 + 2; false;", &object.Integer{3})
		RequireEqualEval(t, "123; return 1 + 2; false;", &object.Integer{3})
	})

	t.Run("nested returns", func(t *testing.T) {
		input := `
			if (10 > 1) {
				if (10 > 1) {
					return 10;
				}
				return 1;
			}
		`
		RequireEqualEval(t, input, &object.Integer{10})
	})

	t.Run("errors", func(t *testing.T) {
		RequireEvalError(t, "-true", "1:1: unknown operator: -BOOLEAN")
		RequireEvalError(t, "5 + true;", "1:3: type mismatch: INTEGER + BOOLEAN")
		RequireEvalError(t, "5 + true; 5;", "1:3: type mismatch: INTEGER + BOOLEAN")
		RequireEvalError(t, "true + false", "1:6: unknown operator: BOOLEAN + BOOLEAN")
		RequireEvalError(t, "5; true + false; 5", "1:9: unknown operator: BOOLEAN + BOOLEAN")
		RequireEvalError(t, "if (10 > 1) { true + false; }", "1:20: unknown operator: BOOLEAN + BOOLEAN")
	})

	t.Run("let statement", func(t *testing.T) {
		RequireEqualEval(t, "let a = 5; a;", &object.Integer{5})
		RequireEqualEval(t, "let a = 2; let b = 5; let c = a * b; c", &object.Integer{10})
	})

	t.Run("functions", func(t *testing.T) {
		RequireEqualEval(t, "let id = fn(x) { x }; id(1)", &object.Integer{1})
		RequireEqualEval(t, "let add = fn(a, b) { return a + b; }; add(2, 8)", &object.Integer{10})
		RequireEqualEval(t, "let twice = fn(f, x) { return f(f(x)) }; let inc = fn(x) { x + 1}; twice(inc, 0)", &object.Integer{2})
	})

	t.Run("strings", func(t *testing.T) {
		RequireEqualEval(t, `"hello" + "world"`, &object.String{"helloworld"})
		RequireEqualEval(t, `"foo" != "bar"`, TRUE)
		RequireEqualEval(t, `"bbb" > "aaa"`, TRUE)
		RequireEqualEval(t, `"bbb" < "aaa"`, FALSE)
		RequireEqualEval(t, `"bbb" >= "aaa"`, TRUE)
		RequireEqualEval(t, `"aaa" >= "aaa"`, TRUE)
		RequireEqualEval(t, `"aaa" == "aaa"`, TRUE)
	})

	t.Run("builtin", func(t *testing.T) {
		RequireEqualEval(t, `len("hello world")`, &object.Integer{11})
		RequireEqualEval(t, `len("")`, &object.Integer{0})
		RequireEvalError(t, `len(1)`, "1:4: len: invalid argument type INTEGER")
		RequireEvalError(t, `len("one", "two")`, "1:4: len: wrong number of arguments")
		RequireEqualEval(t, `len([])`, &object.Integer{0})
		RequireEqualEval(t, `let x = append([], 1, 2); x[(len(x) - 1)]`, &object.Integer{2})
	})

	t.Run("array", func(t *testing.T) {
		RequireEqualEval(t, "[1, 2, 3]", &object.Array{
			Elements: []object.Object{
				&object.Integer{1},
				&object.Integer{2},
				&object.Integer{3},
			},
		})
		RequireEqualEval(t, "let x = [1]; x[0]", &object.Integer{1})
	})

	t.Run("empty hash", func(t *testing.T) {
		obj, err := ParseEval(t, "{}")
		require.NoError(t, err)
		hash, ok := obj.(*object.Hash)
		require.True(t, ok, "should be hash")
		require.Empty(t, hash.Pairs())
	})

	t.Run("hash", func(t *testing.T) {
		expected := []*object.HashPair{
			{
				Key:   &object.Integer{123},
				Value: TRUE,
			},
		}
		obj, err := ParseEval(t, "{ 123: true }")
		require.NoError(t, err)
		hash, ok := obj.(*object.Hash)
		require.True(t, ok, "should be hash")
		require.Equal(t, expected, hash.Pairs())
	})

	t.Run("index", func(t *testing.T) {
		RequireEqualEval(t, "{}[0]", NULL)
		RequireEqualEval(t, "let x = { true: 123, false: 321 }; x[false]", &object.Integer{321})
		RequireEqualEval(t, `"test"[0]`, &object.String{"t"})
	})

	t.Run("function statement", func(t *testing.T) {
		RequireEqualEval(t, "function add(x, y) { x + y }; add(1, 1)", &object.Integer{2})
	})

	t.Run("assignment", func(t *testing.T) {
		RequireEqualEval(t, "let x = 1; x = 2; x", &object.Integer{2})
		RequireEqualEval(t, "let x = [1]; x[0] = 2; x[0]", &object.Integer{2})
		RequireEqualEval(t, "let x = {}; x[true] = 123; x[true]", &object.Integer{123})
	})

	t.Run("while loop", func(t *testing.T) {
		RequireEqualEval(t, "let x = true; while(x) { x = false }; x", FALSE)
		RequireEqualEval(t, `function foo() { let x = true; while (x) { return "hello"; x = false;  }}; foo()`, &object.String{Value: "hello"})
	})

	t.Run("property access", func(t *testing.T) {
		RequireEqualEval(t, `let x = { "foo": 123 }; x.foo`, &object.Integer{123})
		RequireEqualEval(t, `let x = { "foo": { "bar": true } }; x.foo.bar`, TRUE)
		RequireEqualEval(t, `let x = {}; x.foo = 123; x.foo`, &object.Integer{123})
	})

	t.Run("type checking", func(t *testing.T) {
		RequireEqualEval(t, "fn(x: integer){x}(123)", &object.Integer{123})
		RequireEvalError(t, "fn(x: integer){x}(false)", "1:18: wrong type: expected INTEGER, got BOOLEAN")
		RequireEvalError(t, "let x: boolean = 123", "1:1: wrong type: expected BOOLEAN, got INTEGER")
		RequireEqualEval(t, "let x: integer = 123; x", &object.TypedObject{ObjectType: object.INTEGER, Object: &object.Integer{123}})
		RequireEvalError(t, "let x: boolean = false; x = 123", "1:27: wrong type: expected BOOLEAN, got INTEGER")
	})

	t.Run("in expression", func(t *testing.T) {
		RequireEqualEval(t, "1 in [1]", TRUE)
		RequireEqualEval(t, "0 in [1]", FALSE)
		RequireEqualEval(t, "true in { true: 123 }", TRUE)
		RequireEqualEval(t, "true in {}", FALSE)
	})

	t.Run("simple switch", func(t *testing.T) {
		input := `
			let x = null;
			switch "yes" {
			case "yes":
				x = true
			case "no":
				x = false;
			}
			x;
		`
		RequireEqualEval(t, input, TRUE)
	})

	t.Run("switch with default", func(t *testing.T) {
		input := `
			let x = null;
			switch "maybe" {
			case "yes":
			case "no":
			default:
				x = true
			}
			x;
		`
		RequireEqualEval(t, input, TRUE)
	})

	t.Run("return from switch case", func(t *testing.T) {

		input := `
			function foo() {
				switch true {
				case true:
					return "hello";
				}
			}
			foo()
		`
		RequireEqualEval(t, input, &object.String{Value: "hello"})
	})

	t.Run("return from switch default", func(t *testing.T) {

		input := `
			function foo() {
				switch true {
				default:
					return "hello";
				}
			}
			foo()
		`
		RequireEqualEval(t, input, &object.String{Value: "hello"})
	})
}

func ParseEval(t *testing.T, input string) (object.Object, error) {
	t.Helper()
	program, err := parser.Parse(input)
	require.NoError(t, err)
	return Eval(program, object.NewEnv(nil))
}

func RequireEvalError(t *testing.T, input string, message string) {
	t.Helper()
	_, err := ParseEval(t, input)
	require.EqualError(t, err, message)
}

func RequireEqualEval(t *testing.T, input string, expected object.Object) {
	t.Helper()
	actual, err := ParseEval(t, input)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}
