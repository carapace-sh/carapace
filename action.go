package zsh

import (
	"fmt"
	"strings"
)

// Action indicates how to complete the corresponding argument
// https://github.com/zsh-users/zsh-completions/blob/master/zsh-completions-howto.org
// http://zsh.sourceforge.net/Doc/Release/Completion-System.html
type Action struct {
	Value    string
	Callback CompletionCallback
}
type ActionMap map[string]Action
type CompletionCallback func(args []string) Action

// finalize replaces value if a callback function is set
func (a Action) finalize(uid string) Action {
	if a.Callback != nil {
		a.Value = ActionExecute(fmt.Sprintf(`${os_args[1]} _zsh_completion '%v' ${${os_args:1:gs/\"/\\\"}:gs/\'/\\\"}`, uid)).Value
	}
	return a
}

// ActionCallback invokes a go function during completion
func ActionCallback(callback CompletionCallback) Action {
	return Action{Callback: callback}
}

// ActionExecute uses command substitution to invoke a command and evalues it's result as Action
func ActionExecute(command string) Action {
	return Action{Value: fmt.Sprintf(` eval \$(%v)`, command)} // {EVAL-STRING} action did not handle space escaping ('\ ') as expected (separate arguments), this one works
}

// ActionBool completes true/false
func ActionBool() Action {
	return ActionValues("true", "false")
}

// ActionPathFiles completes filepaths
// [http://zsh.sourceforge.net/Doc/Release/Completion-System.html#index-_005fpath_005ffiles]
func ActionPathFiles(pattern string) Action { // TODO support additional options
	if pattern == "" {
		return Action{Value: "_path_files"}
	} else {
		return Action{Value: fmt.Sprintf("_path_files -g '%v'", pattern)}
	}
}

// ActionFiles _path_files with all options except -g and -/. These options depend on file-patterns style setting. // TODO fix doc
// [http://zsh.sourceforge.net/Doc/Release/Completion-System.html#index-_005ffiles]
func ActionFiles(pattern string) Action {
	if pattern == "" {
		return Action{Value: "_files"}
	} else {
		return Action{Value: fmt.Sprintf("_files -g '%v'", pattern)}
	}
}

// ActionNetInterfaces completes network interface names
func ActionNetInterfaces() Action {
	return Action{Value: "_net_interfaces"}
}

// ActionUsers completes user names
func ActionUsers() Action {
	return Action{Value: "_users"}
}

// ActionGroups completes group names
func ActionGroups() Action {
	return Action{Value: "_groups"}
}

// ActionHosts completes host names
func ActionHosts() Action {
	return Action{Value: "_hosts"}
}

// ActionOptions completes the names of shell options
func ActionOptions() Action {
	return Action{Value: "_options"}
}

// ActionValues completes arbitrary keywords (values)
func ActionValues(values ...string) Action {
	if len(strings.TrimSpace(strings.Join(values, ""))) == 0 {
		return ActionMessage("no values to complete")
	}

	vals := make([]string, len(values))
	for index, val := range values {
		// TODO escape special characters
		vals[index] = strings.Replace(val, " ", `\ `, -1)
	}
	return Action{Value: fmt.Sprintf(`_values '' %v`, strings.Join(vals, " "))}
}

// ActionValuesDescribed completes arbitrary key (values) with an additional description (value, description pairs)
func ActionValuesDescribed(values ...string) Action {
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
func ActionMessage(msg string) Action {
	return Action{Value: fmt.Sprintf(" _message -r '%v'", msg)} // space before _message is necessary
}

// ActionMultiParts completes multiple parts of words separately where each part is separated by some char
func ActionMultiParts(separator rune, values ...string) Action {
	return Action{Value: fmt.Sprintf("_multi_parts %v '(%v)'", string(separator), strings.Join(values, " "))}
}
