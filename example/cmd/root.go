package cmd

import (
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "example",
	Short: "example completion",
}

// Execute executes cmd.
func Execute() error {
	carapace.Override(carapace.Opts{
		BridgeCompletion: true,
	})
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringP("persistentFlag", "p", "", "Help message for persistentFlag")
	rootCmd.Flag("persistentFlag").NoOptDefVal = "defaultValue" // no argument required

	rootCmd.Flags().StringArrayP("array", "a", []string{}, "multiflag")

	carapace.Gen(rootCmd).FlagCompletion(carapace.ActionMap{
		"persistentFlag": carapace.ActionValues("p1", "p2", "p3"),
	})

	rootCmd.AddGroup(
		&cobra.Group{ID: "main", Title: "Main Commands"},
		&cobra.Group{ID: "modifier", Title: "Modifier Commands"},
		&cobra.Group{ID: "test", Title: "Test Commands"},
	)
}
