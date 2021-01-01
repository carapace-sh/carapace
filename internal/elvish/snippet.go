package elvish

import (
	"fmt"
	"github.com/rsteube/carapace/internal/uid"
	"github.com/spf13/cobra"
	"strings"
)

var replacer = strings.NewReplacer( // TODO
	`:`, `\:`,
	"\n", ``,
	`"`, `\"`,
	`[`, `\[`,
	`]`, `\]`,
	`'`, `\"`,
)

func Snippet(cmd *cobra.Command, actions map[string]string) string {
	return fmt.Sprintf(`edit:completion:arg-completer[%v] = [@arg]{
    if (eq $arg[-1] "") {
        arg[-1] = "''"
    }
    eval (%v _carapace elvish _ (all $arg) | slurp) &ns=(ns [&arg=$arg])
}
`, cmd.Name(), uid.Executable())
}
