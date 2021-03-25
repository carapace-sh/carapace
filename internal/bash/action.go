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

func Sanitize(values ...string) []string {
	sanitized := make([]string, len(values))
	for index, value := range values {
		sanitized[index] = sanitizer.Replace(value)
	}
	return sanitized
}

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

func ActionRawValues(callbackValue string, values ...common.RawValue) string {
	filtered := make([]common.RawValue, 0)

	lastSegment := callbackValue // last segment of callbackValue split by COMP_WORDBREAKS

	for _, r := range values {
		if strings.HasPrefix(r.Value, callbackValue) {
			// TODO optimize
			if wordbreaks, ok := os.LookupEnv("COMP_WORDBREAKS"); ok {
				wordbreaks = strings.Replace(wordbreaks, " ", "", -1)
				if index := strings.LastIndexAny(callbackValue, wordbreaks); index != -1 {
					r.Value = strings.TrimPrefix(r.Value, callbackValue[:index+1])
					lastSegment = callbackValue[index+1:]
				}
			}
			filtered = append(filtered, r)
		}
	}

	if len(filtered) > 1 && commonDisplayPrefix(filtered...) != "" {
		// When all display values have the same prefix bash will insert is as partial completion (which skips prefixes/formatting).
		if valuePrefix := commonValuePrefix(filtered...); lastSegment != valuePrefix {
			// replace values with common value prefix (`\001` is removed in snippet and compopt nospace will be set)
			filtered = common.RawValuesFrom(commonValuePrefix(filtered...) + "\001")
		} else {
			// prevent insertion of partial display values by prefixing one with space
			filtered[0].Display = " " + filtered[0].Display
		}
	}

	vals := make([]string, len(filtered))
	for index, val := range filtered {
		if len(filtered) == 1 {
			vals[index] = quoter.Replace(sanitizer.Replace(val.Value))
		} else {
			if val.Description != "" {
				vals[index] = fmt.Sprintf("%v (%v)", val.Display, sanitizer.Replace(val.Description))
			} else {
				vals[index] = val.Display
			}
		}
	}
	return strings.Join(vals, "\n")
}
