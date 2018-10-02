package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/icholy/monkey/ast"
	"github.com/icholy/monkey/lexer"
	"github.com/icholy/monkey/token"
)

func TestLetStatement(t *testing.T) {

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
					Value: nil,
				},
				&ast.LetStatement{
					Token: token.Token{token.LET, "let"},
					Name: &ast.Identifier{
						Token: token.Token{token.IDENT, "y"},
						Value: "y",
					},
					Value: nil,
				},
				&ast.LetStatement{
					Token: token.Token{token.LET, "let"},
					Name: &ast.Identifier{
						Token: token.Token{token.IDENT, "foobar"},
						Value: "foobar",
					},
					Value: nil,
				},
			},
		})
	})

	t.Run("return", func(t *testing.T) {
		input := `
			return 5;
			return 10;
			return 993322;
		`
		RequireEqualAST(t, input, &ast.Program{
			Statements: []ast.Statement{
				&ast.ReturnStatement{
					Token:       token.Token{token.RETURN, "return"},
					ReturnValue: nil,
				},
				&ast.ReturnStatement{
					Token:       token.Token{token.RETURN, "return"},
					ReturnValue: nil,
				},
				&ast.ReturnStatement{
					Token:       token.Token{token.RETURN, "return"},
					ReturnValue: nil,
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
