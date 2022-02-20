// Package bash_ble provides bash-ble completion
package bash_ble

import (
	"fmt"

	"github.com/rsteube/carapace/internal/uid"
	"github.com/spf13/cobra"
)

// Snippet creates the bash-ble completion script
func Snippet(cmd *cobra.Command) string {
	result := fmt.Sprintf(`#!/bin/bash
_%v_completion() {
  bleopt complete_menu_style=desc
  #export COMP_WORDBREAKS

  local compline="${COMP_LINE:0:${COMP_POINT}}"
  local IFS=$'\n'
  local c
  mapfile -t c < <(echo "$compline" | sed -e "s/ \$/ ''/" -e 's/"/\"/g' | xargs %v _carapace bash-ble)
  [[ "${COMPREPLY[*]}" == "" ]] && COMPREPLY=() # fix for mapfile creating a non-empty array from empty command output

  for cand in "${c[@]}"; do
    [ ! -z "$cand" ] && ble/complete/cand/yield mandb "${cand%%$'\t'*}" "${cand##*$'\t'}"
  done
}

complete -F _%v_completion %v
`, cmd.Name(), uid.Executable(), cmd.Name(), cmd.Name())

	return result
}
