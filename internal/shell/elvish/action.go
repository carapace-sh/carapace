package elvish

import (
	"encoding/json"
	"strings"

	"github.com/rsteube/carapace/internal/common"
	"github.com/rsteube/carapace/pkg/style"
	"github.com/rsteube/carapace/third_party/github.com/elves/elvish/pkg/ui"
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
	Value            string
	Display          string
	Description      string
	DescriptionStyle string
	CodeSuffix       string
	Style            string
}

// ActionRawValues formats values for elvish.
func ActionRawValues(currentWord string, nospace bool, values common.RawValues) string {
	values = common.AddMessageToValues(currentWord, values)

	suffix := " "
	if nospace {
		suffix = ""
	}

	valueStyle := "default"
	if s := style.Carapace.Value; s != "" && ui.ParseStyling(s) != nil {
		valueStyle = s
	}

	descriptionStyle := "default"
	if s := style.Carapace.Description; s != "" && ui.ParseStyling(s) != nil {
		descriptionStyle = s
	}

	vals := make([]complexCandidate, len(values))
	for index, val := range sanitize(values) {
		if val.Style == "" || ui.ParseStyling(val.Style) == nil {
			val.Style = valueStyle
		}
		vals[index] = complexCandidate{Value: val.Value, Display: val.Display, Description: val.Description, CodeSuffix: suffix, Style: val.Style, DescriptionStyle: descriptionStyle}
	}
	m, _ := json.Marshal(vals)
	return string(m)
}
