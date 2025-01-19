package zsh

import (
	"fmt"
	"strings"

	"github.com/carapace-sh/carapace/internal/common"
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
	for index, value := range values {
		switch value.Tag {
		case "shorthand flags", "longhand flags":
			values[index].Tag = "flags" // join to single tag group for classic zsh side-by-side view
		}
	}

	tagGroup := make([]string, 0)
	values.EachTag(func(tag string, values common.RawValues) {
		vals := make([]string, 0, len(values))
		displays := make([]string, 0, len(values))
		valsNospace := make([]string, 0, len(values))
		displaysNospace := make([]string, 0, len(values))

		for _, val := range values {
			value := sanitizer.Replace(val.Value)
			value = quoteValue(value)
			value = strings.ReplaceAll(value, `\`, `\\`) // TODO find out why `_describe` needs another backslash
			value = strings.ReplaceAll(value, `:`, `\:`) // TODO find out why `_describe` needs another backslash

			display := sanitizer.Replace(val.Display)
			display = strings.ReplaceAll(display, `\`, `\\`) // TODO find out why `_describe` needs another backslash
			display = strings.ReplaceAll(display, `:`, `\:`) // TODO find out why `_describe` needs another backslash

			description := sanitizer.Replace(val.Description)
			if strings.TrimSpace(description) != "" {
				display = fmt.Sprintf("%v:%v", display, description)
			}

			switch {
			case meta.Nospace.Matches(val.Value): // checks unmodified value
				valsNospace = append(valsNospace, value)
				displaysNospace = append(displaysNospace, display)
			default:
				vals = append(vals, value)
				displays = append(displays, display)
			}
		}
		tagGroup = append(tagGroup, strings.Join(
			[]string{
				tag,
				strings.Join(displays, "\n"),
				strings.Join(vals, "\n"),
				strings.Join(displaysNospace, "\n"),
				strings.Join(valsNospace, "\n"),
			},
			"\003"))
	})
	return fmt.Sprintf("%v\001%v\001%v\001", zstyles{values}.Format(), message{meta}.Format(), strings.Join(tagGroup, "\002")+"\002")
}
