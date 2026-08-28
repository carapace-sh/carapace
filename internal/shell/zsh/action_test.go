package zsh

import (
	"strings"
	"testing"

	"github.com/carapace-sh/carapace/internal/common"
	"github.com/carapace-sh/carapace/internal/env"
)

func TestActionRawValuesEscapesDefaultStateOnce(t *testing.T) {
	t.Setenv(env.CARAPACE_COMPLINE, "example ")

	output := ActionRawValues("", common.Meta{}, common.RawValues{
		{Value: "test/Test2&Testaaa.txt", Display: "test/Test2&Testaaa.txt"},
		{Value: "spaceT TT.x", Display: "spaceT TT.x"},
		{Value: "colon:value", Display: "colon:value"},
		{Value: `back\slash`, Display: `back\slash`},
	})

	values := zshActionValues(t, output)
	assertContains(t, values, `test/Test2\&Testaaa.txt `)
	assertNotContains(t, values, `test/Test2\\&Testaaa.txt `)
	assertContains(t, values, `spaceT\ TT.x `)
	assertNotContains(t, values, `spaceT\\ TT.x `)
	assertContains(t, values, `colon\:value `)
	assertContains(t, values, `back\\slash `)
}

func zshActionValues(t *testing.T, output string) string {
	t.Helper()

	sections := strings.Split(output, "\001")
	if len(sections) < 3 {
		t.Fatalf("expected zsh action output to have at least 3 sections, got %d: %q", len(sections), output)
	}

	blocks := strings.Split(strings.TrimSuffix(sections[2], "\002"), "\002")
	if len(blocks) == 0 {
		t.Fatalf("expected zsh action output to have at least one tag block: %q", output)
	}

	fields := strings.Split(blocks[0], "\003")
	if len(fields) != 3 {
		t.Fatalf("expected zsh action tag block to have 3 fields, got %d: %q", len(fields), blocks[0])
	}
	return fields[2]
}

func assertContains(t *testing.T, s string, substr string) {
	t.Helper()

	if !strings.Contains(s, substr) {
		t.Fatalf("expected %q to contain %q", s, substr)
	}
}

func assertNotContains(t *testing.T, s string, substr string) {
	t.Helper()

	if strings.Contains(s, substr) {
		t.Fatalf("expected %q not to contain %q", s, substr)
	}
}
