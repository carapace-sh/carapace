// Package spec provides spec file generation for use with carapace-bin
package spec

import (
	"fmt"
	"strings"

	"github.com/rsteube/carapace/internal/uid"
	"github.com/spf13/cobra"
)

// Snippet generates the spec file.
func Snippet(cmd *cobra.Command) string {
	replacer := strings.NewReplacer( // TODO might need more replacements
		`"`, `\"`,
		`'`, `\'`,
		`[`, `\[`,
		`]`, `\]`,
	)
	return fmt.Sprintf(`name: %v
description: %v
completion:
  positionalany: ["$_bridge.Carapace(%v)"]
`, uid.Executable(), replacer.Replace(cmd.Short), uid.Executable())
}
