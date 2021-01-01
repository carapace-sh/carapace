package powershell

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
	`’`, ``,
	`,`, "`,",
)

func Sanitize(values ...string) []string {
	sanitized := make([]string, len(values))
	for index, value := range values {
		sanitized[index] = sanitizer.Replace(value)
	}
	return sanitized
}

func EscapeSpace(value string) string {
	return strings.Replace(value, " ", "` ", -1)
}

func ActionRawValues(values ...common.RawValue) string {
	vals := make([]string, len(values))
	for index, val := range values {
		if val.Value != "" { // must not be empty - any empty `''` parameter in CompletionResult causes an error
			vals[index] = fmt.Sprintf(`[CompletionResult]::new('%v', '%v ', [CompletionResultType]::ParameterValue, '%v ')`, EscapeSpace(sanitizer.Replace(val.Value)), sanitizer.Replace(val.Display), sanitizer.Replace(val.Description))
		}
	}
	return strings.Join(vals, "\n")
}
