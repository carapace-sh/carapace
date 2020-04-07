package carapace

import (
	"github.com/rsteube/carapace/bash"
	"github.com/rsteube/carapace/fish"
	"github.com/rsteube/carapace/powershell"
	"github.com/rsteube/carapace/zsh"
	"github.com/spf13/cobra"
)

type Action struct {
	Bash       string
	Fish       string
	Zsh        string
	Powershell string
	Callback   CompletionCallback
}
type ActionMap map[string]Action
type CompletionCallback func(args []string) Action

// finalize replaces value if a callback function is set
func (a Action) finalize(cmd *cobra.Command, uid string) Action {
	if a.Callback != nil {
		if a.Bash == "" {
			a.Bash = bash.Callback(cmd.Root().Name(), uid)
		}
		if a.Fish == "" {
			a.Fish = fish.Callback(cmd.Root().Name(), uid)
		}
		if a.Powershell == "" {
			a.Powershell = powershell.Callback(cmd.Root().Name(), uid)
		}
		if a.Zsh == "" {
			a.Zsh = zsh.Callback(uid)
		}
	}
	return a
}

func (a Action) Value(shell string) string {
	switch shell {
	case "bash":
		return a.Bash
	case "fish":
		return a.Fish
	case "powershell":
		return a.Powershell
	case "zsh":
		return a.Zsh
	default:
		return ""
	}
}

// ActionCallback invokes a go function during completion
func ActionCallback(callback CompletionCallback) Action {
	return Action{Callback: callback}
}

// ActionExecute uses command substitution to invoke a command and evalues it's result as Action
func ActionExecute(command string) Action {
	return Action{
		Bash:       bash.ActionExecute(command),
		Fish:       fish.ActionExecute(command),
		Powershell: powershell.ActionExecute(command),
		Zsh:        zsh.ActionExecute(command),
	}
}

// ActionBool completes true/false
func ActionBool() Action {
	return ActionValues("true", "false")
}

func ActionDirectories() Action {
	return Action{
		Bash:       bash.ActionDirectories(),
		Fish:       fish.ActionDirectories(),
		Powershell: powershell.ActionDirectories(),
		Zsh:        zsh.ActionDirectories(),
	}
}

func ActionFiles(suffix string) Action {
	return Action{
		Bash:       bash.ActionFiles(suffix),
		Fish:       fish.ActionFiles(suffix),
		Powershell: powershell.ActionFiles(suffix),
		Zsh:        zsh.ActionFiles("*" + suffix),
	}
}

// ActionNetInterfaces completes network interface names
func ActionNetInterfaces() Action {
	return Action{
		Bash:       bash.ActionNetInterfaces(),
		Fish:       fish.ActionNetInterfaces(),
		Powershell: powershell.ActionNetInterfaces(),
		Zsh:        zsh.ActionNetInterfaces(),
	}
}

// ActionUsers completes user names
func ActionUsers() Action {
	return Action{
		Bash:       bash.ActionUsers(),
		Fish:       fish.ActionUsers(),
		Powershell: powershell.ActionUsers(),
		Zsh:        zsh.ActionUsers(),
	}
}

// ActionGroups completes group names
func ActionGroups() Action {
	return Action{
		Bash:       bash.ActionGroups(),
		Fish:       fish.ActionGroups(),
		Powershell: powershell.ActionGroups(),
		Zsh:        zsh.ActionGroups(),
	}
}

// ActionHosts completes host names
func ActionHosts() Action {
	return Action{
		Bash:       bash.ActionHosts(),
		Fish:       fish.ActionHosts(),
		Powershell: powershell.ActionHosts(),
		Zsh:        zsh.ActionHosts(),
	}
}

// ActionValues completes arbitrary keywords (values)
func ActionValues(values ...string) Action {
	return Action{
		Bash:       bash.ActionValues(values...),
		Fish:       fish.ActionValues(values...),
		Powershell: powershell.ActionValues(values...),
		Zsh:        zsh.ActionValues(values...),
	}
}

// ActionValuesDescribed completes arbitrary key (values) with an additional description (value, description pairs)
func ActionValuesDescribed(values ...string) Action {
	return Action{
		Bash:       bash.ActionValuesDescribed(values...),
		Fish:       fish.ActionValuesDescribed(values...),
		Powershell: powershell.ActionValuesDescribed(values...),
		Zsh:        zsh.ActionValuesDescribed(values...),
	}
}

// ActionMessage displays a help messages in places where no completions can be generated
func ActionMessage(msg string) Action {
	return Action{
		Bash:       bash.ActionMessage(msg),
		Fish:       fish.ActionMessage(msg),
		Powershell: powershell.ActionMessage(msg),
		Zsh:        zsh.ActionMessage(msg),
	}
}

// ActionMultiParts completes multiple parts of words separately where each part is separated by some char
func ActionMultiParts(separator rune, values ...string) Action {
	return Action{
		Bash:       bash.ActionMultiParts(separator, values...),
		Fish:       fish.ActionMultiParts(separator, values...),
		Powershell: powershell.ActionMultiParts(separator, values...),
		Zsh:        zsh.ActionMultiParts(separator, values...),
	}
}
