package parser

import (
	"fmt"
	"strconv"
	"strings"

	"go-interpreter.com/m/ast"
	"go-interpreter.com/m/lexer"
	"go-interpreter.com/m/token"
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
)

var precendences = map[token.TokenType]int{
	token.Product:    PRODUCT,
	token.Slash:      PRODUCT,
	token.Equal:      EQUALS,
	token.LessThan:   LESSGREATER,
	token.BiggerThan: LESSGREATER,
	token.Different:  EQUALS,
	token.Plus:       SUM,
	token.Minus:      SUM,
	token.LParen:     CALL,
}

type prefixParseFn func() ast.Expression
type infixParseFn func(ast.Expression) ast.Expression

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	Errors []error

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, Errors: []error{}}

	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.Ident, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiterals)
	p.registerPrefix(token.Negation, p.parsePrefixExpression)
	p.registerPrefix(token.Minus, p.parsePrefixExpression)
	p.registerPrefix(token.True, p.parseBooleanLiterals)
	p.registerPrefix(token.False, p.parseBooleanLiterals)
	p.registerPrefix(token.LParen, p.parseGroupedExpression)
	p.registerPrefix(token.IfConditional, p.parseIfExpression)
	p.registerPrefix(token.Function, p.parseFunction)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	for tokenType := range precendences {
		if tokenType == token.LParen {
			continue
		}

		p.registerInfix(tokenType, p.parseInfixExpression)
	}
	p.registerInfix(token.LParen, p.parseCallExpressionArguments)

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curIsToken(token.EOF) {
		smt := p.parseStatement()
		if smt != nil {
			program.Statements = append(program.Statements, smt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.Let:
		return p.parseLetStatement()
	case token.ReturnStatement:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stm := &ast.ExpressionStatement{Token: p.curToken}

	stm.Expression = p.parseExpression(LOWEST)

	if p.peekIsToken(token.Semicolon) {
		p.nextToken()
	}

	return stm
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.prefixErrors(p.curToken)
		return nil
	}

	leftExpr := prefix()
	for !p.peekIsToken(token.Semicolon) && precedence < p.peekPrecendence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExpr
		}

		p.nextToken()
		leftExpr = infix(leftExpr)
	}

	return leftExpr
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stm := &ast.LetStatement{Token: p.curToken}
	if !p.expectPeek(token.Ident) {
		return nil
	}

	stm.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.Assign) {
		return nil
	}
	p.nextToken()
	stm.Value = p.parseExpression(LOWEST)

	if p.peekIsToken(token.Semicolon) {
		p.nextToken()
	}

	return stm
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stm := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()
	stm.Value = p.parseExpression(LOWEST)

	if p.peekIsToken(token.Semicolon) {
		p.nextToken()
	}

	return stm
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiterals() ast.Expression {
	num, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Errorf("could not parse %q as integer", p.curToken.Literal)
		p.Errors = append(p.Errors, msg)
		return nil
	}
	return &ast.IntegerLiteral{Token: p.curToken, Value: int(num)}
}

func (p *Parser) parseBooleanLiterals() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curIsToken(token.True)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RParent) {
		return nil
	}

	return exp
}

func (p *Parser) parseFunction() ast.Expression {
	fnExpression := &ast.FunctionExpression{Token: p.curToken}
	if !p.expectPeek(token.LParen) {
		return nil
	}

	fnExpression.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBrace) {
		return nil
	}

	fnExpression.Body = p.parseBlockStatement()
	return fnExpression
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	tokens := make([]*ast.Identifier, 0)

	if p.peekIsToken(token.RParent) {
		p.nextToken()
		return tokens
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	tokens = append(tokens, ident)

	for p.peekIsToken(token.Comma) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		tokens = append(tokens, ident)
	}

	if !p.expectPeek(token.RParent) {
		return nil
	}

	return tokens
}

func (p *Parser) parseIfExpression() ast.Expression {
	ifExp := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.LParen) {
		return nil
	}

	p.nextToken()
	ifExp.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RParent) {
		return nil
	}

	if !p.expectPeek(token.LBrace) {
		return nil
	}

	ifExp.Consequence = p.parseBlockStatement()

	if p.peekIsToken(token.ElseConditional) {
		p.nextToken()

		if !p.expectPeek(token.LBrace) {
			return nil
		}

		ifExp.Alternative = p.parseBlockStatement()
	}

	return ifExp
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curIsToken(token.RBrace) && !p.curIsToken(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseCallExpressionArguments(left ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: left}
	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekIsToken(token.RParent) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))
	for p.peekIsToken(token.Comma) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RParent) {
		return nil
	}

	return args
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	prefixExpression := &ast.PrefixExpression{Token: p.curToken, Operator: p.curToken.Literal}
	p.nextToken()
	prefixExpression.Right = p.parseExpression(PREFIX)

	return prefixExpression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	infixExpression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecendence()
	p.nextToken()
	infixExpression.Right = p.parseExpression(precedence)

	return infixExpression
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) peekIsToken(expectedToken token.TokenType) bool {
	return p.peekToken.Type == expectedToken
}

func (p *Parser) curIsToken(expectedToken token.TokenType) bool {
	return p.curToken.Type == expectedToken
}

func (p *Parser) expectPeek(expectedToken token.TokenType) bool {
	if p.peekIsToken(expectedToken) {
		p.nextToken()
		return true
	} else {
		p.peekErrors(expectedToken)
		return false
	}
}

func (p *Parser) curPrecendence() int {
	if p, ok := precendences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) peekPrecendence() int {
	if p, ok := precendences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) Error() string {
	b := strings.Builder{}
	for _, err := range p.Errors {
		b.WriteString(fmt.Sprintf("%s\n", err.Error()))
	}
	return b.String()
}

func (p *Parser) peekErrors(expectedToken token.TokenType) {
	msg := fmt.Errorf("expected next token to be %s, got %s instead",
		expectedToken, p.peekToken.Type)
	p.Errors = append(p.Errors, msg)
}

func (p *Parser) prefixErrors(t token.Token) {
	msg := fmt.Errorf("no prefix parse function for %s found", t.Type)
	p.Errors = append(p.Errors, msg)
}
