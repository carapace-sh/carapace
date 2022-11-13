package common

import (
	"github.com/rsteube/carapace/pkg/style"
)

// CompletionMessage is a message that is passed by ActionMessage calls.
// This message is either wrapped into an Action (as a raw completion value into it),
// or used by per-shell completion generation code (like ZSH with _message calls).
var CompletionMessage string

// CompletionHint TODO: Should be renamed and refactored, because its fundamentally a completion message.
var CompletionHint string

// AddMessageToValues checks if we should wrap an existing message into a completion action
// (for all shells) or let the per-shell completion code to use it as they see fit.
func AddMessageToValues(currentCallback string, values RawValues) RawValues {
	if CompletionMessage == "" {
		return values
	}

	// If all cases, empty this message, so as not induce
	// further completion calls in error. Do this only after
	// the shells Go code (eg. ZSH) have had to time to use it.
	defer func() { CompletionMessage = "" }()

	msgValue := RawValue{
		Value:       currentCallback + "ERR",
		Description: CompletionMessage,
		Style:       style.Carapace.Error,
	}

	return append([]RawValue{msgValue}, values...)
}
