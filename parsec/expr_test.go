package parsec_test

// Integration test: arithmetic expression parser
//
// Grammar:
//   expr   = term   ('+' | '-') term)*
//   term   = factor (('*' | '/') factor)*
//   factor = '(' expr ')' | natural
//
// This demonstrates recursive parser construction using var + func.

import (
	"testing"

	"github.com/ajiyoshi-vg/goparsec/parsec"
)

func buildExprParser() parsec.Parser[int] {
	var expr parsec.Parser[int]

	factor := func(in parsec.Input) (int, parsec.Input, error) {
		paren := parsec.Between(
			parsec.Lexeme(parsec.Char('(')),
			parsec.Lexeme(parsec.Char(')')),
			parsec.Parser[int](func(in parsec.Input) (int, parsec.Input, error) {
				return expr(in)
			}),
		)
		return parsec.Choice(parsec.Lexeme(parsec.Natural()), paren)(in)
	}

	mulOp := parsec.Choice(
		parsec.Map(parsec.Lexeme(parsec.Char('*')), func(rune) func(int, int) int { return func(a, b int) int { return a * b } }),
		parsec.Map(parsec.Lexeme(parsec.Char('/')), func(rune) func(int, int) int { return func(a, b int) int { return a / b } }),
	)

	term := func(in parsec.Input) (int, parsec.Input, error) {
		return parsec.Map(
			parsec.Bind(factor, func(first int) parsec.Parser[int] {
				return parsec.Map(
					parsec.Many(parsec.Bind(mulOp, func(op func(int, int) int) parsec.Parser[int] {
						return parsec.Map(parsec.Parser[int](factor), func(v int) int { return op(first, v) })
					})),
					func(vs []int) int {
						if len(vs) == 0 {
							return first
						}
						return vs[len(vs)-1]
					},
				)
			}),
			func(v int) int { return v },
		)(in)
	}

	addOp := parsec.Choice(
		parsec.Map(parsec.Lexeme(parsec.Char('+')), func(rune) func(int, int) int { return func(a, b int) int { return a + b } }),
		parsec.Map(parsec.Lexeme(parsec.Char('-')), func(rune) func(int, int) int { return func(a, b int) int { return a - b } }),
	)

	expr = func(in parsec.Input) (int, parsec.Input, error) {
		return parsec.Map(
			parsec.Bind(parsec.Parser[int](term), func(first int) parsec.Parser[int] {
				return parsec.Map(
					parsec.Many(parsec.Bind(addOp, func(op func(int, int) int) parsec.Parser[int] {
						return parsec.Map(parsec.Parser[int](term), func(v int) int { return op(first, v) })
					})),
					func(vs []int) int {
						if len(vs) == 0 {
							return first
						}
						return vs[len(vs)-1]
					},
				)
			}),
			func(v int) int { return v },
		)(in)
	}

	return expr
}

func TestExpr(t *testing.T) {
	expr := buildExprParser()

	tests := []struct {
		input string
		want  int
	}{
		{"1", 1},
		{"42", 42},
		{"1 + 2", 3},
		{"10 - 3", 7},
		{"2 * 3", 6},
		{"8 / 4", 2},
		{"1 + 2 * 3", 7},
		{"(1 + 2) * 3", 9},
		{"(10 - 2) / 4", 2},
		{"2 * (3 + 4)", 14},
	}

	for _, tt := range tests {
		got, err := parsec.RunFull(expr, tt.input)
		if err != nil {
			t.Errorf("Expr(%q): %v", tt.input, err)
			continue
		}
		if got != tt.want {
			t.Errorf("Expr(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}
