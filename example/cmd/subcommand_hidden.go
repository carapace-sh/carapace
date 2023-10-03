package cmd

import (
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

var subcommand_hiddenCmd = &cobra.Command{
	Use:    "hidden",
	Short:  "hidden subcommand",
	Hidden: true,
	Run:    func(cmd *cobra.Command, args []string) {},
}

func init() {
	carapace.Gen(subcommand_hiddenCmd).Standalone()

	subcommandCmd.AddCommand(subcommand_hiddenCmd)
}
