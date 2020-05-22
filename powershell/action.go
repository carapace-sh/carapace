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

func EscapeSpace(value string) string {
	return strings.Replace(value, " ", "` ", -1)
}

func Callback(prefix string, uid string) string {
	return fmt.Sprintf("_%v_callback '%v'", prefix, uid)
}

func ActionExecute(command string) string {
	return fmt.Sprintf(`"%v" | Out-String | InvokeExpression`, strings.Replace(command, "\n", "`n", -1))
}

func ActionDirectories() string {
	return `[CompletionResult]::new('', '', [CompletionResultType]::ParameterValue, '')`
}

func ActionFiles(suffix string) string {
	return `[CompletionResult]::new('', '', [CompletionResultType]::ParameterValue, '')`
}

func ActionNetInterfaces() string {
	return `$(Get-NetAdapter).Name` // TODO test this
}

func ActionUsers() string {
	return `$(Get-LocalUser).Name` // TODO test this
}

func ActionGroups() string {
	return `$(Get-Localgroup).Name` // TODO test this
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
		vals[index] = fmt.Sprintf(`[CompletionResult]::new('%v ', '%v', [CompletionResultType]::ParameterValue, ' ')`, EscapeSpace(val), val)
	}
	return strings.Join(vals, "\n")
}

func ActionValuesDescribed(values ...string) string {
	sanitized := Sanitize(values...)
	// TODO verify length (description always exists)
	vals := make([]string, len(sanitized))
	for index, val := range sanitized {
		if index%2 == 0 {
			vals[index/2] = fmt.Sprintf(`[CompletionResult]::new('%v ', '%v', [CompletionResultType]::ParameterValue, '%v')`, EscapeSpace(val), val, values[index+1])
		}
	}
	return strings.Join(vals, "\n")
}

func ActionMessage(msg string) string {
	return ActionValuesDescribed("_", msg, "ERR", msg)
}

func ActionPrefixValues(prefix string, values ...string) string {
	sanitized := Sanitize(values...)
	if len(strings.TrimSpace(strings.Join(sanitized, ""))) == 0 {
		return ActionMessage("no values to complete")
	}

	vals := make([]string, len(sanitized))
	for index, val := range sanitized {
		vals[index] = fmt.Sprintf(`[CompletionResult]::new('%v', '%v', [CompletionResultType]::ParameterValue, ' ')`, EscapeSpace(prefix+val), val)
	}
	return strings.Join(vals, "\n")
}
