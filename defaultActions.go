package carapace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rsteube/carapace/internal/common"
	"github.com/rsteube/carapace/internal/config"
	"github.com/rsteube/carapace/internal/shell/export"
	"github.com/rsteube/carapace/pkg/style"
	"github.com/rsteube/carapace/third_party/github.com/acarl005/stripansi"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// ActionCallback invokes a go function during completion
func ActionCallback(callback CompletionCallback) Action {
	return Action{callback: callback}
}

// ActionExecCommand invokes given command and transforms its output using given function on success or returns ActionMessage with the first line of stderr if available.
//   carapace.ActionExecCommand("git", "remote")(func(output []byte) carapace.Action {
//     lines := strings.Split(string(output), "\n")
//     return carapace.ActionValues(lines[:len(lines)-1]...)
//   })
func ActionExecCommand(name string, arg ...string) func(f func(output []byte) Action) Action {
	return func(f func(output []byte) Action) Action {
		return ActionCallback(func(c Context) Action {
			var stdout, stderr bytes.Buffer
			cmd := c.Command(name, arg...)
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			if err := cmd.Run(); err != nil {
				if firstLine := strings.SplitN(stderr.String(), "\n", 2)[0]; strings.TrimSpace(firstLine) != "" {
					return ActionMessage(stripansi.Strip(firstLine))
				}
				return ActionMessage(err.Error())
			}
			return f(stdout.Bytes())
		})
	}
}

// ActionImport parses the json output from export as Action
//   carapace.Gen(rootCmd).PositionalAnyCompletion(
//   	carapace.ActionCallback(func(c carapace.Context) carapace.Action {
//   		args := []string{"_carapace", "export", ""}
//   		args = append(args, c.Args...)
//   		args = append(args, c.CallbackValue)
//   		return carapace.ActionExecCommand("command", args...)(func(output []byte) carapace.Action {
//   			return carapace.ActionImport(output)
//   		})
//   	}),
//   )
func ActionImport(output []byte) Action {
	return ActionCallback(func(c Context) Action {
		var e export.Export
		if err := json.Unmarshal(output, &e); err != nil {
			return ActionMessage(err.Error())
		}
		a := actionRawValues(e.RawValues...)
		if e.Nospace {
			a = a.NoSpace()
		}
		return a
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
				// TODO verify
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

// ActionDirectories completes directories
func ActionDirectories() Action {
	return ActionCallback(func(c Context) Action {
		return actionPath([]string{""}, true).Invoke(c).ToMultiPartsA("/").noSpace(true).StyleF(func(s string) string {
			if abs, err := c.Abs(s); err == nil {
				return style.ForPath(abs)
			}
			return ""
		})
	})
}

// ActionFiles completes files with optional suffix filtering
func ActionFiles(suffix ...string) Action {
	return ActionCallback(func(c Context) Action {
		return actionPath(suffix, false).Invoke(c).ToMultiPartsA("/").noSpace(true).StyleF(func(s string) string {
			if abs, err := c.Abs(s); err == nil {
				return style.ForPath(abs)
			}
			return ""
		})
	})
}

func actionPath(fileSuffixes []string, dirOnly bool) Action {
	return ActionCallback(func(c Context) Action {
		abs, err := c.Abs(c.CallbackValue)
		if err != nil {
			return ActionMessage(err.Error())
		}

		displayFolder := filepath.Dir(c.CallbackValue)
		if displayFolder == "." {
			displayFolder = ""
		} else if !strings.HasSuffix(displayFolder, "/") {
			displayFolder = displayFolder + "/"
		}

		actualFolder := filepath.Dir(abs)
		files, err := ioutil.ReadDir(actualFolder)
		if err != nil {
			return ActionMessage(err.Error())
		}

		showHidden := !strings.HasSuffix(abs, "/") && strings.HasPrefix(filepath.Base(abs), ".")

		vals := make([]string, 0, len(files)*2)
		for _, file := range files {
			if !showHidden && strings.HasPrefix(file.Name(), ".") {
				continue
			}

			resolvedFile := file
			if resolved, err := filepath.EvalSymlinks(actualFolder + file.Name()); err == nil {
				if stat, err := os.Stat(resolved); err == nil {
					resolvedFile = stat
				}
			}

			if resolvedFile.IsDir() {
				vals = append(vals, displayFolder+file.Name()+"/", style.ForPath(filepath.Clean(actualFolder+"/"+file.Name()+"/")))
			} else if !dirOnly {
				if len(fileSuffixes) == 0 {
					fileSuffixes = []string{""}
				}
				for _, suffix := range fileSuffixes {
					if strings.HasSuffix(file.Name(), suffix) {
						vals = append(vals, displayFolder+file.Name(), style.ForPath(filepath.Clean(actualFolder+"/"+file.Name())))
						break
					}
				}
			}
		}
		if strings.HasPrefix(c.CallbackValue, "./") {
			return ActionStyledValues(vals...).Invoke(Context{}).Prefix("./").ToA()
		}
		return ActionStyledValues(vals...)
	})
}

// ActionValues completes arbitrary keywords (values)
func ActionValues(values ...string) Action {
	return ActionCallback(func(c Context) Action {
		vals := make([]common.RawValue, 0)
		for _, val := range values {
			vals = append(vals, common.RawValue{Value: val, Display: val, Description: "", Style: style.Default})
		}
		return actionRawValues(vals...)
	})
}

// ActionStyledValues is like ActionValues but also accepts a style
func ActionStyledValues(values ...string) Action {
	return ActionCallback(func(c Context) Action {
		vals := make([]common.RawValue, 0)
		for index, val := range values {
			if index%2 == 0 {
				vals = append(vals, common.RawValue{Value: val, Display: val, Description: "", Style: values[index+1]})
			}
		}
		return actionRawValues(vals...)
	})
}

// ActionValuesDescribed completes arbitrary key (values) with an additional description (value, description pairs)
func ActionValuesDescribed(values ...string) Action {
	return ActionCallback(func(c Context) Action {
		vals := make([]common.RawValue, 0)
		for index, val := range values {
			if index%2 == 0 {
				vals = append(vals, common.RawValue{Value: val, Display: val, Description: values[index+1], Style: style.Default})
			}
		}
		return actionRawValues(vals...)
	})
}

// ActionStyledValuesDescribed is like ActionValues but also accepts a style
func ActionStyledValuesDescribed(values ...string) Action {
	return ActionCallback(func(c Context) Action {
		vals := make([]common.RawValue, 0)
		for index, val := range values {
			if index%3 == 0 {
				vals = append(vals, common.RawValue{Value: val, Display: val, Description: values[index+1], Style: values[index+2]})
			}
		}
		return actionRawValues(vals...)
	})
}

func actionRawValues(rawValues ...common.RawValue) Action {
	return Action{
		rawValues: rawValues,
	}
}

// ActionMessage displays a help messages in places where no completions can be generated
func ActionMessage(msg string) Action {
	return ActionCallback(func(c Context) Action {
		return ActionValuesDescribed("_", "", "ERR", msg).
			Invoke(c).Prefix(c.CallbackValue).ToA(). // needs to be prefixed with current callback value to not be filtered out
			noSpace(true).skipCache(true)
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

		return callback(c).Invoke(c).Prefix(prefix).ToA().noSpace(true)
	})
}

func actionSubcommands(cmd *cobra.Command) Action {
	vals := make([]string, 0)
	for _, subcommand := range cmd.Commands() {
		if !subcommand.Hidden && subcommand.Deprecated == "" {
			vals = append(vals, subcommand.Name(), subcommand.Short)
			for _, alias := range subcommand.Aliases {
				vals = append(vals, alias, subcommand.Short)
			}
		}
	}
	return ActionValuesDescribed(vals...)
}

func actionFlags(cmd *cobra.Command) Action {
	return ActionCallback(func(c Context) Action {
		re := regexp.MustCompile("^-(?P<shorthand>[^-=]+)")
		isShorthandSeries := re.MatchString(c.CallbackValue)

		vals := make([]string, 0)
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			if f.Deprecated != "" {
				return // skip deprecated flags
			}

			if f.Changed &&
				!strings.Contains(f.Value.Type(), "Slice") &&
				!strings.Contains(f.Value.Type(), "Array") &&
				f.Value.Type() != "count" {
				return // don't repeat flag
			}

			if isShorthandSeries {
				if f.Shorthand != "" && f.ShorthandDeprecated == "" {
					for _, shorthand := range c.CallbackValue[1:] {
						if shorthandFlag := cmd.Flags().ShorthandLookup(string(shorthand)); shorthandFlag != nil && shorthandFlag.Value.Type() != "bool" && shorthandFlag.Value.Type() != "count" && shorthandFlag.NoOptDefVal == "" {
							return // abort shorthand flag series if a previous one is not bool or count and requires an argument (no default value)
						}
					}
					vals = append(vals, f.Shorthand, f.Usage)
				}
			} else {
				if !common.IsShorthandOnly(f) {
					if opts.LongShorthand {
						vals = append(vals, "-"+f.Name, f.Usage)
					} else {
						vals = append(vals, "--"+f.Name, f.Usage)
					}
				}
				if f.Shorthand != "" && f.ShorthandDeprecated == "" {
					vals = append(vals, "-"+f.Shorthand, f.Usage)
				}
			}
		})

		if isShorthandSeries {
			return ActionValuesDescribed(vals...).Invoke(c).Prefix(c.CallbackValue).ToA().noSpace(true)
		}
		return ActionValuesDescribed(vals...)
	})
}

// ActionStyleConfig completes style configuration
//   carapace.Value=blue
//   carapace.Description=magenta
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
					return ActionStyledValuesDescribed(fields...).Invoke(c).Suffix("=").ToA()

				default:
					return ActionValues()
				}
			})
		case 1:
			return ActionMultiParts(",", func(c Context) Action {
				return ActionStyles(c.Parts...).Invoke(c).Filter(c.Parts).ToA()
			})
		default:
			return ActionValues()
		}
	})
}

// Actionstyles completes styles
//   blue
//   bg-magenta
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

		if !fg {
			batch = append(batch, ActionStyledValues(
				style.Black, style.Of(append(styles, style.Black)...),
				style.Red, style.Of(append(styles, style.Red)...),
				style.Green, style.Of(append(styles, style.Green)...),
				style.Yellow, style.Of(append(styles, style.Yellow)...),
				style.Blue, style.Of(append(styles, style.Blue)...),
				style.Magenta, style.Of(append(styles, style.Magenta)...),
				style.Cyan, style.Of(append(styles, style.Cyan)...),
				style.White, style.Of(append(styles, style.White)...),
				style.Gray, style.Of(append(styles, style.Gray)...),

				style.BrightBlack, style.Of(append(styles, style.BrightBlack)...),
				style.BrightRed, style.Of(append(styles, style.BrightRed)...),
				style.BrightGreen, style.Of(append(styles, style.BrightGreen)...),
				style.BrightYellow, style.Of(append(styles, style.BrightYellow)...),
				style.BrightBlue, style.Of(append(styles, style.BrightBlue)...),
				style.BrightMagenta, style.Of(append(styles, style.BrightMagenta)...),
				style.BrightCyan, style.Of(append(styles, style.BrightCyan)...),
				style.BrightWhite, style.Of(append(styles, style.BrightWhite)...),
			))

			if strings.HasPrefix(c.CallbackValue, "color") {
				for i := 0; i <= 255; i++ {
					batch = append(batch, ActionStyledValues(
						fmt.Sprintf("color%v", i), style.Of(append(styles, style.XTerm256Color(uint8(i)))...),
					))
				}
			} else {
				batch = append(batch, ActionStyledValues("color", style.Of(styles...)))
			}
		}

		if !bg {
			batch = append(batch, ActionStyledValues(
				style.BgBlack, style.Of(append(styles, style.BgBlack)...),
				style.BgRed, style.Of(append(styles, style.BgRed)...),
				style.BgGreen, style.Of(append(styles, style.BgGreen)...),
				style.BgYellow, style.Of(append(styles, style.BgYellow)...),
				style.BgBlue, style.Of(append(styles, style.BgBlue)...),
				style.BgMagenta, style.Of(append(styles, style.BgMagenta)...),
				style.BgCyan, style.Of(append(styles, style.BgCyan)...),
				style.BgWhite, style.Of(append(styles, style.BgWhite)...),

				style.BgBrightBlack, style.Of(append(styles, style.BgBrightBlack)...),
				style.BgBrightRed, style.Of(append(styles, style.BgBrightRed)...),
				style.BgBrightGreen, style.Of(append(styles, style.BgBrightGreen)...),
				style.BgBrightYellow, style.Of(append(styles, style.BgBrightYellow)...),
				style.BgBrightBlue, style.Of(append(styles, style.BgBrightBlue)...),
				style.BgBrightMagenta, style.Of(append(styles, style.BgBrightMagenta)...),
				style.BgBrightCyan, style.Of(append(styles, style.BgBrightCyan)...),
				style.BgBrightWhite, style.Of(append(styles, style.BgBrightWhite)...),
			))

			if strings.HasPrefix(c.CallbackValue, "bg-color") {
				for i := 0; i <= 255; i++ {
					batch = append(batch, ActionStyledValues(
						fmt.Sprintf("bg-color%v", i), style.Of(append(styles, "bg-"+style.XTerm256Color(uint8(i)))...),
					))
				}
			} else {
				batch = append(batch, ActionStyledValues("bg-color", style.Of(styles...)))
			}
		}

		batch = append(batch, ActionStyledValues(
			style.Bold, style.Of(append(styles, style.Bold)...),
			style.Dim, style.Of(append(styles, style.Dim)...),
			style.Italic, style.Of(append(styles, style.Italic)...),
			style.Underlined, style.Of(append(styles, style.Underlined)...),
			style.Blink, style.Of(append(styles, style.Blink)...),
			style.Inverse, style.Of(append(styles, style.Inverse)...),
		))

		return batch.ToA()
	})
}
