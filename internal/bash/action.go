package bash

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
	return fmt.Sprintf(`eval $(_%v_callback '%v')`, prefix, uid)
}

func ActionDirectories() string {
	return `compgen -S / -d -- "$cur"`
}

func ActionFiles(suffix string) string {
	return fmt.Sprintf(`compgen -S / -d -- "$cur"; compgen -f -X '!*%v' -- "$cur"`, suffix)
}

func ActionCandidates(values ...common.Candidate) string {
	vals := make([]string, len(values))
	for index, val := range values {
		if val.Description == "" {
			vals[index] = sanitizer.Replace(val.Value)
		} else {
			vals[index] = fmt.Sprintf(`%v (%v)`, sanitizer.Replace(val.Value), sanitizer.Replace(val.Description))
		}
	}

	return fmt.Sprintf(`compgen -W $'%v' -- "${cur//\\ / }" | sed "s!^${curprefix//\\ / }!!"`, strings.Join(vals, `\n`))
}
