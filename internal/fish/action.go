package fish

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
		// TODO sanitize
		//vals[index] = strings.Replace(val, " ", `\ `, -1)
		vals[index] = sanitizer.Replace(val.Value) + "\t" + sanitizer.Replace(val.Description)
	}
	return fmt.Sprintf(`echo -e "%v"`, strings.Join(vals, `\n`))
}
