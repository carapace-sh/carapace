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
			if hasDescriptions {
				vals[index] = fmt.Sprintf("%v\t%v\002", val.Value, val.Display)
			} else {
				vals[index] = fmt.Sprintf("%v\t%v", val.Value, val.Display)
			}
		} else {
			vals[index] = fmt.Sprintf("%v\t%v\002 %v-- %v", val.Value, val.Display, strings.Repeat(" ", maxLength-len(val.Display)), val.TrimmedDescription())
		}
	}

	if len(zstyles) > 1000 { // TODO disable styling for large amount of values (bad performance)
		zstyles = make([]string, 0)
	}
	zstyles = append(zstyles, fmt.Sprintf("=(#b)(*)(\002*)=0=%v=%v", style.SGR(style.Carapace.Value), style.SGR(style.Carapace.Description)))
	return fmt.Sprintf("%v\n%v", strings.Join(zstyles, ":"), strings.Join(vals, "\n"))
}

// formatZstyle creates a zstyle matcher for given display stings.
// `compadd -l` (one per line) accepts ansi escape sequences in display value but it seems in tabular view these are removed.
// To ease matching in list mode, the display values have a hidden `\002` suffix.
func formatZstyle(s, _style string) string {
	if sgr := style.SGR(_style); sgr != "" {
		return fmt.Sprintf("=(#b)(%v)(\002*|)=0=%v=%v", strings.Replace(s, "#", `\#`, -1), sgr, style.SGR(style.Carapace.Description))
	}
	return ""
}
