package xonsh

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rsteube/carapace/internal/common"
)

var sanitizer = strings.NewReplacer( // TODO
	"\n", ``,
	"\t", ``,
	`'`, `\'`,
)

type richCompletion struct {
	Value       string
	Display     string
	Description string
}

// ActionRawValues formats values for xonsh.
func ActionRawValues(currentWord string, nospace bool, values common.RawValues) string {
	values = common.AddMessageToValues(currentWord, values)

	filtered := make([]common.RawValue, 0)

	for _, r := range values {
		if strings.HasPrefix(r.Value, currentWord) {
			filtered = append(filtered, r)
		}
	}

	vals := make([]richCompletion, len(filtered))
	for index, val := range filtered {
		val.Value = sanitizer.Replace(val.Value)

		if strings.ContainsAny(val.Value, ` ()[]{}*$?\"|<>&;#`+"`") {
			if strings.Contains(val.Value, `\`) {
				val.Value = fmt.Sprintf("r'%v'", val.Value) // backslash needs raw string
			} else {
				val.Value = fmt.Sprintf("'%v'", val.Value)
			}
		}

		if !nospace {
			val.Value = val.Value + " "
		}

		vals[index] = richCompletion{Value: val.Value, Display: val.Display, Description: val.TrimmedDescription()}
	}
	m, _ := json.Marshal(vals)
	return string(m)
}
