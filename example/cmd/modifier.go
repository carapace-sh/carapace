package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/cache/key"
	"github.com/rsteube/carapace/pkg/condition"
	"github.com/rsteube/carapace/pkg/style"
	"github.com/rsteube/carapace/pkg/traverse"
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
	modifierCmd.Flags().String("chdirf", "", "ChdirF()")
	modifierCmd.Flags().String("filter", "", "Filter()")
	modifierCmd.Flags().String("filterargs", "", "FilterArgs()")
	modifierCmd.Flags().String("filterparts", "", "FilterParts()")
	modifierCmd.Flags().String("invoke", "", "Invoke()")
	modifierCmd.Flags().String("list", "", "List()")
	modifierCmd.Flags().String("multiparts", "", "MultiParts()")
	modifierCmd.Flags().String("multipartsp", "", "MultiPartsP()")
	modifierCmd.Flags().String("nospace", "", "NoSpace()")
	modifierCmd.Flags().String("prefix", "", "Prefix()")
	modifierCmd.Flags().String("retain", "", "Retain()")
	modifierCmd.Flags().String("shift", "", "Shift()")
	modifierCmd.Flags().String("split", "", "Split()")
	modifierCmd.Flags().String("splitp", "", "SplitP()")
	modifierCmd.Flags().String("style", "", "Style()")
	modifierCmd.Flags().String("stylef", "", "StyleF()")
	modifierCmd.Flags().String("styler", "", "StyleR()")
	modifierCmd.Flags().String("suffix", "", "Suffix()")
	modifierCmd.Flags().String("suppress", "", "Suppress()")
	modifierCmd.Flags().String("tag", "", "Tag()")
	modifierCmd.Flags().String("tagf", "", "TagF()")
	modifierCmd.Flags().String("timeout", "", "Timeout()")
	modifierCmd.Flags().String("uniquelist", "", "UniqueList()")
	modifierCmd.Flags().String("uniquelistf", "", "UniqueListF()")
	modifierCmd.Flags().String("unless", "", "Unless()")
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
			)
		}).Cache(5 * time.Second),
		"cache-key": carapace.ActionMultiParts("/", func(c carapace.Context) carapace.Action {
			switch len(c.Parts) {
			case 0:
				return carapace.ActionValues("one", "two").Suffix("/")
			case 1:
				return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
					return carapace.ActionValues(
						time.Now().Format("15:04:05"),
					)
				}).Cache(10*time.Second, key.String(c.Parts[0]))
			default:
				return carapace.ActionValues()
			}
		}),
		"chdir":  carapace.ActionFiles().Chdir(os.TempDir()),
		"chdirf": carapace.ActionFiles().ChdirF(traverse.GitWorkTree),
		"filter": carapace.ActionValuesDescribed(
			"1", "one",
			"2", "two",
			"3", "three",
			"4", "four",
		).Filter("2", "4"),
		"filterargs": carapace.ActionValues(
			"one",
			"two",
			"three",
		).FilterArgs(),
		"filterparts": carapace.ActionMultiParts(",", func(c carapace.Context) carapace.Action {
			return carapace.ActionValues(
				"one",
				"two",
				"three",
			).FilterParts().Suffix(",")
		}),
		"list": carapace.ActionValues("one", "two", "three").List(","),
		"invoke": carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			switch {
			case strings.HasPrefix(c.Value, "file://"):
				c.Value = strings.TrimPrefix(c.Value, "file://")
			case strings.HasPrefix("file://", c.Value):
				c.Value = ""
			default:
				return carapace.ActionValues()
			}
			return carapace.ActionFiles().Invoke(c).Prefix("file://").ToA()
		}),
		"nospace": carapace.ActionValues(
			"one,",
			"two/",
			"three",
		).NoSpace(',', '/'),
		"timeout": carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			time.Sleep(3 * time.Second)
			return carapace.ActionValues("within timeout")
		}).Timeout(2*time.Second, carapace.ActionMessage("timeout exceeded")),
		"multiparts": carapace.ActionValues(
			"dir/subdir1/fileA.txt",
			"dir/subdir1/fileB.txt",
			"dir/subdir2/fileC.txt",
		).MultiParts("/"),
		"multipartsp": carapace.ActionStyledValuesDescribed(
			"keys/<key>", "key example", style.Default,
			"keys/<key>/<value>", "key/value example", style.Default,
			"styles/custom", "custom style", style.Of(style.Blue, style.Blink),
			"styles", "list", style.Yellow,
			"styles/<style>", "details", style.Default,
		).MultiPartsP("/", "<.*>", func(placeholder string, matches map[string]string) carapace.Action {
			switch placeholder {
			case "<key>":
				return carapace.ActionValues("key1", "key2")
			case "<style>":
				return carapace.ActionStyles()
			case "<value>":
				switch matches["<key>"] {
				case "key1":
					return carapace.ActionValues("val1", "val2")
				case "key2":
					return carapace.ActionValues("val3", "val4")
				default:
					return carapace.ActionValues()
				}
			default:
				return carapace.ActionValues()
			}
		}),
		"prefix": carapace.ActionFiles().Prefix("file://"),
		"retain": carapace.ActionValuesDescribed(
			"1", "one",
			"2", "two",
			"3", "three",
			"4", "four",
		).Retain("2", "4"),
		"shift": carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			return carapace.ActionMessage("%#v", c.Args)
		}).Shift(1),
		"split": carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			cmd := &cobra.Command{}
			carapace.Gen(cmd).Standalone()
			cmd.Flags().BoolP("bool", "b", false, "bool flag")
			cmd.Flags().StringP("string", "s", "", "string flag")

			carapace.Gen(cmd).FlagCompletion(carapace.ActionMap{
				"string": carapace.ActionValues("one", "two", "three with space"),
			})

			carapace.Gen(cmd).PositionalCompletion(
				carapace.ActionValues("pos1", "positional1"),
				carapace.ActionFiles(),
			)

			return carapace.ActionExecute(cmd)
		}).Split(),
		"splitp": carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			cmd := &cobra.Command{}
			carapace.Gen(cmd).Standalone()
			cmd.Flags().BoolP("bool", "b", false, "bool flag")
			cmd.Flags().StringP("string", "s", "", "string flag")

			carapace.Gen(cmd).FlagCompletion(carapace.ActionMap{
				"string": carapace.ActionValues("one", "two", "three with space"),
			})

			carapace.Gen(cmd).PositionalCompletion(
				carapace.ActionValues("pos1", "positional1"),
				carapace.ActionFiles(),
			)

			return carapace.ActionExecute(cmd)
		}).SplitP(),
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
		"uniquelistf": carapace.ActionMultiPartsN(":", 2, func(c carapace.Context) carapace.Action {
			switch len(c.Parts) {
			case 0:
				return carapace.ActionValues("one", "two", "three")
			default:
				return carapace.ActionValues("1", "2", "3")
			}
		}).UniqueListF(",", func(s string) string {
			return strings.SplitN(s, ":", 2)[0]
		}),
		"unless": carapace.ActionValues(
			"./local",
			"~/home",
			"/abs",
			"one",
			"two",
			"three",
		).Unless(condition.CompletingPath),
		"usage": carapace.ActionValues().Usage("explicit usage"),
	})

	carapace.Gen(modifierCmd).PositionalCompletion(
		carapace.ActionValues().Usage("explicit positional usage"),
	)
}
