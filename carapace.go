package carapace

import (
	"fmt"
	"os"
	"strings"

	"github.com/rsteube/carapace/bash"
	"github.com/rsteube/carapace/fish"
	"github.com/rsteube/carapace/uid"
	"github.com/rsteube/carapace/zsh"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Completions struct {
	actions ActionMap
}

func (c Completions) invokeCallback(uid string, args []string) Action {
	if action, ok := c.actions[uid]; ok {
		if action.Callback != nil {
			return action.Callback(args)
		}
	}
	return ActionMessage(fmt.Sprintf("callback %v unknown", uid))
}

func (c Completions) GenerateZsh(cmd *cobra.Command) string {
	result := fmt.Sprintf("#compdef %v\n", cmd.Name())
	result += c.GenerateZshFunctions(cmd)

	result += fmt.Sprintf("if compquote '' 2>/dev/null; then _%v; else compdef _%v %v; fi\n", cmd.Name(), cmd.Name(), cmd.Name()) // check if withing completion function and enable direct sourcing
	return result
}

func (c Completions) GenerateZshFunctions(cmd *cobra.Command) string {
	function_pattern := `function %v {
  %v%v  _arguments -C \
%v%v
}
`

	commandsVar := ""
	if cmd.HasSubCommands() {
		commandsVar = "local -a commands\n"
	}

	inheritedArgs := ""
	if !cmd.HasParent() {
		inheritedArgs = "  # shellcheck disable=SC2206\n  local -a -x os_args=(${words})\n\n"
	}

	flags := make([]string, 0)
	for _, flag := range zshCompExtractFlag(cmd) {
		if flagAlreadySet(cmd, flag) {
			continue
		}

		var s string
		if action, ok := c.actions[uid.Flag(cmd, flag)]; ok {
			s = "    " + zsh.SnippetFlagCompletion(flag, &action.Zsh)
		} else {
			s = "    " + zsh.SnippetFlagCompletion(flag, nil)
		}

		flags = append(flags, s)
	}

	positionals := make([]string, 0)
	if cmd.HasSubCommands() {
		positionals = []string{`    "1: :->cmnds"`, `    "*::arg:->args"`}
	} else {
		pos := 1
		for {
			if action, ok := c.actions[uid.Positional(cmd, pos)]; ok {
				positionals = append(positionals, "    "+zsh.SnippetPositionalCompletion(pos, action.Zsh))
				pos++
			} else {
				break // TODO only consisten entriess for now
			}
		}
		if len(positionals) == 0 {
			if cmd.ValidArgs != nil {
				positionals = []string{"    " + zsh.SnippetPositionalCompletion(1, ActionValues(cmd.ValidArgs...).Zsh)}
			}
			positionals = append(positionals, `    "*::arg:->args"`)
		}
	}

	arguments := append(flags, positionals...)

	result := make([]string, 0)
	result = append(result, fmt.Sprintf(function_pattern, uid.Command(cmd), commandsVar, inheritedArgs, strings.Join(arguments, " \\\n"), zsh.SnippetSubcommands(cmd)))
	for _, subcmd := range cmd.Commands() {
		if !subcmd.Hidden {
			result = append(result, c.GenerateZshFunctions(subcmd))
		}
	}

	return strings.Join(result, "\n")
}

//fish
func (c Completions) GenerateFish(cmd *cobra.Command) string {
	result := fmt.Sprintf(`function _state
  set -lx CURRENT (commandline -cp)
  if [ "$LINE" != "$CURRENT" ]
    set -gx LINE (commandline -cp)
    set -gx STATE (commandline -cp | xargs %v _carapace fish state)
  end

  [ "$STATE" = "$argv" ]
end

function _callback
  set -lx CALLBACK (commandline -cp | sed "s/ \$/ _/" | xargs %v _carapace fish $argv )
  eval "$CALLBACK"
end

complete -c %v -f
`, cmd.Name(), cmd.Name(), cmd.Name())
	result += c.GenerateFishFunctions(cmd)

	return result
}

func (c Completions) GenerateFishFunctions(cmd *cobra.Command) string {
	// TODO ensure state is only called oncy per LINE
	function_pattern := `
%v
`

	flags := make([]string, 0)
	for _, flag := range zshCompExtractFlag(cmd) {
		if flagAlreadySet(cmd, flag) {
			continue
		}

		var s string
		if action, ok := c.actions[uid.Flag(cmd, flag)]; ok {
			s = fish.SnippetFlagCompletion(cmd, flag, &action.Fish)
		} else {
			s = fish.SnippetFlagCompletion(cmd, flag, nil)
		}

		flags = append(flags, s)
	}

	positionals := make([]string, 0)
	if cmd.HasSubCommands() {
		positionals = []string{}
		for _, subcmd := range cmd.Commands() {
			positionals = append(positionals, fmt.Sprintf(`complete -c %v -f -n '_state %v ' -a %v -d '%v'`, cmd.Root().Name(), uid.Command(cmd), subcmd.Name(), subcmd.Short))
			// TODO repeat for aliases
			// TODO filter hidden
		}
	} else {
		if len(positionals) == 0 {
			if cmd.ValidArgs != nil {
				//positionals = []string{"    " + snippetPositionalCompletion(1, ActionValues(cmd.ValidArgs...))}
			}
			positionals = append(positionals, fmt.Sprintf(`complete -c %v -f -n '_state %v' -a '(_callback _)'`, cmd.Root().Name(), uid.Command(cmd)))
		}
	}

	arguments := append(flags, positionals...)

	result := make([]string, 0)
	result = append(result, fmt.Sprintf(function_pattern, strings.Join(arguments, "\n")))
	for _, subcmd := range cmd.Commands() {
		if !subcmd.Hidden {
			result = append(result, c.GenerateFishFunctions(subcmd))
		}
	}
	return strings.Join(result, "\n")
}

//fish

// bash
func (c Completions) GenerateBash(cmd *cobra.Command) string {
	result := fmt.Sprintf(`#!/bin/bash
_%v_completions() {
  _callback() {
    local compline="${COMP_LINE:0:${COMP_POINT}}"
    echo "$compline" | sed "s/ \$/ _/" | xargs %v _carapace bash "$1"
  }

  local compline="${COMP_LINE:0:${COMP_POINT}}"
  local state=$(echo "$compline" | sed "s/ \$/ _/" | xargs %v _carapace bash state)
  local last="${COMP_WORDS[${COMP_CWORD}]}"
  local previous="${COMP_WORDS[$((${COMP_CWORD}-1))]}"

  case $state in
%v
  esac
}

complete -F _%v_completions %v
`, cmd.Name(), cmd.Name(), cmd.Name(), c.GenerateBashFunctions(cmd), cmd.Name(), cmd.Name())

	return result
}

func (c Completions) GenerateBashFunctions(cmd *cobra.Command) string {
	// TODO ensure state is only called oncy per LINE
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
			if action, ok := c.actions[uid.Flag(cmd, f)]; ok {
				s = bash.SnippetFlagCompletion(f, action.Bash)
			} else {
				s = bash.SnippetFlagCompletion(f, "")
			}
			flags = append(flags, s)
		}
	})

	result := make([]string, 0)
	// uid.Command, flagList, genflagcompletions, commandargumentcompletion
	result = append(result, fmt.Sprintf(function_pattern, uid.Command(cmd), bash.SnippetFlagList(cmd.LocalFlags()), strings.Join(flags, "\n"), bash.Callback("_")))
	for _, subcmd := range cmd.Commands() {
		if !subcmd.Hidden {
			result = append(result, c.GenerateBashFunctions(subcmd))
		}
	}
	return strings.Join(result, "\n")
}

// bash

func flagAlreadySet(cmd *cobra.Command, flag *pflag.Flag) bool {
	if cmd.LocalFlags().Lookup(flag.Name) != nil {
		return false
	}
	// TODO since it is an inherited flag check for parent command that is not hidden
	return true
}

func zshCompExtractFlag(c *cobra.Command) []*pflag.Flag {
	var flags []*pflag.Flag
	c.LocalFlags().VisitAll(func(f *pflag.Flag) {
		if !f.Hidden {
			flags = append(flags, f)
		}
	})
	c.InheritedFlags().VisitAll(func(f *pflag.Flag) {
		if !f.Hidden {
			flags = append(flags, f)
		}
	})
	return flags
}

type Carapace struct {
	cmd *cobra.Command
}

func Gen(cmd *cobra.Command) *Carapace {
	addCompletionCommand(cmd)
	return &Carapace{
		cmd: cmd,
	}
}

func (zsh Carapace) PositionalCompletion(action ...Action) {
	for index, a := range action {
		completions.actions[uid.Positional(zsh.cmd, index+1)] = a.finalize(uid.Positional(zsh.cmd, index+1))
	}
}

func (zsh Carapace) FlagCompletion(actions ActionMap) {
	for name, action := range actions {
		flag := zsh.cmd.Flag(name) // TODO only allowed for local flags
		completions.actions[uid.Flag(zsh.cmd, flag)] = action.finalize(uid.Flag(zsh.cmd, flag))
	}
}

var completions = Completions{
	actions: make(map[string]Action),
}

func addCompletionCommand(cmd *cobra.Command) {
	for _, c := range cmd.Root().Commands() {
		if c.Name() == "_carapace" {
			return
		}
	}
	cmd.Root().AddCommand(&cobra.Command{
		Use:    "_carapace",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("zsh/fish argument missing")
			} else {
				if args[0] == "zsh" {
					if len(args) <= 1 {
						fmt.Println(completions.GenerateZsh(cmd.Root()))
					} else {
						callback := args[1]
						origArg := []string{}
						if len(os.Args) > 4 {
							origArg = os.Args[5:]
						}
						_, targetArgs := traverse(cmd, origArg)
						fmt.Println(completions.invokeCallback(callback, targetArgs).Zsh)
					}
				} else if args[0] == "bash" {
					if len(args) <= 1 {
						fmt.Println(completions.GenerateBash(cmd.Root()))
					} else {
						callback := args[1]
						origArg := []string{}
						if len(os.Args) > 5 {
							origArg = os.Args[5:]
						}
						targetCmd, targetArgs := traverse(cmd, origArg)
						if callback == "_" {
							if len(targetArgs) == 0 {
								callback = uid.Positional(targetCmd, 1)
							} else {
								lastArg := targetArgs[len(targetArgs)-1]
								if strings.HasSuffix(lastArg, " ") {
									callback = uid.Positional(targetCmd, len(targetArgs)+1)
								} else {
									callback = uid.Positional(targetCmd, len(targetArgs))
								}
							}
							if _, ok := completions.actions[callback]; !ok {
								if targetCmd.HasSubCommands() && len(targetArgs) <= 1 {
									subcommands := make([]string, len(targetCmd.Commands()))
									for i, c := range targetCmd.Commands() {
										subcommands[i] = c.Name() // TODO alias
									}
									fmt.Println(bash.ActionValues(subcommands...))
								}
								os.Exit(0) // ensure no message for missing action on positional completion
							}
						} else if callback == "state" {
							fmt.Println(uid.Command(targetCmd))
							os.Exit(0) // TODO
						}
						fmt.Println(completions.invokeCallback(callback, targetArgs).Bash)
					}
				} else { // fish
					// fish
					if len(args) <= 1 {
						fmt.Println(completions.GenerateFish(cmd.Root()))
					} else {
						callback := args[1]
						origArg := []string{}
						if len(os.Args) > 5 {
							origArg = os.Args[5:]
						}
						targetCmd, targetArgs := traverse(cmd, origArg)
						if callback == "_" {
							if len(targetArgs) == 0 {
								callback = uid.Positional(targetCmd, 1)
							} else {
								lastArg := targetArgs[len(targetArgs)-1]
								if strings.HasSuffix(lastArg, " ") {
									callback = uid.Positional(targetCmd, len(targetArgs)+1)
								} else {
									callback = uid.Positional(targetCmd, len(targetArgs))
								}
							}
							if _, ok := completions.actions[callback]; !ok {
								os.Exit(0) // ensure no message for missing action on positional completion
							}
						} else if callback == "state" {
							fmt.Println(uid.Command(targetCmd))
							os.Exit(0) // TODO
						}
						fmt.Println(completions.invokeCallback(callback, targetArgs).Fish)
					}
					//fish
				}
			}
		},
		FParseErrWhitelist: cobra.FParseErrWhitelist{
			UnknownFlags: true,
		},
		DisableFlagParsing: true,
	})
}

func traverse(cmd *cobra.Command, args []string) (*cobra.Command, []string) {
	// ignore flag parse errors (like a missing argument for the flag currently being completed)
	targetCmd, targetArgs, _ := cmd.Root().Traverse(args)
	targetCmd.ParseFlags(targetArgs)
	return targetCmd, targetCmd.Flags().Args() // TODO check length
}
