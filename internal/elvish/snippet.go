package elvish

import (
	"fmt"
	"github.com/rsteube/carapace/internal/uid"
	"github.com/spf13/cobra"
)

// Snippet creates the elvish completion script
func Snippet(cmd *cobra.Command) string {
	return fmt.Sprintf(`edit:completion:arg-completer[%v] = [@arg]{
    %v _carapace elvish _ (all $arg) | from-json | all (one) | each [c]{ edit:complex-candidate $c[Value] &display=$c[Display] &code-suffix=$c[CodeSuffix] }
}
`, cmd.Name(), uid.Executable())
}
