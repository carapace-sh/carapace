package cmd

import (
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "example",
	Short:   "example completion",
	Version: "example",
}

// Execute executes cmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringP("persistentFlag", "p", "", "Help message for persistentFlag")
	rootCmd.PersistentFlags().String("persistentFlag2", "", "Help message for persistentFlag2")
	rootCmd.Flag("persistentFlag").NoOptDefVal = "defaultValue" // no argument required

	rootCmd.Flags().StringArrayP("array", "a", []string{}, "multiflag")

	carapace.Gen(rootCmd).FlagCompletion(carapace.ActionMap{
		"persistentFlag":  carapace.ActionValues("p1", "p2", "p3"),
		"persistentFlag2": carapace.ActionValues("p4", "p5", "p6"),
	})

	rootCmd.AddGroup(
		&cobra.Group{ID: "main", Title: "Main Commands"},
		&cobra.Group{ID: "modifier", Title: "Modifier Commands"},
		&cobra.Group{ID: "test", Title: "Test Commands"},
	)

	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		println(rootCmd.Flag("persistentFlag").Value.String())
	}
}
