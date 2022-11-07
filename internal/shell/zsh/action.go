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
	// filtered, _ := getCompletionPadding(currentWord, values)
	filtered, maxLength := getCompletionPadding(currentWord, values)

	// Set basic styling things
	valueStyle, descriptionStyle := setDefaultValueStyle()

	// We know the number of completions passed to shell
	// but the number of styles might not be equal.
	vals := make([]string, len(filtered))
	zstyles := make([]string, 0)

	// For each completion candidate
	for index, val := range filtered {
		// Apply sanitizers to each component of the completion (actual/display/description)
		val = sanitizeCompletion(val, valueStyle, nospace)

		if strings.TrimSpace(val.Description) == "" {
			// Candidate used by _describe
			vals[index] = fmt.Sprintf("%v\t%v", val.Value, val.Display)

			// Associated style
			compStyle := formatZstyle(fmt.Sprintf("(%v)()",
				zstyleQuoter.Replace(val.Display)),
				val.Style, descriptionStyle)
			zstyles = append(zstyles, compStyle)
		} else {
			// Candidate used by _describe, note the use of : as separator between completion and description
			vals[index] = fmt.Sprintf("%v\t%v:%v", val.Value, val.Display, val.TrimmedDescription())

			// Associated style
			compStyle := formatZstyle(fmt.Sprintf("(%v)(%v)( -- %v)",
				zstyleQuoter.Replace(val.Display),                 // First (%v)
				strings.Repeat(" ", maxLength-len(val.Display)+1), // Second (%v): padding
				zstyleQuoter.Replace(val.TrimmedDescription())),   // Third (%v): descriptions
				val.Style, descriptionStyle)

			zstyles = append(zstyles, compStyle)
		}
	}

	// TODO disable styling for large amount of values (bad performance)
	if len(zstyles) > 1000 {
		zstyles = make([]string, 0)
	}

	return fmt.Sprintf(":%v\n%v", strings.Join(zstyles, ":"), strings.Join(vals, "\n"))
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
