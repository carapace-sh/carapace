package zsh

import (
	"slices"
	"strings"
	"testing"

	"github.com/carapace-sh/carapace/internal/common"
)

func TestActionRawValuesEscapesSpecialCharactersOnce(t *testing.T) {
	output := ActionRawValues("", common.Meta{}, common.RawValuesFrom("Test2&Testaaa.txt", "spaceT TT.x", "foo:bar"))

	want := []string{
		`Test2\&Testaaa.txt `,
		`spaceT\ TT.x `,
		`foo\:bar `,
	}
	if got := valuesFromOutput(t, output); !slices.Equal(got, want) {
		t.Fatalf("values mismatch:\n got: %#v\nwant: %#v", got, want)
	}
}

func valuesFromOutput(t *testing.T, output string) []string {
	t.Helper()

	fields := strings.Split(output, "\001")
	if len(fields) != 4 {
		t.Fatalf("expected four output fields, got %d in %#v", len(fields), output)
	}

	groups := strings.Split(strings.TrimSuffix(fields[2], "\002"), "\002")
	if len(groups) != 1 {
		t.Fatalf("expected one tag group, got %d in %#v", len(groups), fields[2])
	}

	block := strings.Split(groups[0], "\003")
	if len(block) != 3 {
		t.Fatalf("expected three tag group fields, got %d in %#v", len(block), groups[0])
	}

	return strings.Split(block[2], "\n")
}
