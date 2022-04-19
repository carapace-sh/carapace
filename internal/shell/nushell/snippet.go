// Package nushell provides Nushell completion
package nushell

import (
	"fmt"

	"github.com/rsteube/carapace/internal/common"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

/*

Example of git complete from https://github.com/nushell/nushell/blob/main/docs/sample_config/default_config.nu

# This is a simplified version of completions for git branches and git remotes
def "nu-complete git branches" [] {
  ^git branch | lines | each { |line| $line | str find-replace '\* ' '' | str trim }
}

def "nu-complete git remotes" [] {
  ^git remote | lines | each { |line| $line | str trim }
}

# Check out git branches and files
export extern "git checkout" [
  ...targets: string@"nu-complete git branches"   # name of the branch or files to checkout
  --conflict: string                              # conflict style (merge or diff3)
  --detach(-d)                                    # detach HEAD at named commit
  --force(-f)                                     # force checkout (throw away local modifications)
  --guess                                         # second guess 'git checkout <no-such-branch>' (default)
  --ignore-other-worktrees                        # do not check if another worktree is holding the given ref
  --ignore-skip-worktree-bits                     # do not limit pathspecs to sparse entries only
  --merge(-m)                                     # perform a 3-way merge with the new branch
  --orphan: string                                # new unparented branch
  --ours(-2)                                      # checkout our version for unmerged files
  --overlay                                       # use overlay mode (default)
  --overwrite-ignore                              # update ignored files (default)
  --patch(-p)                                     # select hunks interactively
  --pathspec-from-file: string                    # read pathspec from file
  --progress                                      # force progress reporting
  --quiet(-q)                                     # suppress progress reporting
  --recurse-submodules: string                    # control recursive updating of submodules
  --theirs(-3)                                    # checkout their version for unmerged files
  --track(-t)                                     # set upstream info for new branch
  -b: string                                      # create and checkout a new branch
  -B: string                                      # create/reset and checkout a branch
  -l                                              # create reflog for new branch
]

*/

type command struct {
	Name            string
	Short           string
	Long            string    `json:",omitempty"`
	Aliases         []string  `json:",omitempty"`
	Commands        []command `json:",omitempty"`
	LocalFlags      []flag    `json:",omitempty"`
	PersistentFlags []flag    `json:",omitempty"`
}

type flag struct {
	Longhand    string `json:",omitempty"`
	Shorthand   string `json:",omitempty"`
	Usage       string
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
		Usage:       f.Usage,
		Type:        f.Value.Type(),
		NoOptDefVal: noOptDefVal,
	}
}

func convert(cmd *cobra.Command) command {
	c := command{
		Name:    cmd.Name(),
		Short:   cmd.Short,
		Long:    cmd.Long,
		Aliases: cmd.Aliases,
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

	cmd.HasHelpSubCommands()
	return c
}

func exportCommand(doExport bool, fullCmd string, c command) string {
	str := ""
	if doExport {
		str = str + "export "
	}
	str = str + fmt.Sprintf("extern \"%v\" [\n", fullCmd)
	str = fmt.Sprintf("%v  firstarg?\n", str)

	lflags := c.LocalFlags
	for _, f := range lflags {
		shortFlag := f.Shorthand
		longFlag := f.Longhand
		if shortFlag != "" {
			str = str + fmt.Sprintf("  -%v: %v  #%v\n", shortFlag, f.Type, f.Usage)
		}
		if longFlag != "" {
			str = str + fmt.Sprintf("  --%v: %v  #%v\n", longFlag, f.Type, f.Usage)
		}

	}
	str = str + fmt.Sprintf("]\n")

	for _, subCmd := range c.Commands {
		myCmd := subCmd.Name
		newFullCmd := fmt.Sprintf("%v %v", fullCmd, myCmd)
		str = str + exportCommand(doExport, newFullCmd, subCmd)
	}

	return str
}

// Snippet creates the nushell completion script
func Snippet(cmd *cobra.Command) string {

	//TODO support modules
	defineModule := false

	str := ""
	cmdName := cmd.Name()
	if defineModule {
		str = str + fmt.Sprintf("module %v {\n", cmdName)
	}

	c := convert(cmd)
	str = str + exportCommand(defineModule, cmdName, c)

	if defineModule {
		str = str + fmt.Sprintf("}\n")
	}
	return str
}
