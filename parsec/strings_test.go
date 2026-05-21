package parsec_test

import (
	"testing"

	"github.com/ajiyoshi-vg/goparsec/parsec"
)

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
