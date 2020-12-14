package carapace

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/rsteube/carapace/internal/bash"
	"github.com/rsteube/carapace/internal/common"
	"github.com/rsteube/carapace/internal/elvish"
	"github.com/rsteube/carapace/internal/fish"
	"github.com/rsteube/carapace/internal/powershell"
	"github.com/rsteube/carapace/internal/xonsh"
	"github.com/rsteube/carapace/internal/zsh"
	"github.com/spf13/cobra"
)

// Action indicates how to complete a flag or positional argument
type Action struct {
	rawValues  []common.Candidate
	bash       func() string
	elvish     func() string
	fish       func() string
	powershell func() string
	xonsh      func() string
	zsh        func() string
	callback   CompletionCallback
}

type ActionMap map[string]Action

func (a ActionMap) invokeCallback(uid string, args []string) Action {
	if action, ok := a[uid]; ok {
		if action.callback != nil {
			return action.callback(args)
		}
	}
	return ActionMessage(fmt.Sprintf("callback %v unknown", uid))
}

func (a ActionMap) shell(shell string) map[string]string {
	actions := make(map[string]string, len(a))
	for key, value := range map[string]Action(a) {
		actions[key] = value.Value(shell)
	}
	return actions
}

type CompletionCallback func(args []string) Action

// finalize replaces value if a callback function is set
func (a Action) finalize(cmd *cobra.Command, uid string) Action {
	if a.callback != nil {
		if a.bash == nil {
			a.bash = func() string { return bash.Callback(cmd.Root().Name(), uid) }
		}
		if a.elvish == nil {
			a.elvish = func() string { return elvish.Callback(cmd.Root().Name(), uid) }
		}
		if a.fish == nil {
			a.fish = func() string { return fish.Callback(cmd.Root().Name(), uid) }
		}
		if a.powershell == nil {
			a.powershell = func() string { return powershell.Callback(cmd.Root().Name(), uid) }
		}
		if a.xonsh == nil {
			a.xonsh = func() string { return xonsh.Callback(cmd.Root().Name(), uid) }
		}
		if a.zsh == nil {
			a.zsh = func() string { return zsh.Callback(cmd.Root().Name(), uid) }
		}
	}
	return a
}

type InvokedAction Action

// Invoke executes the callback of an action if it exists (supports nesting)
func (a Action) Invoke(args []string) InvokedAction {
	return InvokedAction(a.nestedAction(args, 5))
}

func (a InvokedAction) Merge(others ...InvokedAction) InvokedAction {
	uniqueCandidates := make(map[string]common.Candidate)
	for _, other := range append([]InvokedAction{a}, others...) {
		for _, c := range other.rawValues {
			uniqueCandidates[c.Value] = c
		}
	}

	candidates := make([]common.Candidate, 0, len(uniqueCandidates))
	for _, c := range uniqueCandidates {
		candidates = append(candidates, c)
	}
	return InvokedAction(actionCandidates(candidates...))
}

func (a InvokedAction) Filter(values []string) InvokedAction {
	toremove := make(map[string]bool)
	for _, v := range values {
		toremove[v] = true
	}
	filtered := make(common.Candidates, 0)
	for _, candidate := range a.rawValues {
		if _, ok := toremove[candidate.Value]; !ok {
			filtered = append(filtered, candidate)
		}
	}
	return InvokedAction(actionCandidates(filtered...))
}

func (a InvokedAction) Prefix(prefix string) InvokedAction {
	for index, val := range a.rawValues {
		a.rawValues[index].Value = prefix + val.Value
	}
	return a
}

func (a InvokedAction) Suffix(suffix string) InvokedAction {
	for index, val := range a.rawValues {
		a.rawValues[index].Value = val.Value + suffix
	}
	return a
}

func (a InvokedAction) ToA() Action {
	return Action(a)
}

func (a InvokedAction) ToMultipartsA(divider string) Action {
	return ActionMultiParts(divider, func(args, parts []string) Action {
		vals := make([]string, 0)
		for _, val := range a.rawValues {
			if strings.HasPrefix(val.Value, CallbackValue) {
				if splitted := strings.Split(val.Value, divider); len(splitted) > len(parts) {
					if len(splitted) == len(parts)+1 {
						vals = append(vals, splitted[len(parts)], val.Description)
					} else {
						vals = append(vals, splitted[len(parts)]+divider, val.Description)
					}
				}
			}
		}
		return ActionValuesDescribed(vals...)
	})
}

func (a Action) nestedAction(args []string, maxDepth int) Action {
	if a.rawValues == nil && a.callback != nil && maxDepth > 0 {
		return a.callback(args).nestedAction(args, maxDepth-1)
	} else {
		return a
	}
}

func (a Action) Value(shell string) string {
	var f func() string
	switch shell {
	case "bash":
		f = a.bash
	case "fish":
		f = a.fish
	case "elvish":
		f = a.elvish
	case "powershell":
		f = a.powershell
	case "xonsh":
		f = a.xonsh
	case "zsh":
		f = a.zsh
	}

	if f == nil {
		// TODO "{}" for xonsh?
		return ""
	} else {
		return f()
	}
}

// ActionCallback invokes a go function during completion
func ActionCallback(callback CompletionCallback) Action {
	return Action{callback: callback}
}

// ActionBool completes true/false
func ActionBool() Action {
	return ActionValues("true", "false")
}

// ActionDirectories completes directories
func ActionDirectories() Action {
	return Action{
		bash:       func() string { return bash.ActionDirectories() },
		elvish:     func() string { return elvish.ActionDirectories() },
		fish:       func() string { return fish.ActionDirectories() },
		powershell: func() string { return powershell.ActionDirectories() },
		xonsh:      func() string { return xonsh.ActionDirectories() },
		zsh:        func() string { return zsh.ActionDirectories() },
		callback:   func(args []string) Action { return actionPath("", true) },
	}
}

// ActionFiles completes files with optional suffix filtering
func ActionFiles(suffix string) Action {
	return Action{
		bash:       func() string { return bash.ActionFiles(suffix) },
		elvish:     func() string { return elvish.ActionFiles(suffix) },
		fish:       func() string { return fish.ActionFiles(suffix) },
		powershell: func() string { return powershell.ActionFiles(suffix) },
		xonsh:      func() string { return xonsh.ActionFiles(suffix) },
		zsh:        func() string { return zsh.ActionFiles("*" + suffix) },
		callback:   func(args []string) Action { return actionPath(suffix, false) },
	}
}

func actionPath(fileSuffix string, dirOnly bool) Action {
	folder := filepath.Dir(CallbackValue)
	if files, err := ioutil.ReadDir(folder); err != nil {
		return ActionMessage(err.Error())
	} else {
		if folder == "." {
			folder = ""
		} else if !strings.HasSuffix(folder, "/") {
			folder = folder + "/"
		}

		vals := make([]string, len(files))
		for index, file := range files {
			if file.IsDir() {
				vals[index] = folder + file.Name() + "/"
			} else if !dirOnly && strings.HasSuffix(file.Name(), fileSuffix) {
				vals[index] = folder + file.Name()
			}
		}
		return ActionValues(vals...)
	}
}

// ActionValues completes arbitrary keywords (values)
func ActionValues(values ...string) Action {
	vals := make([]string, len(values)*2)
	for index, val := range values {
		vals[index*2] = val
		vals[(index*2)+1] = ""
	}
	return ActionValuesDescribed(vals...)
}

// ActionValuesDescribed completes arbitrary key (values) with an additional description (value, description pairs)
func ActionValuesDescribed(values ...string) Action {
	vals := make([]common.Candidate, len(values)/2)
	for index, val := range values {
		if index%2 == 0 {
			vals[index/2] = common.Candidate{Value: val, Display: val, Description: values[index+1]}
		}
	}
	return actionCandidates(vals...)
}

func actionCandidates(candidates ...common.Candidate) Action {
	return Action{
		rawValues:  candidates,
		bash:       func() string { return bash.ActionCandidates(candidates...) },
		elvish:     func() string { return elvish.ActionCandidates(candidates...) },
		fish:       func() string { return fish.ActionCandidates(candidates...) },
		powershell: func() string { return powershell.ActionCandidates(candidates...) },
		xonsh:      func() string { return xonsh.ActionCandidates(candidates...) },
		zsh:        func() string { return zsh.ActionCandidates(candidates...) },
	}
}

// ActionMessage displays a help messages in places where no completions can be generated
func ActionMessage(msg string) Action {
	return ActionValuesDescribed("_", "", "ERR", msg)
}

// CallbackValue is set to the currently completed flag/positional value during callback (note that this is updated during ActionMultiParts)
var CallbackValue string

// ActionMultiParts completes multiple parts of words separately where each part is separated by some char (CallbackValue is set to the currently completed part during invocation)
func ActionMultiParts(divider string, callback func(args []string, parts []string) Action) Action {
	return ActionCallback(func(args []string) Action {
		index := strings.LastIndex(CallbackValue, string(divider))
		prefix := ""
		if len(divider) == 0 {
			prefix = CallbackValue
		} else if index != -1 {
			prefix = CallbackValue[0 : index+len(divider)]
			CallbackValue = CallbackValue[index+len(divider):] // update CallbackValue to only contain the currently completed part
		}
		parts := strings.Split(prefix, string(divider))
		if len(parts) > 0 {
			parts = parts[0 : len(parts)-1]
		}

		return callback(args, parts).Invoke(args).Prefix(prefix).ToA()
	})
}
