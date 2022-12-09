package export

import (
	"encoding/json"
	"runtime/debug"
	"sort"

	"github.com/rsteube/carapace/internal/common"
)

type Export struct {
	Version string
	common.Meta
	RawValues common.RawValues
}

func ActionRawValues(currentWord string, meta common.Meta, values common.RawValues) string {
	sort.Sort(common.ByValue(values))
	m, _ := json.Marshal(Export{
		Version:   version(),
		Meta:      meta,
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
