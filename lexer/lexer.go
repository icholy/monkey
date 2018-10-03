package lexer

import (
	"strings"

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

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
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
	l.skipWhitespace()
	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok.Type = token.EQ
			tok.Literal = "=="
		} else {
			tok = l.charToken(token.ASSIGN)
		}
	case ';':
		tok = l.charToken(token.SEMICOLON)
	case ':':
		tok = l.charToken(token.COLON)
	case '(':
		tok = l.charToken(token.LPAREN)
	case ')':
		tok = l.charToken(token.RPAREN)
	case '{':
		tok = l.charToken(token.LBRACE)
	case '}':
		tok = l.charToken(token.RBRACE)
	case '[':
		tok = l.charToken(token.LBRACKET)
	case ']':
		tok = l.charToken(token.RBRACKET)
	case '+':
		tok = l.charToken(token.PLUS)
	case '-':
		tok = l.charToken(token.MINUS)
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok.Type = token.NE
			tok.Literal = "!="
		} else {
			tok = l.charToken(token.BANG)
		}
	case '*':
		tok = l.charToken(token.ASTERISK)
	case '/':
		tok = l.charToken(token.SLASH)
	case ',':
		tok = l.charToken(token.COMMA)
	case '<':
		tok = l.charToken(token.LT)
	case '>':
		tok = l.charToken(token.GT)
	case 0:
		tok = token.Token{Type: token.EOF}
	default:
		if l.ch == '"' {
			tok.Literal = l.readString()
			tok.Type = token.STRING
			return tok
		}
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		}
		if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		}
		tok = l.charToken(token.ILLEGAL)
	}
	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		l.readChar()
	}
}

func (l *Lexer) readString() string {
	l.readChar()
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
		l.readChar()
	}
	l.readChar()
	return b.String()
}

func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) readNumber() string {
	start := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
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
