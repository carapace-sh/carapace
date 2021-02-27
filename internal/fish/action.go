package fish

import (
	"github.com/rsteube/carapace/internal/common"
	"strings"
)

var sanitizer = strings.NewReplacer(
	`$`, ``,
	"`", ``,
	"\n", ``,
	`\`, ``,
	`"`, ``,
	`'`, ``,
	`(`, `[`,
	`)`, `]`,
)

func Sanitize(values ...string) []string {
	sanitized := make([]string, len(values))
	for index, value := range values {
		sanitized[index] = sanitizer.Replace(value)
	}
	return sanitized
}

func ActionRawValues(callbackValues string, values ...common.RawValue) string {
	vals := make([]string, len(values))
	for index, val := range values {
		vals[index] = sanitizer.Replace(val.Value) + "\t" + sanitizer.Replace(val.Description)
	}
	return strings.Join(vals, "\n")
}
