package zsh

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	shlex "github.com/carapace-sh/carapace-shlex"
	"github.com/carapace-sh/carapace/internal/common"
	"github.com/carapace-sh/carapace/internal/env"
)

// typeSeparatorSuffixes are the trailing characters that zsh can natively
// re-insert and auto-remove through `_describe -S <suffix>`. Restricting
// removable suffixes to this set (rather than "any non-alphanumeric trailing
// character") avoids misfiring on values that legitimately end in other
// punctuation. It mirrors the `separators` list used in the zsh snippet.
const typeSeparatorSuffixes = `/,.':@=`

var sanitizer = strings.NewReplacer(
	"\n", ``,
	"\r", ``,
	"\t", ``,
)

var quotingReplacer = strings.NewReplacer(
	`'`, `'\''`,
)

var quotingEscapingReplacer = strings.NewReplacer(
	`\`, `\\`,
	`"`, `\"`,
	`$`, `\$`,
	"`", "\\`",
)

var defaultReplacer = strings.NewReplacer(
	`\`, `\\`,
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
	`~`, `\~`,
)

// additional replacement for use with `_describe` in shell script
var describeReplacer = strings.NewReplacer(
	`\`, `\\`,
	`:`, `\:`,
)

func quoteValue(s string) string {
	if strings.HasPrefix(s, "~/") || NamedDirectories.Matches(s) {
		return "~" + defaultReplacer.Replace(strings.TrimPrefix(s, "~")) // assume file path expansion
	}
	return defaultReplacer.Replace(s)
}

type state int

const (
	DEFAULT_STATE state = iota
	// Word starts with `"`.
	// Values need to end with `"` as well.
	// Weirdly regardless whether there are additional quotes within the word.
	QUOTING_ESCAPING_STATE
	// Word starts with `'`.
	// Values need to end with `'` as well.
	// Weirdly regardless whether there are additional quotes within the word.
	QUOTING_STATE
	// Word starts and ends with `"`.
	// Space suffix somehow ends up within the quotes.
	//    `"action"<TAB>`
	//    `"action "<CURSOR>`
	// Workaround for now is to force nospace.
	FULL_QUOTING_ESCAPING_STATE
	// Word starts and ends with `'`.
	// Space suffix somehow ends up within the quotes.
	//    `'action'<TAB>`
	//    `'action '<CURSOR>`
	// Workaround for now is to force nospace.
	FULL_QUOTING_STATE
)

// ActionRawValues formats values for zsh
func ActionRawValues(currentWord string, meta common.Meta, values common.RawValues) string {
	splitted, err := shlex.Split(env.Compline())
	state := DEFAULT_STATE
	if err == nil {
		rawValue := splitted.CurrentToken().RawValue
		// TODO use token state to determine actual state (might have mixture).
		switch {
		case regexp.MustCompile(`^'$|^'.*[^']$`).MatchString(rawValue):
			state = QUOTING_STATE
		case regexp.MustCompile(`^"$|^".*[^"]$`).MatchString(rawValue):
			state = QUOTING_ESCAPING_STATE
		case regexp.MustCompile(`^".*"$`).MatchString(rawValue):
			state = FULL_QUOTING_ESCAPING_STATE
		case regexp.MustCompile(`^'.*'$`).MatchString(rawValue):
			state = FULL_QUOTING_STATE
		}
	}

	noprefix := false
	tagGroup := make([]string, 0)
	values.EachTag(func(tag string, values common.RawValues) {
		// Split a tag's values into groups sharing the same removable suffix so
		// that zsh can manage suffix insertion/removal natively (e.g. the `/` of
		// a directory or the `=` of an optarg flag). Values without a dedicated
		// suffix stay in the empty-suffix group and are formatted exactly as
		// before.
		for _, group := range groupValuesBySuffix(values, meta, state) {
			suffix := group.suffix
			vals := make([]string, len(group.values))
			displays := make([]string, len(group.values))
			for index, val := range group.values {
				value := sanitizer.Replace(val.Value)
				// Strip the removable suffix from the value: zsh re-inserts it
				// through `_describe -S <suffix>`. Done before quoting so the
				// (unquoted) separator character is matched reliably.
				if suffix != "" {
					value = strings.TrimSuffix(value, suffix)
				}

				switch state {
				case QUOTING_ESCAPING_STATE:
					value = quotingEscapingReplacer.Replace(value)
					value = describeReplacer.Replace(value)
					value = value + `"`
				case QUOTING_STATE:
					value = quotingReplacer.Replace(value)
					value = describeReplacer.Replace(value)
					value = value + `'`
				case FULL_QUOTING_ESCAPING_STATE:
					value = quotingEscapingReplacer.Replace(value)
					value = describeReplacer.Replace(value)
				case FULL_QUOTING_STATE:
					value = quotingReplacer.Replace(value)
					value = describeReplacer.Replace(value)
				default:
					value = quoteValue(value)
					value = describeReplacer.Replace(value)
				}

				// Values carrying a removable suffix are always nospace (zsh
				// re-inserts the suffix), so the trailing space is only ever
				// appended for the empty-suffix group.
				if suffix == "" && !meta.Nospace.Matches(val.Value) {
					switch state {
					case FULL_QUOTING_ESCAPING_STATE, FULL_QUOTING_STATE: // nospace workaround
					default:
						value += " "
					}
				}

				display := sanitizer.Replace(val.Display)
				display = describeReplacer.Replace(display) // TODO check if this needs to be applied to description as well
				description := sanitizer.Replace(val.Description)

				vals[index] = value
				noprefix = noprefix || meta.NoPrefix.Matches(strings.TrimPrefix(value, currentWord)) // TODO likely not correct at all times (e.g. when value gets quoted)

				if strings.TrimSpace(description) == "" {
					displays[index] = display
				} else {
					displays[index] = fmt.Sprintf("%v:%v", display, description)
				}
			}
			tagGroup = append(tagGroup, strings.Join([]string{tag, suffix, strings.Join(displays, "\n"), strings.Join(vals, "\n")}, "\003"))
		}
	})
	return fmt.Sprintf("%v\001%v\001%v\001%v\001", zstyles{values}.Format(), message{meta}.Format(), noprefix, strings.Join(tagGroup, "\002")+"\002")
}

// suffixGroup is a set of values that share the same removable suffix.
type suffixGroup struct {
	suffix string
	values common.RawValues
}

// groupValuesBySuffix partitions a tag's values by their removable suffix.
//
// A value gets a dedicated suffix group only when all of the following hold:
//   - the completion is unquoted (DEFAULT_STATE) — in quoted states the value
//     gains a trailing quote, so re-inserting a separator after it via
//     `_describe -S` would place it outside the quotes (e.g. `foo'/`);
//   - an active nospace matcher applies to the value (so no space would be
//     appended anyway);
//   - the value ends in one of the recognised type-separator characters.
//
// Every other value stays in the empty-suffix ("") group and is rendered
// exactly as it was before this feature. Groups are returned in a
// deterministic (sorted) order so the emitted completion payload is stable.
func groupValuesBySuffix(values common.RawValues, meta common.Meta, state state) []suffixGroup {
	order := make([]string, 0)
	groups := make(map[string]common.RawValues)
	for _, val := range values {
		suffix := ""
		if state == DEFAULT_STATE && len(val.Value) > 0 && meta.Nospace.Matches(val.Value) {
			if last := val.Value[len(val.Value)-1]; strings.IndexByte(typeSeparatorSuffixes, last) >= 0 {
				suffix = string(last)
			}
		}
		if _, ok := groups[suffix]; !ok {
			order = append(order, suffix)
		}
		groups[suffix] = append(groups[suffix], val)
	}

	sort.Strings(order)
	result := make([]suffixGroup, 0, len(order))
	for _, suffix := range order {
		result = append(result, suffixGroup{suffix: suffix, values: groups[suffix]})
	}
	return result
}
