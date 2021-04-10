package common

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// TraverseLenient traverses the command tree but filters errors regarding arguments currently being completed
func TraverseLenient(cmd *cobra.Command, args []string) (*cobra.Command, []string, error) {
	a := args

	// needed so that completion for positional argument that has no value yet to work
	if len(args) > 0 && args[len(args)-1] == "" {
		a = args[0 : len(args)-1]
	}

	targetCmd, targetArgs, err := cmd.Root().Traverse(a)
	if len(args) > 0 && args[len(args)-1] == "" {
		targetArgs = append(targetArgs, "")
	}
	if err != nil {
		return targetCmd, targetArgs, filterError(args, err)
	}

	err = targetCmd.ParseFlags(targetArgs)
	return targetCmd, targetCmd.Flags().Args(), filterError(args, err)
}

func filterError(args []string, err error) error {
	if err == nil || len(args) == 0 {
		return err
	}

	msg := err.Error()

	current := args[len(args)-1]
	if strings.HasPrefix(current, "--") && msg == fmt.Sprintf("flag needs an argument: %v", current) {
		// ignore long flag argument currently being completed
		return nil
	}

	if strings.HasPrefix(current, "-") && msg == fmt.Sprintf("flag needs an argument: '%v' in -%v", current[len(current)-1:], current[len(current)-1:]) {
		// ignore shorthand flag currently being completed
		return nil
	}

	if strings.HasPrefix(current, "--") && msg == fmt.Sprintf("unknown flag: %v", current) {
		// ignore long flag curently being completed
		return nil
	}

	if len(args) > 1 {
		previous := args[len(args)-2]
		if strings.HasPrefix(previous, "--") && msg == fmt.Sprintf("flag needs an argument: %v", previous) {
			// ignore long flag argument currently being completed
			return nil
		}

		if strings.HasPrefix(previous, "-") && msg == fmt.Sprintf("flag needs an argument: '%v' in -%v", previous[len(previous)-1:], previous[len(previous)-1:]) {
			// ignore shorthand flag argument currently being completed
			return nil
		}
	}
	return err
}
