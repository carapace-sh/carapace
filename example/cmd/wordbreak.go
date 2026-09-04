package cmd

import (
	"github.com/carapace-sh/carapace"
	"github.com/spf13/cobra"
)

var wordbreakCmd = &cobra.Command{
	Use:   "wordbreak",
	Short: "wordbreak example",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	carapace.Gen(specialCmd).Standalone()

	wordbreakCmd.Flags().String("space", "", "space wordbreak")
	wordbreakCmd.Flags().String("tab", "", "tab wordbreak")
	wordbreakCmd.Flags().String("carriage-return", "", "carriage-return wordbreak")
	wordbreakCmd.Flags().String("newline", "", "newline wordbreak")
	wordbreakCmd.Flags().String("double-quote", "", "double-quote wordbreak")
	wordbreakCmd.Flags().String("single-quote", "", "single-quote wordbreak")
	wordbreakCmd.Flags().String("at", "", "at wordbreak")
	wordbreakCmd.Flags().String("greater-than", "", "greater-than wordbreak")
	wordbreakCmd.Flags().String("less-than", "", "less-than wordbreak")
	wordbreakCmd.Flags().String("equals", "", "equals wordbreak")
	wordbreakCmd.Flags().String("semicolon", "", "semicolon wordbreak")
	wordbreakCmd.Flags().String("pipe", "", "pipe wordbreak")
	wordbreakCmd.Flags().String("and", "", "and wordbreak")
	wordbreakCmd.Flags().String("round-bracket", "", "round-bracket wordbreak")
	wordbreakCmd.Flags().String("colon", "", "colon wordbreak")

	rootCmd.AddCommand(wordbreakCmd)

	carapace.Gen(wordbreakCmd).FlagCompletion(carapace.ActionMap{
		"space":           carapace.ActionValues("one", "two", "three").UniqueList(" "),
		"tab":             carapace.ActionValues("one", "two", "three").UniqueList("\t"),
		"carriage-return": carapace.ActionValues("one", "two", "three").UniqueList("\r"),
		"newline":         carapace.ActionValues("one", "two", "three").UniqueList("\n"),
		"double-quote":    carapace.ActionValues("one", "two", "three").UniqueList(`"`),
		"single-quote":    carapace.ActionValues("one", "two", "three").UniqueList("'"),
		"at":              carapace.ActionValues("one", "two", "three").UniqueList("@"),
		"greater-than":    carapace.ActionValues("one", "two", "three").UniqueList(">"),
		"less-than":       carapace.ActionValues("one", "two", "three").UniqueList("<"),
		"equals":          carapace.ActionValues("one", "two", "three").UniqueList("="),
		"semicolon":       carapace.ActionValues("one", "two", "three").UniqueList(";"),
		"pipe":            carapace.ActionValues("one", "two", "three").UniqueList("|"),
		"and":             carapace.ActionValues("one", "two", "three").UniqueList("&"),
		"round-bracket":   carapace.ActionValues("one", "two", "three").UniqueList("("),
		"colon":           carapace.ActionValues("one", "two", "three").UniqueList(":"),
	})
}
