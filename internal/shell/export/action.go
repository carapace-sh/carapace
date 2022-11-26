package export

import (
	"encoding/json"
	"runtime/debug"
	"sort"

	"github.com/rsteube/carapace/internal/common"
)

type Export struct {
	Version   string
	Usage     string
	Nospace   common.SuffixMatcher
	RawValues common.RawValues
}

func ActionRawValues(currentWord string, usage string, nospace common.SuffixMatcher, values common.RawValues) string {
	sort.Sort(common.ByValue(values))
	m, _ := json.Marshal(Export{
		Version:   version(),
		Usage:     usage,
		Nospace:   nospace,
		RawValues: values,
	})
	return string(m)
}

func version() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, dep := range info.Deps {
			if dep.Path == "github.com/rsteube/carapace" {
				return dep.Version
			}
		}
	}
	return "unknown"
}
