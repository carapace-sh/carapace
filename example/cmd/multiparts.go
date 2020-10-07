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
	multipartsCmd.Flags().String("equals", "", "multiparts with = as divider")
	multipartsCmd.Flags().String("slash", "", "multiparts with / as divider")
	multipartsCmd.Flags().String("none", "", "multiparts without divider")

	rootCmd.AddCommand(multipartsCmd)

	carapace.Gen(multipartsCmd).FlagCompletion(carapace.ActionMap{
		"at":     ActionMultipartsTest("@"),
		"colon":  ActionMultipartsTest(":"),
		"comma":  ActionMultipartsTest(","),
		"dot":    ActionMultipartsTest("."),
		"equals": ActionMultipartsTest("="),
		"slash":  ActionMultipartsTest("/"),
		"none": carapace.ActionMultiParts("", func(args, parts []string) carapace.Action {
			return carapace.ActionValuesDescribed("a", "first", "b", "second", "c", "third", "d", "fourth").Invoke(args).Filter(strings.Split(carapace.CallbackValue, "")).ToA()
		}),
	})
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
