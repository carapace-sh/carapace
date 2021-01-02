package xonsh

import (
	"fmt"
	"strings"

	"github.com/rsteube/carapace/internal/common"
)

var sanitizer = strings.NewReplacer( // TODO
	`$`, ``,
	"`", ``,
	"\n", ``,
	`\`, ``,
	`"`, ``,
	`'`, ``,
	`|`, ``,
	`>`, ``,
	`<`, ``,
	`&`, ``,
	`(`, ``,
	`)`, ``,
	`;`, ``,
	`#`, ``,
)

func Sanitize(values ...string) []string {
	sanitized := make([]string, len(values))
	for index, value := range values {
		sanitized[index] = sanitizer.Replace(value)
	}
	return sanitized
}

func ActionRawValues(callbackValue string, values ...common.RawValue) string {
	filtered := make([]common.RawValue, 0)

	for _, r := range values {
		if strings.HasPrefix(r.Value, callbackValue) {
			filtered = append(filtered, r)
		}
	}

	vals := make([]string, len(filtered))
	for index, val := range filtered {
		vals[index] = fmt.Sprintf(`  RichCompletion('%v', display='%v', description='%v', prefix_len=0),`, strings.Replace(sanitizer.Replace(val.Value), " ", `\\ `, -1), sanitizer.Replace(val.Display), sanitizer.Replace(val.Description))
	}
	return fmt.Sprintf("{\n%v\n}", strings.Join(vals, "\n"))
}
