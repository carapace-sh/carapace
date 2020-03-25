package zsh

import (
	"fmt"
	"strings"
)

func Callback(uid string) string {
	return ActionExecute(fmt.Sprintf(`${os_args[1]} _carapace zsh '%v' ${${os_args:1:gs/\"/\\\"}:gs/\'/\\\"}`, uid))
}

// ActionExecute uses command substitution to invoke a command and evalues it's result as Action
func ActionExecute(command string) string {
	return fmt.Sprintf(` eval \$(%v)`, command) // {EVAL-STRING} action did not handle space escaping ('\ ') as expected (separate arguments), this one works
}

// ActionPathFiles completes filepaths
// [http://zsh.sourceforge.net/Doc/Release/Completion-System.html#index-_005fpath_005ffiles]
func ActionPathFiles(pattern string) string { // TODO support additional options
	if pattern == "" {
		return "_path_files"
	} else {
		return fmt.Sprintf("_path_files -g '%v'", pattern)
	}
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

// ActionOptions completes the names of shell options
func ActionOptions() string {
	return "_options"
}

// ActionValues completes arbitrary keywords (values)
func ActionValues(values ...string) string {
	if len(strings.TrimSpace(strings.Join(values, ""))) == 0 {
		return ActionMessage("no values to complete")
	}

	vals := make([]string, len(values))
	for index, val := range values {
		// TODO escape special characters
		vals[index] = strings.Replace(val, " ", `\ `, -1)
	}
	return fmt.Sprintf(`_values '' %v`, strings.Join(vals, " "))
}

// ActionValuesDescribed completes arbitrary key (values) with an additional description (value, description pairs)
func ActionValuesDescribed(values ...string) string {
	// TODO verify length (description always exists)
	vals := make([]string, len(values))
	for index, val := range values {
		if index%2 == 0 {
			vals[index/2] = fmt.Sprintf("'%v[%v]'", val, values[index+1])
		}
	}
	return ActionValues(vals...)
}

// ActionMessage displays a help messages in places where no completions can be generated
func ActionMessage(msg string) string {
	return fmt.Sprintf(" _message -r '%v'", msg) // space before _message is necessary
}

// ActionMultiParts completes multiple parts of words separately where each part is separated by some char
func ActionMultiParts(separator rune, values ...string) string {
	return fmt.Sprintf("_multi_parts %v '(%v)'", string(separator), strings.Join(values, " "))
}
