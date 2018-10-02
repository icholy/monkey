package evaluator

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/icholy/monkey/lexer"
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
		RequireEqualEval(t, "-true", &object.Error{"unknown operator: -BOOLEAN"})
		RequireEqualEval(t, "5 + true;", &object.Error{"type mismatch: INTEGER + BOOLEAN"})
		RequireEqualEval(t, "5 + true; 5;", &object.Error{"type mismatch: INTEGER + BOOLEAN"})
		RequireEqualEval(t, "true + false", &object.Error{"unknown operator: BOOLEAN + BOOLEAN"})
		RequireEqualEval(t, "5; true + false; 5", &object.Error{"unknown operator: BOOLEAN + BOOLEAN"})
		RequireEqualEval(t, "if (10 > 1) { true + false; }", &object.Error{"unknown operator: BOOLEAN + BOOLEAN"})
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
}

func RequireEqualEval(t *testing.T, input string, expected object.Object) {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	require.Empty(t, p.Errors(), "parser error")
	env := object.NewEnv(nil)
	actual := Eval(program, env)
	require.Equal(t, expected, actual)
}
