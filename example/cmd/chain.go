package cmd

import (
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

var chainCmd = &cobra.Command{
	Use:                "chain",
	Short:              "shorthand chain",
	Run:                func(cmd *cobra.Command, args []string) {},
	DisableFlagParsing: true,
}

func init() {
	carapace.Gen(chainCmd).Standalone()

	rootCmd.AddCommand(chainCmd)

	carapace.Gen(chainCmd).PositionalAnyCompletion(
		carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			cmd := &cobra.Command{}
			carapace.Gen(cmd).Standalone()

			cmd.Flags().CountP("c", "c", "")
			cmd.Flags().BoolP("b", "b", false, "")
			cmd.Flags().StringP("v", "v", "", "")
			cmd.Flags().StringP("o", "o", "", "")

			cmd.Flag("o").NoOptDefVal = " "

			carapace.Gen(cmd).FlagCompletion(carapace.ActionMap{
				"v": carapace.ActionValues("val1", "val2"),
				"o": carapace.ActionValues("opt1", "opt2"),
			})

			carapace.Gen(cmd).PositionalCompletion(
				carapace.ActionValues("p1", "positional1"),
			)

			return carapace.ActionExecute(cmd)
		}),
	)

}
