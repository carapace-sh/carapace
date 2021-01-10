package zsh

import (
	"fmt"
	"strings"

	"github.com/rsteube/carapace/internal/common"
)

var sanitizer = strings.NewReplacer(
	`$`, ``,
	"`", ``,
	"\n", ``,
	`\`, ``,
	`"`, ``,
	`'`, ``,
	"`", ``,
	`|`, ``,
	`>`, ``,
	`<`, ``,
	`&`, ``,
	`(`, ``,
	`)`, ``,
	`;`, ``,
	`#`, ``,
	`[`, `\[`,
	`]`, `\]`,
)

func Sanitize(values ...string) []string {
	sanitized := make([]string, len(values))
	for index, value := range values {
		sanitized[index] = sanitizer.Replace(value)
	}
	return sanitized
}

func EscapeSpace(s string) string {
	return strings.Replace(s, " ", `\ `, -1)
}

func ActionRawValues(callbackValue string, values ...common.RawValue) string {
	filtered := make([]common.RawValue, 0)

	for _, r := range values {
		if strings.HasPrefix(r.Value, callbackValue) {
			filtered = append(filtered, r)
		}
	}

	vals := make([]string, len(filtered))
	displays := make([]string, len(filtered))
	for index, val := range filtered {
		vals[index] = fmt.Sprintf("'%v'", EscapeSpace(sanitizer.Replace(val.Value)))
		if strings.TrimSpace(val.Description) == "" {
			displays[index] = fmt.Sprintf("'%v'", sanitizer.Replace(val.Display))
		} else {
			displays[index] = fmt.Sprintf("'%v (%v)'", sanitizer.Replace(val.Display), sanitizer.Replace(val.Description))
		}
	}
	return fmt.Sprintf("{local _comp_desc=(%v);compadd -Q -S '' -d _comp_desc -- %v}", strings.Join(displays, " "), strings.Join(vals, " "))
}
