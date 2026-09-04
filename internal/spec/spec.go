// Package spec provides spec file generation for use with carapace-bin
package spec

import (
	"bytes"
	"sort"
	"strings"

	"github.com/carapace-sh/carapace/internal/pflagfork"
	"github.com/carapace-sh/carapace/pkg/command"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

// Spec generates the spec file.
func Spec(cmd *cobra.Command) string {
	m := &bytes.Buffer{}
	enc := yaml.NewEncoder(m)
	enc.SetIndent(2)
	_ = enc.Encode(commandLine(cmd))
	return "# yaml-language-server: $schema=https://carapace.sh/schemas/command.json\n" + m.String()
}

func commandLine(cmd *cobra.Command) command.Command {
	c := command.Command{
		Name:            cmd.Use,
		Description:     cmd.Short,
		Aliases:         cmd.Aliases,
		Group:           cmd.GroupID,
		Hidden:          cmd.Hidden,
		Flags:           make(command.FlagSet),
		PersistentFlags: make(command.FlagSet),
		Commands:        make([]command.Command, 0),
	}

	if cmd.Long != "" {
		c.Documentation.Command = cmd.Long
	}

	c.ExclusiveFlags = exclusiveFlags(cmd)

	cmd.LocalFlags().VisitAll(func(flag *pflag.Flag) {
		if cmd.PersistentFlags().Lookup(flag.Name) != nil {
			return
		}
		c.Flags[flagDefinition(flag, cmd.Flags())] = toFlag(flag)
	})

	cmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		c.PersistentFlags[flagDefinition(flag, cmd.Flags())] = toFlag(flag)
	})

	for _, subcmd := range cmd.Commands() {
		if subcmd.Name() != "_carapace" && subcmd.Deprecated == "" {
			c.Commands = append(c.Commands, commandLine(subcmd))
		}
	}

	return c
}

func flagDefinition(flag *pflag.Flag, flagSet *pflag.FlagSet) string {
	f := pflagfork.Flag{Flag: flag}
	prefix := pflagfork.FlagSet{FlagSet: flagSet}.Prefix()
	return f.Definition(prefix)
}

func toFlag(flag *pflag.Flag) command.Flag {
	f := pflagfork.Flag{Flag: flag}
	mode := f.GetMode()

	var shorthand, longhand string
	switch mode {
	case pflagfork.ShorthandOnly:
		shorthand = f.Shorthand
	case pflagfork.NameAsShorthand:
		shorthand = f.Shorthand
		longhand = f.Name
	default:
		if f.Shorthand != "" {
			shorthand = f.Shorthand
		}
		longhand = f.Name
	}

	var defaultVal string
	if !isZeroDefault(flag) {
		defaultVal = flag.DefValue
	}

	return command.Flag{
		Longhand:        longhand,
		Shorthand:       shorthand,
		Description:     f.Usage,
		NameAsShorthand: mode == pflagfork.NameAsShorthand,
		Repeatable:      f.IsRepeatable(),
		Optarg:          f.IsOptarg() && f.TakesValue(),
		Value:           f.TakesValue(),
		Hidden:          f.Hidden,
		Required:        f.Required(),
		Nargs:           f.Nargs(),
		Default:         defaultVal,
	}
}

func isZeroDefault(flag *pflag.Flag) bool {
	switch flag.Value.Type() {
	case "bool", "count":
		return flag.DefValue == "false" || flag.DefValue == "0"
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64":
		return flag.DefValue == "0"
	case "string":
		return flag.DefValue == ""
	case "stringSlice", "intSlice", "int32Slice", "int64Slice",
		"uintSlice", "float32Slice", "float64Slice",
		"boolSlice", "durationSlice", "stringArray",
		"ipSlice", "ipNetSlice",
		"stringToInt", "stringToInt64", "stringToString":
		return flag.DefValue == "[]"
	case "duration":
		return flag.DefValue == "0s"
	case "ip", "ipNet", "ipMask":
		return flag.DefValue == "<nil>"
	default:
		return flag.DefValue == ""
	}
}

func exclusiveFlags(cmd *cobra.Command) [][]string {
	groups := make(map[string]bool)
	var result [][]string

	collect := func(flag *pflag.Flag) {
		for _, entry := range flag.Annotations["cobra_annotation_mutually_exclusive"] {
			members := strings.Fields(entry)
			sort.Strings(members)
			key := strings.Join(members, "\x00")
			if !groups[key] {
				groups[key] = true
				result = append(result, members)
			}
		}
	}

	cmd.LocalFlags().VisitAll(collect)
	cmd.PersistentFlags().VisitAll(collect)
	return result
}
