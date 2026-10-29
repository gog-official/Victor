package lexer

type TokenType int

const (
	TokenNumber TokenType = iota
	TokenPlus
	TokenMinus
	TokenMultiply
	TokenDivide
	TokenIdentifier
	TokenComma
	TokenLeftParen
	TokenRightParen
	TokenIllegal
	TokenEOF
)

type Token struct {
	Type    TokenType
	Literal string
}
