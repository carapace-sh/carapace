package cmd

import (
	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/example/cmd/action/net"
	"github.com/rsteube/carapace/example/cmd/action/os"
	"github.com/rsteube/carapace/pkg/style"
	"github.com/spf13/cobra"
)

var actionCmd = &cobra.Command{
	Use:     "action",
	Short:   "action example",
	Aliases: []string{"alias"},
	Run:     func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(actionCmd)

	actionCmd.Flags().CountP("count", "c", "count flag")
	actionCmd.Flags().StringP("files", "f", "", "files flag")
	actionCmd.Flags().String("filtered_files", "", "files flag")
	actionCmd.Flags().String("directories", "", "files flag")
	actionCmd.Flags().StringP("groups", "g", "", "groups flag")
	actionCmd.Flags().StringP("message", "m", "", "message flag")
	actionCmd.Flags().StringP("net_interfaces", "n", "", "net_interfaces flag")
	actionCmd.Flags().String("usergroup", "", "user:group flag")
	actionCmd.Flags().StringP("users", "u", "", "users flag")
	actionCmd.Flags().StringP("values", "v", "", "values flag")
	actionCmd.Flags().StringP("values_described", "d", "", "values with description flag")
	actionCmd.Flags().String("styled_values", "", "styled values flag")
	actionCmd.Flags().String("styled_values_described", "", "styled values with description flag")
	//actionCmd.Flags().StringS("shorthandonly", "s", "", "shorthandonly flag")
	actionCmd.Flags().StringP("kill", "k", "", "kill signals")
	actionCmd.Flags().StringP("optarg", "o", "", "optional arg with default value blue")
	actionCmd.Flag("optarg").NoOptDefVal = "blue"

	carapace.Gen(actionCmd).FlagCompletion(carapace.ActionMap{
		"files": carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			return carapace.ActionFiles().Chdir(actionCmd.Flag("directories").Value.String())
		}),
		"filtered_files": carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			return carapace.ActionFiles(".go", "go.mod", ".txt").Chdir(actionCmd.Flag("directories").Value.String())
		}),
		"directories":      carapace.ActionDirectories(),
		"groups":           os.ActionGroups(),
		"message":          carapace.ActionMessage("message example"),
		"net_interfaces":   net.ActionNetInterfaces(),
		"usergroup":        os.ActionUserGroup(),
		"users":            os.ActionUsers(),
		"values":           carapace.ActionValues("values", "example"),
		"values_described": carapace.ActionValuesDescribed("values", "valueDescription", "example", "exampleDescription"),
		"styled_values": carapace.ActionStyledValues(
			"red", style.Red,
			"green", style.Green,
			"yellow", style.Yellow,
			"blue", style.Blue,
			"magenta", style.Magenta,
			"cyan", style.Cyan,
			"bright-black", style.BrightBlack,
			"default", style.Default,
		),
		"styled_values_described": carapace.ActionStyledValuesDescribed("values", "valueDescription", style.Blue, "example", "exampleDescription", style.Red),
		"kill":                    os.ActionKillSignals(),
		"optarg":                  carapace.ActionValues("blue", "red", "green", "yellow"),
	})

	carapace.Gen(actionCmd).PositionalCompletion(
		carapace.ActionValues("positional1", "p1"),
		carapace.ActionValues("positional2", "p2"),
	)

	carapace.Gen(actionCmd).DashCompletion(
		carapace.ActionValues("dash1", "d1"),
		carapace.ActionValues("dash2", "d2"),
	)

	carapace.Gen(actionCmd).DashAnyCompletion(
		carapace.ActionValues("dashAny", "dAny"),
	)
}
