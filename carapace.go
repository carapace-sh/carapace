// Package carapace is a command argument completion generator for spf13/cobra
package carapace

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/rsteube/carapace/internal/common"
	"github.com/rsteube/carapace/internal/config"
	"github.com/rsteube/carapace/internal/pflagfork"
	"github.com/rsteube/carapace/internal/shell/bash"
	"github.com/rsteube/carapace/internal/shell/bash_ble"
	"github.com/rsteube/carapace/internal/shell/elvish"
	"github.com/rsteube/carapace/internal/shell/export"
	"github.com/rsteube/carapace/internal/shell/fish"
	"github.com/rsteube/carapace/internal/shell/ion"
	"github.com/rsteube/carapace/internal/shell/nushell"
	"github.com/rsteube/carapace/internal/shell/oil"
	"github.com/rsteube/carapace/internal/shell/powershell"
	"github.com/rsteube/carapace/internal/shell/spec"
	"github.com/rsteube/carapace/internal/shell/tcsh"
	"github.com/rsteube/carapace/internal/shell/xonsh"
	"github.com/rsteube/carapace/internal/shell/zsh"
	"github.com/rsteube/carapace/internal/uid"
	"github.com/rsteube/carapace/pkg/ps"
	"github.com/rsteube/carapace/pkg/style"
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

	cobra.OnInitialize(func() {
		if opts.BridgeCompletion {
			registerValidArgsFunction(cmd)
			registerFlagCompletion(cmd)
		}
	})

	return &Carapace{
		cmd: cmd,
	}
}

// PreRun sets a function to be run before completion (use on rootCmd).
func (c Carapace) PreRun(f func(cmd *cobra.Command, args []string)) {
	if completionCmd, _, err := c.cmd.Find([]string{"_carapace"}); err == nil {
		completionCmd.PreRun = func(cmd *cobra.Command, args []string) {
			if len(args) > 2 { // skip script generation
				f(c.cmd, args[2:])
			}
		}
	}
}

// PreInvoke sets a function to alter actions before they are invoked (use on rootCmd).
func (c Carapace) PreInvoke(f func(cmd *cobra.Command, flag *pflag.Flag, action Action) Action) {
	if entry := storage.get(c.cmd); entry.preinvoke != nil {
		_f := entry.preinvoke
		entry.preinvoke = func(cmd *cobra.Command, flag *pflag.Flag, action Action) Action {
			return f(cmd, flag, _f(cmd, flag, action)) // TODO verify if this is correct
		}
	} else {
		entry.preinvoke = f
	}
}

// PositionalCompletion defines completion for positional arguments using a list of Actions
func (c Carapace) PositionalCompletion(action ...Action) {
	storage.get(c.cmd).positional = action
}

// PositionalAnyCompletion defines completion for any positional arguments not already defined
func (c Carapace) PositionalAnyCompletion(action Action) {
	storage.get(c.cmd).positionalAny = action
}

// DashCompletion defines completion for positional arguments after dash (`--`) using a list of Actions
func (c Carapace) DashCompletion(action ...Action) {
	storage.get(c.cmd).dash = action
}

// DashAnyCompletion defines completion for any positional arguments after dash (`--`) not already defined
func (c Carapace) DashAnyCompletion(action Action) {
	storage.get(c.cmd).dashAny = action
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
	c.cmd.CompletionOptions = cobra.CompletionOptions{
		DisableDefaultCmd: true,
	}
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
	if shell == "" {
		shell = ps.DetermineShell()
	}
	shellSnippets := map[string]func(cmd *cobra.Command) string{
		"bash":       bash.Snippet,
		"bash-ble":   bash_ble.Snippet,
		"export":     export.Snippet,
		"fish":       fish.Snippet,
		"elvish":     elvish.Snippet,
		"ion":        ion.Snippet,
		"nushell":    nushell.Snippet,
		"oil":        oil.Snippet,
		"powershell": powershell.Snippet,
		"spec":       spec.Snippet,
		"tcsh":       tcsh.Snippet,
		"xonsh":      xonsh.Snippet,
		"zsh":        zsh.Snippet,
	}
	if s, ok := shellSnippets[shell]; ok {
		return s(c.cmd.Root()), nil
	}

	expected := make([]string, 0)
	for key := range shellSnippets {
		expected = append(expected, key)
	}
	sort.Strings(expected)
	return "", fmt.Errorf("expected one of '%v' [was: %v]", strings.Join(expected, "', '"), shell)
}

func lookupFlag(cmd *cobra.Command, arg string) (flag *pflag.Flag) {
	nameOrShorthand := strings.TrimLeft(strings.SplitN(arg, "=", 2)[0], "-")

	if strings.HasPrefix(arg, "--") {
		flag = cmd.Flags().Lookup(nameOrShorthand)
	} else if strings.HasPrefix(arg, "-") && len(nameOrShorthand) > 0 {
		if pflagfork.IsPosix(cmd.Flags()) {
			flag = cmd.Flags().ShorthandLookup(string(nameOrShorthand[len(nameOrShorthand)-1]))
		} else {
			flag = cmd.Flags().ShorthandLookup(nameOrShorthand)
		}
	}
	return
}

func addCompletionCommand(cmd *cobra.Command) {
	for _, c := range cmd.Commands() {
		if c.Name() == "_carapace" {
			return
		}
	}
	carapaceCmd := &cobra.Command{
		Use:    "_carapace",
		Hidden: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) > 2 && strings.HasPrefix(args[2], "_") {
				cmd.Hidden = false
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			logger.Println(os.Args) // TODO replace last with '' if empty
			if s, err := complete(cmd, args); err != nil {
				fmt.Fprintln(io.MultiWriter(cmd.OutOrStderr(), logger.Writer()), err.Error())
			} else {
				fmt.Fprintln(io.MultiWriter(cmd.OutOrStdout(), logger.Writer()), s)
			}
		},
		FParseErrWhitelist: cobra.FParseErrWhitelist{
			UnknownFlags: true,
		},
		DisableFlagParsing: true,
	}
	cmd.AddCommand(carapaceCmd)
	Carapace{carapaceCmd}.PositionalCompletion(
		ActionStyledValues(
			"bash", "#d35673",
			"bash-ble", "#c2039a",
			"elvish", "#ffd6c9",
			"export", style.Default,
			"fish", "#7ea8fc",
			"ion", "#0e5d6d",
			"nushell", "#29d866",
			"oil", "#373a36",
			"powershell", "#e8a16f",
			"spec", style.Default,
			"tcsh", "#412f09",
			"xonsh", "#a8ffa9",
			"zsh", "#efda53",
		),
		ActionValues(cmd.Root().Name()),
	)
	Carapace{carapaceCmd}.PositionalAnyCompletion(
		ActionCallback(func(c Context) Action {
			args := []string{"_carapace", "export", ""}
			args = append(args, c.Args[2:]...)
			args = append(args, c.CallbackValue)
			return ActionExecCommand(uid.Executable(), args...)(func(output []byte) Action {
				if string(output) == "" {
					return ActionValues()
				}
				return ActionImport(output)
			})
		}),
	)

	styleCmd := &cobra.Command{
		Use:  "style",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
	carapaceCmd.AddCommand(styleCmd)

	styleSetCmd := &cobra.Command{
		Use:  "set",
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			for _, arg := range args {
				if splitted := strings.SplitN(arg, "=", 2); len(splitted) == 2 {
					if err := style.Set(splitted[0], splitted[1]); err != nil {
						fmt.Fprint(cmd.ErrOrStderr(), err.Error())
					}
				} else {
					fmt.Fprintf(cmd.ErrOrStderr(), "invalid format: '%v'", arg)
				}
			}
		},
	}
	styleCmd.AddCommand(styleSetCmd)
	Carapace{styleSetCmd}.PositionalAnyCompletion(
		ActionStyleConfig(),
	)
}

func complete(cmd *cobra.Command, args []string) (string, error) {
	if len(args) == 0 {
		if s, err := Gen(cmd).Snippet(ps.DetermineShell()); err != nil {
			fmt.Fprintln(io.MultiWriter(cmd.OutOrStderr(), logger.Writer()), err.Error())
		} else {
			fmt.Fprintln(io.MultiWriter(cmd.OutOrStdout(), logger.Writer()), s)
		}
	} else {
		if len(args) == 1 {
			s, err := Gen(cmd).Snippet(args[0])
			if err != nil {
				return "", err
			}
			return s, nil
		} else {
			shell := args[0]
			current := args[len(args)-1]
			previous := args[len(args)-2]

			if err := config.Load(); err != nil {
				return ActionMessage("failed to load config: "+err.Error()).Invoke(Context{CallbackValue: current}).value(shell, current), nil
			}

			targetCmd, targetArgs, err := findTarget(cmd, args)
			if err != nil {
				return ActionMessage(err.Error()).Invoke(Context{CallbackValue: current}).value(shell, current), nil
			}

			wd, err := os.Getwd()
			if err != nil {
				return ActionMessage(err.Error()).Invoke(Context{CallbackValue: current}).value(shell, current), nil
			}
			context := Context{
				CallbackValue: current,
				Args:          targetArgs,
				Env:           os.Environ(),
				Dir:           wd,
			}
			if err != nil {
				return ActionMessage(err.Error()).Invoke(context).value(shell, current), nil
			}

			// TODO needs more cleanup and tests
			var targetAction Action
			if flag := lookupFlag(targetCmd, previous); !targetCmd.DisableFlagParsing && flag != nil && flag.NoOptDefVal == "" && !common.IsDash(targetCmd) { // previous arg is a flag and needs a value
				targetAction = storage.getFlag(targetCmd, flag.Name)
			} else if !targetCmd.DisableFlagParsing && strings.HasPrefix(current, "-") && !common.IsDash(targetCmd) { // assume flag
				if strings.Contains(current, "=") { // complete value for optarg flag
					if flag := lookupFlag(targetCmd, current); flag != nil && flag.NoOptDefVal != "" {
						a := storage.getFlag(targetCmd, flag.Name)
						splitted := strings.SplitN(current, "=", 2)
						context.CallbackValue = splitted[1]
						current = strings.Replace(current, "=", opts.OptArgDelimiter, 1)                  // revert (potentially) overridden optarg divider for `.value()` invocation below
						targetAction = a.Invoke(context).Prefix(splitted[0] + opts.OptArgDelimiter).ToA() // prefix with (potentially) overridden optarg delimiter
					}
				} else { // complete flagnames
					targetAction = actionFlags(targetCmd)
				}
			} else {
				if len(context.Args) > 0 {
					context.Args = context.Args[:len(context.Args)-1] // current word being completed is a positional so remove it from context.Args
				}

				if common.IsDash(targetCmd) {
					targetAction = findAction(targetCmd, targetArgs[targetCmd.ArgsLenAtDash():])
				} else {
					targetAction = findAction(targetCmd, targetArgs)
					if targetCmd.HasAvailableSubCommands() && len(targetArgs) <= 1 {
						subcommandA := actionSubcommands(targetCmd).Invoke(context)
						targetAction = targetAction.Invoke(context).Merge(subcommandA).ToA()
					}
				}
			}
			return targetAction.Invoke(context).value(shell, current), nil
		}
	}
	return "", nil // TODO
}

func findAction(targetCmd *cobra.Command, targetArgs []string) Action {
	// TODO handle Action not found
	if len(targetArgs) == 0 {
		return storage.getPositional(targetCmd, 0)
	}
	return storage.getPositional(targetCmd, len(targetArgs)-1)
}

func findTarget(cmd *cobra.Command, args []string) (*cobra.Command, []string, error) {
	origArg := []string{}
	if len(args) > 2 {
		origArg = args[2:]
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

// Test verifies the configuration (e.g. flag name exists)
//
//	func TestCarapace(t *testing.T) {
//	    carapace.Test(t)
//	}
func Test(t interface{ Error(args ...interface{}) }) {
	for _, e := range storage.check() {
		t.Error(e)
	}
}
