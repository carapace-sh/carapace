package elvish

import (
	"fmt"
	"strings"
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
	return fmt.Sprintf(`_%v_callback '%v'`, prefix, uid)
}

func ActionExecute(command string) string {
	return `` // TODO
}

func ActionDirectories() string {
	return `edit:complete-filename $arg[-1]` // TODO
}

func ActionFiles(suffix string) string {
	return `edit:complete-filename $arg[-1]` // TODO
}

func ActionNetInterfaces() string {
	return `` // TODO
}

func ActionUsers() string {
	return `` // TODO
}

func ActionHosts() string {
	return `` // TODO
}

func ActionValues(values ...string) string {
	sanitized := Sanitize(values...)
	if len(strings.TrimSpace(strings.Join(sanitized, ""))) == 0 {
		return ActionMessage("no values to complete")
	}

	vals := make([]string, len(sanitized))
	for index, val := range sanitized {
		// TODO escape special characters
		//vals[index] = fmt.Sprintf(`edit:complex-candidate %v`, val)
		vals[index] = fmt.Sprintf(`%v`, val)
	}
	return fmt.Sprintf(`put %v`, strings.Join(vals, " "))
}

func ActionValuesDescribed(values ...string) string {
	// TODO verify length (description always exists)
	sanitized := Sanitize(values...)
	vals := make([]string, len(values)/2)
	for index, val := range sanitized {
		if index%2 == 0 {
			vals[index/2] = fmt.Sprintf(`edit:complex-candidate '%v' &display='%v (%v)'`, val, val, sanitized[index+1])
		}
	}
	return strings.Join(vals, "\n")
}

func ActionMessage(msg string) string {
	return ActionValuesDescribed("ERR", Sanitize(msg)[0], "_", "")
}

func ActionPrefixValues(prefix string, values ...string) string {
	sanitized := Sanitize(values...)
	if len(strings.TrimSpace(strings.Join(sanitized, ""))) == 0 {
		return ActionMessage("no values to complete")
	}

	vals := make([]string, len(sanitized))
	for index, val := range sanitized {
		// TODO escape special characters
		vals[index] = fmt.Sprintf(`edit:complex-candidate '%v' &display='%v'`, prefix+val, val)
	}
	return strings.Join(vals, "\n")
}
