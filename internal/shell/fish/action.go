package fish

import (
	"fmt"
	"strings"

	"github.com/rsteube/carapace/internal/common"
)

var sanitizer = strings.NewReplacer(
	"\n", ``,
	"\r", ``,
	"\t", ``,
)

// ActionRawValues formats values for fish.
func ActionRawValues(currentWord string, nospace bool, values common.RawValues) string {
	values = common.AddMessageToValues(currentWord, values)

	vals := make([]string, len(values))
	for index, val := range values {
		vals[index] = fmt.Sprintf("%v\t%v", sanitizer.Replace(val.Value), sanitizer.Replace(val.TrimmedDescription()))
	}
	return strings.Join(vals, "\n")
}
