package common

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// TraverseLenient traverses the command tree but filters errors regarding arguments currently being completed.
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

	if targetCmd.DisableFlagParsing {
		return targetCmd, targetArgs, nil // TODO might need to consider logic below regarding parent command as well
	}

	err = targetCmd.ParseFlags(targetArgs)
	for _, name := range append(targetCmd.Aliases, targetCmd.Name()) {
		if len(args) > 0 &&
			name == args[len(args)-1] &&
			len(targetCmd.Flags().Args()) == 0 &&
			!anyFlagChanged(targetCmd) {
			if targetCmd.HasParent() {
				targetCmd = targetCmd.Parent() // when argument currently being completed is fully matching a subcommand it will be returned, so fix this to parent
			}
		}
	}
	return targetCmd, targetCmd.Flags().Args(), filterError(args, err)
}

func anyFlagChanged(cmd *cobra.Command) (changed bool) {
	cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
		changed = changed || f.Changed
	})
	return
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

	if strings.HasPrefix(current, "-") && msg == fmt.Sprintf("flag needs an argument: '%v' in -%v", current[len(current)-1:], current[len(current)-1:]) { // spf13/pflag
		// ignore shorthand flag currently being completed
		return nil
	}

	if strings.HasPrefix(current, "-") && msg == fmt.Sprintf(`flag needs an argument: "%v" in -%v`, current[len(current)-1:], current[len(current)-1:]) { // rsteube/carapace-pflag: shorthand chain
		// ignore shorthand flag currently being completed
		return nil
	}

	if strings.HasPrefix(current, "-") && msg == fmt.Sprintf(`flag needs an argument: "%v" in %v`, current[1:], current) { // rsteube/carapace-pflag: long shorthand
		// ignore shorthand flag currently being completed
		return nil
	}

	if strings.HasPrefix(current, "--") && msg == fmt.Sprintf("unknown flag: %v", current) {
		// ignore long flag curently being completed
		return nil
	}

	if strings.HasPrefix(current, "-") && msg == fmt.Sprintf("unknown shorthand flag: %v", current) { // rsteube/carapace-pflag: long shorthand
		// ignore non-posix shorthand flag currently being completed
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

		if re := regexp.MustCompile(`invalid argument ".*" for "(?P<shorthand>-., )?(?P<flag>.*)" flag:.*`); re.MatchString(msg) && previous == re.FindStringSubmatch(msg)[2] {
			// ignore invalid argument for flag currently being completed (e.g. empty IntSlice)
			return nil
		}
	}
	return err
}
