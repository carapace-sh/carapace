package condition

import (
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/util"
)

// Arch returns true if any of the given strings matches `runtime.GOARCH`.
func Arch(s ...string) func(c carapace.Context) bool {
	return func(c carapace.Context) bool {
		for _, arch := range s {
			if arch == runtime.GOARCH {
				return true
			}
		}
		return false
	}
}

// Arch returns true if any of the given strings matches `runtime.GOOS`.
func Os(s ...string) func(c carapace.Context) bool {
	return func(c carapace.Context) bool {
		for _, os := range s {
			if os == runtime.GOOS {
				return true
			}
		}
		return false
	}
}

// Excutable returns true if any of the given strings matches an executable in PATH.
func Executable(s ...string) func(c carapace.Context) bool {
	return func(c carapace.Context) bool {
		for _, executable := range s {
			// TODO needs to use Context.Env
			if _, err := exec.LookPath(executable); err == nil {
				return true
			}
		}
		return false
	}
}

// Retuns true if given string is a valid file or directory.
func File(s string) func(c carapace.Context) bool {
	return func(c carapace.Context) bool {
		if _, err := os.Stat(s); err == nil {
			return true
		}
		return false
	}
}

// CompletingPath returns true when `Context.Value` has a path prefix.
func CompletingPath(c carapace.Context) bool {
	return util.HasPathPrefix(c.Value)
}

// CompletingPathS is like CompletingPathS but also checks for path separator `/`
func CompletingPathS(c carapace.Context) bool {
	return CompletingPath(c) || strings.Contains(c.Value, "/") // TODO support windows backslash at some point?
}
