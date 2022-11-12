package carapace

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rsteube/carapace/internal/common"
	"github.com/rsteube/carapace/internal/config"
	"github.com/rsteube/carapace/internal/pflagfork"
	"github.com/rsteube/carapace/pkg/ps"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func complete(cmd *cobra.Command, args []string) (string, error) {
	if len(args) == 0 {
		if s, err := Gen(cmd).Snippet(ps.DetermineShell()); err != nil {
			fmt.Fprintln(io.MultiWriter(cmd.OutOrStderr(), logger.Writer()), err.Error())
		} else {
			fmt.Fprintln(io.MultiWriter(cmd.OutOrStdout(), logger.Writer()), s)
		}
	} else {
		if len(args) == 1 {
			s, err := Gen(cmd).Snippet(args[0])
			if err != nil {
				return "", err
			}
			return s, nil
		} else {
			shell := args[0]
			current := args[len(args)-1]
			previous := args[len(args)-2]

			if err := config.Load(); err != nil {
				return ActionMessage("failed to load config: "+err.Error()).Invoke(Context{CallbackValue: current}).value(shell, current), nil
			}

			targetCmd, targetArgs, err := findTarget(cmd, args)
			if err != nil {
				return ActionMessage(err.Error()).Invoke(Context{CallbackValue: current}).value(shell, current), nil
			}

			wd, err := os.Getwd()
			if err != nil {
				return ActionMessage(err.Error()).Invoke(Context{CallbackValue: current}).value(shell, current), nil
			}
			context := Context{
				CallbackValue: current,
				Args:          targetArgs,
				Env:           os.Environ(),
				Dir:           wd,
			}
			if err != nil {
				return ActionMessage(err.Error()).Invoke(context).value(shell, current), nil
			}

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
					targetAction = findAction(targetCmd, targetArgs[targetCmd.ArgsLenAtDash():])
				} else {
					targetAction = findAction(targetCmd, targetArgs)
					if targetCmd.HasAvailableSubCommands() && len(targetArgs) <= 1 {
						subcommandA := actionSubcommands(targetCmd).Invoke(context)
						targetAction = targetAction.Invoke(context).Merge(subcommandA).ToA()
					}
				}
			}
			return targetAction.Invoke(context).value(shell, current), nil
		}
	}
	return "", nil // TODO
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
		if pflagfork.IsPosix(cmd.Flags()) {
			flag = cmd.Flags().ShorthandLookup(string(nameOrShorthand[len(nameOrShorthand)-1]))
		} else {
			flag = cmd.Flags().ShorthandLookup(nameOrShorthand)
		}
	}
	return
}
