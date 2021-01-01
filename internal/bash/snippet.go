package bash

import (
	"fmt"

	"github.com/rsteube/carapace/internal/uid"
	"github.com/spf13/cobra"
)

func Snippet(cmd *cobra.Command, actions map[string]string) string {
	result := fmt.Sprintf(`#!/bin/bash
_%v_completion() {
  export COMP_WORDBREAKS

  local compline="${COMP_LINE:0:${COMP_POINT}}"
  local IFS=$'\n'
  mapfile -t COMPREPLY < <(echo "$compline" | sed -e "s/ \$/ ''/" -e 's/"/\"/g' | xargs %v _carapace bash "_")
  [[ "${COMPREPLY[*]}" == "" ]] && COMPREPLY=() # fix for mapfile creating a non-empty array from empty command output

  [[ ${COMPREPLY[0]} == *[/=@:.,$'\001'] ]] && compopt -o nospace
  # TODO use mapfile
  # shellcheck disable=SC2206
  [[ ${#COMPREPLY[@]} -eq 1 ]] && COMPREPLY=(${COMPREPLY%%$'\001'})
}

complete -F _%v_completion %v
`, cmd.Name(), uid.Executable(), cmd.Name(), cmd.Name())

	return result
}
