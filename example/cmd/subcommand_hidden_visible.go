package cmd

import (
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

var subcommand_hidden_visibleCmd = &cobra.Command{
	Use:   "visible",
	Short: "visible subcommand of a hidden command",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	carapace.Gen(subcommand_hidden_visibleCmd).Standalone()

	subcommand_hiddenCmd.AddCommand(subcommand_hidden_visibleCmd)
}
