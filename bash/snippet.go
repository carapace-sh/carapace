package bash

import (
	"fmt"

	"github.com/spf13/pflag"
)

func SnippetFlagList(flags *pflag.FlagSet) string {
	flagValues := make([]string, 0)

	flags.VisitAll(func(flag *pflag.Flag) {
		if !flag.Hidden {
			flagValues = append(flagValues, "--"+flag.Name)
			if flag.Shorthand != "" {
				flagValues = append(flagValues, "-"+flag.Shorthand)
			}
		}
	})
	return ActionValues(flagValues...)
}

func SnippetFlagCompletion(flag *pflag.Flag, action string) (snippet string) {
	if flag.NoOptDefVal != "" {
		return ""
	}

	var names string
	if flag.Shorthand != "" {
		names = fmt.Sprintf("-%v | --%v", flag.Shorthand, flag.Name)
	} else {
		names = "--" + flag.Name
	}

	return fmt.Sprintf(`          %v)
            COMPREPLY=($(%v))
            ;;
`, names, action)
}
