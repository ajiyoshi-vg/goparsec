package parsec_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/ajiyoshi-vg/goparsec/parsec"
)

func TestTryMap_success(t *testing.T) {
	// ManyChars(Digit()) → string → strconv.Atoi
	p := parsec.TryMap(parsec.ManyChars(parsec.Digit()), strconv.Atoi)
	got, err := parsec.RunString(p, "123")
	if err != nil {
		t.Fatal(err)
	}
	if got != 123 {
		t.Errorf("got %d, want 123", got)
	}
}

func TestTryMap_parserFails(t *testing.T) {
	p := parsec.TryMap(parsec.ManyChars(parsec.Digit()), strconv.Atoi)
	_, err := parsec.RunString(p, "abc")
	// ManyChars succeeds with "" then Atoi("") fails
	if err == nil {
		t.Error("expected error")
	}
}

func TestTryMap_fFails(t *testing.T) {
	mustPositive := func(n int) (int, error) {
		if n <= 0 {
			return 0, fmt.Errorf("must be positive, got %d", n)
		}
		return n, nil
	}
	p := parsec.TryMap(parsec.Integer(), mustPositive)

	got, err := parsec.RunString(p, "42")
	if err != nil {
		t.Fatal(err)
	}
	if got != 42 {
		t.Errorf("got %d, want 42", got)
	}

	_, err = parsec.RunString(p, "-1")
	if err == nil {
		t.Error("expected error for non-positive")
	}
}

func TestTryMap_errorIsParseError(t *testing.T) {
	mustPositive := func(n int) (int, error) {
		if n <= 0 {
			return 0, fmt.Errorf("must be positive")
		}
		return n, nil
	}
	p := parsec.TryMap(parsec.Integer(), mustPositive)

	_, err := parsec.RunString(p, "-1")
	if _, ok := err.(*parsec.ParseError); !ok {
		t.Errorf("expected *ParseError, got %T: %v", err, err)
	}
}

// divide parses "numerator/denominator" and fails on zero denominator.
func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}

func TestTryMap2_success(t *testing.T) {
	slash := parsec.Then(parsec.Spaces(), parsec.Char('/'))
	p := parsec.TryMap2(parsec.Natural(), parsec.Then(slash, parsec.Natural()), divide)

	got, err := parsec.RunString(p, "10/2")
	if err != nil {
		t.Fatal(err)
	}
	if got != 5 {
		t.Errorf("got %d, want 5", got)
	}
}

func TestTryMap2_fFails(t *testing.T) {
	slash := parsec.Then(parsec.Spaces(), parsec.Char('/'))
	p := parsec.TryMap2(parsec.Natural(), parsec.Then(slash, parsec.Natural()), divide)

	_, err := parsec.RunString(p, "10/0")
	if err == nil {
		t.Error("expected error for division by zero")
	}
	if _, ok := err.(*parsec.ParseError); !ok {
		t.Errorf("expected *ParseError, got %T: %v", err, err)
	}
}

type YMD struct{ Year, Month, Day int }

func makeYMD(y, m, d int) (YMD, error) {
	if m < 1 || m > 12 {
		return YMD{}, fmt.Errorf("invalid month: %d", m)
	}
	if d < 1 || d > 31 {
		return YMD{}, fmt.Errorf("invalid day: %d", d)
	}
	return YMD{y, m, d}, nil
}

func TestTryMap3_success(t *testing.T) {
	dash := parsec.Then(parsec.Spaces(), parsec.Char('-'))
	p := parsec.TryMap3(
		parsec.Natural(), parsec.Then(dash, parsec.Natural()), parsec.Then(dash, parsec.Natural()),
		makeYMD,
	)

	got, err := parsec.RunString(p, "2026-5-22")
	if err != nil {
		t.Fatal(err)
	}
	if got != (YMD{2026, 5, 22}) {
		t.Errorf("got %v, want {2026 5 22}", got)
	}
}

func TestTryMap3_fFails(t *testing.T) {
	dash := parsec.Then(parsec.Spaces(), parsec.Char('-'))
	p := parsec.TryMap3(
		parsec.Natural(), parsec.Then(dash, parsec.Natural()), parsec.Then(dash, parsec.Natural()),
		makeYMD,
	)

	_, err := parsec.RunString(p, "2026-13-1")
	if err == nil {
		t.Error("expected error for invalid month")
	}
}

func TestTryMap4_success(t *testing.T) {
	colon := parsec.Then(parsec.Spaces(), parsec.Char(':'))
	p := parsec.TryMap4(
		parsec.Natural(), parsec.Then(colon, parsec.Natural()),
		parsec.Then(colon, parsec.Natural()), parsec.Then(colon, parsec.Natural()),
		func(a, b, c, d int) (int, error) {
			if a < 0 || b < 0 || c < 0 || d < 0 {
				return 0, fmt.Errorf("negative value")
			}
			return a + b + c + d, nil
		},
	)

	got, err := parsec.RunString(p, "1:2:3:4")
	if err != nil {
		t.Fatal(err)
	}
	if got != 10 {
		t.Errorf("got %d, want 10", got)
	}
}
