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
	`'`, `'\''`,
)

func Sanitize(values ...string) []string {
	sanitized := make([]string, len(values))
	for index, value := range values {
		sanitized[index] = sanitizer.Replace(value)
	}
	return sanitized
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
		val.Value = sanitizer.Replace(val.Value)
		val.Display = sanitizer.Replace(val.Display)
		val.Description = sanitizer.Replace(val.Description)

		if strings.TrimSpace(val.Description) == "" {
			vals[index] = fmt.Sprintf("%v\t%v", val.Value, val.Display)
		} else {
			vals[index] = fmt.Sprintf("%v\t%v (%v)", val.Value, val.Display, val.Description)
		}
	}
	return strings.Join(vals, "\n")
}
