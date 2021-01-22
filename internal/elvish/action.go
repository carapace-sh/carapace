package elvish

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rsteube/carapace/internal/common"
)

var sanitizer = strings.NewReplacer(
	"\n", ``,
	"\r", ``,
)

func sanitize(values []common.RawValue) []common.RawValue {
	for index, v := range values {
		(&values[index]).Value = sanitizer.Replace(v.Value)
		(&values[index]).Display = sanitizer.Replace(v.Display)
		(&values[index]).Description = sanitizer.Replace(v.Description)
	}
	return values
}

type complexCandidate struct {
	Value   string
	Display string
}

func ActionRawValues(callbackValue string, values ...common.RawValue) string {
	vals := make([]complexCandidate, len(values))
	for index, val := range sanitize(values) {
		// TODO have a look at this again later: seems elvish does a good job quoting any problematic characterS so the sanitize step was removed
		if val.Description == "" {
			vals[index] = complexCandidate{Value: val.Value, Display: val.Display}
		} else {
			vals[index] = complexCandidate{Value: val.Value, Display: fmt.Sprintf(`%v (%v)`, val.Display, val.Description)}
		}
	}
	m, _ := json.Marshal(vals)
	return string(m)
}
