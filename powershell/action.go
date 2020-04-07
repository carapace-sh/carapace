package powershell

import (
	"fmt"
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
	return fmt.Sprintf("_%v_callback '%v'", prefix, uid)
}

func ActionExecute(command string) string {
	return ActionMessage("TODO") // TODO
}

func ActionDirectories() string {
	return ActionValues("") // TODO
}

func ActionFiles(suffix string) string {
	return ActionValues("") // TODO
}

func ActionNetInterfaces() string {
	return `$(Get-NetAdapter).Name`
}

func ActionUsers() string {
	return `$(Get-LocalUser).Name` // TODO
}

func ActionGroups() string {
	return `$(Get-Localgroup).Name` // TODO
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
		vals[index] = fmt.Sprintf(`[CompletionResult]::new('%v ', '%v', [CompletionResultType]::ParameterValue, ' ')`, val, val)
	}
	return strings.Join(vals, "\n")
}

func ActionValuesDescribed(values ...string) string {
	sanitized := Sanitize(values...)
	// TODO verify length (description always exists)
	vals := make([]string, len(sanitized))
	for index, val := range sanitized {
		if index%2 == 0 {
			vals[index/2] = fmt.Sprintf(`[CompletionResult]::new('%v ', '%v', [CompletionResultType]::ParameterValue, '%v')`, val, val, values[index+1])
		}
	}
	return strings.Join(vals, "\n")
}

func ActionMessage(msg string) string {
	return ActionValues("ERR", msg)
}

func ActionMultiParts(separator rune, values ...string) string {
	return ActionValues(values...)
}
