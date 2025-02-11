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

type prefixParseFn func() ast.Expression
type infixParseFn func(ast.Expression) ast.Expression

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	errors []error

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []error{}}

	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.Ident, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiterals)

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
		return nil
	}

	leftExpr := prefix()

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

	// TODO: Parse let expression
	for !p.curIsToken(token.Semicolon) {
		p.nextToken()
	}

	return stm
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stm := &ast.ReturnStatement{Token: p.curToken}

	// TODO: Parse let expression
	p.nextToken()

	for !p.curIsToken(token.Semicolon) {
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
		p.errors = append(p.errors, msg)
		return nil
	}
	return &ast.IntegerLiteral{Token: p.curToken, Value: int(num)}
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

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) Error() string {
	b := strings.Builder{}
	for _, err := range p.errors {
		b.WriteString(fmt.Sprintf("%s\n", err.Error()))
	}
	return b.String()
}

func (p *Parser) peekErrors(expectedToken token.TokenType) {
	msg := fmt.Errorf("expected next token to be %s, got %s instead",
		expectedToken, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
