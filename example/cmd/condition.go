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

	carapace.Gen(conditionCmd).PositionalCompletion(carapace.ActionCallback(func(args []string) carapace.Action {
		if conditionCmd.Flag("required").Value.String() == "valid" {
			return carapace.ActionValues("condition fulfilled")
		} else {
			return carapace.ActionMessage("flag --required must be set to valid: " + conditionCmd.Flag("required").Value.String())
		}
	}))
}
