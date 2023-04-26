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
	modifierCmd.Flags().String("tomultiparts", "", "ToMultiPartsA()")
	modifierCmd.Flags().String("usage", "", "Usage()")
	modifierCmd.Flags().StringS("x", "x", "", "xxxx")
	modifierCmd.Flags().StringS("y", "y", "", "xxxx")
	modifierCmd.Flags().StringS("z", "z", "", "xxxx")
	modifierCmd.Flags().BoolS("u", "u", false, "xxxx")
	modifierCmd.Flags().BoolS("v", "v", false, "xxxx")
	modifierCmd.Flags().BoolS("w", "w", false, "xxxx")

	rootCmd.AddCommand(modifierCmd)

	carapace.Gen(modifierCmd).FlagCompletion(carapace.ActionMap{
		"batch": carapace.Batch(
			carapace.ActionValuesDescribed(
				"A", "description of A",
				"B", "description of first B",
			),
			carapace.ActionValuesDescribed(
				"B", "description of second B",
				"C", "description of first C",
			),
			carapace.ActionValuesDescribed(
				"C", "description of second C",
				"D", "description of D",
			),
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
		"tomultiparts": carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			return carapace.ActionValuesDescribed(
				"1=1==1/1", "one",
				"1=1==1/2", "two",
				"1=1==2/1", "three",
				"1=1==2/2", "four",
				"1=2==1/1", "five",
				"2=1==1/1", "six",
			).Invoke(c).ToMultiPartsA("==", "=", "/")
		}),
		"usage": carapace.ActionValues().Usage("explicit flag usage"),
	})

	carapace.Gen(modifierCmd).PositionalCompletion(
		carapace.ActionValues().Usage("explicit positional usage"),
	)
}
