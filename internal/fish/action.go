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

func Callback(prefix string, uid string) string {
	return fmt.Sprintf(`_%v_callback %v`, prefix, uid)
}

func ActionDirectories() string {
	return `__fish_complete_directories`
}

func ActionFiles(suffix string) string {
	return fmt.Sprintf(`__fish_complete_suffix "%v"`, suffix)
}

func ActionCandidates(values ...common.Candidate) string {
	vals := make([]string, len(values))
	for index, val := range values {
		// TODO sanitize
		//vals[index] = strings.Replace(val, " ", `\ `, -1)
		vals[index] = sanitizer.Replace(val.Value) + "\t" + sanitizer.Replace(val.Description)
	}
	return fmt.Sprintf(`echo -e "%v"`, strings.Join(vals, `\n`))
}
