package carapace

import (
	"strings"

	"github.com/rsteube/carapace/bash"
	"github.com/rsteube/carapace/common"
	"github.com/rsteube/carapace/elvish"
	"github.com/rsteube/carapace/fish"
	"github.com/rsteube/carapace/powershell"
	"github.com/rsteube/carapace/xonsh"
	"github.com/rsteube/carapace/zsh"
	"github.com/spf13/cobra"
)

type Action struct {
	rawValues  []common.Candidate
	bash       func() string
	elvish     func() string
	fish       func() string
	powershell func() string
	xonsh      func() string
	zsh        func() string
	Callback   CompletionCallback
}
type ActionMap map[string]Action
type CompletionCallback func(args []string) Action

// finalize replaces value if a callback function is set
func (a Action) finalize(cmd *cobra.Command, uid string) Action {
	if a.Callback != nil {
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

// TODO maybe use Invoke(args) and work on []Candidate
func (a Action) Prefix(prefix string, args []string) Action {
	if nestedAction := a.NestedAction(args, 5); nestedAction.rawValues != nil {
		for index, val := range nestedAction.rawValues {
			nestedAction.rawValues[index].Value = prefix + val.Value // TODO check if val.Value can be assigned directly
		}
		return nestedAction
	} else {
		return ActionMessage("TODO Prefix(str) failed")
	}
}

// TODO maybe use Invoke(args) and work on []Candidate
func (a Action) Suffix(suffix string, args []string) Action {
	if nestedAction := a.NestedAction(args, 5); nestedAction.rawValues != nil {
		for index, val := range nestedAction.rawValues {
			nestedAction.rawValues[index].Value = val.Value + suffix // TODO check if val.Value can be assigned directly
		}
		return nestedAction
	} else {
		return ActionMessage("TODO Prefix(str) failed")
	}
}

// TODO maybe rename to Invoke(args)
func (a Action) NestedAction(args []string, maxDepth int) Action {
	if a.rawValues == nil && a.Callback != nil && maxDepth > 0 {
		return a.Callback(args).NestedAction(args, maxDepth-1)
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

func (m *ActionMap) Shell(shell string) map[string]string {
	actions := make(map[string]string, len(completions.actions))
	for key, value := range completions.actions {
		actions[key] = value.Value(shell)
	}
	return actions
}

// ActionCallback invokes a go function during completion
func ActionCallback(callback CompletionCallback) Action {
	return Action{Callback: callback}
}

// ActionBool completes true/false
func ActionBool() Action {
	return ActionValues("true", "false")
}

func ActionDirectories() Action {
	return Action{
		bash:       func() string { return bash.ActionDirectories() },
		elvish:     func() string { return elvish.ActionDirectories() },
		fish:       func() string { return fish.ActionDirectories() },
		powershell: func() string { return powershell.ActionDirectories() },
		xonsh:      func() string { return xonsh.ActionDirectories() },
		zsh:        func() string { return zsh.ActionDirectories() },
		// TODO add Callback so that the action can be used in ActionMultiParts as well
	}
}

func ActionFiles(suffix string) Action {
	return Action{
		bash:       func() string { return bash.ActionFiles(suffix) },
		elvish:     func() string { return elvish.ActionFiles(suffix) },
		fish:       func() string { return fish.ActionFiles(suffix) },
		powershell: func() string { return powershell.ActionFiles(suffix) },
		xonsh:      func() string { return xonsh.ActionFiles(suffix) },
		zsh:        func() string { return zsh.ActionFiles("*" + suffix) },
		// TODO add Callback so that the action can be used in ActionMultiParts as well
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
	return Action{
		rawValues:  vals,
		bash:       func() string { return bash.ActionCandidates(vals...) },
		elvish:     func() string { return elvish.ActionCandidates(vals...) },
		fish:       func() string { return fish.ActionCandidates(vals...) },
		powershell: func() string { return powershell.ActionCandidates(vals...) },
		xonsh:      func() string { return xonsh.ActionCandidates(vals...) },
		zsh:        func() string { return zsh.ActionCandidates(vals...) },
	}
}

// ActionMessage displays a help messages in places where no completions can be generated
func ActionMessage(msg string) Action { // TODO somehow handle this differently for Prefix/Suffix
	return ActionValuesDescribed("_", "", "ERR", msg)
	// TODO muss not be filtered if value already contains a submatch (so that it is always shown)
	// TODO zsh is the only one with actual message function		zsh:        func() string { return zsh.ActionMessage(msg) },
}

// TODO find a better solution for this
var CallbackValue string

// ActionMultiParts completes multiple parts of words separately where each part is separated by some char
func ActionMultiParts(divider string, callback func(args []string, parts []string) Action) Action {
	return ActionCallback(func(args []string) Action {
		// TODO multiple dividers by splitting on each char
		index := strings.LastIndex(CallbackValue, string(divider))
		prefix := ""
		if len(divider) == 0 {
			prefix = CallbackValue
		} else if index != -1 {
			prefix = CallbackValue[0 : index+1]
		}
		parts := strings.Split(prefix, string(divider))
		if len(parts) > 0 {
			parts = parts[0 : len(parts)-1]
		}

		return callback(args, parts).Prefix(prefix, args)
	})
}
