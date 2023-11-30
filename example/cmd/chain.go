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

			cmd.Flags().CountP("count", "c", "")
			cmd.Flags().BoolP("bool", "b", false, "")
			cmd.Flags().StringP("value", "v", "", "")
			cmd.Flags().StringP("optarg", "o", "", "")

			cmd.Flag("optarg").NoOptDefVal = " "

			carapace.Gen(cmd).FlagCompletion(carapace.ActionMap{
				"value":  carapace.ActionValues("val1", "val2"),
				"optarg": carapace.ActionValues("opt1", "opt2"),
			})

			carapace.Gen(cmd).PositionalCompletion(
				carapace.ActionValues("p1", "positional1"),
			)

			return carapace.ActionExecute(cmd)
		}),
	)
}
