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

func ActionRawValues(callbackValue string, nospace bool, values ...common.RawValue) string {
	filtered := make([]common.RawValue, 0)

	for _, r := range values {
		// TODO should rather access callbackvalue (circular dependency) - seems to work though so good enough for now
		if strings.HasPrefix(r.Value, callbackValue) {
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
				vals[index] = fmt.Sprintf("%v (%v)", val.Value, sanitizer.Replace(val.Description))
			} else {
				vals[index] = val.Value
			}
		}
	}
	return strings.Join(vals, "\n")
}
