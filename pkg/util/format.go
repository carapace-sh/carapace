package util

import (
	"fmt"
	"strings"
)

// FormatCmd joins given args to a formatted command.
// TODO experimental
func FormatCmd(args ...string) string {
	replacer := strings.NewReplacer(
		"$", "\\$",
		"`", "\\`",
	)

	formatted := make([]string, 0, len(args))
	for _, arg := range args {
		switch {
		case arg == "",
			strings.ContainsAny(arg, `"' `+"\n\r\t"):
			formatted = append(formatted, replacer.Replace(fmt.Sprintf("%#v", arg)))
		default:
			formatted = append(formatted, arg)
		}
	}
	return strings.Join(formatted, " ")
}
