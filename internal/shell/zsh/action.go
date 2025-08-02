package zsh

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"

	shlex "github.com/carapace-sh/carapace-shlex"
	"github.com/carapace-sh/carapace/internal/common"
	"github.com/carapace-sh/carapace/internal/env"
	"github.com/carapace-sh/carapace/internal/log"
)

var sanitizer = strings.NewReplacer(
	"\n", ``,
	"\r", ``,
	"\t", ``,
)

var quotingReplacer = strings.NewReplacer(
	`'`, `'''`,
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
	start := time.Now()
	defer func() {
		log.LOG.Printf("zsh action processing took %v", time.Since(start))
	}()

	splitted, err := shlex.Split(env.Compline())
	state := DEFAULT_STATE
	if err == nil {
		rawValue := splitted.CurrentToken().RawValue
		// TODO use token state to determine actual state (might have mixture).
		switch {
		case regexp.MustCompile(`^'$|^'.*[^']$`).MatchString(rawValue):
			state = QUOTING_STATE
		case regexp.MustCompile(`^"$|^\".*[^\"]$`).MatchString(rawValue):
			state = QUOTING_ESCAPING_STATE
		case regexp.MustCompile(`^\".*\"$`).MatchString(rawValue):
			state = FULL_QUOTING_ESCAPING_STATE
		case regexp.MustCompile(`^'.*'$`).MatchString(rawValue):
			state = FULL_QUOTING_STATE
		}
	}

	tagGroups := make([]string, 0)
	values.EachTag(func(tag string, tagValues common.RawValues) {
		for suffix, suffixedValues := range groupValuesBySuffix(tagValues, meta, state) {
			vals := make([]string, len(suffixedValues))
			displays := make([]string, len(suffixedValues))

			for index, val := range suffixedValues {
				value := sanitizer.Replace(val.Value)
				if suffix != " " && suffix != "" {
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

				display := sanitizer.Replace(val.Display)
				display = describeReplacer.Replace(display) // TODO check if this needs to be applied to description as well
				description := sanitizer.Replace(val.Description)

				vals[index] = value
				if strings.TrimSpace(description) == "" {
					displays[index] = display
				} else {
					displays[index] = fmt.Sprintf("%v:%v", display, description)
				}
			}
			tagGroups = append(tagGroups, strings.Join([]string{tag, suffix, strings.Join(displays, "\n"), strings.Join(vals, "\n")}, "\003"))
		}
	})
	return fmt.Sprintf("%v\001%v\001%v\002", zstyles{values}.Format(), message{meta}.Format(), strings.Join(tagGroups, "\002"))
}

func groupValuesBySuffix(values common.RawValues, meta common.Meta, state state) map[string]common.RawValues {
	groups := make(map[string]common.RawValues)
	for _, val := range values {
		suffix := ""
		var removableSuffix bool

		// If the last character is not an alphanumeric character, we assume that
		// this character should be removed if the user inserts either a space, the
		// same character or any type separator character ( /,.:@=)
		if len(val.Value) > 0 {
			lastChar := val.Value[len(val.Value)-1:]
			removableSuffix = !isAlphaNumeric(rune(lastChar[len(lastChar)-1]))
		}

		// If we have matched the value suffix against carapace-registered suffixes,
		// and the suffix is not an alphanumeric, then we should register the suffix
		// as removable by ZSH (ie. ZSH will handle automatic insert/erase)
		if meta.Nospace.Matches(val.Value) && len(val.Value) > 0 && removableSuffix {
			lastChar := val.Value[len(val.Value)-1:]
			suffix = lastChar
		}

		if suffix == "" {
			if !meta.Nospace.Matches(val.Value) {
				switch state {
				case FULL_QUOTING_ESCAPING_STATE, FULL_QUOTING_STATE: // nospace workaround
				default:
					suffix = " "
				}
			}
		}

		if _, ok := groups[suffix]; !ok {
			groups[suffix] = make(common.RawValues, 0)
		}
		groups[suffix] = append(groups[suffix], val)
	}
	return groups
}

func isAlphaNumeric(suffix rune) bool {
	return unicode.IsDigit(suffix) || unicode.IsNumber(suffix) || unicode.IsLetter(suffix)
}
