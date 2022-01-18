package export

import (
	"encoding/json"
	"strings"

	"github.com/rsteube/carapace/internal/common"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type command struct {
	Name            string
	Description     string
	Url             string    `json:",omitempty"`
	Aliases         []string  `json:",omitempty"`
	Commands        []command `json:",omitempty"`
	LocalFlags      []flag    `json:",omitempty"`
	PersistentFlags []flag    `json:",omitempty"`
}

type flag struct {
	Longhand    string `json:",omitempty"`
	Shorthand   string `json:",omitempty"`
	Description string
	Type        string
	NoOptDefVal string `json:",omitempty"`
}

func convertFlag(f *pflag.Flag) flag {
	longhand := ""
	if !common.IsShorthandOnly(f) {
		longhand = f.Name
	}

	noOptDefVal := ""
	if f.Value.Type() != "bool" {
		noOptDefVal = f.NoOptDefVal
	}
	return flag{
		Longhand:    longhand,
		Shorthand:   f.Shorthand,
		Description: f.Usage,
		Type:        f.Value.Type(),
		NoOptDefVal: noOptDefVal,
	}
}

func convert(cmd *cobra.Command) command {
	c := command{
		Name:        cmd.Name(),
		Description: cmd.Short,
		Aliases:     cmd.Aliases,
	}

	if strings.HasPrefix(strings.ToLower(cmd.Long), "https://") ||
		strings.HasPrefix(strings.ToLower(cmd.Long), "http://") {
		c.Url = cmd.Long
	}

	lflags := make([]flag, 0)
	cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
		lflags = append(lflags, convertFlag(f))
	})
	c.LocalFlags = lflags

	pflags := make([]flag, 0)
	cmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		pflags = append(pflags, convertFlag(f))
	})
	c.PersistentFlags = pflags

	subcommands := make([]command, 0)
	for _, s := range cmd.Commands() {
		if !s.Hidden {
			subcommands = append(subcommands, convert(s))
		}
	}
	c.Commands = subcommands
	return c
}

// Snippet exports the command structure as json
func Snippet(cmd *cobra.Command) string {
	out, err := json.Marshal(convert(cmd))
	if err == nil {
		return string(out)
	}
	return err.Error()
}
