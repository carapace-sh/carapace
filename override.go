package carapace

import (
	"os"
	"strings"
)

var opts Opts

// Opts contains overrides for completion behaviour.
type Opts struct {
	// OptArgDelimiter changes the delimiter for optional flag arguments
	//   "=" // tail --verbose=descriptor (default)
	//   ":" // java -verbose:class
	OptArgDelimiter string
	// BridgeCompletion registers carapace completions to cobra's default completion
	BridgeCompletion bool
}

func init() {
	opts.OptArgDelimiter = "="
}

// Override changes completion behaviour for non-posix style flags in standalone mode.
// Mostly done by patching `os.Args` before command execution and thus must be called before it.
// These are considered hacks and might undergo changes in future (or replaced by a custom pflag fork).
//
//	func Execute() error {
//	    carapace.Override(carapace.Opts{
//	        LongShorthand: true,
//	        OptArgDelimiter: ":",
//	    })
//		return rootCmd.Execute()
//	}
func Override(o Opts) {
	if o.OptArgDelimiter != "=" &&
		o.OptArgDelimiter != "" {
		opts.OptArgDelimiter = o.OptArgDelimiter
		for index, arg := range os.Args {
			if strings.HasPrefix(arg, "-") {
				// TODO needs solution compatible with ActionInvoke (not changing os.Args)
				os.Args[index] = strings.Replace(arg, o.OptArgDelimiter, `=`, 1)
			}
		}
	}

	opts.BridgeCompletion = o.BridgeCompletion
}
