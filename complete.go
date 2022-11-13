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

	//
	// Flags and their arguments, embedded or space separated -----------------------------------
	//

	// If we want an argument for the previous word (-flag)
	if yes, flag := needsArgument(cmd, previous); yes && flag != nil {
		// First generate the user provided completions, or default ones.
		flagArgAction := storage.getFlag(cmd, flag.Name).Invoke(ctx).ToA()

		// Add a hint if requested
		if len(flagArgAction.rawValues) == 0 {
			flagArgAction = ActionHint(flag.Usage)
		}

		return current, flagArgAction, ctx
	}

	// Else we assume that the argument is a flag, trying to split for embedded args
	if !cmd.DisableFlagParsing && strings.HasPrefix(current, "-") && !common.IsDash(cmd) {
		return completeFlag(cmd, ctx, current)
	}

	//
	// Positionals and commands, compounded ------------------------------------------------------
	//

	// Else, we deal with a positional word (either arg or command)
	if len(ctx.Args) > 0 {
		// current word being completed is a positional so remove it from context.Args
		ctx.Args = ctx.Args[:len(ctx.Args)-1]
	}

	// If the argument is a DoubleDash that has an impact on how the
	// remaining arguments will be parsed (usually all in one list)
	if common.IsDash(cmd) {
		return current, initPositionalAction(cmd, args[cmd.ArgsLenAtDash():], ctx), ctx
	}

	// Else, we must either consume this arg with one of
	// our positional arg completers, and/or complete subcommands.
	targetAction = initPositionalAction(cmd, args, ctx)

	// We only propose subcommands to complete if we have determined
	// that we don't have required arguments anymore. Also take account
	// of if there are still argument completions to provide next.
	if cmd.HasAvailableSubCommands() && len(args) <= 1 {
		subcommandA := actionSubcommands(cmd).Invoke(ctx)
		targetAction = targetAction.Invoke(ctx).Merge(subcommandA).ToA()
	}

	return current, targetAction, ctx
}

func completeFlag(cmd *cobra.Command, ctx Context, current string) (string, Action, Context) {
	// Always try to dected a flag for completion
	flag := lookupFlag(cmd, current)

	if flag == nil || !strings.Contains(current, "=") {
		return current, actionFlags(cmd), ctx
	}

	// Else, process the string and build the completion context.
	action := storage.getFlag(cmd, flag.Name)
	splitted := strings.SplitN(current, "=", 2)
	ctx.CallbackValue = splitted[1]
	current = strings.Replace(current, "=", opts.OptArgDelimiter, 1)             // revert (potentially) overridden optarg divider for `.value()` invocation below
	action = action.Invoke(ctx).Prefix(splitted[0] + opts.OptArgDelimiter).ToA() // prefix with (potentially) overridden optarg delimiter

	return current, action, ctx
}

// initPositionalAction invokes any non-nil action to be found for this positional word.
func initPositionalAction(cmd *cobra.Command, args []string, ctx Context) (action Action) {
	if len(args) == 0 {
		action = storage.getPositional(cmd, 0)
	} else {
		action = storage.getPositional(cmd, len(args)-1)
	}

	// Add a hint if requested
	action = action.Invoke(ctx).ToA()

	return
}

func findTarget(cmd *cobra.Command, args []string) (*cobra.Command, []string, error) {
	origArg := []string{}
	if len(args) > 2 {
		origArg = args[2:]
	}
	return common.TraverseLenient(cmd, origArg)
}

// lookupFlag must actually deal with an arbitrary number of potential short arguments.
func lookupFlag(cmd *cobra.Command, arg string) (flag *pflag.Flag) {
	nameOrShorthand := strings.TrimLeft(strings.SplitN(arg, "=", 2)[0], "-")

	// If we are treating with a long flag request
	if strings.HasPrefix(arg, "--") {
		return cmd.Flags().Lookup(nameOrShorthand)
	}

	if strings.HasPrefix(arg, "-") && len(nameOrShorthand) > 0 {
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

	// We might have that having an embedded argument
	// if flag.Changed {
	// 	return
	// }
	if len(strings.Split(previous, "=")) >= 2 {
		return false, nil
	}

	if !targetCmd.DisableFlagParsing && flag.NoOptDefVal == "" && !common.IsDash(targetCmd) {
		yes = true
	}

	return
}
