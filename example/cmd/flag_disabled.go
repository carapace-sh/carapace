package cmd

import (
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

var flag_disabledCmd = &cobra.Command{
	Use:                "disabled",
	Short:              "flag parsing disabled",
	DisableFlagParsing: true,
	Run:                func(cmd *cobra.Command, args []string) {},
}

func init() {
	carapace.Gen(flag_disabledCmd).Standalone()

	flagCmd.AddCommand(flag_disabledCmd)

	carapace.Gen(flag_disabledCmd).PositionalCompletion(
		carapace.ActionValues("-p1", "positional1"),
		carapace.ActionValues("p2", "--positional2"),
	)
}
