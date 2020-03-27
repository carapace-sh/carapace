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

type Carapace struct {
	cmd *cobra.Command
}

func Gen(cmd *cobra.Command) *Carapace {
	addCompletionCommand(cmd)
	return &Carapace{
		cmd: cmd,
	}
}

func (c Carapace) PositionalCompletion(action ...Action) {
	for index, a := range action {
		completions.actions[uid.Positional(c.cmd, index+1)] = a.finalize(c.cmd, uid.Positional(c.cmd, index+1))
	}
}

func (c Carapace) FlagCompletion(actions ActionMap) {
	for name, action := range actions {
		flag := c.cmd.Flag(name) // TODO only allowed for local flags
		completions.actions[uid.Flag(c.cmd, flag)] = action.finalize(c.cmd, uid.Flag(c.cmd, flag))
	}
}

func (c Carapace) Bash() string {
	actions := make(map[string]string, len(completions.actions))
	for key, value := range completions.actions {
		actions[key] = value.Bash
	}
	return bash.Snippet(c.cmd.Root(), actions)
}

func (c Carapace) Fish() string {
	actions := make(map[string]string, len(completions.actions))
	for key, value := range completions.actions {
		actions[key] = value.Fish
	}
	return fish.Snippet(c.cmd.Root(), actions)
}

func (c Carapace) Zsh() string {
	actions := make(map[string]string, len(completions.actions))
	for key, value := range completions.actions {
		actions[key] = value.Zsh
	}
	return zsh.Snippet(c.cmd.Root(), actions)
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
				if len(args) == 1 {
					switch args[0] {
					case "bash":
						fmt.Println(Gen(cmd).Bash())
					case "fish":
						fmt.Println(Gen(cmd).Fish())
					case "zsh":
						fmt.Println(Gen(cmd).Zsh())
					}
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
								if args[0] == "bash" { // TODO print for other shells as well?
									subcommands := make([]string, 0)
									for _, c := range targetCmd.Commands() {
										if !c.Hidden {
											subcommands = append(subcommands, c.Name())
											for _, alias := range c.Aliases {
												subcommands = append(subcommands, alias)
											}
										}
									}
									fmt.Println(ActionValues(subcommands...).Bash)
								}
							}
							os.Exit(0) // ensure no message for missing action on positional completion // TODO this was only for bash, maybe enable for other shells?
						}
					} else if callback == "state" {
						fmt.Println(uid.Command(targetCmd))
						os.Exit(0) // TODO
					}
					switch args[0] {
					case "bash":
						fmt.Println(completions.invokeCallback(callback, targetArgs).Bash)
					case "fish":
						fmt.Println(completions.invokeCallback(callback, targetArgs).Fish)
					case "zsh":
						fmt.Println(completions.invokeCallback(callback, targetArgs).Zsh)
					}
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
