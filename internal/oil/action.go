package oil

import (
	"fmt"
	"os"
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

func Callback(prefix string, uid string) string {
	return ""
}

func ActionDirectories() string {
	return ""
}

func ActionFiles(suffix string) string {
	return ""
}

func ActionRawValues(values ...common.RawValue) string {
	filtered := make([]common.RawValue, 0)

	for _, r := range values {
		// TODO should rather access callbackvalue (circular dependency) - seems to work though so good enough for now
		prefix := os.Args[len(os.Args)-1]
		if strings.HasPrefix(prefix, "-") && strings.Contains(prefix, "=") {
			prefix = strings.SplitN(prefix, "=", 2)[0]
			r.Value = prefix + "=" + r.Value
		}
		if strings.HasPrefix(r.Value, os.Args[len(os.Args)-1]) {
			filtered = append(filtered, r)
		}
	}

	vals := make([]string, len(filtered))
	for index, val := range filtered {
		if len(filtered) == 1 {
			formattedVal := sanitizer.Replace(val.Value)
			vals[index] = formattedVal
		} else {
			if val.Description != "" {
				vals[index] = fmt.Sprintf("%v (%v)", val.Value, sanitizer.Replace(val.Description))
			} else {
				vals[index] = val.Value
			}
		}
	}
	return strings.Join(vals, "\n")
}
