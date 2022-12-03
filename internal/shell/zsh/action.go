package zsh

import (
	"fmt"
	"strings"

	"github.com/rsteube/carapace/internal/common"
)

var sanitizer = strings.NewReplacer(
	"\n", ``,
	"\r", ``,
	"\t", ``,
)

// TODO verify these are correct/complete (copied from bash)
var quoter = strings.NewReplacer(
	`\`, `\\`,
	`&`, `\&`,
	`<`, `\<`,
	`>`, `\>`,
	"`", "\\`",
	`'`, `\'`,
	`"`, `\"`,
	`{`, `\{`,
	`}`, `\}`,
	`$`, `\$`,
	`#`, `\#`,
	`|`, `\|`,
	`?`, `\?`,
	`(`, `\(`,
	`)`, `\)`,
	`;`, `\;`,
	` `, `\ `,
	`[`, `\[`,
	`]`, `\]`,
	`*`, `\*`,
	`~`, `\~`,
)

func quoteValue(s string) string {
	if strings.HasPrefix(s, "~/") || NamedDirectories.Matches(s) {
		return "~" + quoter.Replace(strings.TrimPrefix(s, "~")) // assume file path expansion
	}
	return quoter.Replace(s)
}

// ActionRawValues formats values for zsh
func ActionRawValues(currentWord string, meta common.Meta, values common.RawValues) string {
	maxLength := 0
	for _, r := range values {
		if length := len(r.Display); length > maxLength {
			maxLength = length
		}
	}

	vals := make([]string, len(values))
	for index, val := range values {
		val.Value = sanitizer.Replace(val.Value)
		val.Value = quoteValue(val.Value)
		val.Value = strings.ReplaceAll(val.Value, `\`, `\\`) // TODO find out why `_describe` needs another backslash
		val.Value = strings.ReplaceAll(val.Value, `:`, `\:`) // TODO find out why `_describe` needs another backslash
		if !meta.Nospace.Matches(val.Value) {
			val.Value = val.Value + " "
		}
		val.Display = sanitizer.Replace(val.Display)
		val.Display = strings.ReplaceAll(val.Display, `\`, `\\`) // TODO find out why `_describe` needs another backslash
		val.Display = strings.ReplaceAll(val.Display, `:`, `\:`) // TODO find out why `_describe` needs another backslash
		val.Description = sanitizer.Replace(val.Description)

		if strings.TrimSpace(val.Description) == "" {
			vals[index] = fmt.Sprintf("%v\t%v", val.Value, val.Display)
		} else {
			vals[index] = fmt.Sprintf("%v\t%v:%v", val.Value, val.Display, val.TrimmedDescription())
		}
	}
	return fmt.Sprintf(":%v\n%v\n%v", zstyles{values}.Format(), message{meta}.Format(), strings.Join(vals, "\n"))
}
