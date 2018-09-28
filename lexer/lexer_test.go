package lexer

import (
	"testing"

	"github.com/icholy/monkey/token"
)

func ExpectTokens(t *testing.T, input string, expected []token.Token) {
	l := New(input)
	for i, e := range expected {
		tok := l.NextToken()
		if tok.Type != e.Type {
			t.Fatalf("test[%d] - wrong TokenType. want=%v, got=%v", i, e.Type, tok.Type)
		}
		if tok.Literal != e.Literal {
			t.Fatalf("test[%d] - wrong Literal. want=%v, got=%v", i, e.Literal, tok.Literal)
		}
	}
}

func TestNextToken(t *testing.T) {

	t.Run("simple", func(t *testing.T) {
		input := `=+(){},;`
		ExpectTokens(t, input, []token.Token{
			{token.ASSIGN, "="},
			{token.PLUS, "+"},
			{token.LPAREN, "("},
			{token.RPAREN, ")"},
			{token.LBRACE, "{"},
			{token.RBRACE, "}"},
			{token.COMMA, ","},
			{token.SEMICOLON, ";"},
			{token.EOF, ""},
		})
	})

}
