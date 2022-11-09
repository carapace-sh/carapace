package zsh

import (
	"fmt"
	"strings"

	"github.com/rsteube/carapace/internal/common"
	"github.com/rsteube/carapace/pkg/style"
	"github.com/rsteube/carapace/third_party/github.com/elves/elvish/pkg/ui"
)

// ActionRawValues formats values, structured by tag:groups, for zsh.
func ActionRawValues(currentWord string, nospace bool, values common.RawValues) string {
	valueStyle, descriptionStyle := setDefaultValueStyle()

	// First group completions according to their tag:group, so that they can compute
	// their paddings and styles independently. Filters candidates not matching the prefix.
	// We do NOT sanitize values yet.
	headers, groups := groupValues(values, currentWord, valueStyle, nospace)

	var zstyles []string
	var groupedComps []string

	for _, header := range headers {
		group := groups[header]

		// Make completion strings for all values, which don't need any padding
		// We perform value sanitization in the formatCompletion function only,
		// since after this, when computing paddings, we need to use values as
		// they will actually be printed, which is the same as not being sanitized.
		var comps []string

		for _, val := range group {
			comps = append(comps, formatCompletion(val, valueStyle, nospace))
		}

		groupedComps = append(groupedComps, fmt.Sprintf("%v %v", header, strings.Join(comps, " ")))

		// Get the correct padding for each value, including those to proposed
		// as completions on the same line (aliases and short/long flags).
		// Note that the values in the group are not sanitized, to get the right padding.
		paddings := getGroupPaddings(group)

		// Get all formatted zstyles for these completions
		zstyles = append(zstyles, formatStyles(group, paddings, descriptionStyle)...)
	}

	// TODO disable styling for large amount of values (bad performance)
	if len(zstyles) > 1000 {
		zstyles = make([]string, 0)
	}

	// The first line is a header containing any message, and an indication to the shell
	// telling it if we want to complete something or not (irrespective of the number of comps)
	return fmt.Sprintf("%v\n%v\n%v", makeHeader(), strings.Join(zstyles, ":"), strings.Join(groupedComps, "\n"))
}

// groupValues groups all completions in their groups, and sanitizes the raw completion values.
func groupValues(vals []common.RawValue, current, style string, nospace bool) ([]string, map[string][]common.RawValue) {
	var headers []string

	groups := make(map[string][]common.RawValue)

	for _, val := range vals {
		if !strings.HasPrefix(val.Value, current) {
			continue
		}

		// val = sanitizeCompletion(val, style, nospace) // Sanitize each part of the completion (actual/display/description)
		groupHeader := setGroupHeader(val) // Generate the tag:group header

		group, exists := groups[groupHeader]
		if !exists {
			group = make([]common.RawValue, 0)
			groups[groupHeader] = group

			headers = append(headers, groupHeader)
		}

		group = append(group, val)
		groups[groupHeader] = group
	}

	return headers, groups
}

// getGroupPaddings returns a map with each value candidate as key,
// and its padding as value.  This padding is what is predicted to be
// used as padding by ZSH, so that we can compute correct format strings.
func getGroupPaddings(vals []common.RawValue) map[string]int {
	paddings := make(map[string]int)
	maxLength := 0
	hasDescriptions := false

	for _, raw := range vals {
		// Since the _describe call in the ZSH snippet does not support unquoting values like
		// with compadd -Q, and since we must use "candidate:description" format, we must use
		// a special sanitizer for display values (an example of this is IPv6 addresses).
		length := len(raw.Display)
		// length := len(sanitizer.Replace(raw.Display))
		if length > maxLength {
			maxLength = length
		}

		hasDescriptions = hasDescriptions || raw.Description != ""
	}

	for _, raw := range vals {
		paddings[raw.Value] = maxLength - len(raw.Value) + 1
	}

	return paddings
}

// formatCompletion generates the completion string to pass to ZSH.
func formatCompletion(val common.RawValue, style string, nospace bool) (comp string) {
	// Sanitize each part of the completion (actual/display/description)
	val = sanitizeCompletion(val, style, nospace)

	// We quote the entire string with single quotes, so that the ZSH script can split
	// them correctly into an array, and also to preserve any special characters.
	if strings.TrimSpace(val.Description) == "" {
		comp = fmt.Sprintf("'%v\t%v'", val.Value, val.Display)
	} else {
		// Note the use of : as separator between completion and description.
		comp = fmt.Sprintf("'%v\t%v:%v'", val.Value, val.Display, val.TrimmedDescription())
	}

	return
}

// formatStyles makes the styles strings for completions in a group, respecting their padding.
func formatStyles(vals []common.RawValue, paddings map[string]int, descStyle string) (zstyles []string) {
	for _, val := range vals {
		padding := paddings[val.Value]

		var compStyle string

		if strings.TrimSpace(val.Description) == "" {
			compStyle = formatZstyle(fmt.Sprintf("(%v)()",
				zstyleQuoter.Replace(val.Display)),
				val.Style, descStyle)
		} else {
			// compStyle = formatZstyle(fmt.Sprintf("(%v%v*)(-- %v)",
			// 	zstyleQuoter.Replace(val.Display),               // First (%v)
			// 	strings.Repeat(" ", maxLength-len(val.Display)), // Second (%v): padding
			// 	zstyleQuoter.Replace(val.TrimmedDescription())), // Third (%v): descriptions
			// 	val.Style, descriptionStyle)

			// compStyle = formatZstyle(fmt.Sprintf("(%v%v*)( -- %v)",
			// 	zstyleQuoter.Replace(val.Display),               // First (%v)
			// 	strings.Repeat(" ", maxLength-len(val.Display)), // Second (%v): padding
			// 	zstyleQuoter.Replace(val.TrimmedDescription())), // Third (%v): descriptions
			// 	val.Style, descriptionStyle)

			compStyle = formatZstyle(fmt.Sprintf("(%v)(%v)( -- %v)",
				zstyleQuoter.Replace(val.Display),               // First (%v)
				strings.Repeat(" ", padding),                    // Second (%v): padding
				zstyleQuoter.Replace(val.TrimmedDescription())), // Third (%v): descriptions
				val.Style, descStyle)
		}

		zstyles = append(zstyles, compStyle)
	}

	return zstyles
}

// setGroupHeader checks that all completions have a group, sets default if needed.
func setGroupHeader(val common.RawValue) string {
	// Set defaults
	if val.Tag == "" {
		val.Tag = "values"
	}

	if val.Group == "" {
		val.Group = "completions"
	}

	tag := quoter.Replace(val.Tag)
	group := quoter.Replace(val.Group)

	// We escape both the tag/group strings, and
	// the entire string itself, like for completions.
	return fmt.Sprintf("'%v:%v'", tag, group)
}

// sanitizeCompletion applies a series of string sanitizers to the completion
// candidate, its display value and its description.
func sanitizeCompletion(val common.RawValue, valueStyle string, nospace bool) common.RawValue {
	// The Value is what will actually be inserted in the command-line.
	val.Value = sanitizer.Replace(val.Value)
	val.Value = quoteValue(val.Value)

	if nospace {
		val.Value += "\001"
	}

	// The display value is used only when displaying the completions,
	// and needs a different sanitizer, in order not to mess up the shell
	// when it tries to know where the description starts (eg. with IPv6)
	val.Display = displaySanitizer.Replace(val.Display)

	// Then sanitize the description only.
	val.Description = sanitizer.Replace(val.Description)

	// Style
	if val.Style == "" || ui.ParseStyling(val.Style) == nil {
		val.Style = valueStyle
	}

	return val
}

func setDefaultValueStyle() (valueStyle, descriptionStyle string) {
	valueStyle = "default"
	if s := style.Carapace.Value; s != "" && ui.ParseStyling(s) != nil {
		valueStyle = s
	}

	descriptionStyle = "default"
	if s := style.Carapace.Description; s != "" && ui.ParseStyling(s) != nil {
		descriptionStyle = s
	}

	return valueStyle, descriptionStyle
}

// formatZstyle creates a zstyle matcher for given display stings.
// `compadd -l` (one per line) accepts ansi escape sequences in display value but it seems in tabular view these are removed.
// To ease matching in list mode, the display values have a hidden `\002` suffix.
func formatZstyle(s, styleValue, styleDescription string) string {
	zstyle := fmt.Sprintf("=(#b)%v=0=%v=%v=%v", s,
		style.SGR(styleValue),
		style.SGR(styleDescription+" bg-default"),
		style.SGR(styleDescription))

	// return fmt.Sprintf("=(#b)%v=0=%v=%v", s, style.SGR(_styleValue), style.SGR(_styleDescription))
	return zstyle
}

// Creates a header line with some indications for the shell caller.
func makeHeader() (header string) {
	// TODO: Find a way to know if actually want to complete something.
	header += "0"

	header += "\t"

	// Format the completion message if needed
	if common.CompletionMessage == "" {
		return
	}

	header += fmt.Sprintf("\x1b[%vm%v\x1b[%vm %v\x1b[%vm",
		style.SGR(style.Carapace.Error),
		"ERR",
		style.SGR("fg-default"),
		sanitizer.Replace(common.CompletionMessage),
		style.SGR("fg-default"),
	)

	return
}

// Notes:

// Completions ========================
// This is used with _describe, but actually takes care of adding padding.
// This is not really useful, as _describe internals will actually take care of computing it,
// including when several candidates have the same description.
// vals[index] = fmt.Sprintf("%v\t%v%v:%v", val.Value, val.Display, strings.Repeat(" ", maxLength-len(val.Display)+1), val.TrimmedDescription())

// Styling ============================

// This one considers we don't need padding, since _describe will compute a new display string anyway.
// So it should work at times, but will not as soon as some padding is applied in the display string.
// formatZstyle(fmt.Sprintf("(%v)(-- %v)", zstyleQuoter.Replace(val.Display), zstyleQuoter.Replace(val.TrimmedDescription())), val.Style, descriptionStyle)

// This one does not work either, probably for several reasons:
// - _describe modifies the candidate padding, so (%v%v) becomes obsolete.
// - For some reason (-- %v) not havin a space first is not recognized, while ( -- %v) works.
// formatZstyle(fmt.Sprintf("(%v%v) (-- %v)", zstyleQuoter.Replace(val.Display), strings.Repeat(" ", maxLength-len(val.Display)), zstyleQuoter.Replace(val.TrimmedDescription())), val.Style, descriptionStyle)

// Originals used with compadd ========
// vals[index] = fmt.Sprintf("%v\t%v%v-- %v", val.Value, val.Display, strings.Repeat(" ", maxLength-len(val.Display)+1), val.TrimmedDescription())
// formatZstyle(fmt.Sprintf("(%v)(%v)(-- %v)", zstyleQuoter.Replace(val.Display), strings.Repeat(" ", maxLength-len(val.Display)+1), zstyleQuoter.Replace(val.TrimmedDescription())), val.Style, descriptionStyle)
