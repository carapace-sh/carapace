package fish

import (
	"fmt"
	"github.com/rsteube/carapace/internal/common"
	"strings"
)

var sanitizer = strings.NewReplacer(
	"\n", ``,
	"\r", ``,
	"\t", ``,
)

func ActionRawValues(callbackValues string, nospace bool, values ...common.RawValue) string {
	vals := make([]string, len(values))
	for index, val := range values {
		vals[index] = fmt.Sprintf("%v\t%v", sanitizer.Replace(val.Value), sanitizer.Replace(val.Description))
	}
	return strings.Join(vals, "\n")
}
