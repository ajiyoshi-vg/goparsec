package parsec_test

// Integration test: arithmetic expression parser
//
// Grammar (all operators left-associative):
//   expr   = chainl1(term,   addOp)
//   term   = chainl1(factor, mulOp)
//   factor = '(' expr ')' | integer

import (
	"testing"

	"github.com/ajiyoshi-vg/goparsec/input"
	"github.com/ajiyoshi-vg/goparsec/parsec"
)

func add(a, b int) int { return a + b }
func sub(a, b int) int { return a - b }
func mul(a, b int) int { return a * b }
func div(a, b int) int { return a / b }

func buildExprParser() parsec.Parser[int] {
	var expr parsec.Parser[int]
	w := parsec.Spaces()

	tok := func(c rune) parsec.Parser[rune] {
		return parsec.Then(w, parsec.Char(c))
	}

	factor := func(in input.Input) (int, input.Input, error) {
		paren := parsec.Between(
			tok('('),
			parsec.Parser[int](func(in input.Input) (int, input.Input, error) {
				return expr(in)
			}),
			tok(')'),
		)
		return parsec.Choice(parsec.Then(w, parsec.Integer()), paren)(in)
	}

	mulOp := parsec.Choice(parsec.Value(tok('*'), mul), parsec.Value(tok('/'), div))
	addOp := parsec.Choice(parsec.Value(tok('+'), add), parsec.Value(tok('-'), sub))

	term := parsec.Chainl1(parsec.Parser[int](factor), mulOp)
	expr = parsec.Chainl1(term, addOp)

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
		{"12 / 3 / 2", 2},  // 左結合: (12/3)/2 = 2
		{"10 - 3 - 2", 5},  // 左結合: (10-3)-2 = 5
		// 負数リテラル
		{"-1", -1},
		{"-2 + 3", 1},
		{"2 * -3", -6},
		{"1 - -2", 3},
		{"(-4 + 1) * -2", 6},
	}

	for _, tt := range tests {
		got, err := parsec.RunStringFull(expr, tt.input)
		if err != nil {
			t.Errorf("Expr(%q): %v", tt.input, err)
			continue
		}
		if got != tt.want {
			t.Errorf("Expr(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func BenchmarkExpr(b *testing.B) {
	expr := buildExprParser()
	inputs := []string{
		"1",
		"1 + 2 * 3",
		"(1 + 2) * (3 - 4)",
		"10 - 3 - 2",
		"(-4 + 1) * -2",
	}
	for b.Loop() {
		for _, in := range inputs {
			parsec.RunStringFull(expr, in)
		}
	}
}
