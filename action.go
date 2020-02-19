package zsh

import (
	"fmt"
	"strings"
)

// https://github.com/zsh-users/zsh-completions/blob/master/zsh-completions-howto.org
// http://zsh.sourceforge.net/Doc/Release/Completion-System.html
type Action struct {
	Value    string
	Callback CompletionCallback
}
type ActionMap map[string]Action
type CompletionCallback func(args []string) Action

// replaces value if a callback function is set
func (a Action) finalize(uid string) Action {
	if a.Callback != nil {
		a.Value = ActionExecute(fmt.Sprintf("${os_args[1]} _zsh_completion '%v' ${os_args:1}", uid)).Value
	}
	return a
}

func ActionCallback(callback CompletionCallback) Action {
	return Action{Callback: callback}
}

// Wraps given command as command substitution and evaluates the output as zsh expression
func ActionExecute(command string) Action {
	return Action{Value: fmt.Sprintf(` eval \$(%v)`, command)} // {EVAL-STRING} action did not handle space escaping ('\ ') as expected (separate arguments), this one works
}

func ActionBool() Action {
	return ActionValues("true", "false")
}

// Used to complete filepaths.
// [http://zsh.sourceforge.net/Doc/Release/Completion-System.html#index-_005fpath_005ffiles]
func ActionPathFiles(pattern string) Action { // TODO support additional options
	if pattern == "" {
		return Action{Value: "_path_files"}
	} else {
		return Action{Value: fmt.Sprintf("_path_files -g '%v'", pattern)}
	}
}

// Calls _path_files with all options except -g and -/. These options depend on file-patterns style setting.
// [http://zsh.sourceforge.net/Doc/Release/Completion-System.html#index-_005ffiles]
func ActionFiles(pattern string) Action {
	if pattern == "" {
		return Action{Value: "_files"}
	} else {
		return Action{Value: fmt.Sprintf("_files -g '%v'", pattern)}
	}
}

// Used for completing network interface names
func ActionNetInterfaces() Action {
	return Action{Value: "_net_interfaces"}
}

// Used for completing user names
func ActionUsers() Action {
	return Action{Value: "_users"}
}

// Used for completing group names
func ActionGroups() Action {
	return Action{Value: "_groups"}
}

// Used for completing hosst names
func ActionHosts() Action {
	return Action{Value: "_hosts"}
}

// Used for completing the names of shell options.
func ActionOptions() Action {
	return Action{Value: "_options"}
}

// used to complete arbitrary keywords (values)
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

// Like ActionValues but with a list of value, description pairs
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

// Used for displaying help messages in places where no completions can be generated.
func ActionMessage(msg string) Action {
	return Action{Value: fmt.Sprintf(" _message -r '%v'", msg)} // space before _message is necessary
}
