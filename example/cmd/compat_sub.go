package cmd

import (
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

var compat_subCmd = &cobra.Command{
	Use:   "sub",
	Short: "",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	carapace.Gen(compat_subCmd).Standalone()

	compatCmd.AddCommand(compat_subCmd)
}
