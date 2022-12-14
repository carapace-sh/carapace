package cmd

import (
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

var modifierCmd = &cobra.Command{
	Use:     "modifier [pos1]",
	Short:   "modifier example",
	GroupID: "modifier",
	Run:     func(cmd *cobra.Command, args []string) {},
}

func init() {
	modifierCmd.Flags().String("batch", "", "Batch()")
	modifierCmd.Flags().String("usage", "", "Usage()")

	rootCmd.AddCommand(modifierCmd)

	carapace.Gen(modifierCmd).FlagCompletion(carapace.ActionMap{
		"batch": carapace.Batch(
			carapace.ActionValues("A", "B"),
			carapace.ActionValues("C", "D"),
			carapace.ActionValues("E", "F"),
		).ToA(),
		"usage": carapace.ActionValues().Usage("explicit flag usage"),
	})

	carapace.Gen(modifierCmd).PositionalCompletion(
		carapace.ActionValues().Usage("explicit positional usage"),
	)
}
