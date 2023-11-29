package cmd

import (
	"strings"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/style"
	"github.com/spf13/cobra"
)

var multipartsCmd = &cobra.Command{
	Use:   "multiparts",
	Short: "multiparts example",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	multipartsCmd.Flags().String("at", "", "multiparts with @ as divider")
	multipartsCmd.Flags().String("colon", "", "multiparts with : as divider")
	multipartsCmd.Flags().String("comma", "", "multiparts with , as divider")
	multipartsCmd.Flags().String("dot", "", "multiparts with . as divider")
	multipartsCmd.Flags().String("dotdotdot", "", "multiparts with ... as divider")
	multipartsCmd.Flags().String("equals", "", "multiparts with = as divider")
	multipartsCmd.Flags().String("none", "", "multiparts without divider")
	multipartsCmd.Flags().String("none-zero", "", "multiparts without divider limited to 0")
	multipartsCmd.Flags().String("none-one", "", "multiparts without divider limited to 1")
	multipartsCmd.Flags().String("none-two", "", "multiparts without divider limited to 2")
	multipartsCmd.Flags().String("none-three", "", "multiparts without divider limited to 3")
	multipartsCmd.Flags().String("slash", "", "multiparts with / as divider")
	multipartsCmd.Flags().String("space", "", "multiparts with space as divider")

	rootCmd.AddCommand(multipartsCmd)

	carapace.Gen(multipartsCmd).FlagCompletion(carapace.ActionMap{
		"at":        actionMultipartsTest("@"),
		"colon":     actionMultipartsTest(":"),
		"comma":     actionMultipartsTest(","),
		"dot":       actionMultipartsTest("."),
		"dotdotdot": actionMultipartsTest("..."),
		"equals":    actionMultipartsTest("="),
		"none":      carapace.ActionValuesDescribed("a", "first", "b", "second", "c", "third", "d", "fourth").UniqueList(""),
		"none-zero": carapace.ActionMultiPartsN("", 0, func(c carapace.Context) carapace.Action {
			return carapace.ActionMessage("unreachable")
		}),
		"none-one": carapace.ActionMultiPartsN("", 1, func(c carapace.Context) carapace.Action {
			return carapace.ActionValues("a", "b")
		}),
		"none-two": carapace.ActionMultiPartsN("", 2, func(c carapace.Context) carapace.Action {
			switch len(c.Parts) {
			case 0:
				return carapace.ActionValuesDescribed(
					"a", "zero",
					"b", "zero",
				).Style(style.Blue)
			default:
				return carapace.ActionValuesDescribed(
					"a", "default",
					"b", "default",
					"c", "default",
				).Style(style.Red).
					UniqueList("")
			}
		}),
		"none-three": carapace.ActionMultiPartsN("", 3, func(c carapace.Context) carapace.Action {
			switch len(c.Parts) {
			case 0:
				return carapace.ActionValuesDescribed(
					"a", "zero",
					"b", "zero",
				).Style(style.Blue)
			case 1:
				return carapace.ActionValuesDescribed(
					"a", "one",
					"b", "one",
					"c", "one",
				).Style(style.Red)
			default:
				return carapace.ActionValuesDescribed(
					"a", "default",
					"b", "default",
					"c", "default",
					"d", "default",
				).Style(style.Green).
					UniqueList("")
			}
		}),
		"slash": actionMultipartsTest("/"),
		"space": carapace.ActionValues(
			"one",
			"two",
			"three",
			"four",
		).UniqueList(" "),
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
					return carapace.ActionValues("FILE", "DIRECTORY", "VALUE").Filter(keys...).Suffix("=")
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
			return actionTestValues().Suffix(divider)
		case 1:
			return actionTestValues().FilterParts().Suffix(divider)
		case 2:
			return actionTestValues().FilterParts()
		default:
			return carapace.ActionValues()
		}
	})
}

func actionTestValues() carapace.Action {
	return carapace.ActionValuesDescribed("first", "first value", "second", "second value", "third with space", "third value", "fourth", "fourth value")
}
