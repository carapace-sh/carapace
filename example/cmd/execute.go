package cmd

import (
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

var executeCmd = &cobra.Command{
	Use:                "execute",
	Short:              "execute example",
	DisableFlagParsing: true,
	Run:                func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(executeCmd)

	cmd := &cobra.Command{
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		Run: func(cmd *cobra.Command, args []string) {},
	}

	cmd.Flags().Bool("test", false, "")

	carapace.Gen(cmd).PositionalCompletion(
		carapace.ActionValues("one", "two"),
		carapace.ActionValues("three", "four"),
	)

	carapace.Gen(executeCmd).PositionalAnyCompletion(
		carapace.ActionExecute(cmd),
	)
}
