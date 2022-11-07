package cmd

import (
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "example-nonposix",
	Short: "nonposix examples",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func Execute() error {
	return rootCmd.Execute()
}
func init() {
	carapace.Gen(rootCmd).Standalone()

	rootCmd.Flags().BoolN("bool-long", "bool-short", false, "BoolN")
	rootCmd.Flags().StringS("delim-colon", "delim-colon", "", "OptargDelimiter ':'")
	rootCmd.Flags().StringS("delim-slash", "delim-slash", "", "OptargDelimiter '/'")
	rootCmd.Flags().CountN("count", "c", "CountN")

	rootCmd.Flag("delim-colon").NoOptDefVal = " "
	rootCmd.Flag("delim-colon").OptargDelimiter = ':'
	rootCmd.Flag("delim-slash").NoOptDefVal = " "
	rootCmd.Flag("delim-slash").OptargDelimiter = '/'

	carapace.Gen(rootCmd).FlagCompletion(carapace.ActionMap{
		"delim-colon": carapace.ActionValues("d1", "d2", "d3"),
		"delim-slash": carapace.ActionValues("d1", "d2", "d3"),
	})

	carapace.Gen(rootCmd).PositionalCompletion(
		carapace.ActionValues("p1", "positional1"),
	)
}
