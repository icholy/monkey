package ast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/icholy/monkey/token"
)

type Node interface {
	TokenText() string
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

func (p *Program) TokenText() string {
	if len(p.Statements) == 0 {
		return ""
	}
	return p.Statements[0].TokenText()
}

type PackageStatement struct {
	Token token.Token
	Name  *Identifier
}

func (p *PackageStatement) String() string {
	return fmt.Sprintf("package %s", p.Name)
}

func (p *PackageStatement) statementNode() {}
func (p *PackageStatement) TokenText() string {
	return p.Token.Text
}

type ImportStatement struct {
	Token   token.Token
	Value   string
	Program *Program
}

func (i *ImportStatement) String() string {
	return fmt.Sprintf("import(%s)", i.Value)
}
func (ImportStatement) statementNode() {}
func (i *ImportStatement) TokenText() string {
	return i.Token.Text
}

type Parameter struct {
	Token token.Token
	Name  *Identifier
	Type  *Identifier
}

func (p *Parameter) expressionNode() {}
func (p *Parameter) TokenText() string {
	return p.Token.Text
}
func (p *Parameter) String() string {
	if p.Type != nil {
		return fmt.Sprintf("%s: %s", p.Name, p.Type)
	}
	return p.Name.Value
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) String() string { return i.Value }
func (Identifier) expressionNode()   {}
func (i *Identifier) TokenText() string {
	return i.Token.Text
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Type  *Identifier
	Value Expression
}

func (l *LetStatement) String() string {
	if l.Type != nil {
		return fmt.Sprintf("let %s: %s = %s;", l.Name, l.Type, l.Value)
	}
	return fmt.Sprintf("let %s = %s;", l.Name, l.Value)
}

func (LetStatement) statementNode() {}
func (l *LetStatement) TokenText() string {
	return l.Token.Text
}

type WhileStatement struct {
	Token     token.Token
	Condition Expression
	Body      *BlockStatement
}

func (w *WhileStatement) String() string {
	return fmt.Sprintf("while (%s) { %s}", w.Condition, w.Body)
}

func (WhileStatement) statementNode() {}
func (w *WhileStatement) TokenText() string {
	return w.Token.Text
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (r *ReturnStatement) String() string {
	if r.ReturnValue == nil {
		return fmt.Sprint("return")
	}
	return fmt.Sprintf("return %s", r.ReturnValue)
}
func (ReturnStatement) statementNode() {}
func (r *ReturnStatement) TokenText() string {
	return r.Token.Text
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
func (e *ExpressionStatement) TokenText() string {
	return e.Token.String()
}

type DebuggerStatement struct {
	Token token.Token
}

func (d *DebuggerStatement) String() string {
	return d.TokenText()
}
func (DebuggerStatement) statementNode() {}
func (d *DebuggerStatement) TokenText() string {
	return d.Token.Text
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) String() string {
	return strconv.FormatInt(i.Value, 10)
}

func (IntegerLiteral) expressionNode() {}
func (i *IntegerLiteral) TokenText() string {
	return i.Token.Text
}

type AssignmentExpression struct {
	Token token.Token
	Left  Expression
	Value Expression
}

func (a *AssignmentExpression) String() string {
	return fmt.Sprintf("%s = %s", a.Left, a.Value)
}

func (AssignmentExpression) expressionNode() {}
func (a *AssignmentExpression) TokenText() string {
	return a.Token.Text
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (StringLiteral) expressionNode() {}
func (s *StringLiteral) TokenText() string {
	return s.Token.Text
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
func (p *PrefixExpression) TokenText() string {
	return p.Token.Text
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
func (i *InfixExpression) TokenText() string {
	return i.Token.Text
}

type BooleanExpression struct {
	Token token.Token
	Value bool
}

func (b *BooleanExpression) String() string {
	return b.Token.Text
}

func (BooleanExpression) expressionNode() {}
func (b *BooleanExpression) TokenText() string {
	return b.Token.Text
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Concequence *BlockStatement
	Alternative *BlockStatement
}

func (i *IfExpression) String() string {
	if i.Alternative == nil {
		return fmt.Sprintf("if (%s) {%s}", i.Condition, i.Concequence)
	}
	return fmt.Sprintf("if (%s) {%s} else {%s}", i.Condition, i.Concequence, i.Alternative)
}
func (IfExpression) expressionNode() {}
func (i *IfExpression) TokenText() string {
	return i.Token.Text
}

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (ArrayLiteral) expressionNode() {}
func (a *ArrayLiteral) TokenText() string {
	return a.Token.Text
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
func (i *IndexExpression) TokenText() string {
	return i.Token.Text
}
func (i *IndexExpression) String() string {
	return fmt.Sprintf("%s[%s]", i.Value, i.Index)
}

type PropertyExpression struct {
	Token token.Token
	Value Expression
	Name  *Identifier
}

func (p *PropertyExpression) String() string {
	return fmt.Sprintf("%s.%s", p.Value, p.Name.Value)
}

func (p *PropertyExpression) expressionNode() {}
func (p *PropertyExpression) TokenText() string {
	return p.Token.Text
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (b *BlockStatement) String() string {
	var sb strings.Builder
	for _, s := range b.Statements {
		fmt.Fprintf(&sb, "%s; ", s)
	}
	return sb.String()
}
func (BlockStatement) statementNode() {}
func (b *BlockStatement) TokenText() string {
	return b.Token.Text
}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Parameter
	ReturnType *Identifier
	Body       *BlockStatement
}

func (f *FunctionLiteral) ParameterNames() []string {
	var names []string
	for _, p := range f.Parameters {
		names = append(names, p.Name.Value)
	}
	return names
}

func (f *FunctionLiteral) String() string {
	return fmt.Sprintf("fn(%s) { %s}", strings.Join(f.ParameterNames(), ", "), f.Body)
}
func (FunctionLiteral) expressionNode() {}
func (f *FunctionLiteral) TokenText() string {
	return f.Token.Text
}

type FunctionStatement struct {
	Token      token.Token
	Name       *Identifier
	Parameters []*Parameter
	ReturnType *Identifier
	Body       *BlockStatement
}

func (f *FunctionStatement) ParameterNames() []string {
	var names []string
	for _, p := range f.Parameters {
		names = append(names, p.Name.Value)
	}
	return names
}

func (f *FunctionStatement) String() string {
	return fmt.Sprintf(
		"function %s(%s) %s",
		f.Name,
		strings.Join(f.ParameterNames(), ", "),
		f.Body,
	)
}
func (FunctionStatement) statementNode() {}
func (f *FunctionStatement) TokenText() string {
	return f.Token.Text
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
func (c *CallExpression) TokenText() string {
	return c.Token.Text
}

type HashLiteral struct {
	Token token.Token
	Pairs []*HashPair
}

type HashPair struct {
	Key   Expression
	Value Expression
}

func (hp *HashPair) String() string {
	return fmt.Sprintf("%s: %s", hp.Key, hp.Value)
}

func (h *HashLiteral) expressionNode() {}
func (h *HashLiteral) TokenText() string {
	return h.Token.Text
}
func (h *HashLiteral) String() string {
	var pairs []string
	for _, p := range h.Pairs {
		pairs = append(pairs, p.String())
	}
	return fmt.Sprintf("{ %s }", strings.Join(pairs, ", "))
}
