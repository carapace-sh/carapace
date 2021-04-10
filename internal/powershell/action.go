package powershell

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/rsteube/carapace/internal/common"
)

var sanitizer = strings.NewReplacer( // TODO
	"\n", ``,
	"\t", ``,
	`'`, "``",
)

type completionResult struct {
	CompletionText string
	ListItemText   string
	ToolTip        string
}

// CompletionResult doesn't like empty parameters, so just replace with space if needed
func ensureNotEmpty(s string) string {
	if s == "" {
		return " "
	}
	return s
}

func ActionRawValues(callbackValue string, nospace bool, values ...common.RawValue) string {
	filtered := common.ByValues(values).Filter(callbackValue)
	sort.Sort(common.ByDisplay(filtered))

	vals := make([]completionResult, 0, len(filtered))
	for _, val := range filtered {
		if val.Value != "" { // must not be empty - any empty `''` parameter in CompletionResult causes an error
			val.Value = sanitizer.Replace(val.Value)

			if strings.ContainsAny(val.Value, ` {}()[]*$?\"|<>&(),;#`+"`") {
				val.Value = fmt.Sprintf("'%v'", val.Value)
			}

			if !nospace {
				val.Value = val.Value + " "
			}

			vals = append(vals, completionResult{
				CompletionText: val.Value,
				ListItemText:   ensureNotEmpty(sanitizer.Replace(val.Display)),
				ToolTip:        ensureNotEmpty(sanitizer.Replace(val.Description)),
			})
		}
	}
	m, _ := json.Marshal(vals)
	return string(m)
}
