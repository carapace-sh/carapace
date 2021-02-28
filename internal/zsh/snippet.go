package zsh

import (
	"fmt"
	"strings"

	"github.com/rsteube/carapace/internal/uid"
	"github.com/spf13/cobra"
)

var replacer = strings.NewReplacer(
	"\n", ``,
	"`", `\"`,
	`:`, `\:`,
	`"`, `\"`,
	`[`, `\[`,
	`]`, `\]`,
)

func Snippet(cmd *cobra.Command) string {
	return fmt.Sprintf(`#compdef %v
function _%v_completion {
  local IFS=$'\n'
  # shellcheck disable=SC2207,SC2086
  local c=($(%v _carapace zsh _ ${^words//\\ / }''))
  # shellcheck disable=SC2034,2206
  local vals=(${c%%%%$'\t'*})
  # shellcheck disable=SC2034,2206
  local descriptions=(${c##*$'\t'})
  compadd -Q -S '' -d descriptions -a -- vals
}
compquote '' 2>/dev/null && _%v_completion
compdef _%v_completion %v
`, cmd.Name(), cmd.Name(), uid.Executable(), cmd.Name(), cmd.Name(), cmd.Name())
}
