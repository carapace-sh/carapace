package cmd

import (
	zsh "github.com/rsteube/cobra-zsh-gen"
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

	zsh.Gen(conditionCmd).FlagCompletion(zsh.ActionMap{
		"required": zsh.ActionValues("valid", "invalid"),
	})

	zsh.Gen(conditionCmd).PositionalCompletion(zsh.ActionCallback(func(args []string) zsh.Action {
		if conditionCmd.Flag("required").Value.String() == "valid" {
			return zsh.ActionValues("condition fulfilled")
		} else {
			return zsh.ActionMessage("flag --required must be set to valid: " + conditionCmd.Flag("required").Value.String())
		}
	}))
}
