package carapace

import (
	"strings"

	"github.com/rsteube/carapace/internal/common"
	"github.com/rsteube/carapace/internal/config"
	"github.com/rsteube/carapace/internal/pflagfork"
	"github.com/rsteube/carapace/pkg/ps"
	"github.com/rsteube/carapace/pkg/style"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Complete can be used by Go programs wishing to produce completions for
// themselves, without passing through shell snippets/output or export formats.
//
// The `onFinalize` function parameter, if non nil, will be called after having
// generated the completions from the given command/tree. This function is generally
// used to reset the command tree, which is needed when the Go program is a shell itself.
// Also, and before calling `onFinalize` if not nil, the completion storage is cleared.
func Complete(cmd *cobra.Command, args []string, onFinalize func()) (common.RawValues, common.Meta) {
	// Generate the completion as normally done for an external system shell
	action, current := generate(cmd, args)

	// And adapt/fetch the results from invoked action
	return internalValues(action, current, onFinalize)
}

func complete(cmd *cobra.Command, args []string) (string, error) {
	switch len(args) {
	case 0:
		return Gen(cmd).Snippet(ps.DetermineShell())
	case 1:
		return Gen(cmd).Snippet(args[0])
	}

	shell := args[0]

	// Generate the completions and get the resulting invoked action.
	action, current := generate(cmd, args)

	// And generate the completions string for the target shell caller.
	return action.value(shell, current), nil
}

func generate(cmd *cobra.Command, args []string) (InvokedAction, string) {
	current := args[len(args)-1]
	previous := ""
	if len(args) > 1 {
		previous = args[len(args)-2]
	}

	if err := config.Load(); err != nil {
		return ActionMessage("failed to load config: " + err.Error()).Invoke(Context{CallbackValue: current}), current
	}

	targetCmd, targetArgs, err := findTarget(cmd, args)
	if err != nil {
		return ActionMessage(err.Error()).Invoke(Context{CallbackValue: current}), current
	}

	context := NewContext(append(targetArgs, current))

	// TODO needs more cleanup and tests
	var targetAction Action
	if flag := lookupFlag(targetCmd, previous); !targetCmd.DisableFlagParsing && flag != nil && flag.NoOptDefVal == "" && !common.IsDash(targetCmd) { // previous arg is a flag and needs a value
		targetAction = storage.getFlag(targetCmd, flag.Name)
	} else if !targetCmd.DisableFlagParsing && strings.HasPrefix(current, "-") && !common.IsDash(targetCmd) { // assume flag
		if strings.Contains(current, "=") { // complete value for optarg flag
			if flag := lookupFlag(targetCmd, current); flag != nil && flag.NoOptDefVal != "" {
				a := storage.getFlag(targetCmd, flag.Name)
				splitted := strings.SplitN(current, "=", 2)
				context.CallbackValue = splitted[1]
				current = strings.Replace(current, "=", opts.OptArgDelimiter, 1)                  // revert (potentially) overridden optarg divider for `.value()` invocation below
				targetAction = a.Invoke(context).Prefix(splitted[0] + opts.OptArgDelimiter).ToA() // prefix with (potentially) overridden optarg delimiter
			}
		} else { // complete flagnames
			targetAction = actionFlags(targetCmd)
		}
	} else {
		if len(context.Args) > 0 {
			context.Args = context.Args[:len(context.Args)-1] // current word being completed is a positional so remove it from context.Args
		}

		if common.IsDash(targetCmd) {
			dashArgs := targetArgs[targetCmd.ArgsLenAtDash() : len(targetArgs)-1]
			context.Args = dashArgs
			targetAction = findAction(targetCmd, dashArgs)
		} else {
			targetAction = findAction(targetCmd, targetArgs)
			if targetCmd.HasAvailableSubCommands() && len(targetArgs) <= 1 {
				subcommandA := actionSubcommands(targetCmd).Invoke(context)
				targetAction = targetAction.Invoke(context).Merge(subcommandA).ToA()
			}
		}
	}

	return targetAction.Invoke(context), current
}

func findAction(targetCmd *cobra.Command, targetArgs []string) Action {
	// TODO handle Action not found
	if len(targetArgs) == 0 {
		return storage.getPositional(targetCmd, 0)
	}
	return storage.getPositional(targetCmd, len(targetArgs)-1)
}

func findTarget(cmd *cobra.Command, args []string) (*cobra.Command, []string, error) {
	origArg := []string{}
	if len(args) > 2 {
		origArg = args[2:]
	}
	return common.TraverseLenient(cmd, origArg)
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

func internalValues(a InvokedAction, current string, onFinalize func()) (common.RawValues, common.Meta) {
	unsorted := common.RawValues(a.rawValues)
	sorted := make(common.RawValues, 0)

	// Ensure values are sorted.
	unsorted.EachTag(func(tag string, values common.RawValues) {
		vals := make(common.RawValues, len(values))
		for index, val := range values {
			if !a.meta.Nospace.Matches(val.Value) {
				val.Value += " "
			}
			if val.Style != "" {
				val.Style = style.SGR(val.Style)
			}

			vals[index] = val
		}
		sorted = append(sorted, vals...)
	})

	// Merge/filter completions and meta stuff.
	filtered := sorted.FilterPrefix(current)
	filtered = a.meta.Messages.Integrate(filtered, current)

	// Reset the storage (empty all commands) and run the finalize function, which is
	// generally in charge of binding new command instances, with blank flags.
	if onFinalize != nil {
		storage = make(_storage)
		onFinalize()
	}

	return filtered, a.meta
}
