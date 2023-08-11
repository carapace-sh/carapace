package bash

import (
	"fmt"
	"os"
	"strings"

	"github.com/rsteube/carapace/internal/common"
)

var sanitizer = strings.NewReplacer(
	"\n", ``,
	"\r", ``,
	"\t", ``,
)

var valueReplacer = strings.NewReplacer(
	`\`, `\\`,
	`"`, `\"`,
	`$`, `\$`,
	"`", "\\`",
)

var displayReplacer = strings.NewReplacer(
	`${`, `\\\${`,
)

func commonPrefix(a, b string) string {
	i := 0
	for i < len(a) && i < len(b) && a[i] == b[i] {
		i++
	}
	return a[0:i]
}

func commonDisplayPrefix(values ...common.RawValue) (prefix string) {
	for index, val := range values {
		if index == 0 {
			prefix = val.Display
		} else {
			prefix = commonPrefix(prefix, val.Display)
		}
	}
	return
}

func commonValuePrefix(values ...common.RawValue) (prefix string) {
	for index, val := range values {
		if index == 0 {
			prefix = val.Value
		} else {
			prefix = commonPrefix(prefix, val.Value)
		}
	}
	return
}

// ActionRawValues formats values for bash.
func ActionRawValues(currentWord string, meta common.Meta, values common.RawValues) string {
	lastSegment := currentWord // last segment of currentWord split by COMP_WORDBREAKS

	// TODO handle wordbreaks correctly with carapace-shlex
	// for valueIndex, value := range values {
	// TODO optimize
	// if wordbreaks, ok := os.LookupEnv("COMP_WORDBREAKS"); ok {
	// 	wordbreaks = strings.Replace(wordbreaks, " ", "", -1)
	// 	if index := strings.LastIndexAny(currentWord, wordbreaks); index != -1 {
	// 		values[valueIndex].Value = strings.TrimPrefix(value.Value, currentWord[:index+1])
	// 		lastSegment = currentWord[index+1:]
	// 	}
	// }
	// }

	if len(values) > 1 && commonDisplayPrefix(values...) != "" {
		// When all display values have the same prefix bash will insert is as partial completion (which skips prefixes/formatting).
		if valuePrefix := commonValuePrefix(values...); lastSegment != valuePrefix {
			// replace values with common value prefix
			values = common.RawValuesFrom(commonValuePrefix(values...))
		} else {
			// prevent insertion of partial display values by prefixing one with space
			values[0].Display = " " + values[0].Display
		}
		meta.Nospace.Add('*')
	}

	nospace := true
	vals := make([]string, len(values))
	for index, val := range values {
		if len(values) == 1 {
			if !meta.Nospace.Matches(val.Value) {
				nospace = false
			}

			vals[index] = sanitizer.Replace(val.Value)
			if requiresQuoting(vals[index]) {
				vals[index] = valueReplacer.Replace(vals[index])
				switch {
				case strings.HasPrefix(vals[index], "~"): // assume homedir expansion
					vals[index] = fmt.Sprintf(`~"%v"`, strings.TrimPrefix(vals[index], "~"))
				default:
					vals[index] = fmt.Sprintf(`"%v"`, vals[index])
				}
			}
		} else {
			val.Display = displayReplacer.Replace(val.Display)
			val.Description = displayReplacer.Replace(val.Description)
			if val.Description != "" {
				vals[index] = fmt.Sprintf("%v (%v)", val.Display, sanitizer.Replace(val.TrimmedDescription()))
			} else {
				vals[index] = val.Display
			}
		}
	}
	return fmt.Sprintf("%v\001%v", nospace, strings.Join(vals, "\n"))
}

func requiresQuoting(s string) bool {
	chars := " \t\r\n`" + `[]{}()<>;|$&:*#`
	chars += os.Getenv("COMP_WORDBREAKS")
	chars += `\`
	return strings.ContainsAny(s, chars)

}
