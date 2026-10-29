package tcsh

import (
	"fmt"
	"strings"

	shlex "github.com/carapace-sh/carapace-shlex/v2"
	"github.com/carapace-sh/carapace/pkg/uid"
	"github.com/spf13/cobra"
)

// Snippet creates the tcsh completion script.
func Snippet(cmd *cobra.Command) string {
	return SnippetSingle(cmd.Name(), false)
}

// SnippetMulti creates a multi-completer tcsh completion script.
func SnippetMulti(names []string, defaultName string, snippetFuncs string) string {
	lines := make([]string, 0, len(names))
	for _, name := range names {
		lines = append(lines, fmt.Sprintf("complete %v 'p@*@`echo \"$COMMAND_LINE'\"''\"'\" | xargs %v %v _carapace tcsh `@@' ;", shlex.Quote(name, shlex.Tcsh), shlex.Quote(uid.Executable(), shlex.Tcsh), shlex.Quote(name, shlex.Tcsh)))
	}
	return strings.Join(lines, "\n")
}

// SnippetSingle creates a single-command tcsh completion script.
// When explicitCommand is true, the command name is included in the invocation
// (for multi-completer subcommands). When false, the executable is invoked
// without an explicit command (standalone mode).
func SnippetSingle(command string, explicitCommand bool) string {
	// TODO initial version - needs to handle open quotes

	invocation := fmt.Sprintf("%v _carapace tcsh", shlex.Quote(uid.Executable(), shlex.Tcsh))
	if explicitCommand {
		invocation = fmt.Sprintf("%v %v _carapace tcsh", shlex.Quote(uid.Executable(), shlex.Tcsh), shlex.Quote(command, shlex.Tcsh))
	}

	return fmt.Sprintf("complete %v 'p@*@`echo \"$COMMAND_LINE'\"''\"'\" | xargs %v `@@' ;", shlex.Quote(command, shlex.Tcsh), invocation)
}
