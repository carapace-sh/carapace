package cmd

import (
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

var batchCmd = &cobra.Command{
	Use:   "batch",
	Short: "batch example",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(batchCmd)

	carapace.Gen(batchCmd).PositionalCompletion(
		carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			return carapace.Batch(
				carapace.ActionValues("A", "B"),
				carapace.ActionValues("C", "D"),
				carapace.ActionValues("E", "F"),
			).Invoke(c).Merge().ToA()
		}),
	)
}
