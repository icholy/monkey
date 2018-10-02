package parser

import (
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
		input := `foobar;`
		RequireEqualAST(t, input, &ast.Program{
			Statements: []ast.Statement{
				&ast.ExpressionStatement{
					Expression: &ast.Identifier{
						Token: token.Token{token.IDENT, "foobar"},
						Value: "foobar",
					},
				},
			},
		})
	})
}

func RequireEqualAST(t *testing.T, input string, expected *ast.Program) {
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	require.Empty(t, p.Errors(), "parser errors")
	require.EqualValues(t, expected, program)
}
