// Package nushell provides Nushell completion
package nushell

import (
	"fmt"

	"github.com/rsteube/carapace/internal/uid"
	"github.com/spf13/cobra"
)

// Snippet creates the nushell completion script
func Snippet(cmd *cobra.Command) string {
	return fmt.Sprintf(`module carapace_%v {
  def "nu-complete %v" [line: string, pos: int] {
    let words = ($line | str substring [0 $pos] | split row " ")
    if ($line | str substring [0 $pos] | str ends-with " ") {
      %v _carapace nushell ($words | append "") | from json
    } else {
      %v _carapace nushell $words | from json
    }
  }
  
  export extern "%v" [
    ...args: string@"nu-complete %v"
  ]
}
use carapace_%v *
`, cmd.Name(), cmd.Name(), uid.Executable(), uid.Executable(), cmd.Name(), cmd.Name(), cmd.Name())
}
