package parser

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/icholy/monkey/ast"
	"github.com/icholy/monkey/lexer"
	"github.com/icholy/monkey/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
	INDEX
	ASSIGN
)

type (
	prefixFn func() ast.Expression
	infixFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l *lexer.Lexer

	cur    token.Token
	peek   token.Token
	errors []string

	precedences map[token.TokenType]int
	prefixFns   map[token.TokenType]prefixFn
	infixFns    map[token.TokenType]infixFn
}

func Parse(input string) (*ast.Program, error) {
	l := lexer.New(input)
	p := New(l)
	prog := p.ParseProgram()
	if errs := p.Errors(); len(errs) != 0 {
		return nil, errors.New(errs[0])
	}
	return prog, nil
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
	}
	p.precedences = map[token.TokenType]int{
		token.EQ:       EQUALS,
		token.NE:       EQUALS,
		token.LT:       LESSGREATER,
		token.GT:       LESSGREATER,
		token.PLUS:     SUM,
		token.MINUS:    SUM,
		token.SLASH:    PRODUCT,
		token.ASTERISK: PRODUCT,
		token.LPAREN:   CALL,
		token.LBRACKET: INDEX,
		token.ASSIGN:   ASSIGN,
	}
	p.prefixFns = map[token.TokenType]prefixFn{
		token.IDENT:    p.identExpr,
		token.INT:      p.integerExpr,
		token.STRING:   p.stringLit,
		token.BANG:     p.prefixExpr,
		token.MINUS:    p.prefixExpr,
		token.TRUE:     p.booleanExpr,
		token.FALSE:    p.booleanExpr,
		token.LPAREN:   p.groupesExpr,
		token.LBRACKET: p.arrayExpr,
		token.IF:       p.ifExpr,
		token.FN:       p.fnExpr,
		token.LBRACE:   p.hashExpr,
	}
	p.infixFns = map[token.TokenType]infixFn{
		token.PLUS:     p.infixExpr,
		token.MINUS:    p.infixExpr,
		token.SLASH:    p.infixExpr,
		token.ASTERISK: p.infixExpr,
		token.EQ:       p.infixExpr,
		token.NE:       p.infixExpr,
		token.LT:       p.infixExpr,
		token.GT:       p.infixExpr,
		token.LPAREN:   p.callExpr,
		token.LBRACKET: p.indexExpr,
		token.ASSIGN:   p.assignExpr,
	}
	p.next()
	p.next()
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) precedence(t token.Token) int {
	if p, ok := p.precedences[t.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) peekError(t token.TokenType) {
	p.errorf("expected %s, got %s instead", t, p.peek)
}

func (p *Parser) next() {
	p.cur = p.peek
	p.peek = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	for !p.cur.Is(token.EOF) {
		if stmt := p.stmt(); stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.next()
	}
	return program
}

func (p *Parser) booleanExpr() ast.Expression {
	return &ast.BooleanExpression{
		Token: p.cur,
		Value: p.cur.Is(token.TRUE),
	}
}

func (p *Parser) whileStmt() *ast.WhileStatement {
	while := &ast.WhileStatement{Token: p.cur}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.next()
	while.Condition = p.expression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	p.next()
	while.Body = p.blockStmt()
	return while
}

func (p *Parser) blockStmt() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.cur}
	p.next()
	for !p.cur.Is(token.RBRACE) && !p.cur.Is(token.EOF) {
		stmt := p.stmt()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.next()
	}
	return block
}

func (p *Parser) hashExpr() ast.Expression {
	hash := &ast.HashLiteral{Token: p.cur}

	for !p.peek.Is(token.RBRACE) {
		p.next()
		key := p.expression(LOWEST)
		if !p.expectPeek(token.COLON) {
			return nil
		}
		p.next()
		value := p.expression(LOWEST)
		hash.Pairs = append(hash.Pairs, &ast.HashPair{Key: key, Value: value})

		if !p.peek.Is(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return hash
}

func (p *Parser) ifExpr() ast.Expression {
	expr := &ast.IfExpression{Token: p.cur}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.next()
	expr.Condition = p.expression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	expr.Concequence = p.blockStmt()

	if p.peek.Is(token.ELSE) {
		p.next()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		expr.Alternative = p.blockStmt()
	}

	return expr
}

func (p *Parser) fnParameterIdents() []*ast.Identifier {
	var params []*ast.Identifier
	for p.peek.Is(token.IDENT) {
		p.next()
		params = append(params, &ast.Identifier{
			Token: p.cur,
			Value: p.cur.Text,
		})
		if p.peek.Is(token.COMMA) {
			p.next()
		}
	}
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return params
}

func (p *Parser) functionStmt() ast.Statement {
	stmt := &ast.FunctionStatement{Token: p.cur}
	p.next()
	stmt.Name = &ast.Identifier{Token: p.cur, Value: p.cur.Text}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	stmt.Parameters = p.fnParameterIdents()
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	stmt.Body = p.blockStmt()
	if p.peek.Is(token.SEMICOLON) {
		p.next()
	}
	return stmt
}

func (p *Parser) fnExpr() ast.Expression {
	expr := &ast.FunctionLiteral{Token: p.cur}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	expr.Parameters = p.fnParameterIdents()
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	expr.Body = p.blockStmt()
	return expr
}

func (p *Parser) groupesExpr() ast.Expression {
	p.next()
	expr := p.expression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return expr
}

func (p *Parser) identExpr() ast.Expression {
	return &ast.Identifier{Token: p.cur, Value: p.cur.Text}
}

func (p *Parser) integerExpr() ast.Expression {
	expr := &ast.IntegerLiteral{Token: p.cur}
	v, err := strconv.ParseInt(p.cur.Text, 10, 64)
	if err != nil {
		p.errorf("invalid integer %s: %v", p.cur, err)
		return nil
	}
	expr.Value = v
	return expr
}

func (p *Parser) stringLit() ast.Expression {
	return &ast.StringLiteral{
		Token: p.cur,
		Value: p.cur.Text,
	}
}

func (p *Parser) prefixExpr() ast.Expression {
	expr := &ast.PrefixExpression{
		Token:    p.cur,
		Operator: p.cur.Text,
	}
	p.next()
	expr.Right = p.expression(PREFIX)
	return expr
}

func (p *Parser) indexExpr(left ast.Expression) ast.Expression {
	expr := &ast.IndexExpression{
		Token: p.cur,
		Value: left,
	}
	p.next()
	expr.Index = p.expression(PREFIX)
	if !p.expectPeek(token.RBRACKET) {
		return nil
	}
	return expr
}

func (p *Parser) infixExpr(left ast.Expression) ast.Expression {
	expr := &ast.InfixExpression{
		Token:    p.cur,
		Left:     left,
		Operator: p.cur.Text,
	}
	precedence := p.precedence(p.cur)
	p.next()
	expr.Right = p.expression(precedence)
	return expr
}

func (p *Parser) delimitedExpr(term token.TokenType) []ast.Expression {
	var args []ast.Expression

	// empty list
	if p.peek.Is(term) {
		p.next()
		return args
	}

	// first element of the list
	p.next()
	args = append(args, p.expression(LOWEST))

	for p.peek.Is(token.COMMA) {
		p.next()
		p.next()
		args = append(args, p.expression(LOWEST))
	}

	if !p.expectPeek(term) {
		return nil
	}

	return args
}

func (p *Parser) arrayExpr() ast.Expression {
	expr := &ast.ArrayLiteral{Token: p.cur}
	expr.Elements = p.delimitedExpr(token.RBRACKET)
	return expr
}

func (p *Parser) callExpr(left ast.Expression) ast.Expression {
	expr := &ast.CallExpression{
		Token:    p.cur,
		Function: left,
	}
	expr.Arguments = p.delimitedExpr(token.RPAREN)
	return expr
}

func (p *Parser) stmt() ast.Statement {
	switch p.cur.Type {
	case token.FUNCTION:
		return p.functionStmt()
	case token.LET:
		return p.letStmt()
	case token.RETURN:
		return p.returnStmt()
	case token.IMPORT:
		return p.importStmt()
	case token.WHILE:
		return p.whileStmt()
	default:
		return p.expressionStmt()
	}
}

func (p *Parser) expressionStmt() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{}
	stmt.Token = p.cur
	stmt.Expression = p.expression(LOWEST)
	if p.peek.Is(token.SEMICOLON) {
		p.next()
	}
	return stmt
}

func (p *Parser) returnStmt() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{}
	stmt.Token = p.cur
	if p.peek.Is(token.SEMICOLON) {
		p.next()
		return stmt
	}
	p.next()
	stmt.ReturnValue = p.expression(LOWEST)
	if p.peek.Is(token.SEMICOLON) {
		p.next()
	}
	return stmt
}

func (p *Parser) assignExpr(left ast.Expression) ast.Expression {
	expr := &ast.AssignmentExpression{
		Token: p.cur,
		Left:  left,
	}
	p.next()
	expr.Value = p.expression(LOWEST)
	return expr
}

func (p *Parser) importStmt() *ast.ImportStatement {
	stmt := &ast.ImportStatement{Token: p.cur}
	if !p.expectPeek(token.STRING) {
		return nil
	}
	stmt.Value = p.cur.Text
	if p.peek.Is(token.SEMICOLON) {
		p.next()
	}
	return stmt
}

func (p *Parser) letStmt() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.cur}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{
		Token: p.cur,
		Value: p.cur.Text,
	}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	p.next()
	stmt.Value = p.expression(LOWEST)
	if p.peek.Is(token.SEMICOLON) {
		p.next()
	}
	return stmt
}

func (p *Parser) expression(precedence int) ast.Expression {
	prefix, ok := p.prefixFns[p.cur.Type]
	if !ok {
		p.errorf("no prefix parse function for: %s", p.cur)
		return nil
	}
	left := prefix()

	for !p.peek.Is(token.SEMICOLON) && precedence < p.precedence(p.peek) {
		infix, ok := p.infixFns[p.peek.Type]
		if !ok {
			return left
		}
		p.next()
		left = infix(left)
	}

	return left
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peek.Is(t) {
		p.next()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) errorf(format string, args ...interface{}) {
	p.errors = append(p.errors, fmt.Sprintf(format, args...))
}
