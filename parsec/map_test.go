package parsec_test

import (
	"testing"

	"github.com/ajiyoshi-vg/goparsec/parsec"
)

func TestMap2_success(t *testing.T) {
	p := parsec.Map2(parsec.Char('a'), parsec.Char('b'), func(a, b rune) string {
		return string([]rune{a, b})
	})
	got, err := parsec.RunString(p, "ab")
	if err != nil {
		t.Fatal(err)
	}
	if got != "ab" {
		t.Errorf("got %q, want \"ab\"", got)
	}
}

func TestMap2_firstFails(t *testing.T) {
	p := parsec.Map2(parsec.Char('a'), parsec.Char('b'), func(a, b rune) string {
		return string([]rune{a, b})
	})
	_, err := parsec.RunString(p, "xb")
	if err == nil {
		t.Error("expected error")
	}
}

func TestMap2_secondFails(t *testing.T) {
	p := parsec.Map2(parsec.Char('a'), parsec.Char('b'), func(a, b rune) string {
		return string([]rune{a, b})
	})
	_, err := parsec.RunString(p, "ax")
	if err == nil {
		t.Error("expected error")
	}
}

func TestMap3_success(t *testing.T) {
	p := parsec.Map3(parsec.Char('a'), parsec.Char('b'), parsec.Char('c'),
		func(a, b, c rune) string { return string([]rune{a, b, c}) })
	got, err := parsec.RunString(p, "abc")
	if err != nil {
		t.Fatal(err)
	}
	if got != "abc" {
		t.Errorf("got %q, want \"abc\"", got)
	}
}

func TestMap4_success(t *testing.T) {
	p := parsec.Map4(parsec.Char('a'), parsec.Char('b'), parsec.Char('c'), parsec.Char('d'),
		func(a, b, c, d rune) string { return string([]rune{a, b, c, d}) })
	got, err := parsec.RunString(p, "abcd")
	if err != nil {
		t.Fatal(err)
	}
	if got != "abcd" {
		t.Errorf("got %q, want \"abcd\"", got)
	}
}

func TestMap2_differentTypes(t *testing.T) {
	p := parsec.Map2(parsec.Char('-'), parsec.Natural(), func(_ rune, n int) int { return -n })
	got, err := parsec.RunString(p, "-42")
	if err != nil {
		t.Fatal(err)
	}
	if got != -42 {
		t.Errorf("got %d, want -42", got)
	}
}
