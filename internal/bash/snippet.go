package bash

import (
	"fmt"
	"strings"

	"github.com/rsteube/carapace/internal/common"
	"github.com/rsteube/carapace/internal/uid"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func snippetLazy(cmd *cobra.Command) string {
	return fmt.Sprintf(`#!/bin/bash
_%v_completions() {
   source <(%v _carapace bash) 
   _%v_completions
}
complete -F _%v_completions %v
`, cmd.Name(), uid.Executable(), cmd.Name(), cmd.Name(), cmd.Name())
}

func Snippet(cmd *cobra.Command, actions map[string]string, lazy bool) string {
	if lazy {
		return snippetLazy(cmd)
	}

	result := fmt.Sprintf(`#!/bin/bash
_%v_callback() {
  local compline="${COMP_LINE:0:${COMP_POINT}}"
  # TODO
  #if [[ $last =~ ^[\"\'] ]] && ! echo "$last" | xargs echo 2>/dev/null >/dev/null ; then
  #    compline="${compline}${last:0:1}"
  #    last="${last// /\\\\ }" 
  #fi

  echo "$compline" | sed -e "s/ $/ ''/" -e 's/"/\"/g' | xargs %v _carapace bash "$1"
}

_%v_completions() {
  local cur prev #words cword split
  _init_completion -n /=:.,
  local curprefix
  curprefix="$(echo "$cur" | sed -r 's_^(.*[:=])?.*_\1_')"
  local compline="${COMP_LINE:0:${COMP_POINT}}"
 
  # TODO
  #if [[ $last =~ ^[\"\'] ]] && ! echo "$last" | xargs echo 2>/dev/null >/dev/null ; then
  #    compline="${compline}${last:0:1}"
  #    last="${last// /\\\\ }" 
  #else
  #    last="${last// /\\\ }" 
  #fi

  local state
  state="$(echo "$compline" | sed -e "s/ \$/ ''/" -e 's/"/\"/g' | xargs %v _carapace bash state)"

  local IFS=$'\n'

  case $state in
%v
  esac

  [[ $cur =~ ^[\"\'] ]] && COMPREPLY=("${COMPREPLY[@]//\\ /\ }")
  
  [[ ${#COMPREPLY[@]} -gt 1 ]] && for entry in "${COMPREPLY[@]}"; do
    value="${entry%%	*}"
    display="${entry#*	}"
    if [[ "${value::1}" != "${display::1}"  ]]; then # inserted value differs from display value
       [[ "$(printf  "%%c\n" "${COMPREPLY[@]#*	}" | uniq | wc -l)" -eq 1 ]] && COMPREPLY=("${COMPREPLY[@]}" "") # prevent insertion if all have same first character (workaround for #164)
      break
    fi
  done

  [[ ${#COMPREPLY[@]} -gt 1 ]] && COMPREPLY=("${COMPREPLY[@]#*	}") # show visual part (all after tab)
  [[ ${#COMPREPLY[@]} -eq 1 ]] && COMPREPLY=( ${COMPREPLY[0]%%	*} ) # show value to insert (all before tab) https://stackoverflow.com/a/10130007
  [[ ${COMPREPLY[0]} == *[/=@:.,] ]] && compopt -o nospace
}

complete -F _%v_completions %v
`, cmd.Name(), uid.Executable(), cmd.Name(), uid.Executable(), snippetFunctions(cmd, actions), cmd.Name(), cmd.Name())

	return result
}

func snippetFunctions(cmd *cobra.Command, actions map[string]string) string {
	function_pattern := `
    '%v' )
      if [[ $cur == -* ]]; then
        case $cur in
%v
          *)
            COMPREPLY=($(%v))
            ;;
        esac
      else
        case $prev in
%v
          *)
            COMPREPLY=($(%v))
            ;;
        esac
      fi
      ;;
`

	optArgflags := make([]string, 0)
	cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
		if !f.Hidden {
			var s string
			if action, ok := actions[uid.Flag(cmd, f)]; ok {
				s = snippetFlagCompletion(f, action, true)
			} else {
				s = snippetFlagCompletion(f, "", true)
			}
			if s != "" {
				optArgflags = append(optArgflags, s)
			}
		}
	})

	flags := make([]string, 0)
	cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
		if !f.Hidden {
			var s string
			if action, ok := actions[uid.Flag(cmd, f)]; ok {
				s = snippetFlagCompletion(f, action, false)
			} else {
				s = snippetFlagCompletion(f, "", false)
			}
			if s != "" {
				flags = append(flags, s)
			}
		}
	})

	var positionalAction string
	if cmd.HasAvailableSubCommands() {
		subcommands := make([]common.Candidate, 0)
		for _, c := range cmd.Commands() {
			if !c.Hidden {
				subcommands = append(subcommands, common.Candidate{Value: c.Name(), Display: c.Name(), Description: c.Short})
				for _, alias := range c.Aliases {
					subcommands = append(subcommands, common.Candidate{Value: alias, Display: c.Name(), Description: c.Short})
				}
			}
		}
		positionalAction = ActionCandidates(subcommands...)
	} else {
		positionalAction = Callback(cmd.Root().Name(), "_")
	}

	result := make([]string, 0)
	result = append(result, fmt.Sprintf(function_pattern, uid.Command(cmd), strings.Join(optArgflags, "\n"), snippetFlagList(cmd.LocalFlags()), strings.Join(flags, "\n"), positionalAction))
	for _, subcmd := range cmd.Commands() {
		if !subcmd.Hidden {
			result = append(result, snippetFunctions(subcmd, actions))
		}
	}
	return strings.Join(result, "\n")
}

func snippetFlagList(flags *pflag.FlagSet) string {
	flagValues := make([]common.Candidate, 0)

	flags.VisitAll(func(flag *pflag.Flag) {
		if !flag.Hidden {
			if !common.IsShorthandOnly(flag) {
				//flagValues = append(flagValues, "--"+flag.Name)
				flagValues = append(flagValues, common.Candidate{Value: "--" + flag.Name, Display: "--" + flag.Name, Description: flag.Usage})
			}
			if flag.Shorthand != "" {
				//flagValues = append(flagValues, "-"+flag.Shorthand)
				flagValues = append(flagValues, common.Candidate{Value: "-" + flag.Shorthand, Display: "-" + flag.Shorthand, Description: flag.Usage})
			}
		}
	})
	if len(flagValues) > 0 {
		return ActionCandidates(flagValues...) // TODO use candidatas
	} else {
		return ""
	}
}

func snippetFlagCompletion(flag *pflag.Flag, action string, optArgFlag bool) (snippet string) {
	// TODO cleanup this mess
	if flag.Value.Type() == "bool" {
		return
	}
	if flag.NoOptDefVal != "" && !optArgFlag {
		return
	}
	if flag.NoOptDefVal == "" && optArgFlag {
		return
	}

	optArgSuffix := ""
	if flag.NoOptDefVal != "" {
		optArgSuffix = "=*"
	}

	var names string
	if flag.Shorthand != "" {
		if common.IsShorthandOnly(flag) {
			names = fmt.Sprintf("-%v", flag.Shorthand+optArgSuffix)
		} else {
			names = fmt.Sprintf("-%v | --%v", flag.Shorthand+optArgSuffix, flag.Name+optArgSuffix)
		}
	} else {
		names = "--" + flag.Name + optArgSuffix
	}

	if optArgFlag {
		return fmt.Sprintf(`          %v)
            cur=${cur#*=}
            curprefix=${curprefix#*=}
            COMPREPLY=($(%v))
            ;;
`, names, action)

	} else {
		return fmt.Sprintf(`          %v)
            COMPREPLY=($(%v))
            ;;
`, names, action)
	}
}
