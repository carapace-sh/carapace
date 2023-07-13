package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/cache"
	"github.com/rsteube/carapace/pkg/style"
	"github.com/spf13/cobra"
)

var modifierCmd = &cobra.Command{
	Use:     "modifier [pos1]",
	Short:   "modifier example",
	GroupID: "modifier",
	Run:     func(cmd *cobra.Command, args []string) {},
}

func init() {
	modifierCmd.Flags().String("batch", "", "Batch()")

	modifierCmd.Flags().String("cache", "", "Cache()")
	modifierCmd.Flags().String("cache-key", "", "Cache()")
	modifierCmd.Flags().String("chdir", "", "Chdir()")
	modifierCmd.Flags().String("filter", "", "Filter()")
	modifierCmd.Flags().String("invoke", "", "Invoke()")
	modifierCmd.Flags().String("list", "", "List()")
	modifierCmd.Flags().String("multiparts", "", "MultiParts()")
	modifierCmd.Flags().String("nospace", "", "NoSpace()")
	modifierCmd.Flags().String("prefix", "", "Prefix()")
	modifierCmd.Flags().String("retain", "", "Retain()")
	modifierCmd.Flags().String("shift", "", "Shift()")
	modifierCmd.Flags().String("style", "", "Style()")
	modifierCmd.Flags().String("stylef", "", "StyleF()")
	modifierCmd.Flags().String("styler", "", "StyleR()")
	modifierCmd.Flags().String("suffix", "", "Suffix()")
	modifierCmd.Flags().String("suppress", "", "Suppress()")
	modifierCmd.Flags().String("tag", "", "Tag()")
	modifierCmd.Flags().String("tagf", "", "TagF()")
	modifierCmd.Flags().String("timeout", "", "Timeout()")
	modifierCmd.Flags().String("uniquelist", "", "UniqueList()")
	modifierCmd.Flags().String("usage", "", "Usage()")

	rootCmd.AddCommand(modifierCmd)

	carapace.Gen(modifierCmd).FlagCompletion(carapace.ActionMap{
		"batch": carapace.Batch(
			carapace.ActionValuesDescribed(
				"A", "description of A",
				"B", "description of first B",
			),
			carapace.ActionValuesDescribed(
				"B", "description of second B",
				"C", "description of first C",
			),
			carapace.ActionValuesDescribed(
				"C", "description of second C",
				"D", "description of D",
			),
		).ToA(),
		"cache": carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			return carapace.ActionValues(
				time.Now().Format("15:04:05"),
			).Cache(5 * time.Second)
		}),
		"cache-key": carapace.ActionMultiParts("/", func(c carapace.Context) carapace.Action {
			switch len(c.Parts) {
			case 0:
				return carapace.ActionValues("one", "two").Suffix("/")
			case 1:
				return carapace.ActionValues(
					time.Now().Format("15:04:05"),
				).Cache(10*time.Second, cache.String(c.Parts[0]))
			default:
				return carapace.ActionValues()
			}
		}),
		"chdir": carapace.ActionFiles().Chdir(os.TempDir()),
		"filter": carapace.ActionValuesDescribed(
			"1", "one",
			"2", "two",
			"3", "three",
			"4", "four",
		).Filter([]string{"2", "4"}),
		"list": carapace.ActionValues("one", "two", "three").List(","),
		"invoke": carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			if !strings.HasPrefix(c.Value, "file://") {
				return carapace.ActionValues("file://").NoSpace()
			}

			c.Value = strings.TrimPrefix(c.Value, "file://")
			return carapace.ActionFiles().Invoke(c).Prefix("file://").ToA()
		}),
		"nospace": carapace.ActionValues(
			"one,",
			"two/",
			"three",
		).NoSpace(',', '/'),
		"timeout": carapace.ActionMultiParts(":", func(c carapace.Context) carapace.Action {
			switch len(c.Parts) {
			case 0:
				return carapace.ActionValuesDescribed(
					"1s", "within timeout",
					"3s", "exceeding timeout",
				).Suffix(":")
			case 1:
				d, err := time.ParseDuration(c.Parts[0])
				if err != nil {
					return carapace.ActionMessage(err.Error())
				}

				return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
					time.Sleep(d)
					return carapace.ActionValues("within timeout")
				}).Timeout(2*time.Second, carapace.ActionMessage("timeout exceeded"))
			default:
				return carapace.ActionValues()
			}
		}),
		"multiparts": carapace.ActionValues(
			"dir/subdir1/fileA.txt",
			"dir/subdir1/fileB.txt",
			"dir/subdir2/fileC.txt",
		).MultiParts("/"),
		"prefix": carapace.ActionValues(
			"melon",
			"drop",
			"fall",
		).Prefix("water"),
		"retain": carapace.ActionValuesDescribed(
			"1", "one",
			"2", "two",
			"3", "three",
			"4", "four",
		).Retain([]string{"2", "4"}),
		"shift": carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			return carapace.ActionMessage("%#v", c.Args)
		}).Shift(1),
		"style": carapace.ActionValues(
			"one",
			"two",
		).Style(style.Green),
		"stylef": carapace.ActionValues(
			"one",
			"two",
			"three",
		).StyleF(func(s string, sc style.Context) string {
			switch s {
			case "one":
				return style.Green
			case "two":
				return style.Red
			default:
				return style.Default
			}
		}),
		"styler": carapace.ActionValues(
			"one",
			"two",
		).StyleR(&style.Carapace.KeywordAmbiguous),
		"suffix": carapace.ActionValues(
			"apple",
			"melon",
			"orange",
		).Suffix("juice"),
		"suppress": carapace.Batch(
			carapace.ActionMessage("unexpected error"),
			carapace.ActionMessage("ignored error"),
		).ToA().Suppress("ignored"),
		"tag": carapace.ActionValues(
			"192.168.1.1",
			"127.0.0.1",
		).Tag("interfaces"),
		"tagf": carapace.ActionValues(
			"one.png",
			"two.gif",
			"three.txt",
			"four.md",
		).StyleF(style.ForPathExt).TagF(func(s string) string {
			switch filepath.Ext(s) {
			case ".png", ".gif":
				return "images"
			case ".txt", ".md":
				return "documents"
			default:
				return ""
			}
		}),
		"uniquelist": carapace.ActionValues("one", "two", "three").UniqueList(","),
		"usage": carapace.ActionMultiParts(":", func(c carapace.Context) carapace.Action {
			switch len(c.Parts) {
			case 0:
				return carapace.ActionValues("explicit", "implicit").Suffix(":")
			case 1:
				if c.Parts[0] == "explicit" {
					return carapace.ActionValues().Usage("explicit usage")
				}
				return carapace.ActionValues()

			default:
				return carapace.ActionValues()
			}
		}),
	})

	carapace.Gen(modifierCmd).PositionalCompletion(
		carapace.ActionValues().Usage("explicit positional usage"),
	)
}
