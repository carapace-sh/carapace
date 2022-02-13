// Package elvish provides elvish completion
package elvish

import (
	"fmt"
	"github.com/rsteube/carapace/internal/uid"
	"github.com/spf13/cobra"
)

// Snippet creates the elvish completion script
func Snippet(cmd *cobra.Command) string {
	return fmt.Sprintf(`set edit:completion:arg-completer[%v] = {|@arg|
    %v _carapace elvish (all $arg) | from-json | all (one) | each {|c| edit:complex-candidate $c[Value] &display=$c[Display] &code-suffix=$c[CodeSuffix] }
}
`, cmd.Name(), uid.Executable())
}
