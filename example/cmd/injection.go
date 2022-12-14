package cmd

import (
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

var injectionCmd = &cobra.Command{
	Use:     "injection",
	Short:   "just trying to break things",
	GroupID: "test",
	Run:     func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(injectionCmd)

	carapace.Gen(injectionCmd).PositionalCompletion(
		carapace.ActionValues("$(echo fail)"),
		carapace.ActionValues(`\$(echo fail)`),
		carapace.ActionValues("`echo fail`"),
		carapace.ActionValues(`"; echo fail #`),
		carapace.ActionValues(`"| echo fail #`),
		carapace.ActionValues(`"&& echo fail #`),
		carapace.ActionValues(`\$(echo fail)`),
		carapace.ActionValues(`\`),
		carapace.ActionValues(`LAST POSITIONAL VALUE`),
	)
}
