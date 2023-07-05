package cmd

import (
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "group example",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	carapace.Gen(groupCmd).Standalone()

	rootCmd.AddCommand(groupCmd)

	groupCmd.AddGroup(
		&cobra.Group{ID: "main", Title: "Main Commands"},
		&cobra.Group{ID: "setup", Title: "Setup Commands"},
	)

	run := func(cmd *cobra.Command, args []string) {}
	groupCmd.AddCommand(
		&cobra.Command{Use: "sub1", GroupID: "main", Run: run},
		&cobra.Command{Use: "sub2", GroupID: "main", Run: run},
		&cobra.Command{Use: "sub3", GroupID: "setup", Run: run},
		&cobra.Command{Use: "sub4", GroupID: "setup", Run: run},
		&cobra.Command{Use: "sub5", Run: run},
	)
}
