package parser

import (
	"testing"

	"github.com/gog-official/victor/internal/lexer"
	"github.com/gog-official/victor/internal/number"
	"github.com/gog-official/victor/internal/vector"
	"github.com/gog-official/victor/internal/workspace"
	"github.com/stretchr/testify/require"
)

var testWorkspace = workspace.NewWorkspace()

func newTestParser(input string) *Parser {
	return NewParser(lexer.NewLexer(input), testWorkspace)
}

func parseExpression(t *testing.T, input string) number.Number {
	t.Helper()

	result, err := newTestParser(input).Parse()
	require.NoError(t, err)

	return result
}

func TestParseInteger(t *testing.T) {
	t.Run("positive integer", func(t *testing.T) {
		result := parseExpression(t, "42")

		require.Equal(t, "42", result.String())
	})

	t.Run("zero", func(t *testing.T) {
		result := parseExpression(t, "0")

		require.Equal(t, "0", result.String())
	})

	t.Run("negative integer", func(t *testing.T) {
		result := parseExpression(t, "-15")

		require.Equal(t, "-15", result.String())
	})
}

func TestParseAddition(t *testing.T) {
	t.Run("two integers", func(t *testing.T) {
		result := parseExpression(t, "2 + 3")

		require.Equal(t, "5", result.String())
	})

	t.Run("multiple additions", func(t *testing.T) {
		result := parseExpression(t, "2 + 3 + 4")

		require.Equal(t, "9", result.String())
	})

	t.Run("addition with negative number", func(t *testing.T) {
		result := parseExpression(t, "5 + -2")

		require.Equal(t, "3", result.String())
	})
}

func TestParseSubtraction(t *testing.T) {
	t.Run("two integers", func(t *testing.T) {
		result := parseExpression(t, "9 - 4")

		require.Equal(t, "5", result.String())
	})

	t.Run("negative result", func(t *testing.T) {
		result := parseExpression(t, "3 - 8")

		require.Equal(t, "-5", result.String())
	})

	t.Run("multiple subtractions", func(t *testing.T) {
		result := parseExpression(t, "10 - 3 - 2")

		require.Equal(t, "5", result.String())
	})
}

func TestParseMultiplication(t *testing.T) {
	t.Run("two integers", func(t *testing.T) {
		result := parseExpression(t, "4 * 5")

		require.Equal(t, "20", result.String())
	})

	t.Run("multiple multiplications", func(t *testing.T) {
		result := parseExpression(t, "2 * 3 * 4")

		require.Equal(t, "24", result.String())
	})
}

func TestParseDivision(t *testing.T) {
	t.Run("exact division", func(t *testing.T) {
		result := parseExpression(t, "10 / 2")

		require.Equal(t, "5", result.String())
	})

	t.Run("fractional result", func(t *testing.T) {
		result := parseExpression(t, "1 / 3")

		require.Equal(t, "1/3", result.String())
	})

	t.Run("fraction arithmetic", func(t *testing.T) {
		result := parseExpression(t, "1 / 3 + 1 / 6")

		require.Equal(t, "1/2", result.String())
	})
}

func TestOperatorPrecedence(t *testing.T) {
	t.Run("multiplication before addition", func(t *testing.T) {
		result := parseExpression(t, "2 + 3 * 4")

		require.Equal(t, "14", result.String())
	})

	t.Run("division before subtraction", func(t *testing.T) {
		result := parseExpression(t, "10 - 6 / 2")

		require.Equal(t, "7", result.String())
	})

	t.Run("addition before multiplication with parentheses", func(t *testing.T) {
		result := parseExpression(t, "(2 + 3) * 4")

		require.Equal(t, "20", result.String())
	})

	t.Run("complex expression", func(t *testing.T) {
		result := parseExpression(t, "2 + 3 * (4 - 1)")

		require.Equal(t, "11", result.String())
	})
}

func TestParseParentheses(t *testing.T) {
	t.Run("grouped addition", func(t *testing.T) {
		result := parseExpression(t, "(2 + 3)")

		require.Equal(t, "5", result.String())
	})

	t.Run("nested parentheses", func(t *testing.T) {
		result := parseExpression(t, "((2 + 3) * (4 - 1))")

		require.Equal(t, "15", result.String())
	})
}

func TestParseUnaryMinus(t *testing.T) {
	t.Run("negative integer", func(t *testing.T) {
		result := parseExpression(t, "-5")

		require.Equal(t, "-5", result.String())
	})

	t.Run("negative grouped expression", func(t *testing.T) {
		result := parseExpression(t, "-(2 + 3)")

		require.Equal(t, "-5", result.String())
	})

	t.Run("double negative", func(t *testing.T) {
		result := parseExpression(t, "--5")

		require.Equal(t, "5", result.String())
	})

	t.Run("negative multiplication", func(t *testing.T) {
		result := parseExpression(t, "-3 * 4")

		require.Equal(t, "-12", result.String())
	})
}

func TestParseErrors(t *testing.T) {
	t.Run("empty input", func(t *testing.T) {
		parser := newTestParser("")

		_, err := parser.Parse()

		require.Error(t, err)
	})

	t.Run("illegal character", func(t *testing.T) {
		parser := newTestParser("@")

		_, err := parser.Parse()

		require.Error(t, err)
	})

	t.Run("trailing operator", func(t *testing.T) {
		parser := newTestParser("2 +")

		_, err := parser.Parse()

		require.Error(t, err)
	})

	t.Run("missing closing parenthesis", func(t *testing.T) {
		parser := newTestParser("(2 + 3")

		_, err := parser.Parse()

		require.Error(t, err)
	})

	t.Run("unexpected closing parenthesis", func(t *testing.T) {
		parser := newTestParser("2 + 3)")

		_, err := parser.Parse()

		require.Error(t, err)
	})

	t.Run("division by zero", func(t *testing.T) {
		parser := newTestParser("5 / 0")

		_, err := parser.Parse()

		require.ErrorIs(t, err, number.ErrDivisionByZero)
	})

	t.Run("two numbers without operator", func(t *testing.T) {
		parser := newTestParser("2 3")

		_, err := parser.Parse()

		require.Error(t, err)
	})
}

func parseValue(t *testing.T, input string) Value {
	t.Helper()

	result, err := newTestParser(input).ParseValue()
	require.NoError(t, err)

	return result
}

func defineTestVector(t *testing.T, name string, x, y int64) {
	t.Helper()

	err := testWorkspace.Define(name, vector.Vector2D{X: number.NewInteger(x), Y: number.NewInteger(y)})
	require.NoError(t, err)

	t.Cleanup(func() {
		testWorkspace.Delete(name)
	})
}

func requireVector(t *testing.T, value Value, x, y string) {
	t.Helper()

	require.True(t, value.IsVector)
	require.Equal(t, x, value.Vector.X.String())
	require.Equal(t, y, value.Vector.Y.String())
}

func TestParseVectorLiteral(t *testing.T) {
	t.Run("integer components", func(t *testing.T) {
		result := parseValue(t, "(3, 4)")
		requireVector(t, result, "3", "4")
	})

	t.Run("fraction components", func(t *testing.T) {
		result := parseValue(t, "(1/2, 3/2)")
		requireVector(t, result, "1/2", "3/2")
	})

	t.Run("negative components", func(t *testing.T) {
		result := parseValue(t, "(-3, 4)")
		requireVector(t, result, "-3", "4")
	})

	t.Run("component expressions", func(t *testing.T) {
		result := parseValue(t, "(2 + 3, 4 * 2)")
		requireVector(t, result, "5", "8")
	})
}

func TestParseVectorExpressions(t *testing.T) {
	defineTestVector(t, "a", 3, 4)
	defineTestVector(t, "b", 2, 1)

	t.Run("addition", func(t *testing.T) {
		requireVector(t, parseValue(t, "a + b"), "5", "5")
	})

	t.Run("subtraction", func(t *testing.T) {
		requireVector(t, parseValue(t, "a - b"), "1", "3")
	})

	t.Run("scale on the right", func(t *testing.T) {
		requireVector(t, parseValue(t, "a * 2"), "6", "8")
	})

	t.Run("scale on the left", func(t *testing.T) {
		requireVector(t, parseValue(t, "2 * a"), "6", "8")
	})

	t.Run("fractional scale", func(t *testing.T) {
		requireVector(t, parseValue(t, "a * 1/2"), "3/2", "2")
	})

	t.Run("division", func(t *testing.T) {
		requireVector(t, parseValue(t, "a / 2"), "3/2", "2")
	})

	t.Run("negation", func(t *testing.T) {
		requireVector(t, parseValue(t, "-a"), "-3", "-4")
	})

	t.Run("dot product", func(t *testing.T) {
		result := parseValue(t, "a * b")

		require.False(t, result.IsVector)
		require.Equal(t, "10", result.Scalar.String())
	})
}

func TestParseVectorPrecedence(t *testing.T) {
	defineTestVector(t, "a", 1, 1)
	defineTestVector(t, "b", 2, 3)

	t.Run("multiplication before addition", func(t *testing.T) {
		requireVector(t, parseValue(t, "a + b * 2"), "5", "7")
	})

	t.Run("parentheses override precedence", func(t *testing.T) {
		requireVector(t, parseValue(t, "(a + b) * 2"), "6", "8")
	})

	t.Run("nested vector expression", func(t *testing.T) {
		requireVector(t, parseValue(t, "a + (b - a) * 2"), "3", "5")
	})
}

func TestParseVectorErrors(t *testing.T) {
	defineTestVector(t, "a", 3, 4)
	defineTestVector(t, "b", 2, 1)

	tests := []struct {
		name  string
		input string
	}{
		{"undefined vector", "missing + a"},
		{"vector plus number", "a + 1"},
		{"number plus vector", "1 + a"},
		{"vector minus number", "a - 1"},
		{"number minus vector", "1 - a"},
		{"vector divided by vector", "a / b"},
		{"number divided by vector", "1 / a"},
		{"vector as a component", "(a, b)"},
		{"missing second component", "(1,)"},
		{"too many components", "(1, 2, 3)"},
		{"missing closing parenthesis", "(1, 2"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := newTestParser(test.input).ParseValue()
			require.Error(t, err)
		})
	}
}

func TestParseReturnsErrorForVector(t *testing.T) {
	defineTestVector(t, "a", 3, 4)

	_, err := newTestParser("a").Parse()
	require.Error(t, err)
}
