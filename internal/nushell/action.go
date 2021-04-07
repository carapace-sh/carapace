package nushell

import (
	"encoding/json"
	"fmt"
	"sort"
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

type suggestion struct {
	Value   string
	Display string
}

func ActionRawValues(callbackValue string, nospace bool, values common.RawValues) string {
	filtered := values.FilterPrefix(callbackValue)
	sort.Sort(common.ByDisplay(filtered))

	vals := make([]suggestion, len(filtered))
	for index, val := range sanitize(filtered) {
		if !nospace {
			val.Value = val.Value + " "
		}

		if val.Description == "" {
			vals[index] = suggestion{Value: val.Value, Display: val.Display}
		} else {
			vals[index] = suggestion{Value: val.Value, Display: fmt.Sprintf(`%v (%v)`, val.Display, val.Description)}
		}
	}
	m, _ := json.Marshal(vals)
	return string(m)
}
