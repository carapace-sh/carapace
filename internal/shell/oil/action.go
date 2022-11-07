package oil

import (
	"fmt"
	"strings"

	"github.com/rsteube/carapace/internal/common"
)

var sanitizer = strings.NewReplacer(
	"\n", ``,
	"\t", ``,
)

const nospaceIndicator = "\001"

// ActionRawValues formats values for oil.
func ActionRawValues(currentWord string, nospace bool, values common.RawValues) string {
	values = common.AddMessageToValues(currentWord, values)

	filtered := make([]common.RawValue, 0)

	for _, r := range values {
		// TODO should rather access callbackvalue (circular dependency) - seems to work though so good enough for now
		if strings.HasPrefix(r.Value, currentWord) {
			filtered = append(filtered, r)
		}
	}

	vals := make([]string, len(filtered))
	for index, val := range filtered {
		if nospace && !strings.HasSuffix(val.Value, nospaceIndicator) {
			val.Value = val.Value + nospaceIndicator
		}

		if len(filtered) == 1 {
			formattedVal := sanitizer.Replace(val.Value)
			vals[index] = formattedVal
		} else {
			if val.Description != "" {
				vals[index] = fmt.Sprintf("%v (%v)", val.Value, sanitizer.Replace(val.TrimmedDescription()))
			} else {
				vals[index] = val.Value
			}
		}
	}
	return strings.Join(vals, "\n")
}
