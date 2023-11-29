// Package assert provides test helpers
package assert

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/rsteube/carapace/third_party/github.com/hexops/gotextdiff"
	"github.com/rsteube/carapace/third_party/github.com/hexops/gotextdiff/myers"
	"github.com/rsteube/carapace/third_party/github.com/hexops/gotextdiff/span"
)

type T interface {
	Cleanup(func())
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Helper()
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Name() string
	Setenv(key, value string)
	Skip(args ...interface{})
	SkipNow()
	Skipf(format string, args ...interface{})
	Skipped() bool
	TempDir() string
}

// Equal calls t.Error if given strings are not equal.
func Equal(t T, expected string, actual string) {
	if expected != actual {
		_, file, line, _ := runtime.Caller(2)
		t.Errorf("%v:%v:\n%v", filepath.Base(file), line, Diff(expected, actual))
	}
}

func Diff(expected, actual string) string {
	edits := myers.ComputeEdits(span.URIFromPath(""), expected, actual)
	diff := fmt.Sprint(gotextdiff.ToUnified("expected", "actual", expected, edits))

	highlighted := make([]string, 0)
	for _, line := range strings.Split(diff, "\n") {
		switch {
		case strings.HasPrefix(line, "-"):
			highlighted = append(highlighted, fmt.Sprintf("\033[0;31m%v\033[0m", line))
		case strings.HasPrefix(line, "+"):
			highlighted = append(highlighted, fmt.Sprintf("\033[0;32m%v\033[0m", line))
		default:
			highlighted = append(highlighted, line)
		}
	}
	return strings.Join(highlighted, "\n")
}
