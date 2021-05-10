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
		(&values[index]).Description = sanitizer.Replace(v.TrimmedDescription())
	}
	return values
}

type complexCandidate struct {
	Value      string
	Display    string
	CodeSuffix string
}

// ActionRawValues formats values for elvish
func ActionRawValues(currentWord string, nospace bool, values ...common.RawValue) string {
	suffix := " "
	if nospace {
		suffix = ""
	}

	vals := make([]complexCandidate, len(values))
	for index, val := range sanitize(values) {
		// TODO have a look at this again later: seems elvish does a good job quoting any problematic characterS so the sanitize step was removed
		if val.Description == "" {
			vals[index] = complexCandidate{Value: val.Value, Display: val.Display, CodeSuffix: suffix}
		} else {
			vals[index] = complexCandidate{Value: val.Value, Display: fmt.Sprintf(`%v (%v)`, val.Display, val.Description), CodeSuffix: suffix}
		}
	}
	m, _ := json.Marshal(vals)
	return string(m)
}
