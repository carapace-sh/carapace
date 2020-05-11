package zsh

import (
	"fmt"
	"github.com/rsteube/carapace/uid"
	"strings"
)

var sanitizer = strings.NewReplacer(
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
	`[`, `\[`,
	`]`, `\]`,
)

func Sanitize(values ...string) []string {
	sanitized := make([]string, len(values))
	for index, value := range values {
		sanitized[index] = sanitizer.Replace(value)
	}
	return sanitized
}

func Callback(cuid string) string {
	return ActionExecute(fmt.Sprintf(`%v _carapace zsh '%v' ${${os_args:1:gs/\"/\\\"}:gs/\'/\\\"}`, uid.Executable(), cuid))
}

// ActionExecute uses command substitution to invoke a command and evalues it's result as Action
func ActionExecute(command string) string {
	return fmt.Sprintf(` eval \$(%v)`, command) // {EVAL-STRING} action did not handle space escaping ('\ ') as expected (separate arguments), this one works
}

func ActionDirectories() string {
	return `_files -/`
}

// ActionFiles _path_files with all options except -g and -/. These options depend on file-patterns style setting. // TODO fix doc
// [http://zsh.sourceforge.net/Doc/Release/Completion-System.html#index-_005ffiles]
func ActionFiles(pattern string) string {
	if pattern == "" {
		return "_files"
	} else {
		return fmt.Sprintf("_files -g '%v'", pattern)
	}
}

// ActionNetInterfaces completes network interface names
func ActionNetInterfaces() string {
	return "_net_interfaces"
}

// ActionUsers completes user names
func ActionUsers() string {
	return "_users"
}

// ActionGroups completes group names
func ActionGroups() string {
	return "_groups"
}

// ActionHosts completes host names
func ActionHosts() string {
	return "_hosts"
}

// ActionValues completes arbitrary keywords (values)
func ActionValues(values ...string) string {
	sanitized := Sanitize(values...)
	if len(strings.TrimSpace(strings.Join(sanitized, ""))) == 0 {
		return ActionMessage("no values to complete")
	}

	vals := make([]string, len(sanitized))
	for index, val := range sanitized {
		// TODO escape special characters
		vals[index] = strings.Replace(val, " ", `\ `, -1)
	}
	return fmt.Sprintf(`_values '' %v`, strings.Join(vals, " "))
}

// ActionValuesDescribed completes arbitrary key (values) with an additional description (value, description pairs)
func ActionValuesDescribed(values ...string) string {
	sanitized := Sanitize(values...)
	if len(strings.TrimSpace(strings.Join(sanitized, ""))) == 0 {
		return ActionMessage("no values to complete")
	}

	// TODO verify length (description always exists)
	vals := make([]string, len(sanitized))
	for index, val := range sanitized {
		if index%2 == 0 {
			vals[index/2] = fmt.Sprintf("'%v[%v]'", strings.Replace(val, " ", `\ `, -1), strings.Replace(sanitized[index+1], " ", `\ `, -1))
		}
	}
	return fmt.Sprintf(`_values '' %v`, strings.Join(vals, " "))
}

// ActionMessage displays a help messages in places where no completions can be generated
func ActionMessage(msg string) string {
	return fmt.Sprintf(" _message -r '%v'", msg) // space before _message is necessary
}

// ActionMultiParts completes multiple parts of words separately where each part is separated by some char
func ActionMultiParts(separator rune, values ...string) string {
	return fmt.Sprintf("_multi_parts %v '(%v)'", string(separator), strings.Join(values, " "))
}
