package carapace

import (
	"fmt"
	"io"
	"strings"

	"github.com/rsteube/carapace/internal/common"
	"github.com/rsteube/carapace/internal/config"
	"github.com/rsteube/carapace/pkg/ps"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func complete(cmd *cobra.Command, args []string) (string, error) {
	// Directly skip to snippet generation if
	// we don't have any arguments to deal with.
	if len(args) == 0 {
		outputSnippet(cmd)

		return "", nil
	}

	// If we have only one argument, its a shell indication/command
	// for which we must output a completion snippet (to be sourced)
	if len(args) == 1 {
		s, err := Gen(cmd).Snippet(args[0])
		if err != nil {
			return "", err
		}
		return s, nil
	}

	shell, current, previous := prepareWords(args)

	if configErr := config.Load(); configErr != nil {
		return ActionMessage("failed to load config: "+configErr.Error()).
				Invoke(Context{CallbackValue: current}).
				value(shell, current),
			nil
	}

	// Else we are completing the command we've been compiled into.
	// Any error arising here is returned after creating a context.
	targetCmd, targetArgs, parseErr := findTarget(cmd, args)

	// Create a new context: it is never nil, even if an error occurs
	// and is returned from this instantation/setup.
	ctx, err := newContext(current, targetArgs)
	if err != nil {
		return ActionMessage(err.Error()).Invoke(ctx).value(shell, current), nil
	}

	var action Action

	// unrecoverable command/arg/flag parsing errors are returned immediately.
	if parseErr != nil {
		return ActionMessage(parseErr.Error()).Invoke(ctx).value(shell, current), nil
	}

	// Or we have at least a positional, a flag or a subcommand word to handle.
	// We simply want to settle on a "root" completion Action, since it can host
	// an arbitrary tree of child completers/directives.
	current, action, ctx = getTargetAction(current, previous, ctx, targetCmd, targetArgs)

	// Invoke the completion actions with a context updated/set for the target command.
	return action.Invoke(ctx).value(shell, current), nil
}

func getTargetAction(current, previous string, ctx Context, cmd *cobra.Command, args []string) (string, Action, Context) {
	var targetAction Action

	// If we want an argument for the previous word (-flag)
	if yes, flag := needsArgument(cmd, previous); yes {
		return current, storage.getFlag(cmd, flag.Name), ctx
	}

	// Else we assume that the argument is a flag, trying to split for embedded args
	if !cmd.DisableFlagParsing && strings.HasPrefix(current, "-") && !common.IsDash(cmd) {
		return completeOption(cmd, ctx, current, targetAction)
	}

	// Else, we deal with a positional word (either arg or command)
	if len(ctx.Args) > 0 {
		// current word being completed is a positional so remove it from context.Args
		ctx.Args = ctx.Args[:len(ctx.Args)-1]
	}

	// If the argument is a DoubleDash that has an impact on how the
	// remaining arguments will be parsed (usually all in one list)
	if common.IsDash(cmd) {
		return current, findAction(cmd, args[cmd.ArgsLenAtDash():]), ctx
	}

	// Else, we must either consume this arg with one of
	// our positional arg completers, and/or complete subcommands.
	targetAction = findAction(cmd, args)

	// We only propose subcommands to complete if we have determined
	// that we don't have required arguments anymore. Also take account
	// of if there are still argument completions to provide next.
	if cmd.HasAvailableSubCommands() && len(args) <= 1 {
		subcommandA := actionSubcommands(cmd).Invoke(ctx)
		targetAction = targetAction.Invoke(ctx).Merge(subcommandA).ToA()
	}

	return current, targetAction, ctx
}

func completeOption(cmd *cobra.Command, ctx Context, current string, targetAction Action) (string, Action, Context) {
	// If there is no embedded argument, just parse the next argument word
	if !strings.Contains(current, "=") {
		return current, actionFlags(cmd), ctx
	}

	// If we have found such a split in the string, proceed.
	// First lookup the flag in the last word.
	// We must find a flag, or we return an empty Action.
	flag := lookupFlag(cmd, current)
	if flag == nil || flag.NoOptDefVal == "" {
		return current, targetAction, ctx
	}

	// Else, process the string and build the completion context.
	action := storage.getFlag(cmd, flag.Name)
	splitted := strings.SplitN(current, "=", 2)
	ctx.CallbackValue = splitted[1]
	current = strings.Replace(current, "=", opts.OptArgDelimiter, 1)                   // revert (potentially) overridden optarg divider for `.value()` invocation below
	targetAction = action.Invoke(ctx).Prefix(splitted[0] + opts.OptArgDelimiter).ToA() // prefix with (potentially) overridden optarg delimiter

	return current, targetAction, ctx
}

func findAction(targetCmd *cobra.Command, targetArgs []string) Action {
	// TODO handle Action not found
	if len(targetArgs) == 0 {
		return storage.getPositional(targetCmd, 0)
	}

	// lastArg := targetArgs[len(targetArgs)-1]
	// if strings.HasSuffix(lastArg, " ") { // TODO is this still correct/needed?
	// 	return storage.getPositional(targetCmd, len(targetArgs))
	// }
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
		flag = cmd.Flags().ShorthandLookup(string(nameOrShorthand[len(nameOrShorthand)-1]))
	}
	return
}

func outputSnippet(cmd *cobra.Command) {
	snip, err := Gen(cmd).Snippet(ps.DetermineShell())
	if err != nil {
		fmt.Fprintln(io.MultiWriter(cmd.OutOrStderr(),
			logger.Writer()),
			err.Error())
	} else {
		fmt.Fprintln(io.MultiWriter(cmd.OutOrStdout(),
			logger.Writer()),
			snip)
	}
}

func prepareWords(args []string) (shell, current, previous string) {
	shell = args[0]
	current = args[len(args)-1]
	previous = args[len(args)-2]

	return
}

func needsArgument(targetCmd *cobra.Command, previous string) (yes bool, flag *pflag.Flag) {
	// We need the flag, non-nil
	flag = lookupFlag(targetCmd, previous)
	if flag == nil {
		return
	}

	if !targetCmd.DisableFlagParsing && flag.NoOptDefVal == "" && !common.IsDash(targetCmd) {
		yes = true
	}

	return
}
