package parsec_test

import (
	"testing"

	"github.com/ajiyoshi-vg/goparsec/parsec"
)

func TestMany_zero(t *testing.T) {
	got, err := parsec.Run(parsec.Many(parsec.Char('a')), "bbb")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Errorf("got %v, want empty slice", got)
	}
}

func TestMany_multiple(t *testing.T) {
	got, err := parsec.Run(parsec.Many(parsec.Char('a')), "aaab")
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "aaa" {
		t.Errorf("got %q, want \"aaa\"", string(got))
	}
}

func TestMany1_one(t *testing.T) {
	got, err := parsec.Run(parsec.Many1(parsec.Digit()), "1abc")
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "1" {
		t.Errorf("got %q, want \"1\"", string(got))
	}
}

func TestMany1_fail(t *testing.T) {
	_, err := parsec.Run(parsec.Many1(parsec.Digit()), "abc")
	if err == nil {
		t.Error("expected error when no match")
	}
}

func TestChoice_first(t *testing.T) {
	p := parsec.Choice(parsec.Char('a'), parsec.Char('b'), parsec.Char('c'))
	got, err := parsec.Run(p, "apple")
	if err != nil {
		t.Fatal(err)
	}
	if got != 'a' {
		t.Errorf("got %q, want 'a'", got)
	}
}

func TestChoice_second(t *testing.T) {
	p := parsec.Choice(parsec.Char('a'), parsec.Char('b'), parsec.Char('c'))
	got, err := parsec.Run(p, "banana")
	if err != nil {
		t.Fatal(err)
	}
	if got != 'b' {
		t.Errorf("got %q, want 'b'", got)
	}
}

func TestChoice_fail(t *testing.T) {
	p := parsec.Choice(parsec.Char('a'), parsec.Char('b'))
	_, err := parsec.Run(p, "xyz")
	if err == nil {
		t.Error("expected error")
	}
}

func TestOption_present(t *testing.T) {
	got, err := parsec.Run(parsec.Option('_', parsec.Char('a')), "abc")
	if err != nil {
		t.Fatal(err)
	}
	if got != 'a' {
		t.Errorf("got %q, want 'a'", got)
	}
}

func TestOption_absent(t *testing.T) {
	got, err := parsec.Run(parsec.Option('_', parsec.Char('a')), "xyz")
	if err != nil {
		t.Fatal(err)
	}
	if got != '_' {
		t.Errorf("got %q, want '_'", got)
	}
}

func TestMap(t *testing.T) {
	p := parsec.Map(parsec.Digit(), func(r rune) int { return int(r - '0') })
	got, err := parsec.Run(p, "7abc")
	if err != nil {
		t.Fatal(err)
	}
	if got != 7 {
		t.Errorf("got %d, want 7", got)
	}
}

func TestBind(t *testing.T) {
	// Parse 'a' then based on that parse 'b'
	p := parsec.Bind(parsec.Char('a'), func(r rune) parsec.Parser[string] {
		return parsec.Map(parsec.Char('b'), func(r2 rune) string {
			return string([]rune{r, r2})
		})
	})
	got, err := parsec.Run(p, "ab")
	if err != nil {
		t.Fatal(err)
	}
	if got != "ab" {
		t.Errorf("got %q, want \"ab\"", got)
	}
}

func TestThen(t *testing.T) {
	p := parsec.Then(parsec.Char('('), parsec.Digit())
	got, err := parsec.Run(p, "(5")
	if err != nil {
		t.Fatal(err)
	}
	if got != '5' {
		t.Errorf("got %q, want '5'", got)
	}
}

func TestSkip(t *testing.T) {
	p := parsec.Skip(parsec.Digit(), parsec.Char(';'))
	got, err := parsec.Run(p, "3;")
	if err != nil {
		t.Fatal(err)
	}
	if got != '3' {
		t.Errorf("got %q, want '3'", got)
	}
}

func TestBetween(t *testing.T) {
	p := parsec.Between(parsec.Char('('), parsec.Char(')'), parsec.Digit())
	got, err := parsec.Run(p, "(5)")
	if err != nil {
		t.Fatal(err)
	}
	if got != '5' {
		t.Errorf("got %q, want '5'", got)
	}
}

func TestSepBy_zero(t *testing.T) {
	p := parsec.SepBy(parsec.Digit(), parsec.Char(','))
	got, err := parsec.Run(p, "abc")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Errorf("got %v, want empty", got)
	}
}

func TestSepBy_one(t *testing.T) {
	p := parsec.SepBy(parsec.Digit(), parsec.Char(','))
	got, err := parsec.Run(p, "3,")
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "3" {
		t.Errorf("got %q, want \"3\"", string(got))
	}
}

func TestSepBy_multiple(t *testing.T) {
	p := parsec.SepBy(parsec.Digit(), parsec.Char(','))
	got, err := parsec.Run(p, "1,2,3")
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "123" {
		t.Errorf("got %q, want \"123\"", string(got))
	}
}

func TestSepBy1_fail(t *testing.T) {
	p := parsec.SepBy1(parsec.Digit(), parsec.Char(','))
	_, err := parsec.Run(p, "abc")
	if err == nil {
		t.Error("expected error on empty match")
	}
}

func TestSpaces(t *testing.T) {
	got, err := parsec.Run(parsec.Spaces(), "   hello")
	if err != nil {
		t.Fatal(err)
	}
	if got != "   " {
		t.Errorf("got %q, want \"   \"", got)
	}
}

func TestLexeme(t *testing.T) {
	p := parsec.Lexeme(parsec.Digit())
	got, err := parsec.Run(p, "3   next")
	if err != nil {
		t.Fatal(err)
	}
	if got != '3' {
		t.Errorf("got %q, want '3'", got)
	}
}

func TestMany_nonConsuming_panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic: Many with non-consuming parser causes infinite loop")
		}
	}()
	// Option always succeeds without consuming on mismatch
	parsec.Run(parsec.Many(parsec.Option('x', parsec.Char('a'))), "bbb")
}

func TestChoice_furthestError(t *testing.T) {
	// "test" reaches pos 3, "true" reaches pos 1 on "tesar"
	// should report the furthest position regardless of parser order
	p := parsec.Choice(parsec.String("test"), parsec.String("true"))
	_, err := parsec.Run(p, "tesar")
	if err == nil {
		t.Fatal("expected error")
	}
	pe, ok := err.(*parsec.ParseError)
	if !ok {
		t.Fatalf("expected *parsec.ParseError, got %T", err)
	}
	if pe.Pos != 3 {
		t.Errorf("want error at pos 3 (furthest reached), got pos %d", pe.Pos)
	}
}

func TestChainl1_single(t *testing.T) {
	addOp := parsec.Map(parsec.Char('+'), func(rune) func(int, int) int { return func(a, b int) int { return a + b } })
	got, err := parsec.Run(parsec.Chainl1(parsec.Natural(), addOp), "1")
	if err != nil {
		t.Fatal(err)
	}
	if got != 1 {
		t.Errorf("got %d, want 1", got)
	}
}

func TestChainl1_multiple(t *testing.T) {
	addOp := parsec.Map(parsec.Char('+'), func(rune) func(int, int) int { return func(a, b int) int { return a + b } })
	got, err := parsec.Run(parsec.Chainl1(parsec.Natural(), addOp), "1+2+3")
	if err != nil {
		t.Fatal(err)
	}
	if got != 6 {
		t.Errorf("got %d, want 6", got)
	}
}

func TestChainl1_leftAssoc(t *testing.T) {
	divOp := parsec.Map(parsec.Char('/'), func(rune) func(int, int) int { return func(a, b int) int { return a / b } })
	got, err := parsec.Run(parsec.Chainl1(parsec.Natural(), divOp), "12/3/2")
	if err != nil {
		t.Fatal(err)
	}
	if got != 2 { // (12/3)/2 = 2、右結合なら 12/(3/2) = 8
		t.Errorf("got %d, want 2 (left-assoc)", got)
	}
}

func TestChainl1_fail(t *testing.T) {
	addOp := parsec.Map(parsec.Char('+'), func(rune) func(int, int) int { return func(a, b int) int { return a + b } })
	_, err := parsec.Run(parsec.Chainl1(parsec.Natural(), addOp), "abc")
	if err == nil {
		t.Error("expected error when p doesn't match")
	}
}

func TestInteger(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"0", 0},
		{"42", 42},
		{"-7", -7},
		{"-0", 0},
		{"123abc", 123},
	}
	for _, tt := range tests {
		got, err := parsec.Run(parsec.Integer(), tt.input)
		if err != nil {
			t.Fatalf("Integer(%q): %v", tt.input, err)
		}
		if got != tt.want {
			t.Errorf("Integer(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestInteger_fail(t *testing.T) {
	_, err := parsec.Run(parsec.Integer(), "abc")
	if err == nil {
		t.Error("expected error for non-integer input")
	}
}

func TestNatural(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"0", 0},
		{"42", 42},
		{"123abc", 123},
	}
	for _, tt := range tests {
		got, err := parsec.Run(parsec.Natural(), tt.input)
		if err != nil {
			t.Fatalf("Natural(%q): %v", tt.input, err)
		}
		if got != tt.want {
			t.Errorf("Natural(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}
