package common

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestIsDash(t *testing.T) {
	cmd := &cobra.Command{}

	if IsDash(cmd) {
		t.Error("IsDash should be false")
	}

	_ = cmd.Flags().Parse([]string{"--", ""})

	if !IsDash(cmd) {
		t.Error("IsDash should be true")
	}
}
