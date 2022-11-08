package zsh

import (
	"fmt"
	"strings"

	"github.com/rsteube/carapace/internal/common"
	"github.com/rsteube/carapace/pkg/style"
	"github.com/rsteube/carapace/third_party/github.com/elves/elvish/pkg/ui"
)

var sanitizer = strings.NewReplacer(
	"\n", ``,
	"\r", ``,
	"\t", ``,
)

var displaySanitizer = strings.NewReplacer(
	"\n", ``,
	"\r", ``,
	"\t", ``,
	`:`, `\:`,
)

// TODO verify these are correct/complete (copied from bash).
var quoter = strings.NewReplacer(
	`&`, `\&`,
	`<`, `\<`,
	`>`, `\>`,
	"`", "\\`",
	`'`, `\'`,
	`"`, `\"`,
	`{`, `\{`,
	`}`, `\}`,
	`$`, `\$`,
	`#`, `\#`,
	`|`, `\|`,
	`?`, `\?`,
	`(`, `\(`,
	`)`, `\)`,
	`;`, `\;`,
	` `, `\ `,
	`[`, `\[`,
	`]`, `\]`,
	`*`, `\*`,
	`\`, `\\`,
	`~`, `\~`,
	`:`, `\:`,
)

func quoteValue(s string) string {
	if strings.HasPrefix(s, "~/") || NamedDirectories.Matches(s) {
		return "~" + quoter.Replace(strings.TrimPrefix(s, "~")) // assume file path expansion
	}

	return quoter.Replace(s)
}

// ActionRawValues formats values for zsh.
func ActionRawValues(currentWord string, nospace bool, values common.RawValues) string {
	// First compute paddings and filter out any completions we don't need.
	filtered, maxLength := getCompletionPadding(currentWord, values)

	// Generate and assemble all completions structured in groups, and zstyles
	groups, zstyles := getGroupedComps(filtered, maxLength, nospace)

	// TODO disable styling for large amount of values (bad performance)
	if len(zstyles) > 1000 {
		zstyles = make([]string, 0)
	}

	// The first line is a header containing any message, and an indication to the shell
	// telling it if we want to complete something or not (irrespective of the number of comps)
	return fmt.Sprintf("%v\n%v\n%v", makeHeader(), strings.Join(zstyles, ":"), strings.Join(groups, "\n"))
}

// getGroupedComps classifies all completions in their respective groups, makes the corresponding header
// and completion strings, and assembles all of them into a single string to be passed to ZSH.
// Also takes care of preparing the zstyles format string.
func getGroupedComps(vals []common.RawValue, maxLength int, nospace bool) (comps, zstyles []string) {
	valueStyle, descriptionStyle := setDefaultValueStyle()

	groups := make(map[string][]string) // tag:group => []formattedComps
	var headers []string                // Keeps track of the order

	// Prepare all completions and put them in their respective groups
	for _, val := range vals {
		// Prepare all values
		val = sanitizeCompletion(val, valueStyle, nospace)         // Sanitize each part of the completion (actual/display/description)
		comp := formatCompletion(val)                              // Generate the completion candidate string, with description if needed
		compStyle := formatStyle(val, descriptionStyle, maxLength) // Generate the style for this completion.
		compGroup := setGroup(val)                                 // Generate the tag:group header

		// Get the group for the completion, creating it if needed.
		group, exists := groups[compGroup]
		if !exists {
			group = make([]string, 0)
			groups[compGroup] = group

			headers = append(headers, compGroup)
		}

		// And store the comp and its zstyle
		group = append(group, comp)
		groups[compGroup] = group

		zstyles = append(zstyles, compStyle)
	}

	// Assemble all groups' headers and completions together
	for _, header := range headers {
		completions := strings.Join(groups[header], " ")
		comps = append(comps, fmt.Sprintf("%v %v", header, completions))
	}

	return comps, zstyles
}

func setGroup(val common.RawValue) string {
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

// formatCompletion generates the completion string to pass to ZSH.
func formatCompletion(val common.RawValue) (comp string) {
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

// formatCompletion generates the style string to pass to ZSH.
func formatStyle(val common.RawValue, descriptionStyle string, maxLength int) (compStyle string) {
	if strings.TrimSpace(val.Description) == "" {
		compStyle = formatZstyle(fmt.Sprintf("(%v)()",
			zstyleQuoter.Replace(val.Display)),
			val.Style, descriptionStyle)
	} else {
		compStyle = formatZstyle(fmt.Sprintf("(%v)(%v)( -- %v)",
			zstyleQuoter.Replace(val.Display),                 // First (%v)
			strings.Repeat(" ", maxLength-len(val.Display)+1), // Second (%v): padding
			zstyleQuoter.Replace(val.TrimmedDescription())),   // Third (%v): descriptions
			val.Style, descriptionStyle)
	}

	return
}

// getCompletionPadding computes the padding needed for candidates, filters the values that we do
// not need as completions, and determines if we must pass descriptions along with values.
func getCompletionPadding(currentWord string, values common.RawValues) ([]common.RawValue, int) {
	filtered := make([]common.RawValue, 0)

	maxLength := 0
	hasDescriptions := false

	for _, raw := range values {
		if strings.HasPrefix(raw.Value, currentWord) {
			filtered = append(filtered, raw)

			// Since the _describe call in the ZSH snippet does not support unquoting values like
			// with compadd -Q, and since we must use "candidate:description" format, we must use
			// a special sanitizer for display values (an example of this is IPv6 addresses).
			if length := len(displaySanitizer.Replace(raw.Display)); length > maxLength {
				maxLength = length
			}

			hasDescriptions = hasDescriptions || raw.Description != ""
		}
	}

	return filtered, maxLength
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

var zstyleQuoter = strings.NewReplacer(
	"#", `\#`,
	"*", `\*`,
	"(", `\(`,
	")", `\)`,
	"[", `\[`,
	"]", `\]`,
	"|", `\|`,
	"~", `\~`,
)

// formatZstyle creates a zstyle matcher for given display stings.
// `compadd -l` (one per line) accepts ansi escape sequences in display value but it seems in tabular view these are removed.
// To ease matching in list mode, the display values have a hidden `\002` suffix.
func formatZstyle(s, _styleValue, _styleDescription string) string {
	// return fmt.Sprintf("=(#b)%v=0=%v=%v", s, style.SGR(_styleValue), style.SGR(_styleDescription))
	return fmt.Sprintf("=(#b)%v=0=%v=%v=%v", s, style.SGR(_styleValue), style.SGR(_styleDescription+" bg-default"), style.SGR(_styleDescription))
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
