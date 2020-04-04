package powershell

import (
	"strings"
)

var sanitizer = strings.NewReplacer( // TODO
	`$`, ``,
	"`", ``,
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
	return `TODO`
}

func ActionExecute(command string) string {
	return `TODO`
}

func ActionDirectories() string {
	return `TODO`
}

func ActionFiles(suffix string) string {
	return `TODO`
}

func ActionNetInterfaces() string {
	return `TODO`
}

func ActionUsers() string {
	return `TODO`
}

func ActionGroups() string {
	return `TODO`
}

func ActionHosts() string {
	return `TODO`
}

func ActionValues(values ...string) string {
	return `TODO`
}

func ActionValuesDescribed(values ...string) string {
	return `TODO`
}

func ActionMessage(msg string) string {
	return `TODO`
}

func ActionMultiParts(separator rune, values ...string) string {
	return `TODO`
}
