package cmd

import (
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
	rootCmd.Flags().StringP("chdir", "C", "", "change work directory")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringP("persistentFlag", "p", "", "Help message for persistentFlag")
	rootCmd.PersistentFlags().String("persistentFlag2", "", "Help message for persistentFlag2")
	rootCmd.Flag("persistentFlag").NoOptDefVal = "defaultValue" // no argument required

	rootCmd.Flags().StringArrayP("array", "a", []string{}, "multiflag")

	carapace.Gen(rootCmd).FlagCompletion(carapace.ActionMap{
		"chdir":           carapace.ActionDirectories(),
		"persistentFlag":  carapace.ActionValues("p1", "p2", "p3"),
		"persistentFlag2": carapace.ActionValues("p4", "p5", "p6"),
	})

	rootCmd.AddGroup(
		&cobra.Group{ID: "main", Title: "Main Commands"},
		&cobra.Group{ID: "modifier", Title: "Modifier Commands"},
		&cobra.Group{ID: "plugin", Title: "Plugin Commands"},
	)

	carapace.Gen(rootCmd).PreRun(func(cmd *cobra.Command, args []string) {
		pluginCmd := &cobra.Command{
			Use:     "plugin",
			Short:   "dynamic plugin command",
			GroupID: "plugin",
			Run:     func(cmd *cobra.Command, args []string) {},
		}

		carapace.Gen(pluginCmd).PositionalCompletion(
			carapace.ActionValues("pl1", "pluginArg1"),
		)

		cmd.AddCommand(pluginCmd)
	})

	carapace.Gen(rootCmd).PreInvoke(func(cmd *cobra.Command, flag *pflag.Flag, action carapace.Action) carapace.Action {
		return action.Chdir(rootCmd.Flag("chdir").Value.String())
	})
}
