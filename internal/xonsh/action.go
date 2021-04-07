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

func Sanitize(values ...string) []string {
	sanitized := make([]string, len(values))
	for index, value := range values {
		sanitized[index] = sanitizer.Replace(value)
	}
	return sanitized
}

type richCompletion struct {
	Value       string
	Display     string
	Description string
}

func ActionRawValues(callbackValue string, nospace bool, values ...common.RawValue) string {
	filtered := make([]common.RawValue, 0)

	for _, r := range values {
		if strings.HasPrefix(r.Value, callbackValue) {
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

		vals[index] = richCompletion{Value: val.Value, Display: val.Display, Description: val.Description}
	}
	m, _ := json.Marshal(vals)
	return string(m)
}
