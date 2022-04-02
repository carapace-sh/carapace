package export

import (
	"encoding/json"
	"runtime/debug"

	"github.com/rsteube/carapace/internal/common"
)

type Export struct {
	Version   string
	Nospace   bool
	RawValues common.RawValues
}

func ActionRawValues(currentWord string, nospace bool, values common.RawValues) string {
	m, _ := json.Marshal(Export{
		Version:   version(),
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
