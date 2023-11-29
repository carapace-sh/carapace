package cmd

import (
	"fmt"

	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var rootCmd = &cobra.Command{
	Use:   "example-nonposix",
	Short: "nonposix examples",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Flags().Visit(func(f *pflag.Flag) {
			fmt.Printf("flag %#v is %#v\n", f.Name, f.Value.String())
		})
	},
}

func Execute() error {
	return rootCmd.Execute()
}
func init() {
	carapace.Gen(rootCmd).Standalone()

	rootCmd.Flags().BoolN("bool-long", "bool-short", false, "BoolN")
	rootCmd.Flags().StringS("delim-colon", "delim-colon", "", "OptargDelimiter ':'")
	rootCmd.Flags().StringS("delim-slash", "delim-slash", "", "OptargDelimiter '/'")
	rootCmd.Flags().CountN("count", "c", "CountN")
	rootCmd.Flags().StringSlice("nargs-any", []string{}, "Nargs")
	rootCmd.Flags().StringSlice("nargs-two", []string{}, "Nargs")

	rootCmd.Flag("delim-colon").NoOptDefVal = " "
	rootCmd.Flag("delim-colon").OptargDelimiter = ':'
	rootCmd.Flag("delim-slash").NoOptDefVal = " "
	rootCmd.Flag("delim-slash").OptargDelimiter = '/'
	rootCmd.Flag("nargs-any").Nargs = -1
	rootCmd.Flag("nargs-two").Nargs = 2

	rootCmd.Flags().SetInterspersed(false)

	carapace.Gen(rootCmd).FlagCompletion(carapace.ActionMap{
		"delim-colon": carapace.ActionValues("d1", "d2", "d3"),
		"delim-slash": carapace.ActionValues("d1", "d2", "d3"),
		"nargs-any": carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			return carapace.ActionValues("na1", "na2", "na3").Invoke(c).Filter(c.Parts...).ToA() // only filters current occurrence
		}),
		"nargs-two": carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			switch len(c.Parts) {
			case 0:
				return carapace.ActionValues("nt1", "nt2", "nt3")
			case 1:
				return carapace.ActionValues("nt4", "nt5", "nt6")
			default:
				return carapace.ActionValues()
			}
		}),
	})

	carapace.Gen(rootCmd).PositionalCompletion(
		carapace.ActionValues("p1", "positional1"),
	)
}
