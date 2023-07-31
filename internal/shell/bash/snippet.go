// Package bash provides bash completion
package bash

import (
	"fmt"

	"github.com/rsteube/carapace/internal/uid"
	"github.com/spf13/cobra"
)

// Snippet creates the bash completion script.
func Snippet(cmd *cobra.Command) string {
	result := fmt.Sprintf(`#!/bin/bash
_%v_completion() {
  export COMP_WORDBREAKS

  local compline="${COMP_LINE:0:${COMP_POINT}}"
  local IFS=$'\n'

  if echo ${compline}"''" | xargs echo 2>/dev/null > /dev/null; then
  	mapfile -t COMPREPLY < <(echo ${compline}"''" | xargs %v _carapace bash )
  elif echo ${compline} | sed "s/\$/'/" | xargs echo 2>/dev/null > /dev/null; then
  	mapfile -t COMPREPLY < <(echo ${compline} | sed "s/\$/'/" | xargs %v _carapace bash)
  else
  	mapfile -t COMPREPLY < <(echo ${compline} | sed 's/$/"/' | xargs %v _carapace bash)
  fi
		
  [[ "${COMPREPLY[*]}" == "" ]] && COMPREPLY=() # fix for mapfile creating a non-empty array from empty command output
}

complete -o nospace -o noquote -F _%v_completion %v
`, cmd.Name(), uid.Executable(), uid.Executable(), uid.Executable(), cmd.Name(), cmd.Name())

	return result
}
