package cmd

import (
	"strings"

	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

var multipartsCmd = &cobra.Command{
	Use:   "multiparts",
	Short: "multiparts example",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	multipartsCmd.Flags().String("at", "", "multiparts with @ as divider")
	multipartsCmd.Flags().String("colon", "", "multiparts with : as divider ")
	multipartsCmd.Flags().String("comma", "", "multiparts with , as divider")
	multipartsCmd.Flags().String("dot", "", "multiparts with . as divider")
	multipartsCmd.Flags().String("dotdotdot", "", "multiparts with ... as divider")
	multipartsCmd.Flags().String("equals", "", "multiparts with = as divider")
	multipartsCmd.Flags().String("none", "", "multiparts without divider")
	multipartsCmd.Flags().String("slash", "", "multiparts with / as divider")

	rootCmd.AddCommand(multipartsCmd)

	carapace.Gen(multipartsCmd).FlagCompletion(carapace.ActionMap{
		"at":        actionMultipartsTest("@"),
		"colon":     actionMultipartsTest(":"),
		"comma":     actionMultipartsTest(","),
		"dot":       actionMultipartsTest("."),
		"dotdotdot": actionMultipartsTest("..."),
		"equals":    actionMultipartsTest("="),
		"none": carapace.ActionMultiParts("", func(c carapace.Context) carapace.Action {
			return carapace.ActionValuesDescribed("a", "first", "b", "second", "c", "third", "d", "fourth").Invoke(c).Filter(c.Parts).ToA()
		}),
		"slash": actionMultipartsTest("/"),
	})

	carapace.Gen(multipartsCmd).PositionalCompletion(
		carapace.ActionMultiParts(",", func(cEntries carapace.Context) carapace.Action {
			return carapace.ActionMultiParts("=", func(c carapace.Context) carapace.Action {
				switch len(c.Parts) {
				case 0:
					keys := make([]string, len(cEntries.Parts))
					for index, entry := range cEntries.Parts {
						keys[index] = strings.Split(entry, "=")[0]
					}
					return carapace.ActionValues("FILE", "DIRECTORY", "VALUE").Invoke(c).Filter(keys).Suffix("=").ToA()
				case 1:
					switch c.Parts[0] {
					case "FILE":
						return carapace.ActionFiles("")
					case "DIRECTORY":
						return carapace.ActionDirectories()
					case "VALUE":
						return carapace.ActionValues("one", "two", "three")
					default:
						return carapace.ActionValues()

					}
				default:
					return carapace.ActionValues()
				}
			})
		}),
	)
}

func actionMultipartsTest(divider string) carapace.Action {
	return carapace.ActionMultiParts(divider, func(c carapace.Context) carapace.Action {
		switch len(c.Parts) {
		case 0:
			return actionTestValues().Invoke(c).Suffix(divider).ToA()
		case 1:
			return actionTestValues().Invoke(c).Filter(c.Parts).Suffix(divider).ToA()
		case 2:
			return actionTestValues().Invoke(c).Filter(c.Parts).ToA()
		default:
			return carapace.ActionValues()
		}
	})
}

func actionTestValues() carapace.Action {
	return carapace.ActionValuesDescribed("first", "first value", "second", "second value", "third with space", "third value", "fourth", "fourth value")
}
