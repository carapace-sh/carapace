package carapace

import (
	"fmt"

	"github.com/spf13/cobra"
)

func registerValidArgsFunction(cmd *cobra.Command) {
	if cmd.ValidArgsFunction == nil {
		cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			action := storage.getPositional(cmd, len(args)).Invoke(Context{Args: args, CallbackValue: toComplete})
			return cobraValuesFor(action), cobraDirectiveFor(action)
		}
	}
}

func registerFlagCompletion(cmd *cobra.Command, actions ActionMap) {
	for name, action := range actions {
		a := action
		cmd.RegisterFlagCompletionFunc(name, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			action := a.Invoke(Context{Args: args, CallbackValue: toComplete})
			return cobraValuesFor(action), cobraDirectiveFor(action)
		})
	}
}

func cobraValuesFor(action InvokedAction) []string {
	result := make([]string, len(action.rawValues))
	for index, r := range action.rawValues {
		if r.Description != "" {
			result[index] = fmt.Sprintf("%v\t%v", r.Value, r.Description)
		} else {
			result[index] = r.Value
		}
	}
	return result
}

func cobraDirectiveFor(action InvokedAction) cobra.ShellCompDirective {
	directive := cobra.ShellCompDirectiveNoFileComp
	if action.nospace {
		directive = directive & cobra.ShellCompDirectiveNoSpace
	}
	if action.skipcache {
		directive = directive & cobra.ShellCompDirectiveError
	}
	return directive
}
