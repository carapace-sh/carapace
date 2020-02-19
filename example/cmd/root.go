package cmd

import (
	"github.com/rsteube/cobra-zsh-gen"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "example",
	Short: "example completion",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringP("persistentFlag", "p", "", "Help message for persistentFlag")
	rootCmd.Flag("persistentFlag").NoOptDefVal = "defaultValue" // no argument required

	rootCmd.Flags().StringArrayP("array", "a", []string{}, "mulitflag")

	zsh.Gen(rootCmd).FlagCompletion(zsh.ActionMap{
		"toggle": zsh.ActionBool(),
	})
}
