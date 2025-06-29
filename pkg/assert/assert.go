package assert

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/carapace-sh/carapace/third_party/github.com/hexops/gotextdiff"
	"github.com/carapace-sh/carapace/third_party/github.com/hexops/gotextdiff/myers"
	"github.com/carapace-sh/carapace/third_party/github.com/hexops/gotextdiff/span"
)

func compare(t *testing.T, expected, actual interface{}, equal bool) {
	var sExpected, sActual string
	var ok bool
	if sExpected, ok = expected.(string); !ok {
		m, err := json.MarshalIndent(expected, "", " ")
		if err != nil {
			t.Error(err.Error())
		}
		sExpected = string(m)
	}
	if sActual, ok = actual.(string); !ok {
		m, err := json.MarshalIndent(actual, "", " ")
		if err != nil {
			t.Error(err.Error())
		}
		sActual = string(m)
	}

	if sExpected != sActual == equal {
		_, file, line, _ := runtime.Caller(2)
		t.Errorf("%v:%v:\n%v", filepath.Base(file), line, diff(sExpected, sActual))
	}
}

func Equal(t *testing.T, expected, actual interface{}) {
	compare(t, expected, actual, true)
}

func NotEqual(t *testing.T, expected, actual interface{}) {
	compare(t, expected, actual, false)
}

func diff(expected, actual string) string {
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
