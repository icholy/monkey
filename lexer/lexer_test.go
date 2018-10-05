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
			{token.FN, "fn"},
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
		input := `<!-/*5>`
		ExpectTokens(t, input, []token.Token{
			{token.LT, "<"},
			{token.BANG, "!"},
			{token.MINUS, "-"},
			{token.SLASH, "/"},
			{token.ASTERISK, "*"},
			{token.INT, "5"},
			{token.GT, ">"},
			{token.EOF, ""},
		})
	})

	t.Run("more keywords", func(t *testing.T) {
		input := `
			if (true) {
				return false;
			} else {
				return 5;
			}
		`
		ExpectTokens(t, input, []token.Token{
			{token.IF, "if"},
			{token.LPAREN, "("},
			{token.TRUE, "true"},
			{token.RPAREN, ")"},
			{token.LBRACE, "{"},
			{token.RETURN, "return"},
			{token.FALSE, "false"},
			{token.SEMICOLON, ";"},
			{token.RBRACE, "}"},
			{token.ELSE, "else"},
			{token.LBRACE, "{"},
			{token.RETURN, "return"},
			{token.INT, "5"},
			{token.SEMICOLON, ";"},
			{token.RBRACE, "}"},
			{token.EOF, ""},
		})
	})

	t.Run("two char operators", func(t *testing.T) {
		input := `if (x == 10) { y != 3; }`
		ExpectTokens(t, input, []token.Token{
			{token.IF, "if"},
			{token.LPAREN, "("},
			{token.IDENT, "x"},
			{token.EQ, "=="},
			{token.INT, "10"},
			{token.RPAREN, ")"},
			{token.LBRACE, "{"},
			{token.IDENT, "y"},
			{token.NE, "!="},
			{token.INT, "3"},
			{token.SEMICOLON, ";"},
			{token.RBRACE, "}"},
			{token.EOF, ""},
		})
	})

	t.Run("strings", func(t *testing.T) {
		ExpectTokens(t, `""`, []token.Token{
			{token.STRING, ""},
			{token.EOF, ""},
		})
		ExpectTokens(t, `"testing testing"`, []token.Token{
			{token.STRING, "testing testing"},
			{token.EOF, ""},
		})
		ExpectTokens(t, `"this is a \" test"`, []token.Token{
			{token.STRING, `this is a " test`},
			{token.EOF, ""},
		})
		ExpectTokens(t, `"\t\n\r"`, []token.Token{
			{token.STRING, "\t\n\r"},
			{token.EOF, ""},
		})
	})

	t.Run("array", func(t *testing.T) {
		ExpectTokens(t, `[1]`, []token.Token{
			{token.LBRACKET, "["},
			{token.INT, "1"},
			{token.RBRACKET, "]"},
			{token.EOF, ""},
		})
	})

	t.Run("hash", func(t *testing.T) {
		ExpectTokens(t, `{ "test": 123 }`, []token.Token{
			token.Token{token.LBRACE, "{"},
			token.Token{token.STRING, "test"},
			token.Token{token.COLON, ":"},
			token.Token{token.INT, "123"},
			token.Token{token.RBRACE, "}"},
			token.Token{token.EOF, ""},
		})
	})

	t.Run("while loop", func(t *testing.T) {
		ExpectTokens(t, `while (true) {}`, []token.Token{
			token.New(token.WHILE, "while"),
			token.New(token.LPAREN, "("),
			token.New(token.TRUE, "true"),
			token.New(token.RPAREN, ")"),
			token.New(token.LBRACE, "{"),
			token.New(token.RBRACE, "}"),
			token.New(token.EOF, ""),
		})
	})

	t.Run("dot access", func(t *testing.T) {
		ExpectTokens(t, `foo.bar()`, []token.Token{
			token.New(token.IDENT, "foo"),
			token.New(token.DOT, "."),
			token.New(token.IDENT, "bar"),
			token.New(token.LPAREN, "("),
			token.New(token.RPAREN, ")"),
			token.New(token.EOF, ""),
		})
	})

}
