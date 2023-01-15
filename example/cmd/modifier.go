package cmd

import (
	"time"

	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

var modifierCmd = &cobra.Command{
	Use:     "modifier [pos1]",
	Short:   "modifier example",
	GroupID: "modifier",
	Run:     func(cmd *cobra.Command, args []string) {},
}

func init() {
	modifierCmd.Flags().String("batch", "", "Batch()")
	modifierCmd.Flags().String("timeout", "", "Timeout()")
	modifierCmd.Flags().String("usage", "", "Usage()")

	rootCmd.AddCommand(modifierCmd)

	carapace.Gen(modifierCmd).FlagCompletion(carapace.ActionMap{
		"batch": carapace.Batch(
			carapace.ActionValues("A", "B"),
			carapace.ActionValues("C", "D"),
			carapace.ActionValues("E", "F"),
		).ToA(),
		"timeout": carapace.ActionMultiParts(":", func(c carapace.Context) carapace.Action {
			switch len(c.Parts) {
			case 0:
				return carapace.ActionValuesDescribed(
					"1s", "within timeout",
					"3s", "exceeding timeout",
				).Suffix(":")
			case 1:
				d, err := time.ParseDuration(c.Parts[0])
				if err != nil {
					return carapace.ActionMessage(err.Error())
				}

				return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
					time.Sleep(d)
					return carapace.ActionValues("within timeout")
				}).Timeout(2*time.Second, carapace.ActionMessage("timeout exceeded"))
			default:
				return carapace.ActionValues()
			}
		}),

		"usage": carapace.ActionValues().Usage("explicit flag usage"),
	})

	carapace.Gen(modifierCmd).PositionalCompletion(
		carapace.ActionValues().Usage("explicit positional usage"),
	)
}
