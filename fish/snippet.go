package fish

import (
	"fmt"
	"strings"

	"github.com/rsteube/carapace/uid"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var replacer = strings.NewReplacer(
	`:`, `\:`,
	`"`, `\"`,
	`[`, `\[`,
	`]`, `\]`,
	`'`, `\"`,
)

func SnippetFlagCompletion(cmd *cobra.Command, flag *pflag.Flag, action *string) (snippet string) {
	var suffix string
	if action == nil {
		if flag.NoOptDefVal != "" {
			suffix = "" // no argument required for flag
		} else {
			suffix = " -r" // require a value
		}
	} else {
		suffix = fmt.Sprintf(" -a '(%v)' -r", *action)
	}

	if flag.Shorthand == "" { // no shorthannd
		snippet = fmt.Sprintf(`complete -c %v -f -n '_%v_state %v' -l %v -d '%v'%v`, cmd.Root().Name(), cmd.Root().Name(), uid.Command(cmd), flag.Name, replacer.Replace(flag.Usage), suffix)
	} else {
		snippet = fmt.Sprintf(`complete -c %v -f -n '_%v_state %v' -l %v -s %v -d '%v'%v`, cmd.Root().Name(), cmd.Root().Name(), uid.Command(cmd), flag.Name, flag.Shorthand, replacer.Replace(flag.Usage), suffix)
	}
	return
}

func SnippetPositionalCompletion(position int, action string) string {
	return fmt.Sprintf(`"%v:: :%v"`, position, action)
}

func zshCompFlagCouldBeSpecifiedMoreThenOnce(f *pflag.Flag) bool {
	return strings.Contains(f.Value.Type(), "Slice") ||
		strings.Contains(f.Value.Type(), "Array")
}

func SnippetSubcommands(cmd *cobra.Command) string {
	if !cmd.HasSubCommands() {
		return ""
	}
	cmnds := make([]string, 0)
	functions := make([]string, 0)
	for _, c := range cmd.Commands() {
		if !c.Hidden {
			cmnds = append(cmnds, fmt.Sprintf(`        "%v:%v"`, c.Name(), c.Short))
			functions = append(functions, fmt.Sprintf(`    %v)
      %v
      ;;`, c.Name(), uid.Command(c)))

			for _, alias := range c.Aliases {
				cmnds = append(cmnds, fmt.Sprintf(`        "%v:%v"`, alias, c.Short))
				functions = append(functions, fmt.Sprintf(`    %v)
      %v
      ;;`, alias, uid.Command(c)))
			}
		}
	}

	templ := `

  # shellcheck disable=SC2154
  case $state in
    cmnds)
      # shellcheck disable=SC2034
      commands=(
%v
      )
      _describe "command" commands
      ;;
  esac
  
  case "${words[1]}" in
%v
  esac`

	return fmt.Sprintf(templ, strings.Join(cmnds, "\n"), strings.Join(functions, "\n"))
}
