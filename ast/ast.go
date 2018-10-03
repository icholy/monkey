package ast

import (
	"fmt"
	"strconv"
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
		fmt.Fprint(&b, s)
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
	if e.Expression == nil {
		return ""
	}
	return e.Expression.String()
}
func (ExpressionStatement) statementNode() {}
func (e *ExpressionStatement) TokenLiteral() string {
	return e.Token.String()
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) String() string {
	return strconv.FormatInt(i.Value, 10)
}

func (IntegerLiteral) expressionNode() {}
func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (StringLiteral) expressionNode() {}
func (s *StringLiteral) TokenLiteral() string {
	return s.Token.Literal
}
func (s *StringLiteral) String() string {
	return fmt.Sprintf("%v", s.Value)
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (p *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", p.Operator, p.Right)
}
func (PrefixExpression) expressionNode() {}
func (p *PrefixExpression) TokenLiteral() string {
	return p.Token.Literal
}

type InfixExpression struct {
	Token    token.Token
	Operator string
	Left     Expression
	Right    Expression
}

func (i *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", i.Left, i.Operator, i.Right)
}
func (InfixExpression) expressionNode() {}
func (i *InfixExpression) TokenLiteral() string {
	return i.Token.Literal
}

type BooleanExpression struct {
	Token token.Token
	Value bool
}

func (b *BooleanExpression) String() string {
	return b.Token.Literal
}

func (BooleanExpression) expressionNode() {}
func (b *BooleanExpression) TokenLiteral() string {
	return b.Token.Literal
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Concequence *BlockStatement
	Alternative *BlockStatement
}

func (i *IfExpression) String() string {
	if i.Alternative == nil {
		return fmt.Sprintf("if (%s) %s", i.Condition, i.Concequence)
	}
	return fmt.Sprintf("if (%s) %s else %s", i.Condition, i.Concequence, i.Alternative)
}
func (IfExpression) expressionNode() {}
func (i *IfExpression) TokenLiteral() string {
	return i.Token.Literal
}

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (ArrayLiteral) expressionNode() {}
func (a *ArrayLiteral) TokenLiteral() string {
	return a.Token.Literal
}
func (a *ArrayLiteral) String() string {
	var values []string
	for _, v := range a.Elements {
		values = append(values, v.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(values, ", "))
}

type IndexExpression struct {
	Token token.Token
	Value Expression
	Index Expression
}

func (i *IndexExpression) expressionNode() {}
func (i *IndexExpression) TokenLiteral() string {
	return i.Token.Literal
}
func (i *IndexExpression) String() string {
	return fmt.Sprintf("%s[%s]", i.Value, i.Index)
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (b *BlockStatement) String() string {
	var sb strings.Builder
	for _, s := range b.Statements {
		fmt.Fprint(&sb, s)
	}
	return sb.String()
}
func (BlockStatement) statementNode() {}
func (b *BlockStatement) TokenLiteral() string {
	return b.Token.Literal
}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (f *FunctionLiteral) ParameterNames() []string {
	var names []string
	for _, p := range f.Parameters {
		names = append(names, p.Value)
	}
	return names
}

func (f *FunctionLiteral) String() string {
	return fmt.Sprintf("fn(%s) %s", strings.Join(f.ParameterNames(), ", "), f.Body)
}
func (FunctionLiteral) expressionNode() {}
func (f *FunctionLiteral) TokenLiteral() string {
	return f.Token.Literal
}

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (c *CallExpression) String() string {
	var args []string
	for _, a := range c.Arguments {
		args = append(args, a.String())
	}
	return fmt.Sprintf("%s(%s)", c.Function, strings.Join(args, ", "))
}
func (CallExpression) expressionNode() {}
func (c *CallExpression) TokenLiteral() string {
	return c.Token.Literal
}

type HashLiteral struct {
	Token token.Token
	Paris map[Expression]Expression
}

func (h *HashLiteral) expressionNode() {}
func (h *HashLiteral) TokenLiteral() string {
	return h.Token.Literal
}
