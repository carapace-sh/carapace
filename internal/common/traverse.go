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

// filterError is written so that it will backward recursively search all arguments and
// try to parse the parts for which some errors must be analyzed. When the current (last)
// word has not yielded anything related to the error, we go up the previous word.
func filterError(args []string, err error) error {
	if err == nil || len(args) == 0 {
		return err
	}

	msg := err.Error()
	current := args[len(args)-1]

	var lastShort string
	if len(current) > 0 {
		lastShort = current[len(current)-1:]
	}

	if strings.HasPrefix(current, "--") && msg == fmt.Sprintf("flag needs an argument: %v", current) {
		// ignore long flag argument currently being completed
		return nil
	}

	if strings.HasPrefix(current, "-") && msg == fmt.Sprintf("flag needs an argument: '%v' in -%v", lastShort, lastShort) { // spf13/pflag
		// ignore shorthand flag currently being completed
		return nil
	}

	if strings.HasPrefix(current, "-") && msg == fmt.Sprintf(`flag needs an argument: "%v" in -%v`, lastShort, lastShort) { // rsteube/carapace-pflag: shorthand chain
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

	// TODO: Mute errors in flags that have a NoOptDefalue, with option-set OptArgDelimiter
	if strings.HasPrefix(current, "-") && msg == fmt.Sprintf(`unknown shorthand flag: "%s" in -%s`, "=", "=") { // rsteube/carapace-pflag: long shorthand
		// ignore non-posix shorthand flag currently being completed
		return nil
	}

	re := regexp.MustCompile(`invalid argument ".*" for "(?P<shorthand>-., )?(?P<flag>.*)" flag:.*`)
	// short := strings.TrimSuffix(re.FindStringSubmatch(msg)[1], ", ") // The space following the comma is important
	// long := re.FindStringSubmatch(msg)[1]
	if re.MatchString(msg) {
		// Ignore invalid argument for flag currently being completed (e.g. empty IntSlice)
		// Match either the short flag, which we must trim from a potential comma, or the long one
		return nil
	}

	// If we did not catch a single error for the current word,
	// (or any flag relevant in this word for that matter), we
	// go on to the previous word, which might be the origin of
	// the error.
	// WARN: This is recursive so that it should lock the completion
	// state with an error about this invalid flag. This seems to be
	// however useful in terms of notifying the user something is going
	// to go wrong in an almost fully persistent way (message printing).
	if len(args) > 1 {
		return filterError(args[len(args)-1:], err)
	}

	return err
}
