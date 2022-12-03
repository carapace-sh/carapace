package cmd

import (
	"strings"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/style"
	"github.com/spf13/cobra"
)

var actionCmd = &cobra.Command{
	Use:     "action [pos1] [pos2] [--] [dashAny]...",
	Short:   "action example",
	Aliases: []string{"alias"},
	GroupID: "main",
	Run:     func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(actionCmd)

	actionCmd.Flags().String("callback", "", "ActionCallback()")
	actionCmd.Flags().String("directories", "", "ActionDirectories()")
	actionCmd.Flags().String("exec-command", "", "ActionExecCommand()")
	actionCmd.Flags().String("files", "", "ActionFiles()")
	actionCmd.Flags().String("files-filtered", "", "ActionFiles(\".md\", \"go.mod\", \"go.sum\")")
	actionCmd.Flags().String("import", "", "ActionImport()")
	actionCmd.Flags().String("message", "", "ActionMessage()")
	actionCmd.Flags().String("message-multiple", "", "ActionMessage()")
	actionCmd.Flags().String("multiparts", "", "ActionMultiParts()")
	actionCmd.Flags().String("multiparts-nested", "", "ActionMultiParts(...ActionMultiParts...)")
	actionCmd.Flags().String("styled-values", "", "ActionStyledValues()")
	actionCmd.Flags().String("styled-values-described", "", "ActionStyledValuesDescribed()")
	actionCmd.Flags().String("values", "", "ActionValues()")
	actionCmd.Flags().String("values-described", "", "ActionValuesDescribed()")

	carapace.Gen(actionCmd).FlagCompletion(carapace.ActionMap{
		"callback": carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			if flag := actionCmd.Flag("values"); flag.Changed {
				return carapace.ActionMessage("values flag is set to: '%v'", flag.Value.String())
			}
			return carapace.ActionMessage("values flag is not set")
		}),
		"directories": carapace.ActionDirectories(),
		"exec-command": carapace.ActionExecCommand("git", "remote")(func(output []byte) carapace.Action {
			lines := strings.Split(string(output), "\n")
			return carapace.ActionValues(lines[:len(lines)-1]...)
		}),
		"files":          carapace.ActionFiles(),
		"files-filtered": carapace.ActionFiles(".md", "go.mod", "go.sum"),
		"import": carapace.ActionImport([]byte(`
{
  "Version": "unknown",
  "Nospace": "",
  "RawValues": [
    {
      "Value": "first",
      "Display": "first"
    },
    {
      "Value": "second",
      "Display": "second"
    },
    {
      "Value": "third",
      "Display": "third"
    }
  ]
}
		`)),
		"message": carapace.ActionMessage("example message"),
		"message-multiple": carapace.Batch(
			carapace.ActionMessage("first message"),
			carapace.ActionMessage("second message"),
			carapace.ActionMessage("third message"),
			carapace.ActionValues("one", "two", "three"),
		).ToA(),
		"multiparts": carapace.ActionMultiParts(":", func(c carapace.Context) carapace.Action {
			switch len(c.Parts) {
			case 0:
				return carapace.ActionValues("userA", "userB").Invoke(c).Suffix(":").ToA()
			case 1:
				return carapace.ActionValues("groupA", "groupB")
			default:
				return carapace.ActionValues()
			}
		}),
		"multiparts-nested": carapace.ActionMultiParts(",", func(cEntries carapace.Context) carapace.Action {
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
						return carapace.ActionFiles("").NoSpace()
					case "DIRECTORY":
						return carapace.ActionDirectories().NoSpace()
					case "VALUE":
						return carapace.ActionValues("one", "two", "three").NoSpace()
					default:
						return carapace.ActionValues()

					}
				default:
					return carapace.ActionValues()
				}
			})
		}),
		"styled-values": carapace.ActionStyledValues(
			"first", style.Default,
			"second", style.Blue,
			"third", style.Of(style.BgBrightBlack, style.Magenta, style.Bold),
		),
		"styled-values-described": carapace.ActionStyledValuesDescribed(
			"first", "description of first", style.Blink,
			"second", "description of second", style.Of("color210", style.Underlined),
			"third", "description of third", style.Of("#112233", style.Italic),
			"thirdalias", "description of third", style.BgBrightMagenta,
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

	carapace.Gen(actionCmd).DashAnyCompletion(
		carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			cmd := &cobra.Command{
				Use: "embedded",
				CompletionOptions: cobra.CompletionOptions{
					DisableDefaultCmd: true,
				},
				Run: func(cmd *cobra.Command, args []string) {},
			}

			cmd.Flags().Bool("embedded-flag", false, "embedded flag")

			carapace.Gen(cmd).PositionalCompletion(
				carapace.ActionValues("embeddedPositional1", "embeddedP1"),
				carapace.ActionValues("embeddedPositional2", "embeddedP2"),
			)

			return carapace.ActionExecute(cmd)
		}),
	)
}
