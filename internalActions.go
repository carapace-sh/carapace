package carapace

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/rsteube/carapace/internal/common"
	"github.com/rsteube/carapace/internal/pflagfork"
	"github.com/rsteube/carapace/pkg/style"
	"github.com/spf13/cobra"
)

func actionPath(fileSuffixes []string, dirOnly bool) Action {
	return ActionCallback(func(c Context) Action {
		abs, err := c.Abs(c.Value)
		if err != nil {
			return ActionMessage(err.Error())
		}

		displayFolder := filepath.Dir(c.Value)
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
				vals = append(vals, displayFolder+file.Name()+"/", style.ForPath(filepath.Clean(actualFolder+"/"+file.Name()+"/"), c))
			} else if !dirOnly {
				if len(fileSuffixes) == 0 {
					fileSuffixes = []string{""}
				}
				for _, suffix := range fileSuffixes {
					if strings.HasSuffix(file.Name(), suffix) {
						vals = append(vals, displayFolder+file.Name(), style.ForPath(filepath.Clean(actualFolder+"/"+file.Name()), c))
						break
					}
				}
			}
		}
		if strings.HasPrefix(c.Value, "./") {
			return ActionStyledValues(vals...).Invoke(Context{}).Prefix("./").ToA()
		}
		return ActionStyledValues(vals...)
	}).Tag("files").NoSpace('/')
}

func actionFlags(cmd *cobra.Command) Action {
	return ActionCallback(func(c Context) Action {
		flagSet := pflagfork.FlagSet{FlagSet: cmd.Flags()}
		isShorthandSeries := flagSet.IsShorthandSeries(c.Value)

		vals := make([]string, 0)
		flagSet.VisitAll(func(f *pflagfork.Flag) {
			switch {
			case f.Deprecated != "":
				return // skip deprecated flags
			case f.Changed && !f.IsRepeatable():
				return // don't repeat flag
			case flagSet.IsMutuallyExclusive(f.Flag):
				return // skip flag of group already set
			}

			if isShorthandSeries {
				if f.Shorthand != "" && f.ShorthandDeprecated == "" {
					for _, shorthand := range c.Value[1:] {
						if shorthandFlag := cmd.Flags().ShorthandLookup(string(shorthand)); shorthandFlag != nil && shorthandFlag.Value.Type() != "bool" && shorthandFlag.Value.Type() != "count" && shorthandFlag.NoOptDefVal == "" {
							return // abort shorthand flag series if a previous one is not bool or count and requires an argument (no default value)
						}
					}
					vals = append(vals, f.Shorthand, f.Usage, f.Style())
				}
			} else {
				switch f.Mode() {
				case pflagfork.NameAsShorthand:
					vals = append(vals, "-"+f.Name, f.Usage, f.Style())
				case pflagfork.Default:
					vals = append(vals, "--"+f.Name, f.Usage, f.Style())
				}

				if f.Shorthand != "" && f.ShorthandDeprecated == "" {
					vals = append(vals, "-"+f.Shorthand, f.Usage, f.Style())
				}
			}
		})

		if isShorthandSeries {
			return ActionStyledValuesDescribed(vals...).Prefix(c.Value).NoSpace('*')
		}
		return ActionStyledValuesDescribed(vals...).MultiParts(".") // multiparts completion for flags grouped with `.`
	}).Tag("flags")
}

func actionSubcommands(cmd *cobra.Command) Action {
	return ActionCallback(func(c Context) Action {
		batch := Batch()
		for _, subcommand := range cmd.Commands() {
			if !subcommand.Hidden && subcommand.Deprecated == "" {
				group := common.Group{Cmd: subcommand}
				batch = append(batch, ActionStyledValuesDescribed(subcommand.Name(), subcommand.Short, group.Style()).Tag(group.Tag()))
				for _, alias := range subcommand.Aliases {
					batch = append(batch, ActionStyledValuesDescribed(alias, subcommand.Short, group.Style()).Tag(group.Tag()))
				}
			}
		}
		return batch.ToA()
	})
}
