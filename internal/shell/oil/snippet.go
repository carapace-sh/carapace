// Package oil provides Oil completion
package oil

import (
	"fmt"

	"github.com/carapace-sh/carapace/pkg/uid"
	"github.com/spf13/cobra"
)

// Snippet creates the oil completion script.
// TODO: Export COMP_WORDBREAKS (Oil supports it at the shell level, same as bash).
// TODO: The sed expression 's/"/\"/g' is effectively a no-op since the shell already consumed quotes in "$compline".
func Snippet(cmd *cobra.Command) string {
	result := fmt.Sprintf(`#!/bin/osh
_%v_completion() {
  local compline="${COMP_LINE:0:${COMP_POINT}}"
  local IFS=$'\n'
  mapfile -t COMPREPLY < <(echo "$compline" | sed -e "s/ \$/ ''/" -e 's/"/\"/g' | xargs %v _carapace oil)
  [[ "${COMPREPLY[@]}" == "" ]] && COMPREPLY=() # fix for mapfile creating a non-empty array from empty command output
  [[ ${COMPREPLY[0]} == *[/=@:.,$'\001'] ]] && compopt -o nospace
  # TODO use mapfile
  # shellcheck disable=SC2206
  [[ ${#COMPREPLY[@]} -eq 1 ]] && COMPREPLY=(${COMPREPLY[@]%%$'\001'})
}

complete -F _%v_completion %v
`, cmd.Name(), uid.Executable(), cmd.Name(), cmd.Name())

	return result
}
