// Package nushell provides Nushell completion
package nushell

import (
	"fmt"

	"github.com/carapace-sh/carapace/pkg/uid"
	"github.com/spf13/cobra"
)

// Snippet creates the nushell completion script.
func Snippet(cmd *cobra.Command) string {
	return fmt.Sprintf(`let %[1]v_completer = {|spans|
    %[2]v _carapace nushell ...$spans | from json
}

let %[1]v_previous_completer = ($env.config.completions.external.completer? | default {|spans| null })
$env.config = ($env.config | upsert completions.external {
    enable: true
    completer: {|spans|
        if $spans.0 == %[1]q {
            do $%[1]v_completer $spans
        } else {
            do $%[1]v_previous_completer $spans
        }
    }
})`, cmd.Name(), uid.Executable())
}
