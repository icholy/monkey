package lexer

import (
	"testing"

	"github.com/icholy/monkey/token"
)

func TestNextToken(t *testing.T) {
	input := `=+(){},;`

	tests := []token.Token{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.Type {
			t.Fatalf("test[%d] - wrong TokenType. want=%v, got=%v", i, tt.Type, tok.Type)
		}
		if tok.Literal != tt.Literal {
			t.Fatalf("test[%d] - wrong Literal. want=%v, got=%v", i, tt.Literal, tok.Literal)
		}
	}
}
