package main

import (
	"fmt"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	SUM
	PRODUCT
	PREFIX
)

var precedences = map[TokenType]int{
	TokenPlus:     SUM,
	TokenMinus:    SUM,
	TokenMultiply: PRODUCT,
	TokenDivide:   PRODUCT,
}

type Value struct {
	Scalar   Number
	Vector   Vector2D
	IsVector bool
}

func scalarValue(n Number) Value {
	return Value{Scalar: n}
}

func vectorValue(v Vector2D) Value {
	return Value{Vector: v, IsVector: true}
}

type (
	prefixParseFunc func() (Value, error)
	infixParseFunc  func(Value) (Value, error)
)

type Parser struct {
	lexer        *Lexer
	currentToken Token
	peekToken    Token
	errors       []string

	prefixParseFuncs map[TokenType]prefixParseFunc
	infixParseFuncs  map[TokenType]infixParseFunc
}

func NewParser(lexer *Lexer) *Parser {
	p := &Parser{
		lexer:            lexer,
		errors:           []string{},
		prefixParseFuncs: make(map[TokenType]prefixParseFunc),
		infixParseFuncs:  make(map[TokenType]infixParseFunc),
	}

	p.registerPrefix(TokenNumber, p.parseIntegerLiteral)
	p.registerPrefix(TokenIdentifier, p.parseIdentifier)
	p.registerPrefix(TokenMinus, p.parsePrefixExpression)
	p.registerPrefix(TokenLeftParen, p.parseGroupedOrVector)

	p.registerInfix(TokenPlus, p.parseInfixExpression)
	p.registerInfix(TokenMinus, p.parseInfixExpression)
	p.registerInfix(TokenMultiply, p.parseInfixExpression)
	p.registerInfix(TokenDivide, p.parseInfixExpression)

	p.readNextToken()
	p.readNextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) registerPrefix(token TokenType, prefixFunc prefixParseFunc) {
	p.prefixParseFuncs[token] = prefixFunc
}

func (p *Parser) registerInfix(token TokenType, infixFunc infixParseFunc) {
	p.infixParseFuncs[token] = infixFunc
}

func (p *Parser) readNextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) nextPrecedence() int {
	if precedence, ok := precedences[p.peekToken.Type]; ok {
		return precedence
	}

	return LOWEST
}

func (p *Parser) currentPrecedence() int {
	if precedence, ok := precedences[p.currentToken.Type]; ok {
		return precedence
	}

	return LOWEST
}

func (p *Parser) Parse() (Number, error) {
	result, err := p.ParseValue()
	if err != nil {
		return Number{}, err
	}

	if result.IsVector {
		return Number{}, fmt.Errorf("expected a number but got a vector")
	}

	return result.Scalar, nil
}

func (p *Parser) ParseValue() (Value, error) {
	result, err := p.parseExpression(LOWEST)
	if err != nil {
		return Value{}, err
	}

	if p.peekToken.Type != TokenEOF {
		return Value{}, fmt.Errorf(
			"unexpected token %q after expression",
			p.peekToken.Literal,
		)
	}

	return result, nil
}

func (p *Parser) parseExpression(precedence int) (Value, error) {
	prefix := p.prefixParseFuncs[p.currentToken.Type]
	if prefix == nil {
		return Value{}, fmt.Errorf(
			"no prefix parser for token %q",
			p.currentToken.Literal,
		)
	}

	left, err := prefix()
	if err != nil {
		return Value{}, err
	}

	for p.peekToken.Type != TokenEOF &&
		precedence < p.nextPrecedence() {
		infix := p.infixParseFuncs[p.peekToken.Type]
		if infix == nil {
			return left, nil
		}

		p.readNextToken()

		left, err = infix(left)
		if err != nil {
			return Value{}, err
		}
	}

	return left, nil
}

func (p *Parser) parseIntegerLiteral() (Value, error) {
	value, err := strconv.ParseInt(p.currentToken.Literal, 10, 64)
	if err != nil {
		return Value{}, fmt.Errorf(
			"invalid integer %q: %w",
			p.currentToken.Literal,
			err,
		)
	}

	return scalarValue(NewInteger(value)), nil
}

func (p *Parser) parseIdentifier() (Value, error) {
	vector, exists := workspace.Get(p.currentToken.Literal)
	if !exists {
		return Value{}, fmt.Errorf(
			"vector %q is not defined",
			p.currentToken.Literal,
		)
	}

	return vectorValue(vector), nil
}

func (p *Parser) parsePrefixExpression() (Value, error) {
	p.readNextToken()

	right, err := p.parseExpression(PREFIX)
	if err != nil {
		return Value{}, err
	}

	if right.IsVector {
		return vectorValue(Scale(right.Vector, NewInteger(-1))), nil
	}

	return scalarValue(SubN(NewInteger(0), right.Scalar)), nil
}

func (p *Parser) parseGroupedOrVector() (Value, error) {
	p.readNextToken()

	first, err := p.parseExpression(LOWEST)
	if err != nil {
		return Value{}, err
	}

	if p.peekToken.Type == TokenComma {
		p.readNextToken()
		p.readNextToken()

		second, err := p.parseExpression(LOWEST)
		if err != nil {
			return Value{}, err
		}

		if first.IsVector || second.IsVector {
			return Value{}, fmt.Errorf("vector components must be numbers")
		}

		if p.peekToken.Type != TokenRightParen {
			return Value{}, fmt.Errorf("expected closing parenthesis")
		}

		p.readNextToken()

		return vectorValue(Vector2D{X: first.Scalar, Y: second.Scalar}), nil
	}

	if p.peekToken.Type != TokenRightParen {
		return Value{}, fmt.Errorf("expected closing parenthesis")
	}

	p.readNextToken()

	return first, nil
}

func (p *Parser) parseInfixExpression(left Value) (Value, error) {
	operator := p.currentToken.Type
	precedence := p.currentPrecedence()

	p.readNextToken()

	right, err := p.parseExpression(precedence)
	if err != nil {
		return Value{}, err
	}

	switch operator {
	case TokenPlus:
		if left.IsVector && right.IsVector {
			return vectorValue(Add(left.Vector, right.Vector)), nil
		}
		if !left.IsVector && !right.IsVector {
			return scalarValue(AddN(left.Scalar, right.Scalar)), nil
		}
		return Value{}, fmt.Errorf("cannot add a vector and a number")
	case TokenMinus:
		if left.IsVector && right.IsVector {
			return vectorValue(Subtract(left.Vector, right.Vector)), nil
		}
		if !left.IsVector && !right.IsVector {
			return scalarValue(SubN(left.Scalar, right.Scalar)), nil
		}
		return Value{}, fmt.Errorf("cannot subtract a vector and a number")
	case TokenMultiply:
		if left.IsVector && right.IsVector {
			return scalarValue(DotProduct(left.Vector, right.Vector)), nil
		}
		if left.IsVector && !right.IsVector {
			return vectorValue(Scale(left.Vector, right.Scalar)), nil
		}
		if !left.IsVector && right.IsVector {
			return vectorValue(Scale(right.Vector, left.Scalar)), nil
		}
		return scalarValue(MulN(left.Scalar, right.Scalar)), nil
	case TokenDivide:
		if left.IsVector && !right.IsVector {
			result, err := Divide(left.Vector, right.Scalar)
			if err != nil {
				return Value{}, err
			}
			return vectorValue(result), nil
		}
		if !left.IsVector && !right.IsVector {
			result, err := DivN(left.Scalar, right.Scalar)
			if err != nil {
				return Value{}, err
			}
			return scalarValue(result), nil
		}
		return Value{}, fmt.Errorf("cannot divide a vector by a vector")
	default:
		return Value{}, fmt.Errorf(
			"unknown operator %q",
			p.currentToken.Literal,
		)
	}
}
