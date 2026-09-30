package nushell

import (
	"testing"
)

func TestPatch(t *testing.T) {
	for _, tc := range []struct {
		arg      string
		expected string
	}{
		{`"double quoted"`, `double quoted`},
		{`'single quoted'`, `single quoted`},
		{"\"esc\\t\"", "esc\t"}, // nushell double quotes support C-style escapes
		{"\"esc\\\"quote\"", `esc"quote`},
		{`'literal \' escape'`, `literal \`}, // single quotes end at the next quote
		{"`backtick`", `backtick`},
		{"plain", `plain`},
		{``, ``},
	} {
		got := Patch([]string{tc.arg})[0]
		if got != tc.expected {
			t.Errorf("%q: got %q, want %q", tc.arg, got, tc.expected)
		}
	}
}
