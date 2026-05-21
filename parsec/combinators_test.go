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
	if got == nil {
		t.Error("Many: expected non-nil empty slice, got nil")
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

func TestReturn(t *testing.T) {
	got, err := parsec.Run(parsec.Return(42), "anything")
	if err != nil {
		t.Fatal(err)
	}
	if got != 42 {
		t.Errorf("got %d, want 42", got)
	}
	// must not consume input
	p := parsec.Then(parsec.Return("prefix"), parsec.String("hello"))
	s, err := parsec.Run(p, "hello")
	if err != nil {
		t.Fatal(err)
	}
	if s != "hello" {
		t.Errorf("got %q, want \"hello\"", s)
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
	// Parse N, then Count(N, 'a') — the continuation depends on the parsed value
	p := parsec.Bind(parsec.Natural(), func(n int) parsec.Parser[[]rune] {
		return parsec.Count(n, parsec.Char('a'))
	})
	got, err := parsec.Run(p, "3aaa")
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "aaa" {
		t.Errorf("got %q, want \"aaa\"", string(got))
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
	if got == nil {
		t.Error("SepBy: expected non-nil empty slice, got nil")
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

func TestSpaces(t *testing.T) {
	got, err := parsec.Run(parsec.Spaces(), "   hello")
	if err != nil {
		t.Fatal(err)
	}
	if got != "   " {
		t.Errorf("got %q, want \"   \"", got)
	}
}

func TestMany_nonConsuming_panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic: Many with non-consuming parser causes infinite loop")
		}
	}()
	// Return always succeeds without consuming input — the clearest trigger for this panic
	parsec.Run(parsec.Many(parsec.Return('x')), "bbb")
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
	addOp := parsec.Map(parsec.Char('+'), func(rune) func(int, int) int { return add })
	got, err := parsec.Run(parsec.Chainl1(parsec.Natural(), addOp), "1")
	if err != nil {
		t.Fatal(err)
	}
	if got != 1 {
		t.Errorf("got %d, want 1", got)
	}
}

func TestChainl1_multiple(t *testing.T) {
	addOp := parsec.Map(parsec.Char('+'), func(rune) func(int, int) int { return add })
	got, err := parsec.Run(parsec.Chainl1(parsec.Natural(), addOp), "1+2+3")
	if err != nil {
		t.Fatal(err)
	}
	if got != 6 {
		t.Errorf("got %d, want 6", got)
	}
}

func TestChainl1_leftAssoc(t *testing.T) {
	divOp := parsec.Map(parsec.Char('/'), func(rune) func(int, int) int { return div })
	got, err := parsec.Run(parsec.Chainl1(parsec.Natural(), divOp), "12/3/2")
	if err != nil {
		t.Fatal(err)
	}
	if got != 2 { // (12/3)/2 = 2、右結合なら 12/(3/2) = 8
		t.Errorf("got %d, want 2 (left-assoc)", got)
	}
}

func TestChainl1_fail(t *testing.T) {
	addOp := parsec.Map(parsec.Char('+'), func(rune) func(int, int) int { return add })
	_, err := parsec.Run(parsec.Chainl1(parsec.Natural(), addOp), "abc")
	if err == nil {
		t.Error("expected error when p doesn't match")
	}
}

func intPow(base, exp int) int {
	result := 1
	for range exp {
		result *= base
	}
	return result
}

func TestNotFollowedBy_success(t *testing.T) {
	// "if" followed by space: OK as keyword
	p := parsec.Then(parsec.String("if"), parsec.NotFollowedBy(parsec.AlphaNum()))
	_, err := parsec.Run(p, "if ")
	if err != nil {
		t.Fatalf("expected success: %v", err)
	}
}

func TestNotFollowedBy_atEOF(t *testing.T) {
	p := parsec.Then(parsec.String("if"), parsec.NotFollowedBy(parsec.AlphaNum()))
	_, err := parsec.Run(p, "if")
	if err != nil {
		t.Fatalf("expected success at EOF: %v", err)
	}
}

func TestNotFollowedBy_fail(t *testing.T) {
	// "ifs" must not match keyword "if"
	p := parsec.Then(parsec.String("if"), parsec.NotFollowedBy(parsec.AlphaNum()))
	_, err := parsec.Run(p, "ifs")
	if err == nil {
		t.Error("expected failure: 'ifs' should not match keyword 'if'")
	}
}

func TestLabel(t *testing.T) {
	p := parsec.Label(parsec.Digit(), "digit")
	_, err := parsec.Run(p, "abc")
	if err == nil {
		t.Fatal("expected error")
	}
	pe, ok := err.(*parsec.ParseError)
	if !ok {
		t.Fatalf("expected *parsec.ParseError, got %T", err)
	}
	if pe.Message != "expected digit" {
		t.Errorf("got message %q, want \"expected digit\"", pe.Message)
	}
	// Line/Col must be set (not zero)
	if pe.Line == 0 || pe.Col == 0 {
		t.Errorf("Label: Line=%d Col=%d, want non-zero", pe.Line, pe.Col)
	}
}

func TestCount(t *testing.T) {
	got, err := parsec.Run(parsec.Count(3, parsec.Digit()), "123abc")
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "123" {
		t.Errorf("got %q, want \"123\"", string(got))
	}
}

func TestCount_tooFew(t *testing.T) {
	_, err := parsec.Run(parsec.Count(4, parsec.Digit()), "123abc")
	if err == nil {
		t.Error("expected error: not enough matches")
	}
}

func TestManyTill(t *testing.T) {
	p := parsec.ManyTill(parsec.AnyChar(), parsec.String("*/"))
	got, err := parsec.Run(p, "hello*/world")
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "hello" {
		t.Errorf("got %q, want \"hello\"", string(got))
	}
}

func TestManyTill_empty(t *testing.T) {
	p := parsec.ManyTill(parsec.AnyChar(), parsec.String("*/"))
	got, err := parsec.Run(p, "*/rest")
	if err != nil {
		t.Fatal(err)
	}
	if got == nil {
		t.Error("ManyTill: expected non-nil empty slice, got nil")
	}
	if len(got) != 0 {
		t.Errorf("got %v, want empty slice", got)
	}
}

func TestManyTill_nonConsuming_panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic: ManyTill with non-consuming body causes infinite loop")
		}
	}()
	parsec.Run(parsec.ManyTill(parsec.Return('x'), parsec.String("*/")), "hello*/")
}

func TestChoice_noArgs_fails(t *testing.T) {
	_, err := parsec.Run(parsec.Choice[rune](), "a")
	if err == nil {
		t.Error("Choice with no parsers should fail, not succeed with zero value")
	}
}

func TestManyTill_noEnd(t *testing.T) {
	p := parsec.ManyTill(parsec.AnyChar(), parsec.String("*/"))
	_, err := parsec.Run(p, "hello world")
	if err == nil {
		t.Error("expected error when end delimiter not found")
	}
}

func TestChainr1_single(t *testing.T) {
	powOp := parsec.Map(parsec.Char('^'), func(rune) func(int, int) int { return intPow })
	got, err := parsec.Run(parsec.Chainr1(parsec.Natural(), powOp), "4")
	if err != nil {
		t.Fatal(err)
	}
	if got != 4 {
		t.Errorf("got %d, want 4", got)
	}
}

func TestChainr1_rightAssoc(t *testing.T) {
	powOp := parsec.Map(parsec.Char('^'), func(rune) func(int, int) int { return intPow })
	// 2^3^2 = 2^(3^2) = 2^9 = 512, not (2^3)^2 = 64
	got, err := parsec.Run(parsec.Chainr1(parsec.Natural(), powOp), "2^3^2")
	if err != nil {
		t.Fatal(err)
	}
	if got != 512 {
		t.Errorf("got %d, want 512 (right-assoc: 2^(3^2))", got)
	}
}

func TestChainr1_fail(t *testing.T) {
	powOp := parsec.Map(parsec.Char('^'), func(rune) func(int, int) int { return intPow })
	_, err := parsec.Run(parsec.Chainr1(parsec.Natural(), powOp), "abc")
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

func TestFloat(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"0", 0},
		{"42", 42},
		{"-7", -7},
		{"3.14", 3.14},
		{"-2.5", -2.5},
		{"1e10", 1e10},
		{"1.5e-3", 1.5e-3},
		{"2.0E+4", 2.0e4},
		{"-1.23e5", -1.23e5},
	}
	for _, tt := range tests {
		got, err := parsec.Run(parsec.Float(), tt.input)
		if err != nil {
			t.Fatalf("Float(%q): %v", tt.input, err)
		}
		if got != tt.want {
			t.Errorf("Float(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestFloat_fail(t *testing.T) {
	_, err := parsec.Run(parsec.Float(), "abc")
	if err == nil {
		t.Error("expected error for non-float input")
	}
}
