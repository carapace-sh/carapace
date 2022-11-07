package nushell

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/rsteube/carapace/internal/common"
)

type record struct {
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
}

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

// ActionRawValues formats values for nushell.
func ActionRawValues(currentWord string, nospace bool, values common.RawValues) string {
	values = common.AddMessageToValues(currentWord, values)

	filtered := values.FilterPrefix(currentWord)
	sort.Sort(common.ByDisplay(filtered))

	vals := make([]record, len(filtered))
	for index, val := range sanitize(filtered) {
		if strings.ContainsAny(val.Value, ` {}()[]<>$&"|;#\`+"`") {
			val.Value = fmt.Sprintf("'%v'", val.Value)
		}

		if !nospace {
			val.Value = val.Value + " "
		}

		vals[index] = record{Value: val.Value, Description: val.TrimmedDescription()}
	}
	m, _ := json.Marshal(vals)
	return string(m)
}
