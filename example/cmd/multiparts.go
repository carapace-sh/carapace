package cmd

import (
	"strings"

	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

var multipartsCmd = &cobra.Command{
	Use:   "multiparts",
	Short: "multiparts example",
}

func init() {
	multipartsCmd.Flags().String("at", "", "multiparts with @ as divider")
	multipartsCmd.Flags().String("colon", "", "multiparts with : as divider ")
	multipartsCmd.Flags().String("comma", "", "multiparts with , as divider")
	multipartsCmd.Flags().String("dot", "", "multiparts with . as divider")
	multipartsCmd.Flags().String("dotdotdot", "", "multiparts with ... as divider")
	multipartsCmd.Flags().String("equals", "", "multiparts with = as divider")
	multipartsCmd.Flags().String("slash", "", "multiparts with / as divider")
	multipartsCmd.Flags().String("none", "", "multiparts without divider")

	rootCmd.AddCommand(multipartsCmd)

	carapace.Gen(multipartsCmd).FlagCompletion(carapace.ActionMap{
		"at":        ActionMultipartsTest("@"),
		"colon":     ActionMultipartsTest(":"),
		"comma":     ActionMultipartsTest(","),
		"dot":       ActionMultipartsTest("."),
		"dotdotdot": ActionMultipartsTest("..."),
		"equals":    ActionMultipartsTest("="),
		"slash":     ActionMultipartsTest("/"),
		"none": carapace.ActionMultiParts("", func(args, parts []string) carapace.Action {
			return carapace.ActionValuesDescribed("a", "first", "b", "second", "c", "third", "d", "fourth").Invoke(args).Filter(strings.Split(carapace.CallbackValue, "")).ToA()
		}),
	})

	carapace.Gen(multipartsCmd).PositionalCompletion(
		carapace.ActionMultiParts(",", func(args, entries []string) carapace.Action {
			return carapace.ActionMultiParts("=", func(args, parts []string) carapace.Action {
				switch len(parts) {
				case 0:
					keys := make([]string, len(entries))
					for index, entry := range entries {
						keys[index] = strings.Split(entry, "=")[0]
					}
					return carapace.ActionValues("FILE", "DIRECTORY", "VALUE").Invoke(args).Filter(keys).Suffix("=").ToA()
				case 1:
					switch parts[0] {
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

func ActionMultipartsTest(divider string) carapace.Action {
	return carapace.ActionMultiParts(divider, func(args, parts []string) carapace.Action {
		switch len(parts) {
		case 0:
			return ActionTestValues().Invoke(args).Suffix(divider).ToA()
		case 1:
			return ActionTestValues().Invoke(args).Filter(parts).Suffix(divider).ToA()
		case 2:
			return ActionTestValues().Invoke(args).Filter(parts).ToA()
		default:
			return carapace.ActionValues()
		}
	})
}

func ActionTestValues() carapace.Action {
	return carapace.ActionValuesDescribed("first", "first value", "second", "second value", "third with space", "third value")
}
