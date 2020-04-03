package bash

import (
	"fmt"
	"strings"

	"github.com/rsteube/carapace/uid"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func Snippet(cmd *cobra.Command, actions map[string]string) string {
	result := fmt.Sprintf(`#!/bin/bash
_%v_callback() {
  local compline="${COMP_LINE:0:${COMP_POINT}}"
  local last="${COMP_WORDS[${COMP_CWORD}]}"
  if [[ $last =~ ^[\"\'] ]] && ! echo "$last" | xargs echo 2>/dev/null >/dev/null ; then
      compline="${compline}${last:0:1}"
      last="${last/ /\\\\ }" 
  fi

  echo "$compline" | sed -e 's/ $/ _/' -e 's/"/\"/g' | xargs %v _carapace bash "$1"
}

_%v_completions() {
  local compline="${COMP_LINE:0:${COMP_POINT}}"
  local last="${COMP_WORDS[${COMP_CWORD}]}"
  
  if [[ $last =~ ^[\"\'] ]] && ! echo "$last" | xargs echo 2>/dev/null >/dev/null ; then
      compline="${compline}${last:0:1}"
      last="${last/ /\\\\ }" 
  else
      last="${last/ /\\\ }" 
  fi

  local state=$(echo "$compline" | sed -e "s/ \$/ _/" -e 's/"/\"/g' | xargs %v _carapace bash state)
  local previous="${COMP_WORDS[$((${COMP_CWORD}-1))]}"
  local IFS=$'\n'

  case $state in
%v
  esac

  [[ $last =~ ^[\"\'] ]] && COMPREPLY=("${COMPREPLY[@]/\\ /\ }")
  [[ $COMPREPLY == */ ]] && compopt -o nospace
}

complete -F _%v_completions %v
`, cmd.Name(), cmd.Name(), cmd.Name(), cmd.Name(), snippetFunctions(cmd, actions), cmd.Name(), cmd.Name())

	return result
}

func snippetFunctions(cmd *cobra.Command, actions map[string]string) string {
	function_pattern := `
    '%v' )
      if [[ $last == -* ]]; then
        COMPREPLY=($(%v))
      else
        case $previous in
%v
          *)
            COMPREPLY=($(%v))
            ;;
        esac
      fi
      ;;
`

	flags := make([]string, 0)
	cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
		if !f.Hidden {
			var s string
			if action, ok := actions[uid.Flag(cmd, f)]; ok {
				s = snippetFlagCompletion(f, action)
			} else {
				s = snippetFlagCompletion(f, "")
			}
			flags = append(flags, s)
		}
	})

	var positionalAction string
	if cmd.HasSubCommands() {
		subcommands := make([]string, 0)
		for _, c := range cmd.Commands() {
			if !c.Hidden {
				subcommands = append(subcommands, c.Name())
				for _, alias := range c.Aliases {
					subcommands = append(subcommands, alias)
				}
			}
		}
		positionalAction = ActionValues(subcommands...)
	} else {
		positionalAction = Callback(cmd.Root().Name(), "_")
	}

	result := make([]string, 0)
	result = append(result, fmt.Sprintf(function_pattern, uid.Command(cmd), snippetFlagList(cmd.LocalFlags()), strings.Join(flags, "\n"), positionalAction))
	for _, subcmd := range cmd.Commands() {
		if !subcmd.Hidden {
			result = append(result, snippetFunctions(subcmd, actions))
		}
	}
	return strings.Join(result, "\n")
}

func snippetFlagList(flags *pflag.FlagSet) string {
	flagValues := make([]string, 0)

	flags.VisitAll(func(flag *pflag.Flag) {
		if !flag.Hidden {
			flagValues = append(flagValues, "--"+flag.Name)
			if flag.Shorthand != "" {
				flagValues = append(flagValues, "-"+flag.Shorthand)
			}
		}
	})
	if len(flagValues) > 0 {
		return ActionValues(flagValues...)
	} else {
		return ""
	}
}

func snippetFlagCompletion(flag *pflag.Flag, action string) (snippet string) {
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
