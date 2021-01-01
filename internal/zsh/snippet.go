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

func Snippet(cmd *cobra.Command, actions map[string]string) string {
	return fmt.Sprintf(`#compdef %v
function _%v_completion {
  # shellcheck disable=SC2086
  eval "$(%v _carapace zsh _ ${^words//\\ / }'')"

}
compquote '' 2>/dev/null && _%v_completion
compdef _%v_completion %v
`, cmd.Name(), cmd.Name(), uid.Executable(), cmd.Name(), cmd.Name(), cmd.Name())
}
