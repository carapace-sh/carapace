package carapace

import (
	"github.com/rsteube/carapace/internal/common"
	"github.com/rsteube/carapace/internal/config"
	"github.com/rsteube/carapace/internal/shell/bash"
	"github.com/rsteube/carapace/internal/shell/library"
	"github.com/rsteube/carapace/internal/shell/nushell"
	"github.com/rsteube/carapace/pkg/ps"
	"github.com/spf13/cobra"
)

// Complete can be used by Go programs wishing to produce completions for
// themselves, without passing through shell snippets/output or export formats.
//
// The `onFinalize` function parameter, if non nil, will be called after having
// generated the completions from the given command/tree. This function is generally
// used to reset the command tree, which is needed when the Go program is a shell itself.
// Also, and before calling `onFinalize` if not nil, the completion storage is cleared.
func Complete(cmd *cobra.Command, args []string, onFinalize func()) (common.RawValues, common.Meta) {
	// Generate the completion as normally done for an external system shell
	initHelpCompletion(cmd)
	action, context := traverse(cmd, args[2:])

	if err := config.Load(); err != nil {
		action = ActionMessage("failed to load config: " + err.Error())
	}

	if onFinalize != nil {
		storage = make(_storage)

		onFinalize()
	}

	invoked := action.Invoke(context)

	return library.ActionRawValues(context.Value, invoked.meta, invoked.rawValues)
}

func complete(cmd *cobra.Command, args []string) (string, error) {
	switch len(args) {
	case 0:
		return Gen(cmd).Snippet(ps.DetermineShell())
	case 1:
		return Gen(cmd).Snippet(args[0])
	default:
		initHelpCompletion(cmd)

		switch ps.DetermineShell() {
		case "nushell":
			args = nushell.Patch(args) // handle open quotes
			LOG.Printf("patching args to %#v", args)
		case "bash": // TODO what about oil and such?
			var err error
			args, err = bash.Patch(args) // handle redirects
			LOG.Printf("patching args to %#v", args)
			if err != nil {
				context := NewContext(args...)
				if _, ok := err.(bash.RedirectError); ok {
					LOG.Printf("completing redirect target for %#v", args)
					return ActionFiles().Invoke(context).value(args[0], args[len(args)-1]), nil
				}
				return ActionMessage(err.Error()).Invoke(context).value(args[0], args[len(args)-1]), nil
			}
		}

		action, context := traverse(cmd, args[2:])
		if err := config.Load(); err != nil {
			action = ActionMessage("failed to load config: " + err.Error())
		}
		return action.Invoke(context).value(args[0], args[len(args)-1]), nil
	}
}
