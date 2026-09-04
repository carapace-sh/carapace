// Package spec provides spec file generation for use with carapace-bin
package spec

import (
	"bytes"

	"github.com/carapace-sh/carapace/internal/pflagfork"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

// Spec generates the spec file.
func Spec(cmd *cobra.Command) string {
	m := &bytes.Buffer{}
	enc := yaml.NewEncoder(m)
	enc.SetIndent(2)
	_ = enc.Encode(command(cmd))
	return "# yaml-language-server: $schema=https://carapace.sh/schemas/command.json\n" + m.String()
}

// Flag holds metadata for a single flag.
type Flag struct {
	Description string
	Nargs       int
	Default     string
}

// Extended is used for flags that have additional metadata beyond a description.
type Extended struct {
	Description string `yaml:"description,omitempty"`
	Nargs       int    `yaml:"nargs,omitempty"`
	Default     string `yaml:"default,omitempty"`
}

// FlagSet is a map of flag definitions to their metadata.
type FlagSet map[string]Flag

func (fs FlagSet) MarshalYAML() (any, error) {
	m := make(map[string]any)
	for k, f := range fs {
		if f.Nargs != 0 || f.Default != "" {
			m[k] = Extended{
				Description: f.Description,
				Nargs:       f.Nargs,
				Default:     f.Default,
			}
		} else {
			m[k] = f.Description
		}
	}
	return m, nil
}

// Command represents a command with its flags, subcommands, and metadata.
type Command struct {
	Name            string            `yaml:"name"`
	Aliases         []string          `yaml:"aliases,omitempty"`
	Description     string            `yaml:"description,omitempty"`
	Group           string            `yaml:"group,omitempty"`
	Hidden          bool              `yaml:"hidden,omitempty"`
	Parsing         string            `yaml:"parsing,omitempty"`
	Flags           FlagSet           `yaml:"flags,omitempty"`
	PersistentFlags FlagSet           `yaml:"persistentflags,omitempty"`
	ExclusiveFlags  [][]string        `yaml:"exclusiveflags,omitempty"`
	Run             string            `yaml:"run,omitempty"`
	Completion      struct {
		Flag          map[string][]string `yaml:"flag,omitempty"`
		Positional    [][]string          `yaml:"positional,omitempty"`
		PositionalAny []string            `yaml:"positionalany,omitempty"`
		Dash          [][]string          `yaml:"dash,omitempty"`
		DashAny       []string            `yaml:"dashany,omitempty"`
	} `yaml:"completion,omitempty"`
	Commands      []Command `yaml:"commands,omitempty"`
	Documentation struct {
		Command       string            `yaml:"command,omitempty"`
		Flag          map[string]string `yaml:"flag,omitempty"`
		Positional    []string          `yaml:"positional,omitempty"`
		PositionalAny string            `yaml:"positionalany,omitempty"`
		Dash          []string          `yaml:"dash,omitempty"`
		DashAny       string            `yaml:"dashany,omitempty"`
	} `yaml:"documentation,omitempty"`
	Examples map[string]string `yaml:"examples,omitempty"`
}

func command(cmd *cobra.Command) Command {
	c := Command{
		Name:            cmd.Use,
		Description:     cmd.Short,
		Aliases:         cmd.Aliases,
		Group:           cmd.GroupID,
		Hidden:          cmd.Hidden,
		Flags:           make(FlagSet),
		PersistentFlags: make(FlagSet),
		Commands:        make([]Command, 0),
	}

	if cmd.Long != "" {
		c.Documentation.Command = cmd.Long
	}

	// TODO mutually exclusive flags

	cmd.LocalFlags().VisitAll(func(flag *pflag.Flag) {
		if cmd.PersistentFlags().Lookup(flag.Name) != nil {
			return
		}

		f := pflagfork.Flag{Flag: flag}
		prefix := pflagfork.FlagSet{FlagSet: cmd.Flags()}.Prefix()
		c.Flags[f.Definition(prefix)] = Flag{
			Description: f.Usage,
			Nargs:       f.Nargs(),
			Default:     flag.DefValue,
		}
	})

	cmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		f := pflagfork.Flag{Flag: flag}
		prefix := pflagfork.FlagSet{FlagSet: cmd.Flags()}.Prefix()
		c.PersistentFlags[f.Definition(prefix)] = Flag{
			Description: f.Usage,
			Nargs:       f.Nargs(),
			Default:     flag.DefValue,
		}
	})

	for _, subcmd := range cmd.Commands() {
		if subcmd.Name() != "_carapace" && subcmd.Deprecated == "" {
			c.Commands = append(c.Commands, command(subcmd))
		}
	}

	return c
}
