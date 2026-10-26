package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNextToken(t *testing.T) {
	t.Run("arithmetic expression", func(t *testing.T) {
		input := "12 + 34"

		expected := []Token{
			{Type: TokenNumber, Literal: "12"},
			{Type: TokenPlus, Literal: "+"},
			{Type: TokenNumber, Literal: "34"},
			{Type: TokenEOF, Literal: ""},
		}

		lex := NewLexer(input)

		for _, want := range expected {
			got := lex.NextToken()

			require.Equal(t, want.Type, got.Type)
			require.Equal(t, want.Literal, got.Literal)
		}
	})

	t.Run("parentheses and operator precedence tokens", func(t *testing.T) {
		input := "8 * (3 - 1)"

		expected := []Token{
			{Type: TokenNumber, Literal: "8"},
			{Type: TokenMultiply, Literal: "*"},
			{Type: TokenLeftParen, Literal: "("},
			{Type: TokenNumber, Literal: "3"},
			{Type: TokenMinus, Literal: "-"},
			{Type: TokenNumber, Literal: "1"},
			{Type: TokenRightParen, Literal: ")"},
			{Type: TokenEOF, Literal: ""},
		}

		lex := NewLexer(input)

		for _, want := range expected {
			got := lex.NextToken()

			require.Equal(t, want.Type, got.Type)
			require.Equal(t, want.Literal, got.Literal)
		}
	})

	t.Run("whitespace", func(t *testing.T) {
		input := "  \t42 \n\r"

		expected := []Token{
			{Type: TokenNumber, Literal: "42"},
			{Type: TokenEOF, Literal: ""},
		}

		lex := NewLexer(input)

		for _, want := range expected {
			got := lex.NextToken()

			require.Equal(t, want.Type, got.Type)
			require.Equal(t, want.Literal, got.Literal)
		}
	})

	t.Run("all operators", func(t *testing.T) {
		input := "+ - * / ( )"

		expected := []Token{
			{Type: TokenPlus, Literal: "+"},
			{Type: TokenMinus, Literal: "-"},
			{Type: TokenMultiply, Literal: "*"},
			{Type: TokenDivide, Literal: "/"},
			{Type: TokenLeftParen, Literal: "("},
			{Type: TokenRightParen, Literal: ")"},
			{Type: TokenEOF, Literal: ""},
		}

		lex := NewLexer(input)

		for _, want := range expected {
			got := lex.NextToken()

			require.Equal(t, want.Type, got.Type)
			require.Equal(t, want.Literal, got.Literal)
		}
	})

	t.Run("multiple digit number", func(t *testing.T) {
		input := "123456"

		lex := NewLexer(input)

		got := lex.NextToken()

		require.Equal(t, TokenNumber, got.Type)
		require.Equal(t, "123456", got.Literal)

		require.Equal(t, TokenEOF, lex.NextToken().Type)
	})

	t.Run("illegal character", func(t *testing.T) {
		input := "@"

		lex := NewLexer(input)

		got := lex.NextToken()

		require.Equal(t, TokenIllegal, got.Type)
		require.Equal(t, "@", got.Literal)
	})

	t.Run("identifier is illegal for now", func(t *testing.T) {
		input := "abc"

		lex := NewLexer(input)

		got := lex.NextToken()

		require.Equal(t, TokenIllegal, got.Type)
		require.Equal(t, "abc", got.Literal)
	})

	t.Run("empty input", func(t *testing.T) {
		lex := NewLexer("")

		got := lex.NextToken()

		require.Equal(t, TokenEOF, got.Type)
		require.Equal(t, "", got.Literal)
	})

	t.Run("repeated EOF", func(t *testing.T) {
		lex := NewLexer("")

		require.Equal(t, TokenEOF, lex.NextToken().Type)
		require.Equal(t, TokenEOF, lex.NextToken().Type)
	})
}
