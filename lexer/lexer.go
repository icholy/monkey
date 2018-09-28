package lexer

import (
	"github.com/icholy/monkey/token"
)

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) charToken(typ token.TokenType) token.Token {
	return token.Token{Type: typ, Literal: string(l.ch)}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	switch l.ch {
	case '=':
		tok = l.charToken(token.ASSIGN)
	case ';':
		tok = l.charToken(token.SEMICOLON)
	case '(':
		tok = l.charToken(token.LPAREN)
	case ')':
		tok = l.charToken(token.RPAREN)
	case '{':
		tok = l.charToken(token.LBRACE)
	case '}':
		tok = l.charToken(token.RBRACE)
	case '+':
		tok = l.charToken(token.PLUS)
	case ',':
		tok = l.charToken(token.COMMA)
	case 0:
		tok = token.Token{Type: token.EOF}
	default:
		tok = l.charToken(token.ILLEGAL)
	}
	l.readChar()
	return tok
}

func newToken(typ token.TokenType, ch byte) token.Token {
	return token.Token{Type: typ, Literal: string(ch)}
}
