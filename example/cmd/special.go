package cmd

import (
	"fmt"
	"os"

	"github.com/carapace-sh/carapace"
	"github.com/spf13/cobra"
)

var special = `p1 & < > ' " { } $ # | ? ( ) ;  [ ] * \ $() ${} ` + "` ``"

var specialCmd = &cobra.Command{
	Use:   "special",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 && args[0] != special {
			fmt.Printf("expected: %#v\n", special)
			fmt.Printf("actual  : %#v\n", args[0])
			os.Exit(1)
		} else {
			fmt.Println("ok")
		}
	},
}

func init() {
	carapace.Gen(specialCmd).Standalone()

	specialCmd.Flags().CountP("count", "c", "count flag")
	specialCmd.Flags().Bool("exclusive", false, "mutually exclusive flag")
	specialCmd.Flags().Count("exclusiveRepeatable", "mutually exclusive repeatable flag")
	specialCmd.Flags().StringP("optarg", "o", "", "optional argument")

	specialCmd.Flag("optarg").NoOptDefVal = " "

	specialCmd.MarkFlagsMutuallyExclusive("exclusive", "exclusiveRepeatable")

	rootCmd.AddCommand(specialCmd)

	carapace.Gen(specialCmd).FlagCompletion(carapace.ActionMap{
		"optarg": carapace.ActionValues("optarg1", "optarg2", "optarg3"),
	})

	carapace.Gen(specialCmd).PositionalCompletion(
		carapace.ActionValues(special),
	)
}
