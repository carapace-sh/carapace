package cmd

import (
	"strings"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/style"
	"github.com/spf13/cobra"
)

var actionCmd = &cobra.Command{
	Use:     "action",
	Short:   "action example",
	Aliases: []string{"alias"},
	GroupID: "main",
	RunE: func(cmd *cobra.Command, args []string) error {
		if f := cmd.Flag("style"); f.Changed {
			if splitted := strings.Split(f.Value.String(), "="); len(splitted) == 2 {
				return style.Set(splitted[0], strings.Replace(splitted[1], ",", " ", -1))
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(actionCmd)

	actionCmd.Flags().String("directories", "", "ActionDirectories()")
	actionCmd.Flags().String("files", "", "ActionFiles()")
	actionCmd.Flags().String("files-filtered", "", "ActionFiles(\".md\", \"go.mod\", \"go.sum\")")
	actionCmd.Flags().String("message", "", "ActionMessage()")
	actionCmd.Flags().String("styled-values", "", "ActionStyledValues()")
	actionCmd.Flags().String("styled-values-described", "", "ActionStyledValuesDescribed()")
	actionCmd.Flags().String("values", "", "ActionValues()")
	actionCmd.Flags().String("values-described", "", "ActionValuesDescribed()")

	carapace.Gen(actionCmd).FlagCompletion(carapace.ActionMap{
		"directories":    carapace.ActionDirectories(),
		"files":          carapace.ActionFiles(),
		"files-filtered": carapace.ActionFiles(".md", "go.mod", "go.sum"),
		"message":        carapace.ActionMessage("example message"),
		"styled-values": carapace.ActionStyledValues(
			"first", style.Default,
			"second", style.Blue,
			"third", style.Of(style.BgBrightBlack, style.Magenta, style.Bold),
		),
		"styled-values-described": carapace.ActionStyledValuesDescribed(
			"first", "description of first", style.Blink,
			"second", "description of second", style.Of("color210", style.Underlined),
			"third", "description of third", style.Of("#112233", style.Italic),
		),
		"values": carapace.ActionValues("first", "second", "third"),
		"values-described": carapace.ActionValuesDescribed(
			"first", "description of first",
			"second", "description of second",
			"third", "description of third",
		),
	})

	carapace.Gen(actionCmd).PositionalCompletion(
		carapace.ActionValues("positional1", "p1", "positional1 with space"),
		carapace.ActionValues("positional2", "p2", "positional2 with space"),
	)

	carapace.Gen(actionCmd).DashCompletion(
		carapace.ActionValues("dash1", "d1", "dash1 with space"),
		carapace.ActionValues("dash2", "d2", "dash2 with space"),
	)

	carapace.Gen(actionCmd).DashAnyCompletion(
		carapace.ActionValues("dashAny", "dAny", "dashAny with space"),
	)
}
