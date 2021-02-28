package oil

import (
	"fmt"

	"github.com/rsteube/carapace/internal/uid"
	"github.com/spf13/cobra"
)

func Snippet(cmd *cobra.Command) string {
	result := fmt.Sprintf(`#!/bin/osh
_%v_completion() {
  local compline="${COMP_LINE:0:${COMP_POINT}}"
  local IFS=$'\n'
  mapfile -t COMPREPLY < <(echo "$compline" | sed -e "s/ \$/ ''/" -e 's/"/\"/g' | xargs %v _carapace oil "_")
  [[ "${COMPREPLY[@]}" == "" ]] && COMPREPLY=() # fix for mapfile creating a non-empty array from empty command output
  [[ ${COMPREPLY[0]} == *[/=@:.,] ]] && compopt -o nospace
}

complete -F _%v_completion %v
`, cmd.Name(), uid.Executable(), cmd.Name(), cmd.Name())

	return result
}
