package main

type TokenType int

const (
	TokenNumber TokenType = iota
	TokenPlus
	TokenMinus
	TokenMultiply
	TokenDivide
	TokenLeftParen
	TokenRightParen
	TokenIllegal
	TokenEOF
)

type Token struct {
	Type    TokenType
	Literal string
}
