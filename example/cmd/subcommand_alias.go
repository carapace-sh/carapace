package cmd

import (
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

var subcommand_aliasCmd = &cobra.Command{
	Use:     "alias",
	Short:   "subcommand with alias",
	Aliases: []string{"a1", "a2"},
	Run:     func(cmd *cobra.Command, args []string) {},
}

func init() {
	carapace.Gen(subcommand_aliasCmd).Standalone()

	subcommandCmd.AddCommand(subcommand_aliasCmd)
}
