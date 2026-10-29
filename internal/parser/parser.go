package parser

import (
	"fmt"
	"strconv"

	"github.com/gog-official/victor/internal/lexer"
	"github.com/gog-official/victor/internal/number"
	"github.com/gog-official/victor/internal/vector"
	"github.com/gog-official/victor/internal/workspace"
)

const (
	_ int = iota
	LOWEST
	SUM
	PRODUCT
	PREFIX
)

var precedences = map[lexer.TokenType]int{
	lexer.TokenPlus:     SUM,
	lexer.TokenMinus:    SUM,
	lexer.TokenMultiply: PRODUCT,
	lexer.TokenDivide:   PRODUCT,
}

type Value struct {
	Scalar   number.Number
	Vector   vector.Vector2D
	IsVector bool
}

func scalarValue(n number.Number) Value {
	return Value{Scalar: n}
}

func vectorValue(v vector.Vector2D) Value {
	return Value{Vector: v, IsVector: true}
}

type (
	prefixParseFunc func() (Value, error)
	infixParseFunc  func(Value) (Value, error)
)

type Parser struct {
	lexer        *lexer.Lexer
	store        *workspace.Workspace
	currentToken lexer.Token
	peekToken    lexer.Token
	errors       []string

	prefixParseFuncs map[lexer.TokenType]prefixParseFunc
	infixParseFuncs  map[lexer.TokenType]infixParseFunc
}

func NewParser(lex *lexer.Lexer, store *workspace.Workspace) *Parser {
	p := &Parser{
		lexer:            lex,
		store:            store,
		errors:           []string{},
		prefixParseFuncs: make(map[lexer.TokenType]prefixParseFunc),
		infixParseFuncs:  make(map[lexer.TokenType]infixParseFunc),
	}

	p.registerPrefix(lexer.TokenNumber, p.parseIntegerLiteral)
	p.registerPrefix(lexer.TokenIdentifier, p.parseIdentifier)
	p.registerPrefix(lexer.TokenMinus, p.parsePrefixExpression)
	p.registerPrefix(lexer.TokenLeftParen, p.parseGroupedOrVector)

	p.registerInfix(lexer.TokenPlus, p.parseInfixExpression)
	p.registerInfix(lexer.TokenMinus, p.parseInfixExpression)
	p.registerInfix(lexer.TokenMultiply, p.parseInfixExpression)
	p.registerInfix(lexer.TokenDivide, p.parseInfixExpression)

	p.readNextToken()
	p.readNextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) registerPrefix(token lexer.TokenType, prefixFunc prefixParseFunc) {
	p.prefixParseFuncs[token] = prefixFunc
}

func (p *Parser) registerInfix(token lexer.TokenType, infixFunc infixParseFunc) {
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

func (p *Parser) Parse() (number.Number, error) {
	result, err := p.ParseValue()
	if err != nil {
		return number.Number{}, err
	}

	if result.IsVector {
		return number.Number{}, fmt.Errorf("expected a number but got a vector")
	}

	return result.Scalar, nil
}

func (p *Parser) ParseValue() (Value, error) {
	result, err := p.parseExpression(LOWEST)
	if err != nil {
		return Value{}, err
	}

	if p.peekToken.Type != lexer.TokenEOF {
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

	for p.peekToken.Type != lexer.TokenEOF &&
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

	return scalarValue(number.NewInteger(value)), nil
}

func (p *Parser) parseIdentifier() (Value, error) {
	vec, exists := p.store.Get(p.currentToken.Literal)
	if !exists {
		return Value{}, fmt.Errorf(
			"vector %q is not defined",
			p.currentToken.Literal,
		)
	}

	return vectorValue(vec), nil
}

func (p *Parser) parsePrefixExpression() (Value, error) {
	p.readNextToken()

	right, err := p.parseExpression(PREFIX)
	if err != nil {
		return Value{}, err
	}

	if right.IsVector {
		return vectorValue(vector.Scale(right.Vector, number.NewInteger(-1))), nil
	}

	return scalarValue(number.SubN(number.NewInteger(0), right.Scalar)), nil
}

func (p *Parser) parseGroupedOrVector() (Value, error) {
	p.readNextToken()

	first, err := p.parseExpression(LOWEST)
	if err != nil {
		return Value{}, err
	}

	if p.peekToken.Type == lexer.TokenComma {
		p.readNextToken()
		p.readNextToken()

		second, err := p.parseExpression(LOWEST)
		if err != nil {
			return Value{}, err
		}

		if first.IsVector || second.IsVector {
			return Value{}, fmt.Errorf("vector components must be numbers")
		}

		if p.peekToken.Type != lexer.TokenRightParen {
			return Value{}, fmt.Errorf("expected closing parenthesis")
		}

		p.readNextToken()

		return vectorValue(vector.Vector2D{X: first.Scalar, Y: second.Scalar}), nil
	}

	if p.peekToken.Type != lexer.TokenRightParen {
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
	case lexer.TokenPlus:
		if left.IsVector && right.IsVector {
			return vectorValue(vector.Add(left.Vector, right.Vector)), nil
		}
		if !left.IsVector && !right.IsVector {
			return scalarValue(number.AddN(left.Scalar, right.Scalar)), nil
		}
		return Value{}, fmt.Errorf("cannot add a vector and a number")
	case lexer.TokenMinus:
		if left.IsVector && right.IsVector {
			return vectorValue(vector.Subtract(left.Vector, right.Vector)), nil
		}
		if !left.IsVector && !right.IsVector {
			return scalarValue(number.SubN(left.Scalar, right.Scalar)), nil
		}
		return Value{}, fmt.Errorf("cannot subtract a vector and a number")
	case lexer.TokenMultiply:
		if left.IsVector && right.IsVector {
			return scalarValue(vector.DotProduct(left.Vector, right.Vector)), nil
		}
		if left.IsVector && !right.IsVector {
			return vectorValue(vector.Scale(left.Vector, right.Scalar)), nil
		}
		if !left.IsVector && right.IsVector {
			return vectorValue(vector.Scale(right.Vector, left.Scalar)), nil
		}
		return scalarValue(number.MulN(left.Scalar, right.Scalar)), nil
	case lexer.TokenDivide:
		if left.IsVector && !right.IsVector {
			result, err := vector.Divide(left.Vector, right.Scalar)
			if err != nil {
				return Value{}, err
			}
			return vectorValue(result), nil
		}
		if !left.IsVector && !right.IsVector {
			result, err := number.DivN(left.Scalar, right.Scalar)
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
