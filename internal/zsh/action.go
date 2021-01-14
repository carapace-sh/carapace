package zsh

import (
	"fmt"
	"strings"

	"github.com/rsteube/carapace/internal/common"
)

var sanitizer = strings.NewReplacer(
	`$`, ``,
	"`", ``,
	"\n", ``,
	`\`, ``,
	`"`, ``,
	`'`, ``,
	"`", ``,
	`|`, ``,
	`>`, ``,
	`<`, ``,
	`&`, ``,
	`(`, ``,
	`)`, ``,
	`;`, ``,
	`#`, ``,
	`[`, `\[`,
	`]`, `\]`,
)

func Sanitize(values ...string) []string {
	sanitized := make([]string, len(values))
	for index, value := range values {
		sanitized[index] = sanitizer.Replace(value)
	}
	return sanitized
}

func EscapeSpace(s string) string {
	return strings.Replace(s, " ", `\ `, -1)
}

func ActionRawValues(callbackValue string, values ...common.RawValue) string {
	filtered := make([]common.RawValue, 0)

	for _, r := range values {
		if strings.HasPrefix(r.Value, callbackValue) {
			filtered = append(filtered, r)
		}
	}

	vals := make([]string, len(filtered))
	for index, val := range filtered {
		if strings.TrimSpace(val.Description) == "" {
			vals[index] = fmt.Sprintf("%v\t%v", EscapeSpace(sanitizer.Replace(val.Value)), sanitizer.Replace(val.Display))
		} else {
			vals[index] = fmt.Sprintf("%v\t%v (%v)", EscapeSpace(sanitizer.Replace(val.Value)), sanitizer.Replace(val.Display), sanitizer.Replace(val.Description))
		}
	}
	return strings.Join(vals, "\n")
}
