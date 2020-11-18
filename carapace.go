// Pacakge carapace provides multi-shell completion script generation for spf13/cobra
package carapace

import (
	"errors"
	"fmt"
	"os"
	"strings"

	ps "github.com/mitchellh/go-ps"
	"github.com/rsteube/carapace/bash"
	"github.com/rsteube/carapace/elvish"
	"github.com/rsteube/carapace/fish"
	"github.com/rsteube/carapace/powershell"
	"github.com/rsteube/carapace/uid"
	"github.com/rsteube/carapace/xonsh"
	"github.com/rsteube/carapace/zsh"
	"github.com/spf13/cobra"
)

type Completions struct {
	actions ActionMap
}

func (c Completions) invokeCallback(uid string, args []string) Action {
	if action, ok := c.actions[uid]; ok {
		if action.callback != nil {
			return action.callback(args)
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

func (c Carapace) Standalone() {
	// TODO probably needs to be done for each subcommand
	if c.cmd.Root().Flag("help") != nil {
		c.cmd.Root().Flags().Bool("help", false, "skip")
		c.cmd.Root().Flag("help").Hidden = true
	}
	c.cmd.Root().SetHelpCommand(&cobra.Command{Hidden: true})
}

func (c Carapace) Snippet(shell string, lazy bool) (string, error) {
	var snippet func(cmd *cobra.Command, actions map[string]string, lazy bool) string

	if shell == "" {
		shell = determineShell()
	}
	switch shell {
	case "bash":
		snippet = bash.Snippet
	case "elvish":
		snippet = elvish.Snippet
	case "fish":
		snippet = fish.Snippet
	case "oil":
		snippet = bash.Snippet
	case "powershell":
		snippet = powershell.Snippet
	case "xonsh":
		snippet = xonsh.Snippet
	case "zsh":
		snippet = zsh.Snippet
	default:
		return "", errors.New(fmt.Sprintf("expected 'bash', 'elvish', 'fish', 'oil', 'powershell', 'xonsh' or 'zsh' [was: %v]", shell))
	}
	return snippet(c.cmd.Root(), completions.actions.Shell(shell), lazy), nil
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
				if s, err := Gen(cmd).Snippet(determineShell(), true); err != nil {
					fmt.Fprintln(os.Stderr, err.Error())
				} else {
					fmt.Println(s)
				}
			} else {
				if len(args) == 1 {
					switch args[0] {
					case "debug":
						for uid, action := range completions.actions {
							fmt.Printf("%v:\t%v\n", uid, action)
						}
					default:
						if s, err := Gen(cmd).Snippet(args[0], false); err != nil {
							fmt.Fprintln(os.Stderr, err.Error())
						} else {
							fmt.Println(s)
						}
					}
				} else {
					targetCmd, targetArgs := findTarget(cmd)

					shell := args[0]
					id := args[1]

					switch id {
					case "_":
						if _uid, action, ok := findAction(targetCmd, targetArgs); ok {
							CallbackValue = uid.Value(targetCmd, targetArgs, _uid)
							if action.callback == nil {
								fmt.Println(action.Value(shell))
							} else {
								fmt.Println(action.callback(targetArgs).nestedAction(targetArgs, 2).Value(shell))
							}
						}
					case "state":
						fmt.Println(uid.Command(targetCmd))
					default:
						CallbackValue = uid.Value(targetCmd, targetArgs, id)
						fmt.Println(completions.invokeCallback(id, targetArgs).nestedAction(targetArgs, 2).Value(shell))
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

func findAction(targetCmd *cobra.Command, targetArgs []string) (id string, action Action, ok bool) {
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
	a := args
	if len(args) > 0 && args[len(args)-1] == "" {
		a = args[0 : len(args)-1]
	}

	targetCmd, targetArgs, _ := cmd.Root().Traverse(a)
	if len(args) > 0 && args[len(args)-1] == "" {
		targetArgs = append(targetArgs, "")
	}
	targetCmd.ParseFlags(targetArgs)
	return targetCmd, targetCmd.Flags().Args() // TODO check length
}

func determineShell() string {
	process, err := ps.FindProcess(os.Getpid())
	for {
		if process, err = ps.FindProcess(process.PPid()); err != nil {
			return ""
		} else {
			switch process.Executable() {
			case "bash":
				return "bash"
			case "elvish":
				return "elvish"
			case "fish":
				return "fish"
			case "osh":
				return "oil"
			case "pwsh":
				return "powershell"
			case "xonsh":
				return "xonsh"
			case "zsh":
				return "zsh"
			default:
				return ""
			}
		}
	}
}

func IsCallback() bool {
	return len(os.Args) > 3 && os.Args[1] == "_carapace" && os.Args[3] != "state"
}
