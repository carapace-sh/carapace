package cmd

import (
	"fmt"

	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

var interspersedCmd = &cobra.Command{
	Use:   "interspersed",
	Short: "interspersed example",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "#%v", args)
	},
}

func init() {
	carapace.Gen(interspersedCmd).Standalone()

	interspersedCmd.Flags().BoolP("bool", "b", false, "bool flag")
	interspersedCmd.Flags().StringP("string", "s", "", "string flag")

	interspersedCmd.Flags().SetInterspersed(false)

	rootCmd.AddCommand(interspersedCmd)

	carapace.Gen(interspersedCmd).PositionalCompletion(
		carapace.ActionValues("p1", "positional1"),
		carapace.ActionValues("p2", "positional2"),
	)

	carapace.Gen(interspersedCmd).DashCompletion(
		carapace.ActionValues("d1", "dash1"),
		carapace.ActionValues("d2", "dash2"),
	)
}
