package cmd

import (
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

var conditionCmd = &cobra.Command{
	Use:   "condition",
	Short: "condition example",
	Long:  `Simple condition examples.`,
}

func init() {
	rootCmd.AddCommand(conditionCmd)

	conditionCmd.Flags().StringP("required", "r", "", "required flag")

	carapace.Gen(conditionCmd).FlagCompletion(carapace.ActionMap{
		"required": carapace.ActionValues("valid", "invalid"),
	})

	carapace.Gen(conditionCmd).PositionalCompletion(
		carapace.ActionCallback(func(c carapace.Context) (result carapace.Action) {
			if conditionCmd.Flag("required").Value.String() == "valid" {
				result = carapace.ActionValues("condition fulfilled")
			} else {
				result = carapace.ActionMessage("flag --required must be set to valid: " + conditionCmd.Flag("required").Value.String())
			}
			return
		}),
	)
}
