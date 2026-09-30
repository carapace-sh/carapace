package zsh

import (
	"strings"
	"testing"

	"github.com/carapace-sh/carapace/internal/common"
	"github.com/carapace-sh/carapace/internal/env"
)

func TestActionRawValuesQuoting(t *testing.T) {
	for _, tc := range []struct {
		compline string
		expected string
	}{
		{`example "ac"`, `ac`},              // FULL_QUOTING_ESCAPING: no escape, no appended quote
		{`example "ac`, `ac" `},             // QUOTING_ESCAPING: closes quote
		{`example "embeddedP2 w"i`, `ac" `}, // QUOTING_ESCAPING despite closed quote within word
		{`example 'ac'`, `ac`},              // FULL_QUOTING
		{`example 'ac`, `ac' `},             // QUOTING
		{`example 'embeddedP2 w'i`, `ac' `},
		{`example ac`, `ac `}, // DEFAULT: space suffix
	} {
		t.Setenv(env.CARAPACE_COMPLINE, tc.compline)
		out := ActionRawValues("", common.Meta{}, common.RawValues{{Value: "ac", Display: "ac"}})
		val := strings.Split(strings.Split(out, "\002")[0], "\003")[2]
		if val != tc.expected {
			t.Errorf("%q: got %q, want %q", tc.compline, val, tc.expected)
		}
	}
}
