package cmd

import (
	"github.com/carapace-sh/carapace"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var argumentStyleCmd = &cobra.Command{
	Use:   "argumentstyle",
	Short: "test ArgumentStyle configurations",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Flags().Visit(func(f *pflag.Flag) {
			cmd.Println(f.Name, f.Value.String())
		})
	},
}

func init() {
	rootCmd.AddCommand(argumentStyleCmd)

	// AcceptAll (0) - accepts all argument styles (default/backward compatible)
	argumentStyleCmd.Flags().String("accept-all", "", "Accept all styles (default)")

	// AcceptNext only - only accepts `-f arg`
	argumentStyleCmd.Flags().String("accept-next", "", "Accept next arg only")
	argumentStyleCmd.Flag("accept-next").ArgumentStyle = pflag.AcceptNext

	// AcceptDelimited only - only accepts `-f=value`
	argumentStyleCmd.Flags().String("accept-delimited", "", "Accept delimited only")
	argumentStyleCmd.Flag("accept-delimited").ArgumentStyle = pflag.AcceptDelimited

	// AcceptAttached only - only accepts `-fvalue` in POSIX mode
	argumentStyleCmd.Flags().String("accept-attached", "", "Accept attached only")
	argumentStyleCmd.Flag("accept-attached").ArgumentStyle = pflag.AcceptAttached

	// AcceptDelimited | AcceptNext - accepts both delimited and next, but not attached
	argumentStyleCmd.Flags().String("accept-delim-or-next", "", "Accept delimited or next")
	argumentStyleCmd.Flag("accept-delim-or-next").ArgumentStyle = pflag.AcceptDelimited | pflag.AcceptNext

	// AcceptDelimited | AcceptAttached - accepts both delimited and attached, but not next
	argumentStyleCmd.Flags().String("accept-delim-or-attached", "", "Accept delimited or attached")
	argumentStyleCmd.Flag("accept-delim-or-attached").ArgumentStyle = pflag.AcceptDelimited | pflag.AcceptAttached

	// AcceptNext | AcceptAttached - accepts both next and attached, but not delimited
	argumentStyleCmd.Flags().String("accept-next-or-attached", "", "Accept next or attached")
	argumentStyleCmd.Flag("accept-next-or-attached").ArgumentStyle = pflag.AcceptNext | pflag.AcceptAttached

	// Demonstrate with shorthand (v)
	argumentStyleCmd.Flags().StringP("accept-all-s", "a", "", "Accept all styles (shorthand)")
	argumentStyleCmd.Flags().StringP("accept-next-s", "n", "", "Accept next only (shorthand)")
	argumentStyleCmd.Flag("accept-next-s").ArgumentStyle = pflag.AcceptNext
	argumentStyleCmd.Flags().StringP("accept-delimited-s", "d", "", "Accept delimited only (shorthand)")
	argumentStyleCmd.Flag("accept-delimited-s").ArgumentStyle = pflag.AcceptDelimited

	carapace.Gen(argumentStyleCmd).FlagCompletion(carapace.ActionMap{
		"accept-all":             carapace.ActionValues("val1", "val2"),
		"accept-next":            carapace.ActionValues("val1", "val2"),
		"accept-delimited":       carapace.ActionValues("val1", "val2"),
		"accept-attached":        carapace.ActionValues("val1", "val2"),
		"accept-delim-or-next":   carapace.ActionValues("val1", "val2"),
		"accept-delim-or-attached": carapace.ActionValues("val1", "val2"),
		"accept-next-or-attached":  carapace.ActionValues("val1", "val2"),
		"accept-all-s":            carapace.ActionValues("val1", "val2"),
		"accept-next-s":          carapace.ActionValues("val1", "val2"),
		"accept-delimited-s":     carapace.ActionValues("val1", "val2"),
	})
}