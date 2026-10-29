package lexer

type Lexer struct {
	input    string
	position int
	readPos  int
	ch       byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (lex *Lexer) readChar() {
	if lex.readPos >= len(lex.input) {
		lex.ch = 0
	} else {
		lex.ch = lex.input[lex.readPos]
	}

	lex.position = lex.readPos
	lex.readPos++
}

func (lex *Lexer) NextToken() Token {
	var tok Token

	for lex.ch == ' ' || lex.ch == '\t' || lex.ch == '\n' || lex.ch == '\r' {
		lex.readChar()
	}

	switch lex.ch {
	case '+':
		tok = Token{Type: TokenPlus, Literal: string(lex.ch)}
	case ',':
		tok = Token{Type: TokenComma, Literal: string(lex.ch)}
	case '(':
		tok = Token{Type: TokenLeftParen, Literal: string(lex.ch)}
	case ')':
		tok = Token{Type: TokenRightParen, Literal: string(lex.ch)}
	case '-':
		tok = Token{Type: TokenMinus, Literal: string(lex.ch)}
	case '/':
		tok = Token{Type: TokenDivide, Literal: string(lex.ch)}
	case '*':
		tok = Token{Type: TokenMultiply, Literal: string(lex.ch)}
	case 0:
		tok = Token{Type: TokenEOF, Literal: ""}
	default:
		if isLetter(lex.ch) {
			tok.Literal = lex.readIdentifier()
			tok.Type = TokenIdentifier
			return tok
		} else if isNumber(lex.ch) {
			tok.Literal = lex.readNumber()
			tok.Type = TokenNumber
			return tok
		} else {
			tok = Token{Type: TokenIllegal, Literal: string(lex.ch)}
		}
	}

	lex.readChar()
	return tok
}

func isLetter(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_'
}

func (lex *Lexer) readIdentifier() string {
	position := lex.position

	for isLetter(lex.ch) {
		lex.readChar()
	}

	return lex.input[position:lex.position]
}

func isNumber(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (lex *Lexer) readNumber() string {
	position := lex.position

	for isNumber(lex.ch) {
		lex.readChar()
	}

	return lex.input[position:lex.position]
}
