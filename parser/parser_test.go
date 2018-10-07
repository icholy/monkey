package parser

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sanity-io/litter"
	"github.com/stretchr/testify/require"

	"github.com/icholy/monkey/ast"
	"github.com/icholy/monkey/token"
)

func TestMonkey(t *testing.T) {

	t.Run("let statements", func(t *testing.T) {

		input := `
			let x = 5;
			let y = 10;
			let foobar = 838383;
			let str = "testing";
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
				&ast.LetStatement{
					Token: token.Token{token.LET, "let"},
					Name: &ast.Identifier{
						Token: token.Token{token.IDENT, "str"},
						Value: "str",
					},
					Value: &ast.StringLiteral{
						Token: token.Token{token.STRING, "testing"},
						Value: "testing",
					},
				},
			},
		})
	})

	t.Run("return", func(t *testing.T) {
		input := `
			return;
			return 5;
			return foo;
			return (true);
		`
		RequireEqualAST(t, input, &ast.Program{
			Statements: []ast.Statement{
				&ast.ReturnStatement{
					Token: token.Token{token.RETURN, "return"},
				},
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
			{">=", token.GT_EQ},
			{"<=", token.LT_EQ},
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
					Token: token.New(token.FN, "fn"),
					Expression: &ast.FunctionLiteral{
						Token: token.Token{token.FN, "fn"},
						Parameters: []*ast.Parameter{
							{
								Token: token.Token{token.IDENT, "x"},
								Name: &ast.Identifier{
									Token: token.Token{token.IDENT, "x"},
									Value: "x",
								},
							},
							{
								Token: token.Token{token.IDENT, "y"},
								Name: &ast.Identifier{
									Token: token.Token{token.IDENT, "y"},
									Value: "y",
								},
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

	t.Run("strings", func(t *testing.T) {
		input := `"hello" + "world"`
		RequireEqualAST(t, input, &ast.Program{
			Statements: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{token.STRING, "hello"},
					Expression: &ast.InfixExpression{
						Token:    token.Token{token.PLUS, "+"},
						Operator: "+",
						Left: &ast.StringLiteral{
							Token: token.Token{token.STRING, "hello"},
							Value: "hello",
						},
						Right: &ast.StringLiteral{
							Token: token.Token{token.STRING, "world"},
							Value: "world",
						},
					},
				},
			},
		})
	})

	t.Run("arrays", func(t *testing.T) {
		input := `["test", 1, hello]`
		RequireEqualAST(t, input, &ast.Program{
			Statements: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{token.LBRACKET, "["},
					Expression: &ast.ArrayLiteral{
						Token: token.Token{token.LBRACKET, "["},
						Elements: []ast.Expression{
							&ast.StringLiteral{
								Token: token.Token{token.STRING, "test"},
								Value: "test",
							},
							&ast.IntegerLiteral{
								Token: token.Token{token.INT, "1"},
								Value: 1,
							},
							&ast.Identifier{
								Token: token.Token{token.IDENT, "hello"},
								Value: "hello",
							},
						},
					},
				},
			},
		})
	})

	t.Run("index", func(t *testing.T) {
		input := `foo[123]`
		RequireEqualAST(t, input, &ast.Program{
			Statements: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{token.IDENT, "foo"},
					Expression: &ast.IndexExpression{
						Token: token.Token{token.LBRACKET, "["},
						Value: &ast.Identifier{
							Token: token.Token{token.IDENT, "foo"},
							Value: "foo",
						},
						Index: &ast.IntegerLiteral{
							Token: token.Token{token.INT, "123"},
							Value: 123,
						},
					},
				},
			},
		})
	})

	t.Run("empty hash", func(t *testing.T) {
		input := `{}`
		RequireEqualAST(t, input, &ast.Program{
			Statements: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{token.LBRACE, "{"},
					Expression: &ast.HashLiteral{
						Token: token.Token{token.LBRACE, "{"},
						Pairs: nil,
					},
				},
			},
		})
	})

	t.Run("hash", func(t *testing.T) {
		input := `{ true: 1, 1:"yes" }`
		RequireEqualAST(t, input, &ast.Program{
			Statements: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{token.LBRACE, "{"},
					Expression: &ast.HashLiteral{
						Token: token.Token{token.LBRACE, "{"},
						Pairs: []*ast.HashPair{
							{
								Key: &ast.BooleanExpression{
									Token: token.Token{token.TRUE, "true"},
									Value: true,
								},
								Value: &ast.IntegerLiteral{
									Token: token.Token{token.INT, "1"},
									Value: 1,
								},
							},
							{
								Key: &ast.IntegerLiteral{
									Token: token.Token{token.INT, "1"},
									Value: 1,
								},
								Value: &ast.StringLiteral{
									Token: token.Token{token.STRING, "yes"},
									Value: "yes",
								},
							},
						},
					},
				},
			},
		})
	})

	t.Run("function statement", func(t *testing.T) {
		input := `
			function foo(x) {}

			x()
		`
		RequireEqualAST(t, input, &ast.Program{
			Statements: []ast.Statement{
				&ast.FunctionStatement{
					Token: token.New(token.FUNCTION, "function"),
					Name: &ast.Identifier{
						Token: token.New(token.IDENT, "foo"),
						Value: "foo",
					},
					Parameters: []*ast.Parameter{
						{
							Token: token.New(token.IDENT, "x"),
							Name: &ast.Identifier{
								Token: token.New(token.IDENT, "x"),
								Value: "x",
							},
						},
					},
					Body: &ast.BlockStatement{
						Token: token.New(token.LBRACE, "{"),
					},
				},
				&ast.ExpressionStatement{
					Token: token.New(token.IDENT, "x"),
					Expression: &ast.CallExpression{
						Token: token.New(token.LPAREN, "("),
						Function: &ast.Identifier{
							Token: token.New(token.IDENT, "x"),
							Value: "x",
						},
					},
				},
			},
		})
	})

	t.Run("import", func(t *testing.T) {
		RequireEqualAST(t, `import "foo.monkey"`, &ast.Program{
			Statements: []ast.Statement{
				&ast.ImportStatement{
					Token: token.New(token.IMPORT, "import"),
					Value: "foo.monkey",
				},
			},
		})
	})

	t.Run("assignment", func(t *testing.T) {
		RequireEqualAST(t, "foo = 1", &ast.Program{
			Statements: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.New(token.IDENT, "foo"),
					Expression: &ast.AssignmentExpression{
						Token: token.New(token.ASSIGN, "="),
						Left: &ast.Identifier{
							Token: token.New(token.IDENT, "foo"),
							Value: "foo",
						},
						Value: &ast.IntegerLiteral{
							Token: token.New(token.INT, "1"),
							Value: 1,
						},
					},
				},
			},
		})

		RequireEqualAST(t, "foo[123] = bar()", &ast.Program{
			Statements: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.New(token.IDENT, "foo"),
					Expression: &ast.AssignmentExpression{
						Token: token.New(token.ASSIGN, "="),
						Left: &ast.IndexExpression{
							Token: token.New(token.LBRACKET, "["),
							Value: &ast.Identifier{
								Token: token.New(token.IDENT, "foo"),
								Value: "foo",
							},
							Index: &ast.IntegerLiteral{
								Token: token.New(token.INT, "123"),
								Value: 123,
							},
						},
						Value: &ast.CallExpression{
							Token: token.New(token.LPAREN, "("),
							Function: &ast.Identifier{
								Token: token.New(token.IDENT, "bar"),
								Value: "bar",
							},
							Arguments: nil,
						},
					},
				},
			},
		})
	})

	t.Run("while loop", func(t *testing.T) {
		RequireEqualAST(t, "while (1 >= x) {}", &ast.Program{
			Statements: []ast.Statement{
				&ast.WhileStatement{
					Token: token.New(token.WHILE, "while"),
					Condition: &ast.InfixExpression{
						Token:    token.New(token.GT_EQ, ">="),
						Operator: ">=",
						Left: &ast.IntegerLiteral{
							Token: token.New(token.INT, "1"),
							Value: 1,
						},
						Right: &ast.Identifier{
							Token: token.New(token.IDENT, "x"),
							Value: "x",
						},
					},
					Body: &ast.BlockStatement{
						Token: token.New(token.LBRACE, "{"),
					},
				},
			},
		})
	})

	t.Run("property access", func(t *testing.T) {
		RequireEqualAST(t, "foo.bar", &ast.Program{
			Statements: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.New(token.IDENT, "foo"),
					Expression: &ast.PropertyExpression{
						Token: token.New(token.DOT, "."),
						Value: &ast.Identifier{
							Token: token.New(token.IDENT, "foo"),
							Value: "foo",
						},
						Name: &ast.Identifier{
							Token: token.New(token.IDENT, "bar"),
							Value: "bar",
						},
					},
				},
			},
		})
	})

	t.Run("package", func(t *testing.T) {
		RequireEqualAST(t, "package foo", &ast.Program{
			Statements: []ast.Statement{
				&ast.PackageStatement{
					Token: token.New(token.PACKAGE, "package"),
					Name: &ast.Identifier{
						Token: token.New(token.IDENT, "foo"),
						Value: "foo",
					},
				},
			},
		})
	})

}

func RequireEqualString(t *testing.T, input, expected string) {
	program, err := Parse(input)
	require.NoError(t, err)
	require.Equal(t, expected, program.String())
}

func RequireEqualAST(t *testing.T, input string, expected *ast.Program) {
	actual, err := Parse(input)
	require.NoError(t, err)
	if !cmp.Equal(expected, actual) {
		litter.Dump(expected)
		litter.Dump(actual)
		t.Fatal(cmp.Diff(expected, actual))
	}
}
