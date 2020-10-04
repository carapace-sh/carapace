package zsh

import (
	"fmt"
	"strings"

	"github.com/rsteube/carapace/common"
	"github.com/rsteube/carapace/uid"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var replacer = strings.NewReplacer(
	"\n", ``,
	"`", `\"`,
	`:`, `\:`,
	`"`, `\"`,
	`[`, `\[`,
	`]`, `\]`,
)

func Snippet(cmd *cobra.Command, actions map[string]string) string {
	result := fmt.Sprintf("#compdef %v\n", cmd.Name())
	result += fmt.Sprintf(`function _%v_callback {
  # shellcheck disable=SC2086
  eval "$(%v _carapace zsh "$5" ${os_args})"
}
`, cmd.Name(), uid.Executable())
	result += snippetFunctions(cmd, actions)

	result += fmt.Sprintf("if compquote '' 2>/dev/null; then _%v; else compdef _%v %v; fi\n", cmd.Name(), cmd.Name(), cmd.Name()) // check if withing completion function and enable direct sourcing
	return result
}

func snippetFunctions(cmd *cobra.Command, actions map[string]string) string {
	function_pattern := `function %v {
  %v%v  _arguments -C \
%v%v
}
`

	commandsVar := ""
	if cmd.HasAvailableSubCommands() {
		commandsVar = "local -a commands\n"
	}

	inheritedArgs := ""
	if !cmd.HasParent() {
		inheritedArgs = "  # shellcheck disable=SC2206\n  local -a -x os_args=(${words})\n\n"
	}

	flags := make([]string, 0)
	cmd.LocalFlags().VisitAll(func(flag *pflag.Flag) {
		if !flag.Hidden {
			var s string
			if action, ok := actions[uid.Flag(cmd, flag)]; ok {
				s = "    " + snippetFlagCompletion(flag, &action)
			} else {
				s = "    " + snippetFlagCompletion(flag, nil)
			}

			flags = append(flags, s)
		}
	})

	positionals := make([]string, 0)
	if cmd.HasAvailableSubCommands() {
		positionals = []string{`    "1: :->cmnds"`, `    "*::arg:->args"`}
	} else {
		pos := 1
		for {
			if action, ok := actions[uid.Positional(cmd, pos)]; ok {
				positionals = append(positionals, "    "+snippetPositionalCompletion(pos, action))
				pos++
			} else {
				if action, ok := actions[uid.Positional(cmd, 0)]; ok {
					positionals = append(positionals, "    "+snippetPositionalAnyCompletion(action))
				}
				break // TODO only consisten entriess for now
			}
		}
		if len(positionals) == 0 {
			if cmd.ValidArgs != nil {
				positionals = []string{"    " + snippetPositionalCompletion(1, ActionCandidates(common.CandidateFromValues(cmd.ValidArgs...)...))}
			}
			positionals = append(positionals, `    "*::arg:->args"`)
		}
	}

	arguments := append(flags, positionals...)

	result := make([]string, 0)
	result = append(result, fmt.Sprintf(function_pattern, uid.Command(cmd), commandsVar, inheritedArgs, strings.Join(arguments, " \\\n"), snippetSubcommands(cmd)))
	for _, subcmd := range cmd.Commands() {
		if !subcmd.Hidden {
			result = append(result, snippetFunctions(subcmd, actions))
		}
	}

	return strings.Join(result, "\n")
}

func snippetFlagCompletion(flag *pflag.Flag, action *string) (snippet string) {
	var suffix, multimark, multimarkEscaped string
	if action == nil {
		if flag.NoOptDefVal != "" {
			suffix = "::" // no argument required for flag
		} else {
			suffix = ": :" // require a value
		}
	} else {
		if flag.NoOptDefVal != "" {
			suffix = fmt.Sprintf(":: :%v", *action)
		} else {
			suffix = fmt.Sprintf(": :%v", *action)
		}
	}

	if zshCompFlagCouldBeSpecifiedMoreThenOnce(flag) {
		multimark = "*"
		multimarkEscaped = "\\*"
	}

	// TODO cleanup structure so that it can be read easier
	optArgSuffix := ""
	if flag.NoOptDefVal != "" {
		optArgSuffix = "=-"
	}

	if flag.Shorthand == "" { // no shorthannd
		snippet = fmt.Sprintf(`"%v--%v[%v]%v"`, multimark, flag.Name+optArgSuffix, replacer.Replace(flag.Usage), suffix)
	} else if common.IsShorthandOnly(flag) {
		snippet = fmt.Sprintf(`"%v-%v[%v]%v"`, multimark, flag.Shorthand+optArgSuffix, replacer.Replace(flag.Usage), suffix)
	} else {
		snippet = fmt.Sprintf(`"(%v-%v %v--%v)"{%v-%v,%v--%v}"[%v]%v"`, multimark, flag.Shorthand, multimark, flag.Name, multimarkEscaped, flag.Shorthand+optArgSuffix, multimarkEscaped, flag.Name+optArgSuffix, replacer.Replace(flag.Usage), suffix)
	}
	return
}

func snippetPositionalCompletion(position int, action string) string {
	return fmt.Sprintf(`"%v: :%v"`, position, action)
}

func snippetPositionalAnyCompletion(action string) string {
	return fmt.Sprintf(`"*: :%v"`, action)
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
