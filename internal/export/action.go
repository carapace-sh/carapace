package export

import (
	"encoding/json"

	"github.com/rsteube/carapace/internal/common"
)

type Export struct {
	Nospace   bool
	RawValues common.RawValues
}

func ActionRawValues(currentWord string, nospace bool, values common.RawValues) string {
	m, _ := json.Marshal(Export{
		Nospace:   nospace,
		RawValues: values,
	})
	return string(m)
}
