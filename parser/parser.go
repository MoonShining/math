package parser

import (
	"github.com/MoonShining/math/lexer"
	"strconv"
)

var precedences = map[lexer.TokenType]int{
	lexer.Add:       1,
	lexer.Sub:       1,
	lexer.Mul:       2,
	lexer.Div:       2,
	lexer.LeftParen: 4,
}

type Expression interface {
}

type PrefixExpression struct {
	Prefix lexer.Token
	Expr   Expression
}

type InfixExpression struct {
	Operator lexer.Token
	Left     Expression
	Right    Expression
}

type NumberExpression struct {
	Value int64
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l,
		infixParseFns:  make(map[lexer.TokenType]func(Expression) Expression),
		prefixParseFns: make(map[lexer.TokenType]func() Expression),
	}

	p.registerPrefix(lexer.Sub, p.parsePrefixExpression)
	p.registerPrefix(lexer.LeftParen, p.parseGroupedExpression)
	p.registerPrefix(lexer.Number, p.parseNumber)

	p.registerInfix(lexer.Sub, p.parseInfixExpression)
	p.registerInfix(lexer.Add, p.parseInfixExpression)
	p.registerInfix(lexer.Mul, p.parseInfixExpression)
	p.registerInfix(lexer.Div, p.parseInfixExpression)

	p.nextToken()
	p.nextToken()

	return p
}

type Parser struct {
	l         *lexer.Lexer
	curToken  lexer.Token
	peekToken lexer.Token

	prefixParseFns map[lexer.TokenType]func() Expression
	infixParseFns  map[lexer.TokenType]func(Expression) Expression
}

func (p *Parser) registerPrefix(t lexer.TokenType, f func() Expression) {
	p.prefixParseFns[t] = f
}

func (p *Parser) registerInfix(t lexer.TokenType, f func(Expression) Expression) {
	p.infixParseFns[t] = f
}

func (p *Parser) ParseExpression(precedence int) Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	leftExp := prefix()

	for precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return 0
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return 0
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) parseNumber() Expression {
	value, _ := strconv.ParseInt(p.curToken.Literal, 10, 64)
	return NumberExpression{Value: value}
}

func (p *Parser) parsePrefixExpression() Expression {
	pe := PrefixExpression{Prefix: lexer.SUB}
	p.nextToken()
	pe.Expr = p.ParseExpression(3)
	return pe
}

func (p *Parser) parseGroupedExpression() Expression {
	p.nextToken()
	exp := p.ParseExpression(0)
	p.nextToken()
	return exp
}

func (p *Parser) parseInfixExpression(left Expression) Expression {
	exp := InfixExpression{Left: left, Operator: p.curToken}
	precedence := p.curPrecedence()
	p.nextToken()
	exp.Right = p.ParseExpression(precedence)
	return exp
}
