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
		RequireEqualEval(t, "2 != 2", FALSE)
		RequireEqualEval(t, "1 == 2", FALSE)
		RequireEqualEval(t, "1 + 2 == 3", TRUE)
		RequireEqualEval(t, "true == true", TRUE)
		RequireEqualEval(t, "true == false", FALSE)
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
		RequireEvalError(t, "-true", "unknown operator: -BOOLEAN")
		RequireEvalError(t, "5 + true;", "type mismatch: INTEGER + BOOLEAN")
		RequireEvalError(t, "5 + true; 5;", "type mismatch: INTEGER + BOOLEAN")
		RequireEvalError(t, "true + false", "unknown operator: BOOLEAN + BOOLEAN")
		RequireEvalError(t, "5; true + false; 5", "unknown operator: BOOLEAN + BOOLEAN")
		RequireEvalError(t, "if (10 > 1) { true + false; }", "unknown operator: BOOLEAN + BOOLEAN")
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
	})

	t.Run("builtin", func(t *testing.T) {
		RequireEqualEval(t, `len("hello world")`, &object.Integer{11})
		RequireEqualEval(t, `len("")`, &object.Integer{0})
		RequireEvalError(t, `len(1)`, "len: invalid argument type INTEGER")
		RequireEvalError(t, `len("one", "two")`, "len: wrong number of arguments")
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

	t.Run("hash index", func(t *testing.T) {
		RequireEqualEval(t, "{}[0]", NULL)
		RequireEqualEval(t, "let x = { true: 123, false: 321 }; x[false]", &object.Integer{321})
	})

	t.Run("function statement", func(t *testing.T) {
		RequireEqualEval(t, "function add(x, y) { x + y }; add(1, 1)", &object.Integer{2})
	})

}

func ParseEval(t *testing.T, input string) (object.Object, error) {
	program, err := parser.Parse(input)
	require.NoError(t, err)
	return Eval(program, object.NewEnv(nil))
}

func RequireEvalError(t *testing.T, input string, message string) {
	_, err := ParseEval(t, input)
	require.EqualError(t, err, message)
}

func RequireEqualEval(t *testing.T, input string, expected object.Object) {
	actual, err := ParseEval(t, input)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}
