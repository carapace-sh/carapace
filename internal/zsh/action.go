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

// ActionRawValues formats values for zsh
func ActionRawValues(currentWord string, nospace bool, values common.RawValues) string {
	filtered := make([]common.RawValue, 0)

	maxLength := 0
	for _, r := range values {
		if strings.HasPrefix(r.Value, currentWord) {
			filtered = append(filtered, r)
			if length := len(r.Display); length > maxLength {
				maxLength = length
			}
		}
	}

	vals := make([]string, len(filtered))
	for index, val := range filtered {
		val.Value = sanitizer.Replace(val.Value)
		if nospace {
			val.Value = val.Value + "\001"
		}
		val.Display = sanitizer.Replace(val.Display)
		val.Description = sanitizer.Replace(val.Description)

		if strings.TrimSpace(val.Description) == "" {
			vals[index] = fmt.Sprintf("%v\t%v", val.Value, val.Display)
		} else {
			vals[index] = fmt.Sprintf("%v\t%v %v-- %v", val.Value, val.Display, strings.Repeat(" ", maxLength-len(val.Display)), val.TrimmedDescription())
		}
	}
	return strings.Join(vals, "\n")
}
