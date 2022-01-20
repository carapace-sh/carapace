// Package tcsh provides tcsh completion
package tcsh

import (
	"fmt"

	"github.com/rsteube/carapace/internal/uid"
	"github.com/spf13/cobra"
)

func Snippet(cmd *cobra.Command) string {
	// TODO initial version - needs to handle open quotes
	return fmt.Sprintf("complete \"%v\" 'p@*@`echo \"$COMMAND_LINE'\"''\"'\" | xargs %v _carapace tcsh _ `@@' ;", cmd.Name(), uid.Executable())
}
