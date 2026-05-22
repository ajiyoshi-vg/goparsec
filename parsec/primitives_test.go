package parsec_test

import (
	"testing"

	"github.com/ajiyoshi-vg/goparsec/parsec"
)

func TestSatisfy_match(t *testing.T) {
	p := parsec.Satisfy(func(r rune) bool { return r == 'a' })
	got, err := parsec.Run(p, "abc")
	if err != nil {
		t.Fatal(err)
	}
	if got != 'a' {
		t.Errorf("got %q, want 'a'", got)
	}
}

func TestSatisfy_noMatch(t *testing.T) {
	p := parsec.Satisfy(func(r rune) bool { return r == 'a' })
	_, err := parsec.Run(p, "xyz")
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestSatisfy_EOF(t *testing.T) {
	p := parsec.Satisfy(func(r rune) bool { return true })
	_, err := parsec.Run(p, "")
	if err == nil {
		t.Error("expected error on empty input")
	}
}

func TestChar(t *testing.T) {
	got, err := parsec.Run(parsec.Char('h'), "hello")
	if err != nil {
		t.Fatal(err)
	}
	if got != 'h' {
		t.Errorf("got %q, want 'h'", got)
	}
}

func TestChar_fail(t *testing.T) {
	_, err := parsec.Run(parsec.Char('z'), "hello")
	if err == nil {
		t.Error("expected error")
	}
}

func TestLiteral(t *testing.T) {
	got, err := parsec.Run(parsec.Literal("hello"), "hello world")
	if err != nil {
		t.Fatal(err)
	}
	if got != "hello" {
		t.Errorf("got %q, want \"hello\"", got)
	}
}

func TestString_fail(t *testing.T) {
	_, err := parsec.Run(parsec.Literal("world"), "hello")
	if err == nil {
		t.Error("expected error")
	}
}

func TestString_partialEOF(t *testing.T) {
	_, err := parsec.Run(parsec.Literal("hello"), "hel")
	if err == nil {
		t.Error("expected error on short input")
	}
}

func TestEOF_success(t *testing.T) {
	_, err := parsec.Run(parsec.EOF(), "")
	if err != nil {
		t.Fatal(err)
	}
}

func TestEOF_fail(t *testing.T) {
	_, err := parsec.Run(parsec.EOF(), "a")
	if err == nil {
		t.Error("expected error when input remains")
	}
}

func TestDigit(t *testing.T) {
	for _, c := range "0123456789" {
		got, err := parsec.Run(parsec.Digit(), string(c))
		if err != nil {
			t.Fatalf("Digit(%q): %v", c, err)
		}
		if got != c {
			t.Errorf("Digit(%q) = %q, want %q", c, got, c)
		}
	}
}

func TestDigit_fail(t *testing.T) {
	_, err := parsec.Run(parsec.Digit(), "a")
	if err == nil {
		t.Error("expected error for non-digit")
	}
}

func TestLetter(t *testing.T) {
	for _, c := range "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		got, err := parsec.Run(parsec.Letter(), string(c))
		if err != nil {
			t.Fatalf("Letter(%q): %v", c, err)
		}
		if got != c {
			t.Errorf("Letter(%q) = %q, want %q", c, got, c)
		}
	}
}

func TestLetter_fail(t *testing.T) {
	_, err := parsec.Run(parsec.Letter(), "1")
	if err == nil {
		t.Error("expected error for non-letter")
	}
}

func TestAlphaNum(t *testing.T) {
	for _, c := range "abcXYZ019" {
		got, err := parsec.Run(parsec.AlphaNum(), string(c))
		if err != nil {
			t.Fatalf("AlphaNum(%q): %v", c, err)
		}
		if got != c {
			t.Errorf("AlphaNum(%q) = %q, want %q", c, got, c)
		}
	}
}

func TestAlphaNum_fail(t *testing.T) {
	for _, c := range "_ !@" {
		_, err := parsec.Run(parsec.AlphaNum(), string(c))
		if err == nil {
			t.Errorf("AlphaNum(%q): expected error", c)
		}
	}
}

func TestHexDigit(t *testing.T) {
	for _, c := range "0123456789abcdefABCDEF" {
		got, err := parsec.Run(parsec.HexDigit(), string(c))
		if err != nil {
			t.Fatalf("HexDigit(%q): %v", c, err)
		}
		if got != c {
			t.Errorf("HexDigit(%q) = %q, want %q", c, got, c)
		}
	}
}

func TestHexDigit_fail(t *testing.T) {
	for _, c := range "ghGH_z" {
		_, err := parsec.Run(parsec.HexDigit(), string(c))
		if err == nil {
			t.Errorf("HexDigit(%q): expected error", c)
		}
	}
}

func TestRunFull_success(t *testing.T) {
	_, err := parsec.RunFull(parsec.Char('a'), "a")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRunFull_extraInput(t *testing.T) {
	_, err := parsec.RunFull(parsec.Char('a'), "ab")
	if err == nil {
		t.Error("expected error when input not fully consumed")
	}
}

// TestCustomParser_politeErrors verifies that user-written custom parsers can
// use ErrNoMatch for soft failures and NewError for positioned hard failures,
// and that these interact correctly with Choice and error reporting.
func TestCustomParser_politeErrors(t *testing.T) {
	// Custom parser: matches "42", returns ErrNoMatch on any other input.
	fortyTwo := func(in parsec.Input) (int, parsec.Input, error) {
		c, ok := in.Head()
		if !ok || c != '4' {
			return 0, in, parsec.ErrNoMatch // soft failure: caller may try alternatives
		}
		cur := in.Advance()
		c, ok = cur.Head()
		if !ok || c != '2' {
			return 0, in, parsec.NewError(cur, "expected '2' after '4'") // hard failure
		}
		return 42, cur.Advance(), nil
	}

	// ErrNoMatch lets Choice fall through to the next alternative.
	p := parsec.Choice(parsec.Parser[int](fortyTwo), parsec.Natural())
	got, err := parsec.Run(p, "99")
	if err != nil {
		t.Fatalf("Choice fallthrough: unexpected error: %v", err)
	}
	if got != 99 {
		t.Errorf("Choice fallthrough: got %d, want 99", got)
	}

	// Hard failure from NewError propagates with correct position.
	_, err = parsec.Run(parsec.Parser[int](fortyTwo), "49")
	if err == nil {
		t.Fatal("expected error for '49'")
	}
	pe, ok := err.(*parsec.ParseError)
	if !ok {
		t.Fatalf("expected *ParseError, got %T: %v", err, err)
	}
	if pe.Col != 2 {
		t.Errorf("error col = %d, want 2 (position of failed '2')", pe.Col)
	}
}

// TestNewErrorf verifies that NewErrorf produces a *ParseError with the formatted
// message at the correct position, matching NewError(in, fmt.Sprintf(format, args...)).
func TestNewErrorf(t *testing.T) {
	p := func(in parsec.Input) (rune, parsec.Input, error) {
		c, ok := in.Head()
		if !ok {
			return 0, in, parsec.NewErrorf(in, "expected digit, got EOF")
		}
		if c < '0' || c > '9' {
			return 0, in, parsec.NewErrorf(in, "expected digit, got %q", c)
		}
		return c, in.Advance(), nil
	}

	// Success path unchanged.
	got, err := parsec.Run(parsec.Parser[rune](p), "5abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != '5' {
		t.Errorf("got %q, want '5'", got)
	}

	// Failure: formatted message and correct position.
	_, err = parsec.Run(parsec.Parser[rune](p), "abc")
	if err == nil {
		t.Fatal("expected error")
	}
	pe, ok := err.(*parsec.ParseError)
	if !ok {
		t.Fatalf("expected *ParseError, got %T: %v", err, err)
	}
	if pe.Message != `expected digit, got 'a'` {
		t.Errorf("message = %q, want %q", pe.Message, `expected digit, got 'a'`)
	}
	if pe.Col != 1 {
		t.Errorf("col = %d, want 1", pe.Col)
	}
}

// TestChoice_allocsOnSoftFail verifies that failed alternatives in Choice
// do not allocate heap objects (*ParseError).
// The minimum unavoidable allocs are 3: []rune (NewInput), stringInput interface
// boxing (NewInput), and stringInput interface boxing (Advance on success).
func TestChoice_allocsOnSoftFail(t *testing.T) {
	p := parsec.Choice(parsec.Char('+'), parsec.Char('-'), parsec.Char('*'), parsec.Char('/'))
	allocs := testing.AllocsPerRun(100, func() {
		parsec.Run(p, "/") // last alternative matches; 3 alternatives fail first
	})
	// 3 = minimum with immutable Input (NewInput×2 + Advance×1).
	// Before this optimization: 9 allocs (3 failing *ParseError × 2 each via fmt.Sprintf).
	if allocs > 3 {
		t.Errorf("Choice soft-fail allocs = %.0f, want ≤3", allocs)
	}
}

func BenchmarkChar(b *testing.B) {
	p := parsec.Char('a')
	for b.Loop() {
		parsec.Run(p, "abc")
	}
}

func BenchmarkSatisfy(b *testing.B) {
	p := parsec.Satisfy(func(r rune) bool { return r >= 'a' && r <= 'z' })
	for b.Loop() {
		parsec.Run(p, "hello")
	}
}

func BenchmarkLiteral(b *testing.B) {
	p := parsec.Literal("hello")
	for b.Loop() {
		parsec.Run(p, "hello world")
	}
}
