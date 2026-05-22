package parsec_test

import (
	"encoding/json"
	"strconv"
	"testing"
	"unicode/utf8"

	"github.com/ajiyoshi-vg/goparsec/parsec"
)

func TestManyChars(t *testing.T) {
	got, err := parsec.Run(parsec.ManyChars(parsec.Letter()), "hello123")
	if err != nil {
		t.Fatal(err)
	}
	if got != "hello" {
		t.Errorf("got %q, want \"hello\"", got)
	}
}

func TestManyChars_zero(t *testing.T) {
	got, err := parsec.Run(parsec.ManyChars(parsec.Letter()), "123")
	if err != nil {
		t.Fatal(err)
	}
	if got != "" {
		t.Errorf("got %q, want \"\"", got)
	}
}

func TestMany1Chars(t *testing.T) {
	got, err := parsec.Run(parsec.Many1Chars(parsec.AlphaNum()), "abc123 rest")
	if err != nil {
		t.Fatal(err)
	}
	if got != "abc123" {
		t.Errorf("got %q, want \"abc123\"", got)
	}
}

func TestMany1Chars_fail(t *testing.T) {
	_, err := parsec.Run(parsec.Many1Chars(parsec.Letter()), "123")
	if err == nil {
		t.Error("expected error when no match")
	}
}

func TestGoString(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		// basic
		{`""`, ""},
		{`"hello"`, "hello"},
		// simple escapes
		{`"hello\nworld"`, "hello\nworld"},
		{`"tab\there"`, "tab\there"},
		{`"back\\slash"`, `back\slash`},
		{`"say \"hi\""`, `say "hi"`},
		{`"\a\b\f\r\v"`, "\a\b\f\r\v"},
		// hex escape \xNN
		{`"\x41"`, "A"},
		{`"\x61\x62\x63"`, "abc"},
		// unicode escape \uNNNN
		{`"中文"`, "中文"},
		// unicode escape \UNNNNNNNN
		{`"\U0001F600"`, "😀"},
	}

	for _, tt := range tests {
		got, err := parsec.Run(parsec.GoString(), tt.input)
		if err != nil {
			t.Errorf("GoString(%q): %v", tt.input, err)
			continue
		}
		if got != tt.want {
			t.Errorf("GoString(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

// FuzzGoString verifies the round-trip invariant:
// GoString(strconv.Quote(s)) == s for any string s.
func FuzzGoString(f *testing.F) {
	seeds := []string{
		"",
		"hello",
		`say "hi"`,
		"line1\nline2",
		"tab\there",
		`back\slash`,
		"\a\b\f\r\v",
		"中文",
		"😀",
		"\x00\x01",
	}
	for _, s := range seeds {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		if !utf8.ValidString(s) {
			return
		}
		quoted := strconv.Quote(s)
		got, err := parsec.Run(parsec.GoString(), quoted)
		if err != nil {
			t.Errorf("GoString(%q): unexpected error: %v", quoted, err)
			return
		}
		if got != s {
			t.Errorf("GoString(Quote(%q))\n  got  %q\n  want %q", s, got, s)
		}
	})
}

func TestGoString_invalid(t *testing.T) {
	tests := []struct {
		desc  string
		input string
	}{
		{"unterminated", `"hello`},
		{"unknown escape", `"\z"`},
		{"short hex", `"\x4"`},
		{"short unicode4", `"\u4e2"`},
		{"short unicode8", `"\U0001F60"`},
	}

	for _, tt := range tests {
		_, err := parsec.Run(parsec.GoString(), tt.input)
		if err == nil {
			t.Errorf("GoString(%s=%q): expected error, got nil", tt.desc, tt.input)
		}
	}
}

func TestJSONString(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		// basic
		{`""`, ""},
		{`"hello"`, "hello"},
		// simple escapes
		{`"hello\nworld"`, "hello\nworld"},
		{`"tab\there"`, "tab\there"},
		{`"back\\slash"`, `back\slash`},
		{`"say \"hi\""`, `say "hi"`},
		{`"\b\f\r"`, "\b\f\r"},
		// JSON-specific: \/ is a valid escape
		{`"\/"`, "/"},
		// unicode escape \uNNNN
		{`"中文"`, "中文"},
		{`"A"`, "A"},
		// surrogate pair for U+1F600 (😀)
		{`"😀"`, "😀"},
	}

	for _, tt := range tests {
		got, err := parsec.Run(parsec.JSONString(), tt.input)
		if err != nil {
			t.Errorf("JSONString(%q): %v", tt.input, err)
			continue
		}
		if got != tt.want {
			t.Errorf("JSONString(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestJSONString_invalid(t *testing.T) {
	tests := []struct {
		desc  string
		input string
	}{
		{"unterminated", `"hello`},
		{"unknown escape \\a", `"\a"`},
		{"unknown escape \\x", `"\x41"`},
		{"short unicode", `"\u4e2"`},
		{"lone high surrogate", `"\uD83D"`},
		{"high surrogate + non-low", `"\uD83DA"`},
		{"unescaped control", "\"" + "\x00" + "\""},
	}

	for _, tt := range tests {
		_, err := parsec.Run(parsec.JSONString(), tt.input)
		if err == nil {
			t.Errorf("JSONString(%s=%q): expected error, got nil", tt.desc, tt.input)
		}
	}
}

// FuzzJSONString verifies the round-trip invariant:
// JSONString(json.Marshal(s)) == s for any valid UTF-8 string s.
func FuzzJSONString(f *testing.F) {
	seeds := []string{
		"",
		"hello",
		`say "hi"`,
		"line1\nline2",
		"tab\there",
		`back\slash`,
		"\b\f\r",
		"中文",
		"😀",
		"\x00\x01",
	}
	for _, s := range seeds {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		if !utf8.ValidString(s) {
			return
		}
		b, err := json.Marshal(s)
		if err != nil {
			t.Fatalf("json.Marshal(%q): %v", s, err)
		}
		quoted := string(b)
		got, err := parsec.Run(parsec.JSONString(), quoted)
		if err != nil {
			t.Errorf("JSONString(%q): unexpected error: %v", quoted, err)
			return
		}
		if got != s {
			t.Errorf("JSONString(json.Marshal(%q))\n  got  %q\n  want %q", s, got, s)
		}
	})
}

// TestGoString_allocs verifies that GoString avoids the []rune intermediate slice.
// A strings.Builder-based loop should eliminate the []rune allocation and the
// rune-to-string conversion.
func TestGoString_allocs(t *testing.T) {
	p := parsec.GoString()
	allocs := testing.AllocsPerRun(100, func() {
		parsec.Run(p, `"hello"`)
	})
	// 2 (NewInput) + 7 (Advance per char incl. quotes) + ~2 (Builder buf) = ~11
	if allocs > 13 {
		t.Errorf("GoString allocs = %.0f, want ≤13", allocs)
	}
}

// TestJSONString_allocs mirrors TestGoString_allocs for the JSON variant.
func TestJSONString_allocs(t *testing.T) {
	p := parsec.JSONString()
	allocs := testing.AllocsPerRun(100, func() {
		parsec.Run(p, `"hello"`)
	})
	if allocs > 13 {
		t.Errorf("JSONString allocs = %.0f, want ≤13", allocs)
	}
}

func BenchmarkGoString(b *testing.B) {
	p := parsec.GoString()
	inputs := []string{
		`"hello"`,
		`"hello\nworld\ttab"`,
		`"back\\slash and \"quotes\""`,
		`"中文"`,
	}
	for b.Loop() {
		for _, in := range inputs {
			parsec.Run(p, in)
		}
	}
}

func BenchmarkJSONString(b *testing.B) {
	p := parsec.JSONString()
	inputs := []string{
		`"hello"`,
		`"hello\nworld\ttab"`,
		`"back\\slash and \"quotes\""`,
		`"中文"`,
	}
	for b.Loop() {
		for _, in := range inputs {
			parsec.Run(p, in)
		}
	}
}
