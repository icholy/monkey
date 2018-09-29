package parser

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/icholy/monkey/ast"
	"github.com/icholy/monkey/lexer"
	"github.com/icholy/monkey/token"
)

func TestLetStatement(t *testing.T) {
	input := `
		let x = 5;
		let y = 10;
		let foobar = 838383;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	require.EqualValues(t, program, &ast.Program{
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
}
