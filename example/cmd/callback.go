package cmd

import (
	"fmt"
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

var callbackCmd = &cobra.Command{
	Use:   "callback",
	Short: "callback example",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("callbackCmd called with args: %v\n", args)
		fmt.Printf("callback flag is: %v\n", cmd.Flag("callback").Value)
	},
}

func init() {
	rootCmd.AddCommand(callbackCmd)
	callbackCmd.Flags().StringP("callback", "c", "", "Help message for callback")

	carapace.Gen(callbackCmd).FlagCompletion(carapace.ActionMap{
      "callback": carapace.ActionCallback(func(args []string) carapace.Action {
			return carapace.ActionValues("cb1", "cb2", "cb3")
		}),
	})

	carapace.Gen(callbackCmd).PositionalCompletion(
      carapace.ActionCallback(func(args []string) carapace.Action {
			return carapace.ActionValues("callback1", "callback2")
		}),
	)
}
