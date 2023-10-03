package cmd

import (
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

var subcommand_groupCmd = &cobra.Command{
	Use:     "group",
	Short:   "subcommand with group",
	GroupID: "group",
	Run:     func(cmd *cobra.Command, args []string) {},
}

func init() {
	carapace.Gen(subcommand_groupCmd).Standalone()

	subcommandCmd.AddCommand(subcommand_groupCmd)
}
