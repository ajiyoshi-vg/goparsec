package parsec_test

// Integration test: JSON parser
//
// Grammar:
//   value  = null | true | false | number | string | array | object
//   number = '-'? digit+ ('.' digit+)?
//   string = '"' [^"\\]* '"'
//   array  = '[' (value (',' value)*)? ']'
//   object = '{' (key ':' value (',' key ':' value)*)? '}'

import (
	"reflect"
	"testing"

	"github.com/ajiyoshi-vg/goparsec/parsec"
)

// Conversion functions: parsed text → Go values
func jsonNull(_ string) any   { return nil }
func jsonTrue(_ string) any   { return true }
func jsonFalse(_ string) any  { return false }
func jsonString(s string) any { return s }
func jsonArray(vs []any) any  { return vs }
func jsonObject(pairs [][2]any) any {
	m := make(map[string]any, len(pairs))
	for _, p := range pairs {
		m[p[0].(string)] = p[1]
	}
	return m
}

func newJSONParser() parsec.Parser[any] {
	var jsonValue parsec.Parser[any]
	lazy := parsec.Parser[any](func(in parsec.Input) (any, parsec.Input, error) { return jsonValue(in) })

	w := parsec.Spaces()
	// keyword and tok close over w — no need to pass it at each call site
	keyword := func(s string, f func(string) any) parsec.Parser[any] {
		return parsec.Then(w, parsec.Map(parsec.String(s), f))
	}
	tok := func(c rune) parsec.Parser[rune] {
		return parsec.Then(w, parsec.Char(c))
	}

	jnull   := keyword("null",  jsonNull)
	jtrue   := keyword("true",  jsonTrue)
	jfalse  := keyword("false", jsonFalse)
	jnumber := parsec.Then(w, parsec.Map(parsec.Float(), func(f float64) any { return f }))
	jstring := parsec.Then(w, parsec.Map(parsec.JSONString(), jsonString))

	comma := tok(',')
	colon := tok(':')

	// array: '[' ws (value (',' value)*)? ws ']'
	jarray := parsec.Map(
		parsec.Between(tok('['), tok(']'), parsec.SepBy(lazy, comma)),
		jsonArray,
	)

	// object: '{' ws (key ':' value (',' key ':' value)*)? ws '}'
	key  := parsec.Then(w, parsec.JSONString())
	pair := parsec.Bind(key, func(k string) parsec.Parser[[2]any] {
		return parsec.Map(parsec.Then(colon, lazy), func(v any) [2]any { return [2]any{k, v} })
	})
	jobject := parsec.Map(
		parsec.Between(tok('{'), tok('}'), parsec.SepBy(pair, comma)),
		jsonObject,
	)

	jsonValue = parsec.Choice(jnull, jtrue, jfalse, jnumber, jstring, jarray, jobject)
	return jsonValue
}

func TestJSON(t *testing.T) {
	p := newJSONParser()

	tests := []struct {
		input string
		want  any
	}{
		// literals
		{`null`, nil},
		{`true`, true},
		{`false`, false},
		// numbers
		{`0`, float64(0)},
		{`42`, float64(42)},
		{`-7`, float64(-7)},
		{`3.14`, float64(3.14)},
		{`-2.5`, float64(-2.5)},
		// strings
		{`""`, ""},
		{`"hello"`, "hello"},
		{`"say \"hi\""`, `say "hi"`},
		{`"line1\nline2"`, "line1\nline2"},
		// arrays
		{`[]`, []any{}},
		{`[1,2,3]`, []any{float64(1), float64(2), float64(3)}},
		{`[true,null,"x"]`, []any{true, nil, "x"}},
		// objects
		{`{}`, map[string]any{}},
		{`{"key":"value"}`, map[string]any{"key": "value"}},
		{`{"a":1,"b":true}`, map[string]any{"a": float64(1), "b": true}},
		// nested
		{`[1,[2,3]]`, []any{float64(1), []any{float64(2), float64(3)}}},
		{`{"x":{"y":42}}`, map[string]any{"x": map[string]any{"y": float64(42)}}},
		// whitespace
		{` null`, nil},
		{` [ 1 , 2 ] `, []any{float64(1), float64(2)}},
		{` { "k" : "v" } `, map[string]any{"k": "v"}},
	}

	for _, tt := range tests {
		got, err := parsec.Run(p, tt.input)
		if err != nil {
			t.Errorf("JSON(%q): %v", tt.input, err)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("JSON(%q)\n  got  %#v\n  want %#v", tt.input, got, tt.want)
		}
	}
}

func TestJSON_invalid(t *testing.T) {
	p := newJSONParser()

	inputs := []string{
		``,
		`tru`,
		`[1,]`,
		`{"key"}`,
	}

	for _, input := range inputs {
		_, err := parsec.Run(p, input)
		if err == nil {
			t.Errorf("JSON(%q): expected error, got nil", input)
		}
	}
}

func BenchmarkJSON(b *testing.B) {
	p := newJSONParser()
	inputs := []string{
		`null`,
		`42`,
		`"hello\nworld"`,
		`[1, 2, 3]`,
		`{"key": "value", "n": 42}`,
		`{"a": [1, true, null], "b": {"c": 3.14}}`,
	}
	for b.Loop() {
		for _, in := range inputs {
			parsec.Run(p, in)
		}
	}
}
