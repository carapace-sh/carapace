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
	"github.com/rsteube/carapace/internal/common"
	"github.com/rsteube/carapace/internal/elvish"
	"github.com/rsteube/carapace/internal/fish"
	"github.com/rsteube/carapace/internal/oil"
	"github.com/rsteube/carapace/internal/powershell"
	"github.com/rsteube/carapace/internal/uid"
	"github.com/rsteube/carapace/internal/xonsh"
	"github.com/rsteube/carapace/internal/zsh"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

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
	storage.get(c.cmd).positional = action
}

func (c Carapace) PositionalAnyCompletion(action Action) {
	storage.get(c.cmd).positionalAny = action
}

func (c Carapace) FlagCompletion(actions ActionMap) {
	storage.get(c.cmd).flag = actions
}

func (c Carapace) Standalone() {
	// TODO probably needs to be done for each subcommand
	if c.cmd.Root().Flag("help") != nil {
		c.cmd.Root().Flags().Bool("help", false, "skip")
		c.cmd.Root().Flag("help").Hidden = true
	}
	c.cmd.Root().SetHelpCommand(&cobra.Command{Hidden: true})
}

func (c Carapace) Snippet(shell string) (string, error) {
	var snippet func(cmd *cobra.Command) string

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
		snippet = oil.Snippet
	case "powershell":
		snippet = powershell.Snippet
	case "xonsh":
		snippet = xonsh.Snippet
	case "zsh":
		snippet = zsh.Snippet
	default:
		return "", errors.New(fmt.Sprintf("expected 'bash', 'elvish', 'fish', 'oil', 'powershell', 'xonsh' or 'zsh' [was: %v]", shell))
	}
	return snippet(c.cmd.Root()), nil
}

func lookupFlag(cmd *cobra.Command, arg string) (flag *pflag.Flag) {
	nameOrShorthand := strings.TrimLeft(strings.SplitN(arg, "=", 2)[0], "-")

	if strings.HasPrefix(arg, "--") {
		flag = cmd.Flags().Lookup(nameOrShorthand)
	} else if strings.HasPrefix(arg, "-") {
		flag = cmd.Flags().ShorthandLookup(string(nameOrShorthand[len(nameOrShorthand)-1]))
	}
	return
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
				if s, err := Gen(cmd).Snippet(determineShell()); err != nil {
					fmt.Fprintln(io.MultiWriter(os.Stderr, logger.Writer()), err.Error())
				} else {
					fmt.Fprintln(io.MultiWriter(os.Stdout, logger.Writer()), s)
				}
			} else {
				if len(args) == 1 {
					switch args[0] {
					case "debug":
						// TODO broken, not indexed by uid anymore
						for uid, action := range storage {
							fmt.Printf("%v:\t%v\n", uid, action)
						}
					default:
						if s, err := Gen(cmd).Snippet(args[0]); err != nil {
							fmt.Fprintln(io.MultiWriter(os.Stderr, logger.Writer()), err.Error())
						} else {
							fmt.Fprintln(io.MultiWriter(os.Stdout, logger.Writer()), s)
						}
					}
				} else {
					shell := args[0]
					id := args[1]
					current := args[len(args)-1]
					previous := args[len(args)-2]

					targetCmd, targetArgs, err := findTarget(cmd, args)
					context := Context{CallbackValue: current, Args: targetArgs}
					if err != nil {
						fmt.Fprintln(io.MultiWriter(os.Stdout, logger.Writer()), ActionMessage(err.Error()).Invoke(context).value(shell, context.CallbackValue))
						return
					}

					switch id {
					case "_":
						// TODO needs more cleanup and tests
						var targetAction Action
						if flag := lookupFlag(targetCmd, previous); flag != nil && flag.NoOptDefVal == "" { // previous arg is a flag and needs a value
							targetAction = storage.getFlag(targetCmd, flag.Name)
						} else if strings.HasPrefix(current, "-") { // assume flag
							if strings.Contains(current, "=") { // complete value for optarg flag
								if flag := lookupFlag(targetCmd, current); flag != nil && flag.NoOptDefVal != "" {
									a := storage.getFlag(targetCmd, flag.Name)
									// TODO no value for oil (elvish works)
									splitted := strings.SplitN(current, "=", 2)
									context.CallbackValue = splitted[1]
									targetAction = a.Invoke(context).Prefix(splitted[0] + "=").ToA()
								}
							} else { // complete flagnames
								targetAction = actionFlags(targetCmd)
							}
						} else if targetCmd.HasAvailableSubCommands() && len(targetArgs) <= 1 {
							subcommandA := actionSubcommands(targetCmd)
							a := findAction(targetCmd, targetArgs)
							subcommandA = a.Invoke(context).Merge(subcommandA.Invoke(context)).ToA()
							targetAction = subcommandA
						} else {
							targetAction = findAction(targetCmd, targetArgs)
						}
						fmt.Fprintln(io.MultiWriter(os.Stdout, logger.Writer()), targetAction.Invoke(context).value(shell, context.CallbackValue))
					default:
						// TODO disable support for direct uid invocation
						//fmt.Fprintln(io.MultiWriter(os.Stdout, logger.Writer()), actionMap.invokeCallback(id, context).Invoke(context).value(shell, context.CallbackValue))
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

func findAction(targetCmd *cobra.Command, targetArgs []string) Action {
	// TODO handle Action not found
	if len(targetArgs) == 0 {
		return storage.getPositional(targetCmd, 0)
	} else {
		lastArg := targetArgs[len(targetArgs)-1]
		if strings.HasSuffix(lastArg, " ") { // TODO is this still correct/needed?
			return storage.getPositional(targetCmd, len(targetArgs))
		} else {
			return storage.getPositional(targetCmd, len(targetArgs)-1)
		}
	}
}

func findTarget(cmd *cobra.Command, args []string) (*cobra.Command, []string, error) {
	origArg := []string{}
	if len(args) > 3 {
		origArg = args[3:]
	}
	return common.TraverseLenient(cmd, origArg)
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
		case "nu":
			return "nushell"
		case "oil":
			return "oil"
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

type testing interface {
	Error(args ...interface{})
}

// Test verifies the configuration (e.g. flag name exists)
//   func TestCarapace(t *testing.T) {
//       carapace.Test(t)
//   }
func Test(t testing) {
	for _, e := range storage.check() {
		testing(t).Error(e)
	}
}
