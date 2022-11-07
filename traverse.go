package carapace

import (
	"strings"

	"github.com/rsteube/carapace/internal/pflagfork"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type InFlag struct {
	*pflag.Flag
	Args []string
}

func (f InFlag) Consumes(arg string) bool {
	pflagfork.Flag(f.Flag)

	return false // TODO
}

func actionTraverse(c *cobra.Command, args []string) (Action, Context) {
	preInvoke(c, args)

	inArgs := []string{} // args consumed by current command
	var inFlag *InFlag
	fs := pflagfork.FlagSet(c.Flags())

	context := newContext(args)
	for i, arg := range context.Args {
		switch {
		// flag argument
		case inFlag != nil && inFlag.Consumes(arg):
			inArgs = append(inArgs, arg)
			inFlag.Args = append(inFlag.Args, arg)

			if !inFlag.Consumes("") {
				inFlag = nil // no more args expected
			}
			continue

		// flag
		case strings.HasPrefix(arg, "-"):
			inFlag = &InFlag{
				Flag: fs.LookupArg(arg).Flag,
				Args: []string{},
			}
			inArgs = append(inArgs, arg)
			continue

		// subcommand
		case subcommand(c, arg) != nil:
			if err := c.ParseFlags(inArgs); err != nil {
				return ActionMessage(err.Error()), context
			}
			return actionTraverse(subcommand(c, arg), args[i+1:])

		// positional
		default:
			inArgs = append(inArgs, arg)
			inFlag = nil
		}
	}

	if err := c.ParseFlags(inArgs); err != nil { // TODO filter error
		return ActionMessage(err.Error()), context
	}

	switch {
	// flag argument
	case inFlag != nil && inFlag.Consumes(context.CallbackValue):
		return storage.getFlag(c, inFlag.Name), context

	// flag
	case strings.HasPrefix(context.CallbackValue, "-"):
		// TODO handle optargflags with their value
		return actionFlags(c), context

	// positional or subcommand
	default:
		return Batch(
				storage.getPositional(c, len(context.Args)),
				actionSubcommands(c),
			).ToA(),
			context
	}
}

func subcommand(cmd *cobra.Command, arg string) (subcommand *cobra.Command) {
	subcommand, _, _ = cmd.Find([]string{arg})
	return
}

func preInvoke(cmd *cobra.Command, args []string) {
	if subCmd := subcommand(cmd, "_carapace"); subCmd != nil && subCmd.PreRun != nil {
		subCmd.PreRun(cmd, args)
	}
}