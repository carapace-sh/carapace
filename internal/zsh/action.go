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
	result := make([]string, 0)
	for _, word := range strings.Split(_style, " ") {
		switch word {
		case style.Red:
			result = append(result, "\033[31m")
		case style.Green:
			result = append(result, "\033[32m")
		case style.Yellow:
			result = append(result, "\033[33m")
		case style.Blue:
			result = append(result, "\033[34m")
		case style.Magenta:
			result = append(result, "\033[35m")
		case style.Cyan:
			result = append(result, "\033[36m")

		case style.BrightBlack:
			result = append(result, "\033[90m")
		case style.BrightRed:
			result = append(result, "\033[91m")
		case style.BrightGreen:
			result = append(result, "\033[92m")
		case style.BrightYellow:
			result = append(result, "\033[93m")
		case style.BrightBlue:
			result = append(result, "\033[94m")
		case style.BrightMagenta:
			result = append(result, "\033[95m")
		case style.BrightCyan:
			result = append(result, "\033[96m")

		case style.BgRed:
			result = append(result, "\033[41m")
		case style.BgGreen:
			result = append(result, "\033[42m")
		case style.BgYellow:
			result = append(result, "\033[43m")
		case style.BgBlue:
			result = append(result, "\033[44m")
		case style.BgMagenta:
			result = append(result, "\033[45m")
		case style.BgCyan:
			result = append(result, "\033[46m")

		case style.BgBrightBlack:
			result = append(result, "\033[100m")
		case style.BgBrightRed:
			result = append(result, "\033[101m")
		case style.BgBrightGreen:
			result = append(result, "\033[102m")
		case style.BgBrightYellow:
			result = append(result, "\033[103m")
		case style.BgBrightBlue:
			result = append(result, "\033[104m")
		case style.BgBrightMagenta:
			result = append(result, "\033[105m")
		case style.BgBrightCyan:
			result = append(result, "\033[106m")

		case style.Bold:
			result = append(result, "\033[1m")
		case style.Dim:
			result = append(result, "\033[2m")
		case style.Italic:
			result = append(result, "\033[3m")
		case style.Underlined:
			result = append(result, "\033[4m")
		case style.Blink:
			result = append(result, "\033[5m")
		case style.Inverse:
			result = append(result, "\033[7m")
		default:
		}
		result = append(result, s)
		result = append(result, "\033[0m")
	}
	return strings.Join(result, "")
}
