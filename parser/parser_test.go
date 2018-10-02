package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/icholy/monkey/ast"
	"github.com/icholy/monkey/lexer"
	"github.com/icholy/monkey/token"
)

func TestMonkey(t *testing.T) {

	t.Run("let statements", func(t *testing.T) {

		input := `
			let x = 5;
			let y = 10;
			let foobar = 838383;
		`

		RequireEqualAST(t, input, &ast.Program{
			Statements: []ast.Statement{
				&ast.LetStatement{
					Token: token.Token{token.LET, "let"},
					Name: &ast.Identifier{
						Token: token.Token{token.IDENT, "x"},
						Value: "x",
					},
					Value: &ast.IntegerLiteral{
						Token: token.Token{token.INT, "5"},
						Value: 5,
					},
				},
				&ast.LetStatement{
					Token: token.Token{token.LET, "let"},
					Name: &ast.Identifier{
						Token: token.Token{token.IDENT, "y"},
						Value: "y",
					},
					Value: &ast.IntegerLiteral{
						Token: token.Token{token.INT, "10"},
						Value: 10,
					},
				},
				&ast.LetStatement{
					Token: token.Token{token.LET, "let"},
					Name: &ast.Identifier{
						Token: token.Token{token.IDENT, "foobar"},
						Value: "foobar",
					},
					Value: &ast.IntegerLiteral{
						Token: token.Token{token.INT, "838383"},
						Value: 838383,
					},
				},
			},
		})
	})

	t.Run("return", func(t *testing.T) {
		input := `
			return 5;
			return foo;
			return (true);
		`
		RequireEqualAST(t, input, &ast.Program{
			Statements: []ast.Statement{
				&ast.ReturnStatement{
					Token: token.Token{token.RETURN, "return"},
					ReturnValue: &ast.IntegerLiteral{
						Token: token.Token{token.INT, "5"},
						Value: 5,
					},
				},
				&ast.ReturnStatement{
					Token: token.Token{token.RETURN, "return"},
					ReturnValue: &ast.Identifier{
						Token: token.Token{token.IDENT, "foo"},
						Value: "foo",
					},
				},
				&ast.ReturnStatement{
					Token: token.Token{token.RETURN, "return"},
					ReturnValue: &ast.BooleanExpression{
						Token: token.Token{token.TRUE, "true"},
						Value: true,
					},
				},
			},
		})
	})

	t.Run("identifier expression", func(t *testing.T) {
		input := `foobar; 5;`
		RequireEqualAST(t, input, &ast.Program{
			Statements: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{token.IDENT, "foobar"},
					Expression: &ast.Identifier{
						Token: token.Token{token.IDENT, "foobar"},
						Value: "foobar",
					},
				},
				&ast.ExpressionStatement{
					Token: token.Token{token.INT, "5"},
					Expression: &ast.IntegerLiteral{
						Token: token.Token{token.INT, "5"},
						Value: 5,
					},
				},
			},
		})
	})

	t.Run("prefix expression", func(t *testing.T) {
		input := "!5; - foo;"
		RequireEqualAST(t, input, &ast.Program{
			Statements: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{token.BANG, "!"},
					Expression: &ast.PrefixExpression{
						Token:    token.Token{token.BANG, "!"},
						Operator: "!",
						Right: &ast.IntegerLiteral{
							Token: token.Token{token.INT, "5"},
							Value: 5,
						},
					},
				},
				&ast.ExpressionStatement{
					Token: token.Token{token.MINUS, "-"},
					Expression: &ast.PrefixExpression{
						Token:    token.Token{token.MINUS, "-"},
						Operator: "-",
						Right: &ast.Identifier{
							Token: token.Token{token.IDENT, "foo"},
							Value: "foo",
						},
					},
				},
			},
		})
	})

	t.Run("infix expressions", func(t *testing.T) {

		tests := []struct {
			operator    string
			opTokenType token.TokenType
		}{
			{"+", token.PLUS},
			{"-", token.MINUS},
			{">", token.GT},
			{"<", token.LT},
			{"!=", token.NE},
			{"==", token.EQ},
		}

		for _, tt := range tests {
			five := token.Token{token.INT, "5"}
			input := fmt.Sprintf("5 %s 5", tt.operator)
			t.Run(input, func(t *testing.T) {
				RequireEqualAST(t, input, &ast.Program{
					Statements: []ast.Statement{
						&ast.ExpressionStatement{
							Token: five,
							Expression: &ast.InfixExpression{
								Token: token.Token{tt.opTokenType, tt.operator},
								Left: &ast.IntegerLiteral{
									Token: five,
									Value: 5,
								},
								Operator: tt.operator,
								Right: &ast.IntegerLiteral{
									Token: five,
									Value: 5,
								},
							},
						},
					},
				})
			})
		}
	})

	t.Run("operator precedence", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{"-a * b", "((-a) * b)"},
			{"!-a", "(!(-a))"},
			{"a + b + c", "((a + b) + c)"},
			{"a + b - c", "((a + b) - c)"},
			{"a + b * c", "(a + (b * c))"},
			{"3 > 5 == true", "((3 > 5) == true)"},
			{"true != false", "(true != false)"},
			{"(3 + b) * foo", "((3 + b) * foo)"},
		}

		for _, tt := range tests {
			t.Run(tt.input, func(t *testing.T) {
				RequireEqualString(t, tt.input, tt.expected)
			})
		}
	})

	t.Run("boolean expression", func(t *testing.T) {
		input := "true; false"
		RequireEqualAST(t, input, &ast.Program{
			Statements: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{token.TRUE, "true"},
					Expression: &ast.BooleanExpression{
						Token: token.Token{token.TRUE, "true"},
						Value: true,
					},
				},
				&ast.ExpressionStatement{
					Token: token.Token{token.FALSE, "false"},
					Expression: &ast.BooleanExpression{
						Token: token.Token{token.FALSE, "false"},
						Value: false,
					},
				},
			},
		})
	})

	t.Run("if expressions", func(t *testing.T) {
		input := "if (true) { x } else { foo }"
		RequireEqualAST(t, input, &ast.Program{
			Statements: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{token.IF, "if"},
					Expression: &ast.IfExpression{
						Token: token.Token{token.IF, "if"},
						Condition: &ast.BooleanExpression{
							Token: token.Token{token.TRUE, "true"},
							Value: true,
						},
						Concequence: &ast.BlockStatement{
							Token: token.Token{token.LBRACE, "{"},
							Statements: []ast.Statement{
								&ast.ExpressionStatement{
									Token: token.Token{token.IDENT, "x"},
									Expression: &ast.Identifier{
										Token: token.Token{token.IDENT, "x"},
										Value: "x",
									},
								},
							},
						},
						Alternative: &ast.BlockStatement{
							Token: token.Token{token.LBRACE, "{"},
							Statements: []ast.Statement{
								&ast.ExpressionStatement{
									Token: token.Token{token.IDENT, "foo"},
									Expression: &ast.Identifier{
										Token: token.Token{token.IDENT, "foo"},
										Value: "foo",
									},
								},
							},
						},
					},
				},
			},
		})
	})

	t.Run("function literal", func(t *testing.T) {
		input := "fn(x, y) { x }"
		RequireEqualAST(t, input, &ast.Program{
			Statements: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{token.FUNCTION, "fn"},
					Expression: &ast.FunctionLiteral{
						Token: token.Token{token.FUNCTION, "fn"},
						Parameters: []*ast.Identifier{
							{
								Token: token.Token{token.IDENT, "x"},
								Value: "x",
							},
							{
								Token: token.Token{token.IDENT, "y"},
								Value: "y",
							},
						},
						Body: &ast.BlockStatement{
							Token: token.Token{token.LBRACE, "{"},
							Statements: []ast.Statement{
								&ast.ExpressionStatement{
									Token: token.Token{token.IDENT, "x"},
									Expression: &ast.Identifier{
										Token: token.Token{token.IDENT, "x"},
										Value: "x",
									},
								},
							},
						},
					},
				},
			},
		})
	})

	t.Run("call expression", func(t *testing.T) {
		input := "foo(x, 1)"
		RequireEqualAST(t, input, &ast.Program{
			Statements: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{token.IDENT, "foo"},
					Expression: &ast.CallExpression{
						Token: token.Token{token.LPAREN, "("},
						Function: &ast.Identifier{
							Token: token.Token{token.IDENT, "foo"},
							Value: "foo",
						},
						Arguments: []ast.Expression{
							&ast.Identifier{
								Token: token.Token{token.IDENT, "x"},
								Value: "x",
							},
							&ast.IntegerLiteral{
								Token: token.Token{token.INT, "1"},
								Value: 1,
							},
						},
					},
				},
			},
		})
	})
}

func RequireEqualString(t *testing.T, input, expected string) {
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	require.Empty(t, p.Errors(), "parser errors")
	require.Equal(t, expected, program.String())
}

func RequireEqualAST(t *testing.T, input string, expected *ast.Program) {
	l := lexer.New(input)
	p := New(l)
	actual := p.ParseProgram()
	require.Empty(t, p.Errors(), "parser errors")
	require.Equal(t, expected, actual)
}