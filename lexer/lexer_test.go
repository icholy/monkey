package lexer

import (
	"testing"

	"github.com/icholy/monkey/token"
)

func ExpectTokens(t *testing.T, input string, expected []token.Token) {
	l := New(input)
	for i, e := range expected {
		tok := l.NextToken()
		tok.Pos = token.Pos{}
		if tok != e {
			t.Fatalf("test[%d] - wrong token. want=%s, got=%s", i, e, tok)
		}
	}
}

func TestNextToken(t *testing.T) {

	t.Run("single char", func(t *testing.T) {
		input := `=+(){},;`
		ExpectTokens(t, input, []token.Token{
			token.New(token.ASSIGN, "="),
			token.New(token.PLUS, "+"),
			token.New(token.LPAREN, "("),
			token.New(token.RPAREN, ")"),
			token.New(token.LBRACE, "{"),
			token.New(token.RBRACE, "}"),
			token.New(token.COMMA, ","),
			token.New(token.SEMICOLON, ";"),
			token.New(token.EOF, ""),
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
			token.New(token.LET, "let"),
			token.New(token.IDENT, "five"),
			token.New(token.ASSIGN, "="),
			token.New(token.INT, "5"),
			token.New(token.SEMICOLON, ";"),
			token.New(token.LET, "let"),
			token.New(token.IDENT, "ten"),
			token.New(token.ASSIGN, "="),
			token.New(token.INT, "10"),
			token.New(token.SEMICOLON, ";"),
			token.New(token.LET, "let"),
			token.New(token.IDENT, "add"),
			token.New(token.ASSIGN, "="),
			token.New(token.FN, "fn"),
			token.New(token.LPAREN, "("),
			token.New(token.IDENT, "x"),
			token.New(token.COMMA, ","),
			token.New(token.IDENT, "y"),
			token.New(token.RPAREN, ")"),
			token.New(token.LBRACE, "{"),
			token.New(token.IDENT, "x"),
			token.New(token.PLUS, "+"),
			token.New(token.IDENT, "y"),
			token.New(token.RBRACE, "}"),
			token.New(token.SEMICOLON, ";"),
			token.New(token.LET, "let"),
			token.New(token.IDENT, "result"),
			token.New(token.ASSIGN, "="),
			token.New(token.IDENT, "add"),
			token.New(token.LPAREN, "("),
			token.New(token.IDENT, "five"),
			token.New(token.COMMA, ","),
			token.New(token.IDENT, "ten"),
			token.New(token.RPAREN, ")"),
			token.New(token.SEMICOLON, ";"),
			token.New(token.EOF, ""),
		})
	})

	t.Run("one character operators", func(t *testing.T) {
		input := `<!-/*5>`
		ExpectTokens(t, input, []token.Token{
			token.New(token.LT, "<"),
			token.New(token.BANG, "!"),
			token.New(token.MINUS, "-"),
			token.New(token.SLASH, "/"),
			token.New(token.ASTERISK, "*"),
			token.New(token.INT, "5"),
			token.New(token.GT, ">"),
			token.New(token.EOF, ""),
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
			token.New(token.IF, "if"),
			token.New(token.LPAREN, "("),
			token.New(token.TRUE, "true"),
			token.New(token.RPAREN, ")"),
			token.New(token.LBRACE, "{"),
			token.New(token.RETURN, "return"),
			token.New(token.FALSE, "false"),
			token.New(token.SEMICOLON, ";"),
			token.New(token.RBRACE, "}"),
			token.New(token.ELSE, "else"),
			token.New(token.LBRACE, "{"),
			token.New(token.RETURN, "return"),
			token.New(token.INT, "5"),
			token.New(token.SEMICOLON, ";"),
			token.New(token.RBRACE, "}"),
			token.New(token.EOF, ""),
		})
	})

	t.Run("two char operators", func(t *testing.T) {
		input := `if (x == 10) { y != 3; }`
		ExpectTokens(t, input, []token.Token{
			token.New(token.IF, "if"),
			token.New(token.LPAREN, "("),
			token.New(token.IDENT, "x"),
			token.New(token.EQ, "=="),
			token.New(token.INT, "10"),
			token.New(token.RPAREN, ")"),
			token.New(token.LBRACE, "{"),
			token.New(token.IDENT, "y"),
			token.New(token.NE, "!="),
			token.New(token.INT, "3"),
			token.New(token.SEMICOLON, ";"),
			token.New(token.RBRACE, "}"),
			token.New(token.EOF, ""),
		})
	})

	t.Run("strings", func(t *testing.T) {
		ExpectTokens(t, `""`, []token.Token{
			token.New(token.STRING, ""),
			token.New(token.EOF, ""),
		})
		ExpectTokens(t, `"testing testing"`, []token.Token{
			token.New(token.STRING, "testing testing"),
			token.New(token.EOF, ""),
		})
		ExpectTokens(t, `"this is a \" test"`, []token.Token{
			token.New(token.STRING, `this is a " test`),
			token.New(token.EOF, ""),
		})
		ExpectTokens(t, `"\t\n\r"`, []token.Token{
			token.New(token.STRING, "\t\n\r"),
			token.New(token.EOF, ""),
		})
	})

	t.Run("array", func(t *testing.T) {
		ExpectTokens(t, `[1]`, []token.Token{
			token.New(token.LBRACKET, "["),
			token.New(token.INT, "1"),
			token.New(token.RBRACKET, "]"),
			token.New(token.EOF, ""),
		})
	})

	t.Run("hash", func(t *testing.T) {
		ExpectTokens(t, `{ "test": 123 }`, []token.Token{
			token.New(token.LBRACE, "{"),
			token.New(token.STRING, "test"),
			token.New(token.COLON, ":"),
			token.New(token.INT, "123"),
			token.New(token.RBRACE, "}"),
			token.New(token.EOF, ""),
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

	t.Run(">= AND <=", func(t *testing.T) {
		ExpectTokens(t, ">=+<=", []token.Token{
			token.New(token.GT_EQ, ">="),
			token.New(token.PLUS, "+"),
			token.New(token.LT_EQ, "<="),
		})
		ExpectTokens(t, "1 <= 2", []token.Token{
			token.New(token.INT, "1"),
			token.New(token.LT_EQ, "<="),
			token.New(token.INT, "2"),
			token.New(token.EOF, ""),
		})
	})

	t.Run("package", func(t *testing.T) {
		ExpectTokens(t, "package foo", []token.Token{
			token.New(token.PACKAGE, "package"),
			token.New(token.IDENT, "foo"),
			token.New(token.EOF, ""),
		})
	})

	t.Run("|| and &&", func(t *testing.T) {
		ExpectTokens(t, "true || false && true", []token.Token{
			token.New(token.TRUE, "true"),
			token.New(token.OR, "||"),
			token.New(token.FALSE, "false"),
			token.New(token.AND, "&&"),
			token.New(token.TRUE, "true"),
			token.New(token.EOF, ""),
		})
	})

	t.Run("x in []", func(t *testing.T) {
		ExpectTokens(t, "x in []", []token.Token{
			token.New(token.IDENT, "x"),
			token.New(token.IN, "in"),
			token.New(token.LBRACKET, "["),
			token.New(token.RBRACKET, "]"),
			token.New(token.EOF, ""),
		})
	})

}
