package bash

import (
	"fmt"
	"os"
	"strings"

	"github.com/rsteube/carapace/internal/common"
)

var sanitizer = strings.NewReplacer(
	"\n", ``,
	"\r", ``,
	"\t", ``,
)

var quoter = strings.NewReplacer(
	// seems readline provides quotation only for the filename completion (which would add suffixes) so do that here
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
	`\`, `\\`,
)

func commonPrefix(a, b string) string {
	i := 0
	for i < len(a) && i < len(b) && a[i] == b[i] {
		i++
	}
	return a[0:i]
}

func commonDisplayPrefix(values ...common.RawValue) (prefix string) {
	for index, val := range values {
		if index == 0 {
			prefix = val.Display
		} else {
			prefix = commonPrefix(prefix, val.Display)
		}
	}
	return
}

func commonValuePrefix(values ...common.RawValue) (prefix string) {
	for index, val := range values {
		if index == 0 {
			prefix = val.Value
		} else {
			prefix = commonPrefix(prefix, val.Value)
		}
	}
	return
}

const nospaceIndicator = "\001"

// ActionRawValues formats values for bash.
func ActionRawValues(currentWord string, nospace common.SuffixMatcher, values common.RawValues) string {
	lastSegment := currentWord // last segment of currentWord split by COMP_WORDBREAKS

	for _, r := range values {
		// TODO optimize
		if wordbreaks, ok := os.LookupEnv("COMP_WORDBREAKS"); ok {
			wordbreaks = strings.Replace(wordbreaks, " ", "", -1)
			if index := strings.LastIndexAny(currentWord, wordbreaks); index != -1 {
				r.Value = strings.TrimPrefix(r.Value, currentWord[:index+1])
				lastSegment = currentWord[index+1:]
			}
		}
	}

	if len(values) > 1 && commonDisplayPrefix(values...) != "" {
		// When all display values have the same prefix bash will insert is as partial completion (which skips prefixes/formatting).
		if valuePrefix := commonValuePrefix(values...); lastSegment != valuePrefix {
			// replace values with common value prefix (`\001` is removed in snippet and compopt nospace will be set)
			values = common.RawValuesFrom(commonValuePrefix(values...) + nospaceIndicator)
		} else {
			// prevent insertion of partial display values by prefixing one with space
			values[0].Display = " " + values[0].Display
		}
	}

	vals := make([]string, len(values))
	for index, val := range values {
		if nospace.Matches(val.Value) {
			val.Value = val.Value + nospaceIndicator
		}

		if len(values) == 1 {
			vals[index] = quoter.Replace(sanitizer.Replace(val.Value))
		} else {
			if val.Description != "" {
				vals[index] = fmt.Sprintf("%v (%v)", val.Display, sanitizer.Replace(val.TrimmedDescription()))
			} else {
				vals[index] = val.Display
			}
		}
	}
	return strings.Join(vals, "\n")
}
