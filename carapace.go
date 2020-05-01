package carapace

import (
	"fmt"
	"os"
	"strings"

	"github.com/rsteube/carapace/bash"
	"github.com/rsteube/carapace/elvish"
	"github.com/rsteube/carapace/fish"
	"github.com/rsteube/carapace/powershell"
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

func (c Carapace) PositionalAnyCompletion(action Action) {
	completions.actions[uid.Positional(c.cmd, 0)] = action.finalize(c.cmd, uid.Positional(c.cmd, 0))
}

func (c Carapace) FlagCompletion(actions ActionMap) {
	for name, action := range actions {
		if flag := c.cmd.LocalFlags().Lookup(name); flag == nil {
			fmt.Fprintf(os.Stderr, "unknown flag: %v\n", name)
		} else {
			completions.actions[uid.Flag(c.cmd, flag)] = action.finalize(c.cmd, uid.Flag(c.cmd, flag))
		}
	}
}

func (c Carapace) Bash() string {
	actions := make(map[string]string, len(completions.actions))
	for key, value := range completions.actions {
		actions[key] = value.Bash
	}
	return bash.Snippet(c.cmd.Root(), actions)
}

func (c Carapace) Elvish() string {
	actions := make(map[string]string, len(completions.actions))
	for key, value := range completions.actions {
		actions[key] = value.Elvish
	}
	return elvish.Snippet(c.cmd.Root(), actions)
}

func (c Carapace) Fish() string {
	actions := make(map[string]string, len(completions.actions))
	for key, value := range completions.actions {
		actions[key] = value.Fish
	}
	return fish.Snippet(c.cmd.Root(), actions)
}

func (c Carapace) Powershell() string {
	actions := make(map[string]string, len(completions.actions))
	for key, value := range completions.actions {
		actions[key] = value.Powershell
	}
	return powershell.Snippet(c.cmd.Root(), actions)
}

func (c Carapace) Zsh() string {
	actions := make(map[string]string, len(completions.actions))
	for key, value := range completions.actions {
		actions[key] = value.Zsh
	}
	return zsh.Snippet(c.cmd.Root(), actions)
}

func (c Carapace) Standalone() {
	// TODO probably needs to be done for each subcommand
	if c.cmd.Root().Flag("help") != nil {
		c.cmd.Root().Flags().Bool("help", false, "skip")
		c.cmd.Root().Flag("help").Hidden = true
	}
	c.cmd.Root().SetHelpCommand(&cobra.Command{Hidden: true})
}

func (c Carapace) Snippet(shell string) string {
	switch shell {
	case "bash":
		return c.Bash()
	case "elvish":
		return c.Elvish()
	case "fish":
		return c.Fish()
	case "powershell":
		return c.Powershell()
	case "zsh":
		return c.Zsh()
	default:
		return fmt.Sprintf("expected 'bash', 'fish', 'powershell' or 'zsh' [was: %v]", shell)
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
				if len(args) == 1 {
					switch args[0] {
					case "debug":
						for uid, action := range completions.actions {
							fmt.Printf("%v:\t%v\n", uid, action)
						}
					default:
						fmt.Println(Gen(cmd).Snippet(args[0]))
					}
				} else {
					targetCmd, targetArgs := findTarget(cmd)

					shell := args[0]
					id := args[1]

					switch id {
					case "_":
						if action, ok := findAction(targetCmd, targetArgs); ok {
							if action.Callback == nil {
								fmt.Println(action.Value(shell))
							} else {
								fmt.Println(action.Callback(targetArgs).Value(shell))
							}
						}
					case "state":
						fmt.Println(uid.Command(targetCmd))
					default:
						fmt.Println(completions.invokeCallback(id, targetArgs).Value(shell))
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

func findAction(targetCmd *cobra.Command, targetArgs []string) (action Action, ok bool) {
	var id string
	if len(targetArgs) == 0 {
		id = uid.Positional(targetCmd, 1)
	} else {
		lastArg := targetArgs[len(targetArgs)-1]
		if strings.HasSuffix(lastArg, " ") {
			id = uid.Positional(targetCmd, len(targetArgs)+1)
		} else {
			id = uid.Positional(targetCmd, len(targetArgs))
		}
	}
	if action, ok = completions.actions[id]; !ok {
		id = uid.Positional(targetCmd, 0)
		action, ok = completions.actions[id]
	}
	return
}

func findTarget(cmd *cobra.Command) (*cobra.Command, []string) {
	origArg := []string{}
	if len(os.Args) > 5 {
		origArg = os.Args[5:]
	}
	return traverse(cmd, origArg)
}

func traverse(cmd *cobra.Command, args []string) (*cobra.Command, []string) {
	// ignore flag parse errors (like a missing argument for the flag currently being completed)
	targetCmd, targetArgs, _ := cmd.Root().Traverse(args)
	targetCmd.ParseFlags(targetArgs)
	return targetCmd, targetCmd.Flags().Args() // TODO check length
}

func IsCallback() bool {
	return len(os.Args) > 3 && os.Args[1] == "_carapace" && os.Args[3] != "state"
}
