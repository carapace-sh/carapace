package zsh

import (
	"fmt"
	"regexp"
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
	hasDescriptions := false
	for _, r := range values {
		if strings.HasPrefix(r.Value, currentWord) {
			filtered = append(filtered, r)
			if length := len(r.Display); length > maxLength {
				maxLength = length
			}
			hasDescriptions = hasDescriptions || r.Description != ""
		}
	}

	zstyles := make([]string, 0)
	vals := make([]string, len(filtered))
	for index, val := range filtered {
		val.Value = sanitizer.Replace(val.Value)
		if nospace {
			val.Value = val.Value + "\001"
		}
		val.Display = sanitizer.Replace(val.Display)
		val.Description = sanitizer.Replace(val.Description)

		if zstyle := formatZstyle(val.Display, val.Style); zstyle != "" {
			zstyles = append(zstyles, zstyle)
		}

		if strings.TrimSpace(val.Description) == "" {
			vals[index] = fmt.Sprintf("%v\t%v", val.Value, val.Display)
		} else {
			vals[index] = fmt.Sprintf("%v\t%v\002 %v-- %v", val.Value, val.Display, strings.Repeat(" ", maxLength-len(val.Display)), val.TrimmedDescription())
		}
	}
	return fmt.Sprintf(":%v\n%v", strings.Join(zstyles, ":"), strings.Join(vals, "\n")) // first line is intentionally never empty (single `:`) for snippet
}

// formatZstyle creates a zstyle matcher for given display stings.
// `compadd -l` (one per line) accepts ansi escape sequences in display value but it seems in tabular view these are removed.
// To ease matching in list mode, the display values have a hidden `\002` suffix.
func formatZstyle(s, _style string) string {
	result := make([]string, 0)

	reColor256 := regexp.MustCompile(`^color(?P<number>\d+)$`)
	for _, word := range strings.Split(_style, " ") {
		switch word {
		case style.Black:
			result = append(result, "30")
		case style.Red:
			result = append(result, "31")
		case style.Green:
			result = append(result, "32")
		case style.Yellow:
			result = append(result, "33")
		case style.Blue:
			result = append(result, "34")
		case style.Magenta:
			result = append(result, "35")
		case style.Cyan:
			result = append(result, "36")
		case style.White:
			result = append(result, "37")

		case style.BrightBlack:
			result = append(result, "90")
		case style.BrightRed:
			result = append(result, "91")
		case style.BrightGreen:
			result = append(result, "92")
		case style.BrightYellow:
			result = append(result, "93")
		case style.BrightBlue:
			result = append(result, "94")
		case style.BrightMagenta:
			result = append(result, "95")
		case style.BrightCyan:
			result = append(result, "96")
		case style.BrightWhite:
			result = append(result, "97")

		case style.BgBlack:
			result = append(result, "40")
		case style.BgRed:
			result = append(result, "41")
		case style.BgGreen:
			result = append(result, "42")
		case style.BgYellow:
			result = append(result, "43")
		case style.BgBlue:
			result = append(result, "44")
		case style.BgMagenta:
			result = append(result, "45")
		case style.BgCyan:
			result = append(result, "46")
		case style.BgWhite:
			result = append(result, "47")

		case style.BgBrightBlack:
			result = append(result, "100")
		case style.BgBrightRed:
			result = append(result, "101")
		case style.BgBrightGreen:
			result = append(result, "102")
		case style.BgBrightYellow:
			result = append(result, "103")
		case style.BgBrightBlue:
			result = append(result, "104")
		case style.BgBrightMagenta:
			result = append(result, "105")
		case style.BgBrightCyan:
			result = append(result, "106")
		case style.BgBrightWhite:
			result = append(result, "107")

		case style.Bold:
			result = append(result, "1")
		case style.Dim:
			result = append(result, "2")
		case style.Italic:
			result = append(result, "3")
		case style.Underlined:
			result = append(result, "4")
		case style.Blink:
			result = append(result, "5")
		case style.Inverse:
			result = append(result, "7")
		default:
			if reColor256.MatchString(word) {
				result = append(result, fmt.Sprintf("38;5;%v", reColor256.FindStringSubmatch(word)[1]))
			}
		}
	}

	if len(result) > 0 {
		return fmt.Sprintf("=(#b)(%v)((\002*|))=0=%v", strings.Replace(regexp.QuoteMeta(s), "#", `\#`, -1), strings.Join(result, ";"))
	}
	return ""
}
