package lexer

import (
	"strings"

	"github.com/icholy/monkey/token"
)

func New(input string) *Lexer {
	return &Lexer{
		input: input,
		ch:    input[0],
	}
}

type Lexer struct {
	input string
	pos   int
	ch    byte
}

func (l *Lexer) peek() byte {
	next := l.pos + 1
	if next >= len(l.input) {
		return 0
	}
	return l.input[next]
}

func (l *Lexer) read() {
	l.pos++
	if l.pos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.pos]
	}
}

func (l *Lexer) charToken(typ token.TokenType) token.Token {
	return token.Token{Type: typ, Text: string(l.ch)}
}

var bytetokens = map[byte]token.TokenType{
	';': token.SEMICOLON,
	':': token.COLON,
	'(': token.LPAREN,
	')': token.RPAREN,
	'{': token.LBRACE,
	'}': token.RBRACE,
	'[': token.LBRACKET,
	']': token.RBRACKET,
	'+': token.PLUS,
	'-': token.MINUS,
	'*': token.ASTERISK,
	'/': token.SLASH,
	',': token.COMMA,
	'<': token.LT,
	'>': token.GT,
	0:   token.EOF,
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.whitespace()

	if typ, ok := bytetokens[l.ch]; ok {
		if typ == token.EOF {
			tok = token.Token{Type: token.EOF}
		} else {
			tok = token.NewByte(typ, l.ch)
		}
		l.read()
		return tok
	}

	switch l.ch {
	case '=':
		if l.peek() == '=' {
			l.read()
			tok.Type = token.EQ
			tok.Text = "=="
		} else {
			tok = l.charToken(token.ASSIGN)
		}
	case '!':
		if l.peek() == '=' {
			l.read()
			tok.Type = token.NE
			tok.Text = "!="
		} else {
			tok = l.charToken(token.BANG)
		}
	case '"':
		tok.Text = l.str()
		tok.Type = token.STRING
		return tok
	default:
		if isLetter(l.ch) {
			tok.Text = l.ident()
			tok.Type = token.LookupIdent(tok.Text)
			return tok
		}
		if isDigit(l.ch) {
			tok.Text = l.number()
			tok.Type = token.INT
			return tok
		}
		tok = l.charToken(token.ILLEGAL)
	}
	l.read()
	return tok
}

func (l *Lexer) whitespace() {
	for isWhitespace(l.ch) {
		l.read()
	}
}

func (l *Lexer) str() string {
	l.read()
	var escaped bool
	var b strings.Builder
	for l.ch != 0 {
		if escaped {
			switch l.ch {
			case 't':
				b.WriteByte('\t')
			case 'r':
				b.WriteByte('\r')
			case 'n':
				b.WriteByte('\n')
			default:
				b.WriteByte(l.ch)
			}
			escaped = false
		} else {
			if l.ch == '"' {
				break
			}
			if l.ch == '\\' {
				escaped = true
			} else {
				b.WriteByte(l.ch)
			}
		}
		l.read()
	}
	l.read()
	return b.String()
}

func (l *Lexer) ident() string {
	start := l.pos
	for isLetter(l.ch) {
		l.read()
	}
	return l.input[start:l.pos]
}

func (l *Lexer) number() string {
	start := l.pos
	for isDigit(l.ch) {
		l.read()
	}
	return l.input[start:l.pos]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
