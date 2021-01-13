package elvish

import (
	"encoding/json"
	"fmt"
	"github.com/rsteube/carapace/internal/common"
)

type complexCandidate struct {
	Value   string
	Display string
}

func ActionRawValues(callbackValue string, values ...common.RawValue) string {
	vals := make([]complexCandidate, len(values))
	for index, val := range values {
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
