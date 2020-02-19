package cmd

import (
	zsh "github.com/rsteube/cobra-zsh-gen"
	"github.com/spf13/cobra"
)

var actionCmd = &cobra.Command{
	Use:   "action",
	Short: "action example",
}

func init() {
	rootCmd.AddCommand(actionCmd)

	actionCmd.Flags().StringP("files", "f", "", "files flag")
	actionCmd.Flags().StringP("groups", "g", "", "groups flag")
	actionCmd.Flags().String("hosts", "", "hosts flag")
	actionCmd.Flags().StringP("message", "m", "", "message flag")
	actionCmd.Flags().StringP("net_interfaces", "n", "", "net_interfaces flag")
	actionCmd.Flags().StringP("options", "o", "", "options flag")
	actionCmd.Flags().String("path_files", "", "path_files flag")
	actionCmd.Flags().StringP("users", "u", "", "users flag")
	actionCmd.Flags().StringP("values", "v", "", "values flag")
	actionCmd.Flags().StringP("values_described", "d", "", "values with description flag")
	actionCmd.Flags().StringP("custom", "c", "", "custom flag")

	zsh.Gen(actionCmd).FlagCompletion(zsh.ActionMap{
		"files":            zsh.ActionFiles("*.go"),
		"groups":           zsh.ActionGroups(),
		"hosts":            zsh.ActionHosts(),
		"message":          zsh.ActionMessage("message example"),
		"net_interfaces":   zsh.ActionNetInterfaces(),
		"options":          zsh.ActionOptions(),
		"path_files":       zsh.ActionPathFiles(""),
		"users":            zsh.ActionUsers(),
		"values":           zsh.ActionValues("values", "example"),
		"values_described": zsh.ActionValuesDescribed("values", "valueDescription", "example", "exampleDescription"),
		"custom":           zsh.Action{Value: "_most_recent_file 2"},
	})
}
