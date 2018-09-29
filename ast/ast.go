package ast

import (
	"fmt"
	"strings"

	"github.com/icholy/monkey/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var b strings.Builder
	for _, s := range p.Statements {
		fmt.Fprintln(&b, s)
	}
	return b.String()
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) == 0 {
		return ""
	}
	return p.Statements[0].TokenLiteral()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) String() string {
	return i.Value
}

func (Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (l *LetStatement) String() string {
	return fmt.Sprintf("let %s = %s;", l.Name, l.Value)
}

func (LetStatement) statementNode() {}
func (l *LetStatement) TokenLiteral() string {
	return l.Token.Literal
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (r *ReturnStatement) String() string {
	return fmt.Sprintf("return %s;", r.ReturnValue)
}
func (ReturnStatement) statementNode() {}
func (r *ReturnStatement) TokenLiteral() string {
	return r.Token.Literal
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (e *ExpressionStatement) String() string {
	return e.Expression.String()
}
func (ExpressionStatement) statementNode() {}
func (e *ExpressionStatement) TokenLiteral() string {
	return e.Token.String()
}
