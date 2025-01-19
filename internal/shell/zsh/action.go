package zsh

import (
	"fmt"
	"regexp"
	"strings"

	shlex "github.com/carapace-sh/carapace-shlex"
	"github.com/carapace-sh/carapace/internal/common"
	"github.com/carapace-sh/carapace/internal/env"
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

type state int

const (
	DEFAULT_STATE state = iota
	// Word starts with `"`.
	// Values need to end with `"` as well.
	// Weirdly regardless whether there are additional quotes within the word.
	QUOTING_ESCAPING_STATE
	// Word starts with `'`.
	// Values need to end with `'` as well.
	// Weirdly regardless whether there are additional quotes within the word.
	QUOTING_STATE
	// Word starts and ends with quotes.
	// Space suffix somehow ends up within the quotes.
	//    `"action"<TAB>`
	//    `"action "<CURSOR>`
	// Workaround for now is to force nospace.
	FULLY_QUOTED_STATE
)

// ActionRawValues formats values for zsh
func ActionRawValues(currentWord string, meta common.Meta, values common.RawValues) string {
	splitted, err := shlex.Split(env.Compline())
	state := DEFAULT_STATE
	if err == nil {
		rawValue := splitted.CurrentToken().RawValue
		switch {
		case regexp.MustCompile(`^'$|^'.*[^']$`).MatchString(rawValue):
			state = QUOTING_STATE
		case regexp.MustCompile(`^"$|^".*[^"]$`).MatchString(rawValue):
			state = QUOTING_ESCAPING_STATE
		case regexp.MustCompile(`^".*"$|^'.*'$`).MatchString(rawValue):
			state = FULLY_QUOTED_STATE
		}
	}

	for index, value := range values {
		switch value.Tag {
		case "shorthand flags", "longhand flags":
			values[index].Tag = "flags" // join to single tag group for classic zsh side-by-side view
		}
	}

	tagGroup := make([]string, 0)
	values.EachTag(func(tag string, values common.RawValues) {
		vals := make([]string, len(values))
		displays := make([]string, len(values))
		for index, val := range values {
			value := sanitizer.Replace(val.Value)
			value = quoteValue(value)
			value = strings.ReplaceAll(value, `\`, `\\`) // TODO find out why `_describe` needs another backslash
			value = strings.ReplaceAll(value, `:`, `\:`) // TODO find out why `_describe` needs another backslash

			switch state {
			// TODO depending on state value needs to be formatted differently
			// TODO backspace strings are currently an issue
			case QUOTING_STATE:
				value = value + `'`
			case QUOTING_ESCAPING_STATE:
				value = value + `"`
			}

			if !meta.Nospace.Matches(val.Value) && state != FULLY_QUOTED_STATE {
				value += " "
			}

			display := sanitizer.Replace(val.Display)
			display = strings.ReplaceAll(display, `\`, `\\`) // TODO find out why `_describe` needs another backslash
			display = strings.ReplaceAll(display, `:`, `\:`) // TODO find out why `_describe` needs another backslash
			description := sanitizer.Replace(val.Description)

			vals[index] = value

			if strings.TrimSpace(description) == "" {
				displays[index] = display
			} else {
				displays[index] = fmt.Sprintf("%v:%v", display, description)
			}
		}
		tagGroup = append(tagGroup, strings.Join([]string{tag, strings.Join(displays, "\n"), strings.Join(vals, "\n")}, "\003"))
	})
	return fmt.Sprintf("%v\001%v\001%v\001", zstyles{values}.Format(), message{meta}.Format(), strings.Join(tagGroup, "\002")+"\002")
}
