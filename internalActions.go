package carapace

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rsteube/carapace/internal/common"
	"github.com/rsteube/carapace/internal/pflagfork"
	"github.com/rsteube/carapace/pkg/style"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

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

func actionFlags(cmd *cobra.Command) Action {
	return ActionCallback(func(c Context) Action {
		re := regexp.MustCompile("^-(?P<shorthand>[^-=]+)")
		isShorthandSeries := re.MatchString(c.CallbackValue) && pflagfork.FlagSet(cmd.Flags()).IsPosix()

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
				if flagstyle := pflagfork.Flag(f).Style(); flagstyle != pflagfork.ShorthandOnly {
					if flagstyle == pflagfork.NameAsShorthand {
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
			return ActionValuesDescribed(vals...).Invoke(c).Prefix(c.CallbackValue).ToA().noSpace("*")
		}
		for i := 0; i < len(vals); i = i + 2 { // TODO experimental - hardcoded multiparts completion if flags are "grouped" with `.`
			if strings.Contains(vals[i], ".") {
				return ActionValuesDescribed(vals...).Invoke(c).ToMultiPartsA(".")
			}
		}
		return ActionValuesDescribed(vals...)
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

func actionRawValues(rawValues ...common.RawValue) Action {
	return Action{
		rawValues: rawValues,
	}
}
