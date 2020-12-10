// Pacakge carapace provides multi-shell completion script generation for spf13/cobra
package carapace

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/rsteube/carapace/internal/bash"
	"github.com/rsteube/carapace/internal/elvish"
	"github.com/rsteube/carapace/internal/fish"
	"github.com/rsteube/carapace/internal/powershell"
	"github.com/rsteube/carapace/internal/uid"
	"github.com/rsteube/carapace/internal/xonsh"
	"github.com/rsteube/carapace/internal/zsh"
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
			logger.Println(os.Args) // TODO replace last with '' if empty

			if len(args) == 0 {
				if s, err := Gen(cmd).Snippet(determineShell(), true); err != nil {
					fmt.Fprintln(io.MultiWriter(os.Stderr, logger.Writer()), err.Error())
				} else {
					fmt.Fprintln(io.MultiWriter(os.Stdout, logger.Writer()), s)
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
							fmt.Fprintln(io.MultiWriter(os.Stderr, logger.Writer()), err.Error())
						} else {
							fmt.Fprintln(io.MultiWriter(os.Stdout, logger.Writer()), s)
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
								fmt.Fprintln(io.MultiWriter(os.Stdout, logger.Writer()), action.Value(shell))
							} else {
								fmt.Fprintln(io.MultiWriter(os.Stdout, logger.Writer()), action.callback(targetArgs).nestedAction(targetArgs, 2).Value(shell))
							}
						}
					case "state":
						fmt.Fprintln(io.MultiWriter(os.Stdout, logger.Writer()), uid.Command(targetCmd))
					default:
						CallbackValue = uid.Value(targetCmd, targetArgs, id)
						fmt.Fprintln(io.MultiWriter(os.Stdout, logger.Writer()), completions.invokeCallback(id, targetArgs).nestedAction(targetArgs, 2).Value(shell))
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
	for _, executable := range processExecutables() {
		switch executable {
		case "bash":
			return "bash"
		case "elvish":
			return "elvish"
		case "fish":
			return "fish"
		case "osh":
			return "oil"
		case "powershell.exe":
			return "powershell"
		case "pwsh":
			return "powershell"
		case "pwsh.exe":
			return "powershell"
		case "xonsh":
			return "xonsh"
		case "zsh":
			return "zsh"
		}
	}
	return ""
}

func processExecutables() []string {
	if runtime.GOOS == "windows" {
		return []string{"powershell.exe"} // TODO hardcoded for now, but there might be elvish or sth. else on window
	} else {
		if output, err := exec.Command("ps", "-o", "comm").Output(); err != nil {
			return []string{}
		} else {
			lines := strings.Split(string(output), "\n")[1:]      // skip header
			for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 { // reverse slice
				lines[i], lines[j] = lines[j], lines[i]
			}
			return lines
		}
	}
}

func IsCallback() bool {
	return len(os.Args) > 3 && os.Args[1] == "_carapace" && os.Args[3] != "state"
}

var logger = log.New(ioutil.Discard, "", log.Flags())

func init() {
	if _, enabled := os.LookupEnv("CARAPACE_LOG"); enabled {
		if err := initLogger(); err != nil {
			log.Fatal(err.Error())
		}
	}
}

func initLogger() (err error) {
	tmpdir := fmt.Sprintf("%v/carapace", os.TempDir())
	if err = os.MkdirAll(tmpdir, os.ModePerm); err == nil {
		var logfileWriter io.Writer
		if logfileWriter, err = os.OpenFile(fmt.Sprintf("%v/%v.log", tmpdir, uid.Executable()), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666); err == nil {
			Lmsgprefix := 1 << 6
			logger = log.New(logfileWriter, determineShell()+" ", log.Flags()|Lmsgprefix)
			//logger = log.New(logfileWriter, determineShell()+" ", log.Flags()|log.Lmsgprefix)
		}
	}
	return
}
