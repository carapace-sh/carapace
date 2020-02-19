package cmd

import (
	"fmt"
	"github.com/rsteube/cobra-zsh-gen"
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

	zsh.Gen(callbackCmd).FlagCompletion(zsh.ActionMap{
		"callback": zsh.ActionCallback(func(args []string) zsh.Action {
			return zsh.ActionValues("cb1", "cb2", "cb3")
		}),
	})

	zsh.Gen(callbackCmd).PositionalCompletion(
		zsh.ActionCallback(func(args []string) zsh.Action {
			return zsh.ActionValues("callback1", "callback2")
		}),
	)
}
