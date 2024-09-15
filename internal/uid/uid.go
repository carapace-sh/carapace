// Package uid provides unique identifiers
package uid

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/carapace-sh/carapace/internal/pflagfork"
	"github.com/spf13/cobra"
)

// Command creates a uid for given command.
func Command(cmd *cobra.Command) *url.URL {
	path := []string{cmd.Name()}
	for parent := cmd.Parent(); parent != nil; parent = parent.Parent() {
		path = append(path, parent.Name())
	}
	reverse(path) // TODO slices.Reverse
	return &url.URL{
		Scheme: "cmd",
		Host:   path[0],
		Path:   strings.Join(path[1:], "/"),
	}
}

// reverse reverses the elements of the slice in place.
func reverse(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// Flag creates a uid for given flag.
func Flag(cmd *cobra.Command, flag *pflagfork.Flag) *url.URL {
	uid := Command(cmd)
	values := uid.Query()
	values.Set("flag", flag.Name)
	uid.RawQuery = values.Encode()
	return uid
}

// Executable returns the name of the executable.
func Executable() string {
	if executable, err := os.Executable(); err != nil {
		return "echo" // safe fallback that should never happen
	} else if filepath.Base(executable) == "cmd.test" {
		return "example" // for `go test -v ./...`
	} else {
		return filepath.Base(executable)
	}
}

// Map maps values to uids to simplify testing.
//
//	Map(
//	    "go.mod", "file://path/to/go.mod",
//	    "go.sum", "file://path/to/go.sum",
//	)
func Map(uids ...string) func(s string) (*url.URL, error) {
	return func(s string) (*url.URL, error) {
		for i := 0; i < len(uids); i += 2 {
			if uids[i] == s {
				return url.Parse(uids[i+1])
			}
		}
		return &url.URL{}, nil
	}
}
