package bash_ble

import (
	"fmt"
	"strings"

	"github.com/rsteube/carapace/internal/common"
)

// ActionRawValues formats values for bash_ble
func ActionRawValues(currentWord string, nospace common.SuffixMatcher, values common.RawValues) string {
	filtered := values.FilterPrefix(currentWord)
	vals := make([]string, len(filtered))
	for index, val := range filtered {
		suffix := " "
		if nospace.Matches(val.Value) {
			suffix = ""
		}
		vals[index] = fmt.Sprintf("%v\t%v\x1c%v\x1c%v\x1c%v", val.Value, val.Display, "", suffix, val.TrimmedDescription())
	}
	return strings.Join(vals, "\n")
}
