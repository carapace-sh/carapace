package zsh

import (
	"fmt"
	"strings"

	"github.com/rsteube/carapace/internal/common"
)

// values structures an arbitrary number of completions into more structured groups,
// and builds appropriate suffixes and removable patterns for them.
type values struct {
	headers []string
	vals    map[string][]common.RawValue

	suffix       string
	suffixRemove string
}

// scanValues makes a first pass at all the values and gathers as much information as possible,
// to set various default stuff, gather completions in groups, and compute suffixes depending on
// the nature of the completion values.
func scanValues(vals []common.RawValue, current string, nospace bool) (values, bool, int) {
	comps := values{
		vals: make(map[string][]common.RawValue),
	}

	maxLength := 0
	onlyCommands := true

	for _, val := range vals {
		if !strings.HasPrefix(val.Value, current) {
			continue
		}

		// Generate the tag:group header and store value
		groupHeader := setGroupHeader(val)

		group, exists := comps.vals[groupHeader]
		if !exists {
			group = make([]common.RawValue, 0)
			comps.vals[groupHeader] = group

			comps.headers = append(comps.headers, groupHeader)
		}

		group = append(group, val)
		comps.vals[groupHeader] = group

		// Update maximum global display padding
		length := len(val.Display)
		if length > maxLength {
			maxLength = length
		}

		if val.Tag != "command" {
			onlyCommands = false
		}

		// Set suffix modifiers, considering any tags found.
		if !nospace {
			comps.suffix = " "
		}

		comps.suffixRemove = suffixRemovePatterns(val, nospace)
	}

	return comps, onlyCommands, maxLength
}

// setGroupHeader checks that all completions have a group, sets default if needed.
func setGroupHeader(val common.RawValue) string {
	// Set defaults
	if val.Tag == "" {
		val.Tag = string(common.Value)
		if val.Group == "" {
			val.Group = val.Tag + "s"
		}
	}

	if val.Group == "" {
		if val.Group == "" {
			val.Group = val.Tag + "s"
		}
	}

	tag := quoter.Replace(val.Tag)
	group := quoter.Replace(val.Group)

	// We escape both the tag/group strings, and
	// the entire string itself, like for completions.
	return fmt.Sprintf("%v:%v", tag, group)
}

// getPadding computes the required padding (global or per group) with safeguards.
func getPadding(valueLen, maxLenGroup, maxLenAll int) (padding string) {
	var paddingLen int

	if maxLenAll != 0 {
		paddingLen = maxLenAll - valueLen + 1
	} else {
		paddingLen = maxLenGroup - valueLen + 1
	}
	if paddingLen < 0 {
		paddingLen = 0
	}

	return strings.Repeat(" ", paddingLen)
}

// getMaxLength returns the length of the longest completion value.
func getMaxLength(vals []common.RawValue) int {
	maxLength := 0

	for _, raw := range vals {
		length := len(raw.Display)
		if length > maxLength {
			maxLength = length
		}
	}

	return maxLength
}

// TODO: Anylize the given value suffixes with regexps here.
func suffixRemovePatterns(val common.RawValue, nospace bool) string {
	if val.Tag == "option" && val.SuffixRemovable == "=" {
		return "^\\-" // ZSH special pattern matcher for any non-nil character
	}

	if val.SuffixRemovable != "" && !nospace {
		return "^\\-"
	}

	// Something useful here ?
	// \"'/=

	return ""
}

func formatSuffixMatchers(suffix, removeSuffix string) (string, string) {
	if suffix == " " && removeSuffix == "" {
		removeSuffix = " "
	}

	// We either have a specific suffix, other than a space.
	if suffix != "" {
		suffix = fmt.Sprintf("%v", suffix)
	}

	// if removeSuffix != "" {
	// 	removeSuffix = fmt.Sprintf("%v", quoter.Replace(removeSuffix), quoter.Replace(suffix))
	// }

	return suffix, removeSuffix
}

func hasAliasedCompletions(vals []common.RawValue) bool {
	allKeys := make(map[string]bool)
	for _, item := range vals {
		if _, value := allKeys[item.Description]; value {
			return true
		}

		allKeys[item.Description] = true
	}

	return false
}
