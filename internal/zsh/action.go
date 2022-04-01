package zsh

import (
	"fmt"
	"strings"

	"github.com/rsteube/carapace/internal/common"
	"github.com/rsteube/carapace/pkg/style"
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
			vals[index] = fmt.Sprintf("%v\t%v", val.Value, format(val.Display, val.Style))
		} else {
			vals[index] = fmt.Sprintf("%v\t%v %v-- %v", val.Value, format(val.Display, val.Style), strings.Repeat(" ", maxLength-len(val.Display)), val.TrimmedDescription())
		}
	}
	return strings.Join(vals, "\n")
}

func format(s, _style string) string {
	switch _style {
	case style.Red:
		return fmt.Sprintf("\033[31m%v\033[0m", s)
	case style.Green:
		return fmt.Sprintf("\033[32m%v\033[0m", s)
	case style.Yellow:
		return fmt.Sprintf("\033[33m%v\033[0m", s)
	case style.Blue:
		return fmt.Sprintf("\033[34m%v\033[0m", s)
	case style.Magenta:
		return fmt.Sprintf("\033[35m%v\033[0m", s)
	case style.Cyan:
		return fmt.Sprintf("\033[36m%v\033[0m", s)
	case style.BrightBlack:
		return fmt.Sprintf("\033[90m%v\033[0m", s)
	default:
		return s
	}
}
