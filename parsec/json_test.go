package parsec_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/ajiyoshi-vg/goparsec/parsec"
)

// newJSONParser builds a JSON value parser using goparsec combinators.
// Returns any: nil | bool | float64 | string | []any | map[string]any
// String values do not support escape sequences.
func newJSONParser() parsec.Parser[any] {
	var jsonValue parsec.Parser[any]

	w := parsec.Spaces()

	jnull  := parsec.Then(w, parsec.Map(parsec.String("null"),  func(string) any { return nil }))
	jtrue  := parsec.Then(w, parsec.Map(parsec.String("true"),  func(string) any { return true }))
	jfalse := parsec.Then(w, parsec.Map(parsec.String("false"), func(string) any { return false }))

	// number: -? digit+ (. digit+)?
	frac := parsec.Option("", parsec.Map(
		parsec.Bind(parsec.Char('.'), func(rune) parsec.Parser[[]rune] { return parsec.Many1(parsec.Digit()) }),
		func(ds []rune) string { return "." + string(ds) },
	))
	jnumber := parsec.Then(w, parsec.Map(
		parsec.Bind(
			parsec.Option(false, parsec.Map(parsec.Char('-'), func(rune) bool { return true })),
			func(neg bool) parsec.Parser[string] {
				return parsec.Bind(
					parsec.Map(parsec.Many1(parsec.Digit()), func(ds []rune) string { return string(ds) }),
					func(i string) parsec.Parser[string] {
						return parsec.Map(frac, func(f string) string {
							s := i + f
							if neg {
								s = "-" + s
							}
							return s
						})
					},
				)
			},
		),
		func(s string) any {
			f, _ := strconv.ParseFloat(s, 64)
			return f
		},
	))

	// string (no escape sequences)
	rawString := parsec.Between(
		parsec.Char('"'),
		parsec.Char('"'),
		parsec.Map(
			parsec.Many(parsec.Satisfy(func(r rune) bool { return r != '"' && r != '\\' })),
			func(rs []rune) string { return string(rs) },
		),
	)
	jstring := parsec.Then(w, parsec.Map(rawString, func(s string) any { return s }))

	comma := parsec.Then(w, parsec.Char(','))
	colon := parsec.Then(w, parsec.Char(':'))

	// array: '[' ws (value (',' value)*)? ws ']'
	jarray := parsec.Map(
		parsec.Between(
			parsec.Then(w, parsec.Char('[')),
			parsec.Then(w, parsec.Char(']')),
			parsec.SepBy(
				parsec.Parser[any](func(in parsec.Input) (any, parsec.Input, error) { return jsonValue(in) }),
				comma,
			),
		),
		func(vs []any) any { return vs },
	)

	// object: '{' ws (key ':' value (',' key ':' value)*)? ws '}'
	key  := parsec.Then(w, rawString)
	pair := parsec.Bind(key, func(k string) parsec.Parser[[2]any] {
		return parsec.Map(
			parsec.Then(colon, parsec.Parser[any](func(in parsec.Input) (any, parsec.Input, error) { return jsonValue(in) })),
			func(v any) [2]any { return [2]any{k, v} },
		)
	})
	jobject := parsec.Map(
		parsec.Between(
			parsec.Then(w, parsec.Char('{')),
			parsec.Then(w, parsec.Char('}')),
			parsec.SepBy(pair, comma),
		),
		func(pairs [][2]any) any {
			m := make(map[string]any, len(pairs))
			for _, p := range pairs {
				m[p[0].(string)] = p[1]
			}
			return m
		},
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
