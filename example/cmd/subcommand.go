package cmd

import (
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

var subcommandCmd = &cobra.Command{
	Use:   "subcommand",
	Short: "subcommand example",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	carapace.Gen(subcommandCmd).Standalone()

	subcommandCmd.AddGroup(
		&cobra.Group{ID: "group", Title: ""},
	)

	rootCmd.AddCommand(subcommandCmd)
}
