package lexer

import (
	"testing"

	"github.com/icholy/monkey/token"
)

func ExpectTokens(t *testing.T, input string, expected []token.Token) {
	l := New(input)
	for i, e := range expected {
		tok := l.NextToken()
		if tok != e {
			t.Fatalf("test[%d] - wrong token. want=%s, got=%s", i, e, tok)
		}
	}
}

func TestNextToken(t *testing.T) {

	t.Run("single char", func(t *testing.T) {
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

	t.Run("keywords", func(t *testing.T) {
		input := `
			let five = 5;
			let ten = 10;
			let add = fn(x, y) {
				x + y
			};
			let result = add(five, ten);
		`
		ExpectTokens(t, input, []token.Token{
			{token.LET, "let"},
			{token.IDENT, "five"},
			{token.ASSIGN, "="},
			{token.INT, "5"},
			{token.SEMICOLON, ";"},
			{token.LET, "let"},
			{token.IDENT, "ten"},
			{token.ASSIGN, "="},
			{token.INT, "10"},
			{token.SEMICOLON, ";"},
			{token.LET, "let"},
			{token.IDENT, "add"},
			{token.ASSIGN, "="},
			{token.FUNCTION, "fn"},
			{token.LPAREN, "("},
			{token.IDENT, "x"},
			{token.COMMA, ","},
			{token.IDENT, "y"},
			{token.RPAREN, ")"},
			{token.LBRACE, "{"},
			{token.IDENT, "x"},
			{token.PLUS, "+"},
			{token.IDENT, "y"},
			{token.RBRACE, "}"},
			{token.SEMICOLON, ";"},
			{token.LET, "let"},
			{token.IDENT, "result"},
			{token.ASSIGN, "="},
			{token.IDENT, "add"},
			{token.LPAREN, "("},
			{token.IDENT, "five"},
			{token.COMMA, ","},
			{token.IDENT, "ten"},
			{token.RPAREN, ")"},
			{token.SEMICOLON, ";"},
			{token.EOF, ""},
		})
	})

	t.Run("one character operators", func(t *testing.T) {
		input := `!-/*5`
		ExpectTokens(t, input, []token.Token{
			{token.BANG, "!"},
			{token.MINUS, "-"},
			{token.SLASH, "/"},
			{token.ASTERISK, "*"},
			{token.INT, "5"},
		})
	})

}
