package fish

import (
	"fmt"
	"strings"

	"github.com/rsteube/carapace/common"
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

func Snippet(cmd *cobra.Command, actions map[string]string) string {
	result := fmt.Sprintf(`function _%v_quote_suffix
  if not commandline -cp | xargs echo 2>/dev/null >/dev/null
    if commandline -cp | sed 's/$/"/'| xargs echo 2>/dev/null >/dev/null
      echo '"'
    else if commandline -cp | sed "s/\$/'/"| xargs echo 2>/dev/null >/dev/null
      echo "'"
    end
  else 
    echo ""
  end
end

function _%v_state
  set -lx CURRENT (commandline -cp)
  if [ "$LINE" != "$CURRENT" ]
    set -gx LINE (commandline -cp)
    set -gx STATE (commandline -cp | sed "s/\$/"(_%v_quote_suffix)"/" | xargs %v _carapace fish state)
  end

  [ "$STATE" = "$argv" ]
end

function _%v_callback
  set -lx CALLBACK (commandline -cp | sed "s/\$/"(_%v_quote_suffix)"/" | sed "s/ \$/ ''/" | xargs %v _carapace fish $argv )
  eval "$CALLBACK"
end

complete -c %v -f
`, cmd.Name(), cmd.Name(), cmd.Name(), uid.Executable(), cmd.Name(), cmd.Name(), uid.Executable(), cmd.Name())
	result += snippetFunctions(cmd, actions)

	return result
}

func snippetFunctions(cmd *cobra.Command, actions map[string]string) string {
	function_pattern := `
%v
`

	flags := make([]string, 0)
	cmd.LocalFlags().VisitAll(func(flag *pflag.Flag) {
		if !flag.Hidden {
			var s string
			if action, ok := actions[uid.Flag(cmd, flag)]; ok {
				s = snippetFlagCompletion(cmd, flag, &action)
			} else {
				s = snippetFlagCompletion(cmd, flag, nil)
			}

			flags = append(flags, s)
		}
	})

	positionals := make([]string, 0)
	if cmd.HasAvailableSubCommands() {
		positionals = []string{}
		for _, subcmd := range cmd.Commands() {
			if !subcmd.Hidden {
				positionals = append(positionals, fmt.Sprintf(`complete -c '%v' -f -n '_%v_state %v ' -a '%v' -d '%v'`, cmd.Root().Name(), cmd.Root().Name(), uid.Command(cmd), subcmd.Name()+" "+strings.Join(subcmd.Aliases, " "), replacer.Replace(subcmd.Short)))
			}
		}
	} else {
		if len(positionals) == 0 {
			if cmd.ValidArgs != nil {
				//positionals = []string{"    " + snippetPositionalCompletion(1, ActionValues(cmd.ValidArgs...))}
			}
			positionals = append(positionals, fmt.Sprintf(`complete -c '%v' -f -n '_%v_state %v' -a '(_%v_callback _)'`, cmd.Root().Name(), cmd.Root().Name(), uid.Command(cmd), cmd.Root().Name()))
		}
	}

	arguments := append(flags, positionals...)

	result := make([]string, 0)
	result = append(result, fmt.Sprintf(function_pattern, strings.Join(arguments, "\n")))
	for _, subcmd := range cmd.Commands() {
		if !subcmd.Hidden {
			result = append(result, snippetFunctions(subcmd, actions))
		}
	}
	return strings.Join(result, "\n")
}

func snippetFlagCompletion(cmd *cobra.Command, flag *pflag.Flag, action *string) (snippet string) {
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
		snippet = fmt.Sprintf(`complete -c '%v' -f -n '_%v_state %v' -l '%v' -d '%v'%v`, cmd.Root().Name(), cmd.Root().Name(), uid.Command(cmd), flag.Name, replacer.Replace(flag.Usage), suffix)
	} else if common.IsShorthandOnly(flag) {
		snippet = fmt.Sprintf(`complete -c '%v' -f -n '_%v_state %v' -s '%v' -d '%v'%v`, cmd.Root().Name(), cmd.Root().Name(), uid.Command(cmd), flag.Shorthand, replacer.Replace(flag.Usage), suffix)
	} else {
		snippet = fmt.Sprintf(`complete -c '%v' -f -n '_%v_state %v' -l '%v' -s '%v' -d '%v'%v`, cmd.Root().Name(), cmd.Root().Name(), uid.Command(cmd), flag.Name, flag.Shorthand, replacer.Replace(flag.Usage), suffix)
	}
	return
}

func snippetPositionalCompletion(position int, action string) string {
	return fmt.Sprintf(`"%v:: :%v"`, position, action)
}

func zshCompFlagCouldBeSpecifiedMoreThenOnce(f *pflag.Flag) bool {
	return strings.Contains(f.Value.Type(), "Slice") ||
		strings.Contains(f.Value.Type(), "Array")
}

func snippetSubcommands(cmd *cobra.Command) string {
	if !cmd.HasAvailableSubCommands() {
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
