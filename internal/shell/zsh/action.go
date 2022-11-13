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
	// Compute and adjust all defaults first.
	valueStyle, descStyle := setDefaultValueStyle()

	// First go over all values, in order to:
	// - Filter candidates not matching the current prefix callback
	// - Group them according to their tag:group specifications,so that they
	//   can compute their paddings and styles independently.
	// - Compute a global max padding length, irrespectively of groups
	groups, onlyCommands, maxLen := scanValues(values, currentWord, nospace)

	// We actually only pad values globally when all groups are commands.
	if !onlyCommands {
		maxLen = 0
	}

	var zstyles []string
	var groupValues []string

	for _, header := range groups.headers {
		group := groups.vals[header]

		// Generate all formatted completion strings and associated zstyles for the group.
		values, styles := formatGroup(group, valueStyle, descStyle, maxLen)

		// Generate the complete group string, with tag:group header and its completions,
		groupValues = append(groupValues, fmt.Sprintf("%v\n%v", header, strings.Join(values, "\n")))

		// And append the styles, passed all at once, irrespectively of their groups.
		zstyles = append(zstyles, styles...)
	}

	if len(zstyles) > maxZstyles {
		zstyles = make([]string, 0)
	}

	// Header line : The header contains any message that has to be printed, and computed suffix matchers/removers.
	ret, message := formatMessage()
	suffix, removePatterns := formatSuffixMatchers(groups.suffix, groups.suffixRemove)
	header := strings.Join([]string{ret, message, suffix, removePatterns}, "\t")

	return fmt.Sprintf("%v\n%v\n%v", header, strings.Join(zstyles, ":"), strings.Join(groupValues, "\n\n"))
}

// formatGroup generates all strings (completions and styles) for a given group of completions.
// This function optimizes the number of iterations performed on the group's values (2 passes).
func formatGroup(group []common.RawValue, valueStyle, descStyle string, maxLen int) ([]string, []string) {
	completions := make([]string, len(group))
	zstyles := make([]string, len(group)+1)

	// We want to know if some completions will be displayed on the same line
	// If yes, we must add a default style pattern in order to match descriptions.
	hasAliases := hasAliasedCompletions(group)
	if hasAliases {
		zstyles[len(zstyles)-1] = formatZstyleValue("(-- *)", descStyle)
	}

	// Get the correct padding for each value, including those to proposed
	// as completions on the same line (aliases and short/long flags).
	// Values in the group are not sanitized yet, to get the right padding.
	maxLenGrp := getMaxLength(group)

	for idx, val := range group {
		// Generate completion string for this value, respecting/considering a few things:
		// - If some values are to be displayed next to the same description (eg. -f/--file)
		// - If we must use global padding or per-group padding.
		completions[idx] = formatValue(val, valueStyle, hasAliases, maxLenGrp, maxLen)

		// Generate the corresponding zstyle string.
		zstyles[idx] = formatStyle(val, descStyle, hasAliases, maxLenGrp, maxLen)
	}

	return completions, zstyles
}

// formatValues generates a completion string from a value, taking into account various requirements and parameters.
// Those parameters are mostly here for us to generate values that are conform to their associated zstyles.
func formatValue(val common.RawValue, style string, hasAliases bool, maxLenGrp, maxLenAll int) string {
	// Any padding, if used, must be computed before sanitizing the value
	valueLen := len(val.Value)

	// Sanitize each part of the completion (actual/display/description)
	val = sanitizeCompletion(val, style)

	// Shorthands
	comp, display, desc := val.Value+val.SuffixRemovable, val.Display, val.TrimmedDescription()

	// When the completion has no description, we don't need to take any
	// parameters and constraints into account.
	if strings.TrimSpace(val.Description) == "" {
		return fmt.Sprintf("%v\t%v", comp, display)
	}

	// Else we have a description, and then requirements make the actual string need to vary in format.
	// First, if there are completions to be printed on the same line, we format as needed.
	if hasAliases {
		return fmt.Sprintf("%v\t%v:%v", comp, display, desc)
	}

	// Else, we have a description but no two values in this group are aliases of each other,
	// in which case we can use 'custom display strings', in which we are responsible for padding.
	// We either must use an global padding (all groups), or per-group padding.
	padding := getPadding(valueLen, maxLenGrp, maxLenAll)

	return fmt.Sprintf("%v\t%v%v -- %v", comp, display, padding, desc)
}

// formatStyle makes the style strings for a completion, respecting its padding and considering
// the various parameters also considered when generating the completion strings.
func formatStyle(val common.RawValue, descStyle string, hasAliases bool, maxLenGrp, maxLen int) string {
	// Any padding, if used, must be computed before sanitizing the value
	valueDisplayLen := len(val.Display)

	// Shorthands
	display := zstyleQuoter.Replace(val.Display)
	desc := zstyleQuoter.Replace(val.TrimmedDescription())

	// When the completion has no description, we don't need to take any
	// parameters and constraints into account.
	if strings.TrimSpace(val.Description) == "" {
		return formatZstyle(fmt.Sprintf("(%v)()", display), val.Style, descStyle)
	}

	// Else we have a description, and then requirements make the actual string need to vary.
	// If there are completions to be printed on the same line, we don't care neither about
	// the padding, nor the description (which is going to be set later).
	if hasAliases {
		return formatZstyleValue(fmt.Sprintf("(%v)", display), val.Style)
	}

	// Else, we have a description but no two values in this group are aliases
	// of each other, so we make zstyle that takes all components into account.
	// We either must use an global padding (all groups), or per-group padding.
	padding := getPadding(valueDisplayLen, maxLenGrp, maxLen)

	return formatZstyle(fmt.Sprintf("(%v)(%v)(*-- %v)", display, padding, desc), val.Style, descStyle)
}

// sanitizeCompletion applies a series of string sanitizers to the completion
// candidate, its display value and its description.
func sanitizeCompletion(val common.RawValue, valueStyle string) common.RawValue {
	// The Value is what will actually be inserted in the command-line.
	val.Value = sanitizer.Replace(val.Value)
	val.Value = quoteValue(val.Value)

	// The display value is used only when displaying the completions,
	// and needs a different sanitizer, in order not to mess up the shell
	// when it tries to know where the description starts (eg. with IPv6)
	val.Display = displaySanitizer.Replace(val.Display)
	val.Description = displaySanitizer.Replace(val.Description)

	// Style
	if val.Style == "" || ui.ParseStyling(val.Style) == nil {
		val.Style = valueStyle
	}

	return val
}

// formatZstyle creates a zstyle matcher for given display stings.
// `compadd -l` (one per line) accepts ansi escape sequences in display value but it seems in tabular view these are removed.
// To ease matching in list mode, the display values have a hidden `\002` suffix.
func formatZstyle(s, styleValue, styleDescription string) string {
	zstyle := fmt.Sprintf("=(#b)%v=0=%v=%v=%v", s,
		style.SGR(styleValue),
		style.SGR(styleDescription+" bg-default"),
		style.SGR(styleDescription))

	return zstyle
}

// formatZstyleValue creates a zstyle matcher for given display stings.
// The difference with formatZstyle is that we only match one elemet.
// This is used when matching values that are displayed 'stacked' like short/long flags.
func formatZstyleValue(s, styleValue string) string {
	zstyle := fmt.Sprintf("=(#b)%v=%v", s,
		style.SGR(styleValue))

	return zstyle
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

// builds the first part of the header line, checking either errors or hints.
func formatMessage() (retcode, message string) {
	if common.CompletionMessage != "" {
		// Format the completion message if needed
		retcode = "1" // TODO VERY VERY UGLY

		msg := messageSanitizer.Replace(common.CompletionMessage)
		msg = quoter.Replace(msg)

		message = fmt.Sprintf("\x1b[%vm%v\x1b[%vm %v\x1b[%vm",
			style.SGR(style.Carapace.Error),
			"ERR",
			style.SGR(style.Dim),
			msg,
			style.SGR("fg-default"))

	} else if common.CompletionHint != "" {
		// Format the completion message if needed
		msg := messageSanitizer.Replace(common.CompletionHint)
		msg = quoter.Replace(msg)

		message = fmt.Sprintf("\x1b[%vm%v\x1b[%vm %v\x1b[%vm",
			style.SGR(style.Dim),
			"",
			style.SGR(style.Dim),
			msg,
			style.SGR("fg-default"))
	}

	return retcode, message
}
