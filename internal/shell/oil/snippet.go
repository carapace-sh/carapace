package oil

import (
	"fmt"
	"strings"

	shlex "github.com/carapace-sh/carapace-shlex/v2"
	"github.com/carapace-sh/carapace/pkg/uid"
	"github.com/spf13/cobra"
)

// Snippet creates the oil completion script.
func Snippet(cmd *cobra.Command) string {
	return SnippetSingle(cmd.Name(), false)
}

// SnippetMulti creates a multi-completer oil completion script.
func SnippetMulti(names []string, defaultName string, snippetFuncs string) string {
	quoted := make([]string, 0, len(names))
	for _, name := range names {
		quoted = append(quoted, shlex.Quote(name, shlex.Bash))
	}
	return fmt.Sprintf(`#!/bin/osh
%[4]v_%[1]v_completer() {
  local command="${COMP_WORDS[0]}"
  local compline="${COMP_LINE:0:${COMP_POINT}}"
  local IFS=$'\n'
  mapfile -t COMPREPLY < <(echo "$compline" | sed -e "s/ \$/ ''/" -e 's/"/\"/g' | xargs %[2]v "${command}" _carapace oil)
  [[ "${COMPREPLY[@]}" == "" ]] && COMPREPLY=() # fix for mapfile creating a non-empty array from empty command output
  [[ ${COMPREPLY[0]} == *[/=@:.,$'\001'] ]] && compopt -o nospace
  # shellcheck disable=SC2206
  [[ ${#COMPREPLY[@]} -eq 1 ]] && COMPREPLY=(${COMPREPLY[@]%%$'\001'})
}

complete -F _%[1]v_completer %[3]v
`, defaultName, shlex.Quote(uid.Executable(), shlex.Bash), strings.Join(quoted, " "), snippetFuncs)
}

// SnippetSingle creates a single-command oil completion script.
// When explicitCommand is true, the command name is included in the invocation
// (for multi-completer subcommands). When false, the executable is invoked
// without an explicit command (standalone mode).
func SnippetSingle(command string, explicitCommand bool) string {

	invocation := fmt.Sprintf("%v _carapace oil", shlex.Quote(uid.Executable(), shlex.Bash))
	if explicitCommand {
		invocation = fmt.Sprintf("%v %v _carapace oil", shlex.Quote(uid.Executable(), shlex.Bash), shlex.Quote(command, shlex.Bash))
	}

	return fmt.Sprintf(`#!/bin/osh
_%v_completion() {
  local compline="${COMP_LINE:0:${COMP_POINT}}"
  local IFS=$'\n'
  mapfile -t COMPREPLY < <(echo "$compline" | sed -e "s/ \$/ ''/" -e 's/"/\"/g' | xargs %v)
  [[ "${COMPREPLY[@]}" == "" ]] && COMPREPLY=() # fix for mapfile creating a non-empty array from empty command output
  [[ ${COMPREPLY[0]} == *[/=@:.,$'\001'] ]] && compopt -o nospace
  # TODO use mapfile
  # shellcheck disable=SC2206
  [[ ${#COMPREPLY[@]} -eq 1 ]] && COMPREPLY=(${COMPREPLY[@]%%$'\001'})
}

complete -F _%v_completion %v
`, command, invocation, shlex.Quote(command, shlex.Bash), command)
}
