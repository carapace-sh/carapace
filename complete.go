package carapace

import (
	"github.com/rsteube/carapace/internal/config"
	"github.com/rsteube/carapace/pkg/ps"
	"github.com/spf13/cobra"
)

func complete(cmd *cobra.Command, args []string) (string, error) {
	switch len(args) {
	case 0:
		return Gen(cmd).Snippet(ps.DetermineShell())
	case 1:
		return Gen(cmd).Snippet(args[0])
	default:
		initHelpCompletion(cmd)
		action, context := traverse(cmd, args[2:])
		if err := config.Load(); err != nil {
			action = ActionMessage("failed to load config: " + err.Error())
		}
		return action.Invoke(context).value(args[0], args[len(args)-1]), nil
	}
}
