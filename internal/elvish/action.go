package elvish

import (
	"encoding/json"
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
		(&values[index]).Description = sanitizer.Replace(v.TrimmedDescription())
	}
	return values
}

type complexCandidate struct {
	Value       string
	Display     string
	Description string
	CodeSuffix  string
}

// ActionRawValues formats values for elvish
func ActionRawValues(currentWord string, nospace bool, values common.RawValues) string {
	suffix := " "
	if nospace {
		suffix = ""
	}

	vals := make([]complexCandidate, len(values))
	for index, val := range sanitize(values) {
		vals[index] = complexCandidate{Value: val.Value, Display: val.Display, Description: val.Description, CodeSuffix: suffix}
	}
	m, _ := json.Marshal(vals)
	return string(m)
}
