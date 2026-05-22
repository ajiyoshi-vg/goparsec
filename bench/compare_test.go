// Package bench provides side-by-side benchmarks comparing goparsec and parseg.
//
// Design differences:
//   - goparsec: immutable Input (rune slice + position), value-based (T, Input, error)
//   - parseg:   mutable stream (io.ReadSeeker), pointer-based (*T, int, error)
package bench

import (
	"bytes"
	"testing"

	"github.com/ajiyoshi-vg/goparsec/parsec"
	"github.com/ajiyoshi-vg/parseg"
	parsegexpr "github.com/ajiyoshi-vg/parseg/expr"
)

// --- Natural ---

func BenchmarkGoparsec_Natural(b *testing.B) {
	p := parsec.Natural()
	for b.Loop() {
		parsec.Run(p, "12345")
	}
}

func BenchmarkParseg_Natural(b *testing.B) {
	p := parseg.Natural()
	input := []byte("12345")
	for b.Loop() {
		r := bytes.NewReader(input)
		p.Parse(r)
	}
}

// --- Many(Digit) ---

func BenchmarkGoparsec_ManyDigit(b *testing.B) {
	p := parsec.Many(parsec.Digit())
	for b.Loop() {
		parsec.Run(p, "1234567890")
	}
}

func BenchmarkParseg_ManyDigit(b *testing.B) {
	p := parseg.Many(parseg.Digit())
	input := []byte("1234567890")
	for b.Loop() {
		r := bytes.NewReader(input)
		p.Parse(r)
	}
}

// --- Choice/OneOf (hit first alternative) ---

func BenchmarkGoparsec_Choice_first(b *testing.B) {
	p := parsec.Choice(parsec.Char('+'), parsec.Char('-'), parsec.Char('*'), parsec.Char('/'))
	for b.Loop() {
		parsec.Run(p, "+1")
	}
}

func BenchmarkParseg_OneOf_first(b *testing.B) {
	p := parseg.OneOf(parseg.Rune('+'), parseg.Rune('-'), parseg.Rune('*'), parseg.Rune('/'))
	input := []byte("+1")
	for b.Loop() {
		r := bytes.NewReader(input)
		p.Parse(r)
	}
}

// --- Choice/OneOf (hit last alternative) ---

func BenchmarkGoparsec_Choice_last(b *testing.B) {
	p := parsec.Choice(parsec.Char('+'), parsec.Char('-'), parsec.Char('*'), parsec.Char('/'))
	for b.Loop() {
		parsec.Run(p, "/1")
	}
}

func BenchmarkParseg_OneOf_last(b *testing.B) {
	p := parseg.OneOf(parseg.Rune('+'), parseg.Rune('-'), parseg.Rune('*'), parseg.Rune('/'))
	input := []byte("/1")
	for b.Loop() {
		r := bytes.NewReader(input)
		p.Parse(r)
	}
}

// --- String ---

func BenchmarkGoparsec_String(b *testing.B) {
	p := parsec.String("hello")
	for b.Loop() {
		parsec.Run(p, "hello world")
	}
}

func BenchmarkParseg_String(b *testing.B) {
	p := parseg.String("hello")
	input := []byte("hello world")
	for b.Loop() {
		r := bytes.NewReader(input)
		p.Parse(r)
	}
}

// --- Expression parser (no spaces: parseg does not skip whitespace) ---

func buildGoparsecExpr() parsec.Parser[int] {
	var expr parsec.Parser[int]

	add := func(a, b int) int { return a + b }
	sub := func(a, b int) int { return a - b }
	mul := func(a, b int) int { return a * b }
	div := func(a, b int) int { return a / b }

	opParser := func(c rune, fn func(int, int) int) parsec.Parser[func(int, int) int] {
		return parsec.Map(parsec.Char(c), func(rune) func(int, int) int { return fn })
	}

	factor := parsec.Parser[int](func(in parsec.Input) (int, parsec.Input, error) {
		paren := parsec.Between(
			parsec.Char('('),
			parsec.Char(')'),
			parsec.Parser[int](func(in parsec.Input) (int, parsec.Input, error) {
				return expr(in)
			}),
		)
		return parsec.Choice(parsec.Integer(), paren)(in)
	})

	mulOp := parsec.Choice(opParser('*', mul), opParser('/', div))
	addOp := parsec.Choice(opParser('+', add), opParser('-', sub))
	term := parsec.Chainl1(factor, mulOp)
	expr = parsec.Chainl1(term, addOp)
	return expr
}

func BenchmarkGoparsec_Expr(b *testing.B) {
	expr := buildGoparsecExpr()
	inputs := []string{
		"1",
		"1+2*3",
		"(1+2)*(3-4)",
		"10-3-2",
		"1+2*6/(10-7)",
	}
	for b.Loop() {
		for _, in := range inputs {
			parsec.Run(expr, in)
		}
	}
}

func BenchmarkParseg_Expr(b *testing.B) {
	expr := parsegexpr.Parser()
	inputs := [][]byte{
		[]byte("1"),
		[]byte("1+2*3"),
		[]byte("(1+2)*(3-4)"),
		[]byte("10-3-2"),
		[]byte("1+2*6/(10-7)"),
	}
	for b.Loop() {
		for _, in := range inputs {
			r := bytes.NewReader(in)
			expr.Parse(r)
		}
	}
}
