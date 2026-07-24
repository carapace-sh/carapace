package cmd

import (
	"github.com/carapace-sh/carapace"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var customPrefixCmd = &cobra.Command{
	Use:   "customprefix",
	Short: "test custom flag prefix (e.g. '&')",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Flags().Visit(func(f *pflag.Flag) {
			cmd.Println(f.Name, f.Value.String())
		})
	},
}

func init() {
	carapace.Gen(customPrefixCmd).Standalone()
	customPrefixCmd.Flags().SetPrefix('&')

	customPrefixCmd.Flags().BoolN("verbose", "v", false, "verbose output")
	customPrefixCmd.Flags().StringN("name", "n", "", "name flag")
	customPrefixCmd.Flags().CountN("count", "c", "count flag")
	rootCmd.AddCommand(customPrefixCmd)

	carapace.Gen(customPrefixCmd).FlagCompletion(carapace.ActionMap{
		"name": carapace.ActionValues("alpha", "beta", "gamma"),
	})

	carapace.Gen(customPrefixCmd).PositionalCompletion(
		carapace.ActionValues("pos1", "positional1"),
	)
}
