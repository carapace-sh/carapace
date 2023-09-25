package cmd

import (
	"os/exec"
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
	actionCmd.Flags().String("cobra", "", "ActionCobra()")
	actionCmd.Flags().String("commands", "", "ActionCommands()")
	actionCmd.Flags().String("directories", "", "ActionDirectories()")
	actionCmd.Flags().String("execcommand", "", "ActionExecCommand()")
	actionCmd.Flags().String("execcommandE", "", "ActionExecCommand()")
	actionCmd.Flags().String("executables", "", "ActionExecutables()")
	actionCmd.Flags().String("files", "", "ActionFiles()")
	actionCmd.Flags().String("files-filtered", "", "ActionFiles(\".md\", \"go.mod\", \"go.sum\")")
	actionCmd.Flags().String("import", "", "ActionImport()")
	actionCmd.Flags().String("message", "", "ActionMessage()")
	actionCmd.Flags().String("message-multiple", "", "ActionMessage()")
	actionCmd.Flags().String("multiparts", "", "ActionMultiParts()")
	actionCmd.Flags().String("multiparts-nested", "", "ActionMultiParts(...ActionMultiParts...)")
	actionCmd.Flags().String("multipartsn", "", "ActionMultiPartsN()")
	actionCmd.Flags().String("multipartsn-empty", "", "ActionMultiPartsN()")
	actionCmd.Flags().String("styles", "", "ActionStyles()")
	actionCmd.Flags().String("styleconfig", "", "ActionStyleConfig()")
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
		"cobra": carapace.ActionCobra(func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{"one", "two"}, cobra.ShellCompDirectiveNoSpace
		}),
		"commands":    carapace.ActionCommands(rootCmd).Split(),
		"directories": carapace.ActionDirectories(),
		"execcommand": carapace.ActionExecCommand("git", "remote")(func(output []byte) carapace.Action {
			lines := strings.Split(string(output), "\n")
			return carapace.ActionValues(lines[:len(lines)-1]...)
		}),
		"execcommandE": carapace.ActionExecCommandE("false")(func(output []byte, err error) carapace.Action {
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					return carapace.ActionMessage("failed with %v", exitErr.ExitCode())
				}
				return carapace.ActionMessage(err.Error())
			}
			return carapace.ActionValues()
		}),
		"executables":    carapace.ActionExecutables(),
		"files":          carapace.ActionFiles(),
		"files-filtered": carapace.ActionFiles(".md", "go.mod", "go.sum"),
		"import": carapace.ActionImport([]byte(`
{
  "version": "unknown",
  "nospace": "",
  "values": [
    {
      "value": "first",
      "display": "first"
    },
    {
      "value": "second",
      "display": "second"
    },
    {
      "value": "third",
      "display": "third"
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
					return carapace.ActionValues("FILE", "DIRECTORY", "VALUE").Invoke(c).Filter(keys...).Suffix("=").ToA()
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
		"multipartsn": carapace.ActionMultiPartsN("=", 2, func(c carapace.Context) carapace.Action {
			switch len(c.Parts) {
			case 0:
				return carapace.ActionValues("one", "two").Suffix("=")
			case 1:
				return carapace.ActionMultiParts("=", func(c carapace.Context) carapace.Action {
					switch len(c.Parts) {
					case 0:
						return carapace.ActionValues("three", "four").Suffix("=")
					case 1:
						return carapace.ActionValues("five", "six")
					default:
						return carapace.ActionValues()
					}
				})
			default:
				return carapace.ActionMessage("should never happen")
			}
		}),
		"multipartsn-empty": carapace.ActionMultiPartsN("", 2, func(c carapace.Context) carapace.Action {
			switch len(c.Parts) {
			case 0:
				return carapace.ActionValues("a", "b")
			case 1:
				return carapace.ActionValues("c", "d", "e").UniqueList("")
			default:
				return carapace.ActionMessage("should never happen")
			}
		}),
		"styles":      carapace.ActionStyles(),
		"styleconfig": carapace.ActionStyleConfig(),
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

	carapace.Gen(actionCmd).PositionalAnyCompletion(
		carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			cmd := &cobra.Command{
				Use: "embedded",
				CompletionOptions: cobra.CompletionOptions{
					DisableDefaultCmd: true,
				},
				Run: func(cmd *cobra.Command, args []string) {},
			}

			cmd.Flags().Bool("embedded-bool", false, "embedded bool flag")
			cmd.Flags().String("embedded-string", "", "embedded string flag")
			cmd.Flags().String("embedded-optarg", "", "embedded optarg flag")

			cmd.Flag("embedded-optarg").NoOptDefVal = " "

			carapace.Gen(cmd).FlagCompletion(carapace.ActionMap{
				"embedded-string": carapace.ActionValues("es1", "es2", "es3"),
				"embedded-optarg": carapace.ActionValues("eo1", "eo2", "eo3"),
			})

			carapace.Gen(cmd).PositionalCompletion(
				carapace.ActionValues("embeddedPositional1", "embeddedP1"),
				carapace.ActionValues("embeddedPositional2 with space", "embeddedP2 with space"),
			)

			return carapace.ActionExecute(cmd)
		}),
	)

	carapace.Gen(actionCmd).DashAnyCompletion(
		carapace.ActionPositional(actionCmd),
	)
}
