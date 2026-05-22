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

func TestString(t *testing.T) {
	got, err := parsec.Run(parsec.String("hello"), "hello world")
	if err != nil {
		t.Fatal(err)
	}
	if got != "hello" {
		t.Errorf("got %q, want \"hello\"", got)
	}
}

func TestString_fail(t *testing.T) {
	_, err := parsec.Run(parsec.String("world"), "hello")
	if err == nil {
		t.Error("expected error")
	}
}

func TestString_partialEOF(t *testing.T) {
	_, err := parsec.Run(parsec.String("hello"), "hel")
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

func BenchmarkString(b *testing.B) {
	p := parsec.String("hello")
	for b.Loop() {
		parsec.Run(p, "hello world")
	}
}
