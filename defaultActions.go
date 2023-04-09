package carapace

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/rsteube/carapace/internal/common"
	"github.com/rsteube/carapace/internal/config"
	"github.com/rsteube/carapace/internal/export"
	"github.com/rsteube/carapace/pkg/style"
	"github.com/rsteube/carapace/third_party/github.com/acarl005/stripansi"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// ActionCallback invokes a go function during completion.
func ActionCallback(callback CompletionCallback) Action {
	return Action{callback: callback}
}

// ActionExecCommand invokes given command and transforms its output using given function on success or returns ActionMessage with the first line of stderr if available.
//
//	carapace.ActionExecCommand("git", "remote")(func(output []byte) carapace.Action {
//	  lines := strings.Split(string(output), "\n")
//	  return carapace.ActionValues(lines[:len(lines)-1]...)
//	})
func ActionExecCommand(name string, arg ...string) func(f func(output []byte) Action) Action {
	return func(f func(output []byte) Action) Action {
		return ActionExecCommandE(name, arg...)(func(output []byte, err error) Action {
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					if firstLine := strings.SplitN(string(exitErr.Stderr), "\n", 2)[0]; strings.TrimSpace(firstLine) != "" {
						err = errors.New(firstLine)
					}
				}
				return ActionMessage(err.Error())
			}
			return f(output)
		})
	}
}

// ActionExecCommandE is like ActionExecCommand but with custom error handling.
//
//	carapace.ActionExecCommandE("supervisorctl", "--configuration", path, "status")(func(output []byte, err error) carapace.Action {
//		if err != nil {
//			const NOT_RUNNING = 3
//			if exitErr, ok := err.(*exec.ExitError); !ok || exitErr.ExitCode() != NOT_RUNNING {
//				return carapace.ActionMessage(err.Error())
//			}
//		}
//		return carapace.ActionValues("success")
//	})
func ActionExecCommandE(name string, arg ...string) func(f func(output []byte, err error) Action) Action {
	return func(f func(output []byte, err error) Action) Action {
		return ActionCallback(func(c Context) Action {
			var stdout, stderr bytes.Buffer
			cmd := c.Command(name, arg...)
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			if err := cmd.Run(); err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					exitErr.Stderr = stderr.Bytes() // seems this needs to be set manually due to stdout being collected?
				}
				return f(stdout.Bytes(), err)
			}
			return f(stdout.Bytes(), nil)
		})
	}
}

// ActionImport parses the json output from export as Action
//
//	carapace.Gen(rootCmd).PositionalAnyCompletion(
//		carapace.ActionCallback(func(c carapace.Context) carapace.Action {
//			args := []string{"_carapace", "export", ""}
//			args = append(args, c.Args...)
//			args = append(args, c.CallbackValue)
//			return carapace.ActionExecCommand("command", args...)(func(output []byte) carapace.Action {
//				return carapace.ActionImport(output)
//			})
//		}),
//	)
func ActionImport(output []byte) Action {
	return ActionCallback(func(c Context) Action {
		var e export.Export
		if err := json.Unmarshal(output, &e); err != nil {
			return ActionMessage(err.Error())
		}
		return Action{
			rawValues: e.Values,
			meta:      e.Meta,
		}
	})
}

// ActionExecute executes completion on an internal command
// TODO example
func ActionExecute(cmd *cobra.Command) Action {
	return ActionCallback(func(c Context) Action {
		args := []string{"_carapace", "export", cmd.Name()}
		args = append(args, c.Args...)
		args = append(args, c.CallbackValue)
		cmd.SetArgs(args)

		Gen(cmd).PreInvoke(func(cmd *cobra.Command, flag *pflag.Flag, action Action) Action {
			return ActionCallback(func(_c Context) Action {
				_c.Env = c.Env
				_c.Dir = c.Dir
				return action.Invoke(_c).ToA()
			})
		})

		var stdout, stderr bytes.Buffer
		cmd.SetOut(&stdout)
		cmd.SetErr(&stderr)

		if err := cmd.Execute(); err != nil {
			return ActionMessage(err.Error())
		}
		return ActionImport(stdout.Bytes())
	})
}

// ActionDirectories completes directories.
func ActionDirectories() Action {
	return ActionCallback(func(c Context) Action {
		return actionPath([]string{""}, true).Invoke(c).ToMultiPartsA("/").StyleF(style.ForPath)
	}).Tag("directories")
}

// ActionFiles completes files with optional suffix filtering.
func ActionFiles(suffix ...string) Action {
	return ActionCallback(func(c Context) Action {
		return actionPath(suffix, false).Invoke(c).ToMultiPartsA("/").StyleF(style.ForPath)
	}).Tag("files")
}

// ActionValues completes arbitrary keywords (values).
func ActionValues(values ...string) Action {
	return ActionCallback(func(c Context) Action {
		vals := make([]common.RawValue, 0, len(values))
		for _, val := range values {
			vals = append(vals, common.RawValue{Value: val, Display: val})
		}
		return Action{rawValues: vals}
	})
}

// ActionStyledValues is like ActionValues but also accepts a style.
func ActionStyledValues(values ...string) Action {
	return ActionCallback(func(c Context) Action {
		if length := len(values); length%2 != 0 {
			return ActionMessage("invalid amount of arguments [ActionStyledValues]: %v", length)
		}

		vals := make([]common.RawValue, 0, len(values)/2)
		for i := 0; i < len(values); i += 2 {
			vals = append(vals, common.RawValue{Value: values[i], Display: values[i], Style: values[i+1]})
		}
		return Action{rawValues: vals}
	})
}

// ActionValuesDescribed completes arbitrary key (values) with an additional description (value, description pairs).
func ActionValuesDescribed(values ...string) Action {
	return ActionCallback(func(c Context) Action {
		if length := len(values); length%2 != 0 {
			return ActionMessage("invalid amount of arguments [ActionValuesDescribed]: %v", length)
		}

		vals := make([]common.RawValue, 0, len(values)/2)
		for i := 0; i < len(values); i += 2 {
			vals = append(vals, common.RawValue{Value: values[i], Display: values[i], Description: values[i+1]})
		}
		return Action{rawValues: vals}
	})
}

// ActionStyledValuesDescribed is like ActionValues but also accepts a style.
func ActionStyledValuesDescribed(values ...string) Action {
	return ActionCallback(func(c Context) Action {
		if length := len(values); length%3 != 0 {
			return ActionMessage("invalid amount of arguments [ActionStyledValuesDescribed]: %v", length)
		}

		vals := make([]common.RawValue, 0, len(values)/3)
		for i := 0; i < len(values); i += 3 {
			vals = append(vals, common.RawValue{Value: values[i], Display: values[i], Description: values[i+1], Style: values[i+2]})
		}
		return Action{rawValues: vals}
	})
}

// ActionMessage displays a help messages in places where no completions can be generated.
func ActionMessage(msg string, args ...interface{}) Action {
	return ActionCallback(func(c Context) Action {
		if len(args) > 0 {
			msg = fmt.Sprintf(msg, args...)
		}
		a := ActionValues().NoSpace()
		a.meta.Messages.Add(stripansi.Strip(msg))
		return a
	})
}

// ActionMultiParts completes multiple parts of words separately where each part is separated by some char (CallbackValue is set to the currently completed part during invocation)
func ActionMultiParts(divider string, callback func(c Context) Action) Action {
	return ActionCallback(func(c Context) Action {
		index := strings.LastIndex(c.CallbackValue, string(divider))
		prefix := ""
		if len(divider) == 0 {
			prefix = c.CallbackValue
			c.CallbackValue = ""
		} else if index != -1 {
			prefix = c.CallbackValue[0 : index+len(divider)]
			c.CallbackValue = c.CallbackValue[index+len(divider):] // update CallbackValue to only contain the currently completed part
		}
		parts := strings.Split(prefix, string(divider))
		if len(parts) > 0 && len(divider) > 0 {
			parts = parts[0 : len(parts)-1]
		}
		c.Parts = parts

		nospace := '*'
		if runes := []rune(divider); len(runes) > 0 {
			nospace = runes[len(runes)-1]
		}
		return callback(c).Invoke(c).Prefix(prefix).ToA().NoSpace(nospace)
	})
}

// ActionStyleConfig completes style configuration
//
//	carapace.Value=blue
//	carapace.Description=magenta
func ActionStyleConfig() Action {
	return ActionMultiParts("=", func(c Context) Action {
		switch len(c.Parts) {
		case 0:
			return ActionMultiParts(".", func(c Context) Action {
				switch len(c.Parts) {
				case 0:
					return ActionValues(config.GetStyleConfigs()...).Invoke(c).Suffix(".").ToA()

				case 1:
					fields, err := config.GetStyleFields(c.Parts[0])
					if err != nil {
						return ActionMessage(err.Error())
					}
					batch := Batch()
					for _, field := range fields {
						batch = append(batch, ActionStyledValuesDescribed(field.Name, field.Description, field.Style).Tag(field.Tag))
					}
					return batch.Invoke(c).Merge().Suffix("=").ToA()

				default:
					return ActionValues()
				}
			})
		case 1:
			return ActionMultiParts(",", func(c Context) Action {
				return ActionStyles(c.Parts...).Invoke(c).Filter(c.Parts).ToA().NoSpace()
			})
		default:
			return ActionValues()
		}
	})
}

// Actionstyles completes styles
//
//	blue
//	bg-magenta
func ActionStyles(styles ...string) Action {
	return ActionCallback(func(c Context) Action {
		fg := false
		bg := false

		for _, s := range styles {
			if strings.HasPrefix(s, "bg-") {
				bg = true
			}
			if s == style.Black ||
				s == style.Red ||
				s == style.Green ||
				s == style.Yellow ||
				s == style.Blue ||
				s == style.Magenta ||
				s == style.Cyan ||
				s == style.White ||
				s == style.Gray ||
				s == style.BrightBlack ||
				s == style.BrightRed ||
				s == style.BrightGreen ||
				s == style.BrightYellow ||
				s == style.BrightBlue ||
				s == style.BrightMagenta ||
				s == style.BrightCyan ||
				s == style.BrightWhite ||
				strings.HasPrefix(s, "#") ||
				strings.HasPrefix(s, "color") ||
				strings.HasPrefix(s, "fg-") {
				fg = true
			}
		}

		batch := Batch()
		_s := func(s string) string {
			return style.Of(append(styles, s)...)
		}

		if !fg {
			batch = append(batch, ActionStyledValues(
				style.Black, _s(style.Black),
				style.Red, _s(style.Red),
				style.Green, _s(style.Green),
				style.Yellow, _s(style.Yellow),
				style.Blue, _s(style.Blue),
				style.Magenta, _s(style.Magenta),
				style.Cyan, _s(style.Cyan),
				style.White, _s(style.White),
				style.Gray, _s(style.Gray),

				style.BrightBlack, _s(style.BrightBlack),
				style.BrightRed, _s(style.BrightRed),
				style.BrightGreen, _s(style.BrightGreen),
				style.BrightYellow, _s(style.BrightYellow),
				style.BrightBlue, _s(style.BrightBlue),
				style.BrightMagenta, _s(style.BrightMagenta),
				style.BrightCyan, _s(style.BrightCyan),
				style.BrightWhite, _s(style.BrightWhite),
			))

			if strings.HasPrefix(c.CallbackValue, "color") {
				for i := 0; i <= 255; i++ {
					batch = append(batch, ActionStyledValues(
						fmt.Sprintf("color%v", i), _s(style.XTerm256Color(uint8(i))),
					))
				}
			} else {
				batch = append(batch, ActionStyledValues("color", style.Of(styles...)))
			}
		}

		if !bg {
			batch = append(batch, ActionStyledValues(
				style.BgBlack, _s(style.BgBlack),
				style.BgRed, _s(style.BgRed),
				style.BgGreen, _s(style.BgGreen),
				style.BgYellow, _s(style.BgYellow),
				style.BgBlue, _s(style.BgBlue),
				style.BgMagenta, _s(style.BgMagenta),
				style.BgCyan, _s(style.BgCyan),
				style.BgWhite, _s(style.BgWhite),

				style.BgBrightBlack, _s(style.BgBrightBlack),
				style.BgBrightRed, _s(style.BgBrightRed),
				style.BgBrightGreen, _s(style.BgBrightGreen),
				style.BgBrightYellow, _s(style.BgBrightYellow),
				style.BgBrightBlue, _s(style.BgBrightBlue),
				style.BgBrightMagenta, _s(style.BgBrightMagenta),
				style.BgBrightCyan, _s(style.BgBrightCyan),
				style.BgBrightWhite, _s(style.BgBrightWhite),
			))

			if strings.HasPrefix(c.CallbackValue, "bg-color") {
				for i := 0; i <= 255; i++ {
					batch = append(batch, ActionStyledValues(
						fmt.Sprintf("bg-color%v", i), _s("bg-"+style.XTerm256Color(uint8(i))),
					))
				}
			} else {
				batch = append(batch, ActionStyledValues("bg-color", style.Of(styles...)))
			}
		}

		batch = append(batch, ActionStyledValues(
			style.Bold, _s(style.Bold),
			style.Dim, _s(style.Dim),
			style.Italic, _s(style.Italic),
			style.Underlined, _s(style.Underlined),
			style.Blink, _s(style.Blink),
			style.Inverse, _s(style.Inverse),
		))

		return batch.ToA()
	}).Tag("styles")
}

// ActionExecutables completes executable files from PATH
//
//	nvim
//	chmod
func ActionExecutables() Action {
	return ActionCallback(func(c Context) Action {
		// TODO allow additional descriptions to be registered somewhere for carapace-bin (key, value,...)
		batch := Batch()
		manDescriptions := manDescriptions(c)
		dirs := strings.Split(os.Getenv("PATH"), string(os.PathListSeparator))
		for i := len(dirs) - 1; i >= 0; i-- {
			batch = append(batch, actionDirectoryExecutables(dirs[i], c.CallbackValue, manDescriptions))
		}
		return batch.ToA()
	}).Tag("path executables")
}

func actionDirectoryExecutables(dir string, prefix string, manDescriptions map[string]string) Action {
	return ActionCallback(func(c Context) Action {
		if files, err := os.ReadDir(dir); err == nil {
			vals := make([]string, 0)
			for _, f := range files {
				if strings.HasPrefix(f.Name(), prefix) {
					if info, err := f.Info(); err == nil && !f.IsDir() && isExecAny(info.Mode()) {
						vals = append(vals, f.Name(), manDescriptions[f.Name()], style.ForPath(dir+"/"+f.Name(), c))
					}
				}
			}
			return ActionStyledValuesDescribed(vals...)
		}
		return ActionValues()
	})
}

func isExecAny(mode os.FileMode) bool {
	return mode&0111 != 0
}

func manDescriptions(c Context) (descriptions map[string]string) {
	descriptions = make(map[string]string)
	if !strings.HasPrefix(c.CallbackValue, ".") && !strings.HasPrefix(c.CallbackValue, "~") && !strings.HasPrefix(c.CallbackValue, "/") {
		if output, err := c.Command("man", "--names-only", "-k", "^"+c.CallbackValue).Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			r := regexp.MustCompile(`^(?P<name>[^ ]+) [^-]+- (?P<description>.*)$`)
			for _, line := range lines {
				if matches := r.FindStringSubmatch(line); len(matches) > 2 {
					descriptions[matches[1]] = matches[2]
				}
			}
		}
	}
	return
}

// ActionPositional completes positional arguments for given command ignoring `--` (dash).
// TODO: experimental - likely gives issues with preinvoke (does not have the full args)
//
//	carapace.Gen(cmd).DashAnyCompletion(
//		carapace.ActionPositional(cmd),
//	)
func ActionPositional(cmd *cobra.Command) Action {
	return ActionCallback(func(c Context) Action {
		if cmd.ArgsLenAtDash() < 0 {
			return ActionMessage("only allowed for dash arguments [ActionPositional]")
		}

		c.Args = cmd.Flags().Args()
		entry := storage.get(cmd)

		a := entry.positionalAny
		if index := len(c.Args); index < len(entry.positional) {
			a = entry.positional[len(c.Args)]
		}
		return a.Invoke(c).ToA()
	})
}
