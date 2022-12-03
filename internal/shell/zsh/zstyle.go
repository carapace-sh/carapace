package zsh

import (
	"fmt"
	"strings"

	"github.com/rsteube/carapace/internal/common"
	"github.com/rsteube/carapace/pkg/style"
	"github.com/rsteube/carapace/third_party/github.com/elves/elvish/pkg/ui"
)

type zstyles struct {
	rawValues common.RawValues
}

func (z zstyles) descriptionSGR() string {
	if s := style.Carapace.Description; s != "" && ui.ParseStyling(s) != nil {
		return style.SGR(s)
	}
	return style.SGR(style.Default)
}

func (z zstyles) valueSGR(val common.RawValue) string {
	if val.Style != "" && ui.ParseStyling(val.Style) != nil {
		return style.SGR(val.Style)
	}

	if ui.ParseStyling(style.Carapace.Value) != nil {
		return style.SGR(style.Carapace.Value)
	}
	return style.SGR(style.Default)

}

func (z zstyles) hasAliases() bool {
	descriptions := make(map[string]bool)
	for _, val := range z.rawValues {
		if _, exists := descriptions[val.Description]; exists && val.Description != "" {
			return true
		}
		descriptions[val.Description] = true
	}
	return false
}

func (z zstyles) Format() string {
	replacer := strings.NewReplacer(
		"#", `\#`,
		"*", `\*`,
		"(", `\(`,
		")", `\)`,
		"[", `\[`,
		"]", `\]`,
		"|", `\|`,
		"~", `\~`,
	)

	hasAliases := z.hasAliases()
	formatted := make([]string, 0)
	if len(z.rawValues) < 1000 { // disable styling for large amount of values (bad performance)
		for _, val := range z.rawValues {
			// TODO this might need to be handled differently regarding tags/groups (e.g. done for each tag)
			pattern := "=(#b)(%v)( * -- *)=0=%v=%v"  // match value with description
			if val.Description == "" || hasAliases { // different behaviour in `_describe` when values are on the same line
				pattern = "=(#b)(%v)()=0=%v=%v" // only match value
			}

			formatted = append(formatted, fmt.Sprintf(pattern, replacer.Replace(val.Display), z.valueSGR(val), z.descriptionSGR()))
		}
	}
	formatted = append(formatted, fmt.Sprintf("=(#b)(%v)=0=%v", "-- *", z.descriptionSGR()))
	return strings.Join(formatted, ":")
}
