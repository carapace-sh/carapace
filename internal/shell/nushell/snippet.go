// Package nushell provides Nushell completion
package nushell

import (
	"fmt"

	"github.com/carapace-sh/carapace/pkg/uid"
	"github.com/spf13/cobra"
)

// Snippet creates the nushell completion script.
func Snippet(cmd *cobra.Command) string {
	return fmt.Sprintf(`let %v_completer = {|spans|
    # if the current command is an alias, get it's expansion
    let expanded_alias = (scope aliases | where name == $spans.0 | $in.0?.expansion?)

    # overwrite
    let spans = (if $expanded_alias != null  {
      # put the first word of the expanded alias first in the span
      $spans | skip 1 | prepend ($expanded_alias | split row " " | take 1)
    } else {
      $spans | skip 1 | prepend ($spans.0)
    })

    %v _carapace nushell ...$spans | from json
}`, cmd.Name(), uid.Executable())
}
