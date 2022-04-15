package carapace

import (
)

var opts Opts

// Opts contains overrides for completion behaviour
type Opts struct {
	// BridgeCompletion registers carapace completions to cobra's default completion
	BridgeCompletion bool
}

// Override changes completion behaviour for non-posix style flags in standalone mode.
// Mostly done by patching `os.Args` before command execution and thus must be called before it.
// These are considered hacks and might undergo changes in future (or replaced by a custom pflag fork).
//
//   func Execute() error {
//       carapace.Override(carapace.Opts{
//           LongShorthand: true,
//           OptArgDelimiter: ":",
//       })
//   	return rootCmd.Execute()
//   }
func Override(o Opts) {
	opts.BridgeCompletion = o.BridgeCompletion
}
