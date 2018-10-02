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
		RequireEqualEval(t, "true", TRUE)
		RequireEqualEval(t, "false", FALSE)
		RequireEqualEval(t, "!true", FALSE)
		RequireEqualEval(t, "!!true", TRUE)
		RequireEqualEval(t, "!false", TRUE)
		RequireEqualEval(t, "!!false", FALSE)
	})
}

func RequireEqualEval(t *testing.T, input string, expected object.Object) {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	require.Empty(t, p.Errors(), "parser error")
	actual := Eval(program)
	require.Equal(t, expected, actual)
}
