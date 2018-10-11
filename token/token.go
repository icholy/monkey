package token

import (
	"fmt"
)

type TokenType string

type Pos struct {
	Line   int
	Offset int
}

func (p Pos) String() string {
	return fmt.Sprintf("%d:%d", p.Line, p.Offset)
}

type Token struct {
	Pos
	Type TokenType
	Text string
}

func New(typ TokenType, text string) Token {
	return Token{Type: typ, Text: text}
}

func NewByte(typ TokenType, text byte) Token {
	return New(typ, string(text))
}

func (t Token) Is(typ TokenType) bool {
	return t.Type == typ
}

func (t Token) String() string {
	return fmt.Sprintf("%s(\"%s\")", t.Type, t.Text)
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT    = "IDENT"
	INT      = "INT"
	ASSIGN   = "ASSIGN"
	PLUS     = "PLUS"
	MINUS    = "MINUS"
	BANG     = "BANG"
	ASTERISK = "ASTERISK"
	SLASH    = "SLASH"
	GT       = "GT"
	LT       = "LT"
	EQ       = "EQ"
	NE       = "NE"
	GT_EQ    = "GT_EQ"
	LT_EQ    = "LT_EQ"
	DOT      = "DOT"
	OR       = "OR"
	AND      = "AND"

	// Delimiters
	COMMA     = "COMMA"
	SEMICOLON = "SEMICOLON"
	COLON     = "COLON"

	LPAREN   = "LPAREN"
	RPAREN   = "RPAREN"
	LBRACE   = "LBRACE"
	RBRACE   = "RBRACE"
	LBRACKET = "LBRACKET"
	RBRACKET = "RBRACKET"
	STRING   = "STRING"

	// Keywords
	FN       = "FN"
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	IMPORT   = "IMPORT"
	WHILE    = "WHILE"
	PACKAGE  = "PACKAGE"
	DEBUGGER = "DEBUGGER"
	NULL     = "NULL"
)

var keywords = map[string]TokenType{
	"fn":       FN,
	"let":      LET,
	"true":     TRUE,
	"false":    FALSE,
	"if":       IF,
	"else":     ELSE,
	"return":   RETURN,
	"function": FUNCTION,
	"import":   IMPORT,
	"while":    WHILE,
	"package":  PACKAGE,
	"debugger": DEBUGGER,
	"null":     NULL,
}

func LookupIdent(ident string) TokenType {
	if typ, ok := keywords[ident]; ok {
		return typ
	}
	return IDENT
}
