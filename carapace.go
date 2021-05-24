package carapace

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/rsteube/carapace/internal/bash"
	"github.com/rsteube/carapace/internal/common"
	"github.com/rsteube/carapace/internal/elvish"
	"github.com/rsteube/carapace/internal/fish"
	"github.com/rsteube/carapace/internal/ion"
	"github.com/rsteube/carapace/internal/nushell"
	"github.com/rsteube/carapace/internal/oil"
	"github.com/rsteube/carapace/internal/powershell"
	"github.com/rsteube/carapace/internal/uid"
	"github.com/rsteube/carapace/internal/xonsh"
	"github.com/rsteube/carapace/internal/zsh"
	"github.com/rsteube/carapace/pkg/ps"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Carapace wraps cobra.Command to define completions
type Carapace struct {
	cmd *cobra.Command
}

// Gen initialized Carapace for given command
func Gen(cmd *cobra.Command) *Carapace {
	addCompletionCommand(cmd)
	return &Carapace{
		cmd: cmd,
	}
}

// PositionalCompletion defines completion for positional arguments using a list of Actions
func (c Carapace) PositionalCompletion(action ...Action) {
	storage.get(c.cmd).positional = action
}

// PositionalAnyCompletion defines completion for any positional argument not already defined
func (c Carapace) PositionalAnyCompletion(action Action) {
	storage.get(c.cmd).positionalAny = action
}

// FlagCompletion defines completion for flags using a map consisting of name and Action
func (c Carapace) FlagCompletion(actions ActionMap) {
	if e := storage.get(c.cmd); e.flag == nil {
		e.flag = actions
	} else {
		for name, action := range actions {
			e.flag[name] = action
		}
	}
}

// Standalone prevents cobra defaults interfering with standalone mode (e.g. implicit help command)
func (c Carapace) Standalone() {
	// TODO probably needs to be done for each subcommand
	// TODO still needed?
	if c.cmd.Flag("help") != nil {
		c.cmd.Flags().Bool("help", false, "skip")
		c.cmd.Flag("help").Hidden = true
	}
	c.cmd.SetHelpCommand(&cobra.Command{Hidden: true})
}

// Snippet creates completion script for given shell
func (c Carapace) Snippet(shell string) (string, error) {
	var snippet func(cmd *cobra.Command) string

	if shell == "" {
		shell = ps.DetermineShell()
	}
	switch shell {
	case "bash":
		snippet = bash.Snippet
	case "elvish":
		snippet = elvish.Snippet
	case "fish":
		snippet = fish.Snippet
	case "ion":
		snippet = ion.Snippet
	case "nushell":
		snippet = nushell.Snippet
	case "oil":
		snippet = oil.Snippet
	case "powershell":
		snippet = powershell.Snippet
	case "xonsh":
		snippet = xonsh.Snippet
	case "zsh":
		snippet = zsh.Snippet
	default:
		return "", fmt.Errorf("expected 'bash', 'elvish', 'fish', 'ion', 'nushell','oil', 'powershell', 'xonsh' or 'zsh' [was: %v]", shell)
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
	for _, c := range cmd.Commands() {
		if c.Name() == "_carapace" {
			return
		}
	}
	cmd.AddCommand(&cobra.Command{
		Use:    "_carapace",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			logger.Println(os.Args) // TODO replace last with '' if empty

			if len(args) == 0 {
				if s, err := Gen(cmd).Snippet(ps.DetermineShell()); err != nil {
					fmt.Fprintln(io.MultiWriter(os.Stderr, logger.Writer()), err.Error())
				} else {
					fmt.Fprintln(io.MultiWriter(os.Stdout, logger.Writer()), s)
				}
			} else {
				if len(args) == 1 {
					if s, err := Gen(cmd).Snippet(args[0]); err != nil {
						fmt.Fprintln(io.MultiWriter(os.Stderr, logger.Writer()), err.Error())
					} else {
						fmt.Fprintln(io.MultiWriter(os.Stdout, logger.Writer()), s)
					}
				} else {
					shell := args[0]
					id := args[1]
					current := args[len(args)-1]
					previous := args[len(args)-2]

					targetCmd, targetArgs, err := findTarget(cmd, args)
					context := Context{CallbackValue: current, Args: targetArgs}
					if err != nil {
						fmt.Fprintln(io.MultiWriter(os.Stdout, logger.Writer()), ActionMessage(err.Error()).Invoke(context).value(shell, current))
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
						} else {
							if len(context.Args) > 0 {
								context.Args = context.Args[:len(context.Args)-1] // current word being completed is a positional so remove it from context.Args
							}

							targetAction = findAction(targetCmd, targetArgs)
							if targetCmd.HasAvailableSubCommands() && len(targetArgs) <= 1 {
								subcommandA := actionSubcommands(targetCmd).Invoke(context)
								targetAction = targetAction.Invoke(context).Merge(subcommandA).ToA()
							}
						}
						fmt.Fprintln(io.MultiWriter(os.Stdout, logger.Writer()), targetAction.Invoke(context).value(shell, current))
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
	}

	lastArg := targetArgs[len(targetArgs)-1]
	if strings.HasSuffix(lastArg, " ") { // TODO is this still correct/needed?
		return storage.getPositional(targetCmd, len(targetArgs))
	}
	return storage.getPositional(targetCmd, len(targetArgs)-1)
}

func findTarget(cmd *cobra.Command, args []string) (*cobra.Command, []string, error) {
	origArg := []string{}
	if len(args) > 3 {
		origArg = args[3:]
	}
	return common.TraverseLenient(cmd, origArg)
}

// IsCallback returns true if current program invocation is a callback
func IsCallback() bool {
	return len(os.Args) > 1 && os.Args[1] == "_carapace"
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
			logger = log.New(logfileWriter, ps.DetermineShell()+" ", log.Flags()|Lmsgprefix)
			//logger = log.New(logfileWriter, determineShell()+" ", log.Flags()|log.Lmsgprefix)
		}
	}
	return
}

type testingT interface {
	Error(args ...interface{})
}

// Test verifies the configuration (e.g. flag name exists)
//   func TestCarapace(t *testing.T) {
//       carapace.Test(t)
//   }
func Test(t testingT) {
	for _, e := range storage.check() {
		testingT(t).Error(e)
	}
}
