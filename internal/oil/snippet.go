package oil

import (
	"fmt"

	"github.com/rsteube/carapace/internal/uid"
	"github.com/spf13/cobra"
)

func Snippet(cmd *cobra.Command, actions map[string]string) string {
	result := fmt.Sprintf(`#!/bin/oil
_%v_completion() {
  local compline="${COMP_LINE:0:${COMP_POINT}}"
  local IFS=$'\n'
  mapfile -t COMPREPLY < <(echo "$compline" | sed -e "s/ \$/ ''/" -e 's/"/\"/g' | xargs %v _carapace oil "_")
  [[ ""${COMPREPLY[@]}"" == "" ]] && setvar COMPREPLY = %%() # fix for mapfile creating a non-empty array from empty command output
  [[ ${COMPREPLY[@]} == *[/=@:.,] ]] && compopt -o nospace
}

complete -F _%v_completion %v
`, cmd.Name(), uid.Executable(), cmd.Name(), cmd.Name())

	return result
}
