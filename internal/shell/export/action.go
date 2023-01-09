package export

import (
	"encoding/json"
	"runtime/debug"
	"sort"

	"github.com/rsteube/carapace/internal/common"
)

type Export struct {
	Version string `json:"version"`
	common.Meta
	Values common.RawValues `json:"values"`
}

func ActionRawValues(currentWord string, meta common.Meta, values common.RawValues) string {
	sort.Sort(common.ByValue(values))
	m, _ := json.Marshal(Export{
		Version: version(),
		Meta:    meta,
		Values:  values,
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
