package zsh

import (
	"fmt"
	"strings"

	"github.com/rsteube/carapace/internal/common"
	"github.com/rsteube/carapace/pkg/style"
	"github.com/rsteube/carapace/third_party/github.com/elves/elvish/pkg/ui"
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

	valueStyle := "default"
	if s := style.Carapace.Value; s != "" && ui.ParseStyling(s) != nil {
		valueStyle = s
	}

	descriptionStyle := "default"
	if s := style.Carapace.Description; s != "" && ui.ParseStyling(s) != nil {
		descriptionStyle = s
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

		if val.Style == "" || ui.ParseStyling(val.Style) == nil {
			val.Style = valueStyle
		}

		if strings.TrimSpace(val.Description) == "" {
			vals[index] = fmt.Sprintf("%v\t%v", val.Value, val.Display)
			zstyles = append(zstyles, formatZstyle(fmt.Sprintf("(%v)()", zstyleQuoter.Replace(val.Display)), val.Style, descriptionStyle))
		} else {
			vals[index] = fmt.Sprintf("%v\t%v %v-- %v", val.Value, val.Display, strings.Repeat(" ", maxLength-len(val.Display)), val.TrimmedDescription())
			zstyles = append(zstyles, formatZstyle(fmt.Sprintf("(%v)[ ]*-- (%v)", zstyleQuoter.Replace(val.Display), zstyleQuoter.Replace(val.TrimmedDescription())), val.Style, descriptionStyle))
		}
	}

	if len(zstyles) > 1000 { // TODO disable styling for large amount of values (bad performance)
		zstyles = make([]string, 0)
	}
	return fmt.Sprintf(":%v\n%v", strings.Join(zstyles, ":"), strings.Join(vals, "\n"))
}

var zstyleQuoter = strings.NewReplacer(
	"#", `\#`,
	"*", `\*`,
	"(", `\(`,
	")", `\)`,
	"~", `\~`,
	"[", `\[`,
	"]", `\]`,
)

// formatZstyle creates a zstyle matcher for given display stings.
// `compadd -l` (one per line) accepts ansi escape sequences in display value but it seems in tabular view these are removed.
// To ease matching in list mode, the display values have a hidden `\002` suffix.
func formatZstyle(s, _styleValue, _styleDescription string) string {
	return fmt.Sprintf("=(#b)%v=0=%v=%v", s, style.SGR(_styleValue), style.SGR(_styleDescription))
}
