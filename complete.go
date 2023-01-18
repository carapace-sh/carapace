package carapace

import (
	"strings"

	"github.com/rsteube/carapace/internal/pflagfork"
	"github.com/rsteube/carapace/pkg/ps"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func complete(cmd *cobra.Command, args []string) (string, error) {
	switch len(args) {
	case 0:
		return Gen(cmd).Snippet(ps.DetermineShell())
	case 1:
		return Gen(cmd).Snippet(args[0])
	default:
		action, context := traverse(cmd, args[2:])
		return action.Invoke(context).value(args[0], args[len(args)-1]), nil
	}
}

func lookupFlag(cmd *cobra.Command, arg string) (flag *pflag.Flag) {
	nameOrShorthand := strings.TrimLeft(strings.SplitN(arg, "=", 2)[0], "-")

	if strings.HasPrefix(arg, "--") {
		flag = cmd.Flags().Lookup(nameOrShorthand)
	} else if strings.HasPrefix(arg, "-") && len(nameOrShorthand) > 0 {
		if (pflagfork.FlagSet{FlagSet: cmd.Flags()}).IsPosix() {
			flag = cmd.Flags().ShorthandLookup(string(nameOrShorthand[len(nameOrShorthand)-1]))
		} else {
			flag = cmd.Flags().ShorthandLookup(nameOrShorthand)
		}
	}
	return
}
