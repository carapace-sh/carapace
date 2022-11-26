package cmd

import (
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

var batchCmd = &cobra.Command{
	Use:     "batch",
	Short:   "batch example",
	GroupID: "main",
	Run:     func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(batchCmd)

	carapace.Gen(batchCmd).PositionalCompletion(
		carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			return carapace.Batch(
				carapace.ActionValues("A", "B").Tag("first"),
				carapace.ActionValues("C", "D").Tag("second"),
				carapace.ActionValues("E", "F").TagF(func(value string) string { return "third" }),
			).Invoke(c).Merge().ToA()
		}),
	)
}
