package token

import (
	"fmt"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

func (t Token) String() string {
	return fmt.Sprintf("%s(\"%s\")", t.Type, t.Literal)
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
	BANG     = "!"
	ASTERISK = "ASTERISK"
	SLASH    = "/"
	GT       = "GT"
	LT       = "LT"

	// Delimiters
	COMMA     = "COMMA"
	SEMICOLON = "SEMICOLON"

	LPAREN = "LPAREN"
	RPAREN = "RPAREN"
	LBRACE = "LBRACE"
	RBRACE = "RBRACE"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"

	// Operators
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func LookupIdent(ident string) TokenType {
	if typ, ok := keywords[ident]; ok {
		return typ
	}
	return IDENT
}
