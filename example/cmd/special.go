package cmd

import (
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

var specialCmd = &cobra.Command{
	Use:   "special",
	Short: "",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	carapace.Gen(specialCmd).Standalone()

	specialCmd.Flags().CountP("count", "c", "count flag")
	specialCmd.Flags().Bool("exclusive1", false, "mutually exclusive flag")
	specialCmd.Flags().Bool("exclusive2", false, "mutually exclusive flag")
	specialCmd.Flags().StringP("optarg", "o", "", "optional argument")

	specialCmd.Flag("optarg").NoOptDefVal = " "

	specialCmd.MarkFlagsMutuallyExclusive("exclusive1", "exclusive2")

	rootCmd.AddCommand(specialCmd)

	carapace.Gen(specialCmd).FlagCompletion(carapace.ActionMap{
		"optarg": carapace.ActionValues("optarg1", "optarg2", "optarg3"),
	})

	carapace.Gen(specialCmd).PositionalCompletion(
		carapace.ActionValues(`p1 & < > ' " { } $ # | ? ( ) ;  [ ] * \ `+"`", "positional1"),
	)
}
