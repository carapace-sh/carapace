package carapace

import (
	"github.com/rsteube/carapace/fish"
	"github.com/rsteube/carapace/zsh"
)

type Action struct {
	Fish     string
	Zsh      string
	Callback CompletionCallback
}
type ActionMap map[string]Action
type CompletionCallback func(args []string) Action

// finalize replaces value if a callback function is set
func (a Action) finalize(uid string) Action {
	if a.Callback != nil {
		// TODO only set to callback if no value is set (one shell might not need the callback)
		a.Fish = fish.Callback(uid)
		a.Zsh = zsh.Callback(uid)
	}
	return a
}

// ActionCallback invokes a go function during completion
func ActionCallback(callback CompletionCallback) Action {
	return Action{Callback: callback}
}

// ActionExecute uses command substitution to invoke a command and evalues it's result as Action
func ActionExecute(command string) Action {
	return Action{
		Fish: fish.ActionExecute(command),
		Zsh:  zsh.ActionExecute(command),
	}
}

// ActionBool completes true/false
func ActionBool() Action {
	return Action{
		Fish: fish.ActionBool(),
		Zsh:  zsh.ActionBool(),
	}
}

// ActionPathFiles completes filepaths
func ActionPathFiles(suffix string) Action {
	return Action{
		Fish: fish.ActionPathFiles(suffix),
		Zsh:  zsh.ActionPathFiles("*" + suffix),
	}
}

func ActionFiles(suffix string) Action {
	return Action{
		Fish: fish.ActionFiles(suffix),
		Zsh:  zsh.ActionFiles("*" + suffix),
	}
}

// ActionNetInterfaces completes network interface names
func ActionNetInterfaces() Action {
	return Action{
		Fish: fish.ActionNetInterfaces(),
		Zsh:  zsh.ActionNetInterfaces(),
	}
}

// ActionUsers completes user names
func ActionUsers() Action {
	return Action{
		Fish: fish.ActionUsers(),
		Zsh:  zsh.ActionUsers(),
	}
}

// ActionGroups completes group names
func ActionGroups() Action {
	return Action{
		Fish: fish.ActionGroups(),
		Zsh:  zsh.ActionGroups(),
	}
}

// ActionHosts completes host names
func ActionHosts() Action {
	return Action{
		Fish: fish.ActionHosts(),
		Zsh:  zsh.ActionHosts(),
	}
}

// ActionOptions completes the names of shell options
func ActionOptions() Action {
	return Action{
		Fish: fish.ActionOptions(),
		Zsh:  zsh.ActionOptions(),
	}
}

// ActionValues completes arbitrary keywords (values)
func ActionValues(values ...string) Action {
	return Action{
		Fish: fish.ActionValues(values...),
		Zsh:  zsh.ActionValues(values...),
	}
}

// ActionValuesDescribed completes arbitrary key (values) with an additional description (value, description pairs)
func ActionValuesDescribed(values ...string) Action {
	return Action{
		Fish: fish.ActionValuesDescribed(values...),
		Zsh:  zsh.ActionValuesDescribed(values...),
	}
}

// ActionMessage displays a help messages in places where no completions can be generated
func ActionMessage(msg string) Action {
	return Action{
		Fish: fish.ActionMessage(msg),
		Zsh:  zsh.ActionMessage(msg),
	}
}

// ActionMultiParts completes multiple parts of words separately where each part is separated by some char
func ActionMultiParts(separator rune, values ...string) Action {
	return Action{
		Fish: fish.ActionMultiParts(separator, values...),
		Zsh:  zsh.ActionMultiParts(separator, values...),
	}
}
