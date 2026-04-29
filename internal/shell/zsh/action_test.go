package zsh

import (
	"strings"
	"testing"

	"github.com/carapace-sh/carapace/internal/common"
)

func TestActionRawValuesEscapesSpecialCharactersOnce(t *testing.T) {
	output := ActionRawValues("", common.Meta{}, common.RawValues{
		{Value: "Test2&Testaaa.txt", Display: "Test2&Testaaa.txt"},
		{Value: "spaceT TT.x", Display: "spaceT TT.x"},
	})

	if !strings.Contains(output, `Test2\&Testaaa.txt`) {
		t.Fatalf("expected ampersand to be escaped once, got %q", output)
	}
	if strings.Contains(output, `Test2\\&Testaaa.txt`) {
		t.Fatalf("expected ampersand not to be double-escaped, got %q", output)
	}
	if !strings.Contains(output, `spaceT\ TT.x`) {
		t.Fatalf("expected space to be escaped once, got %q", output)
	}
	if strings.Contains(output, `spaceT\\ TT.x`) {
		t.Fatalf("expected space not to be double-escaped, got %q", output)
	}
}
