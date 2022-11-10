package zsh

import (
	"strings"
)

// TODO disable styling for large amount of values (bad performance).
var maxZstyles = 1000

// sanitizer is used only for descriptions and messages.
var sanitizer = strings.NewReplacer(
	"\n", ``,
	"\r", ``,
	"\t", ``,
)

// displaySanitizer is used for candidates and/or their display strings.
// It also escapes colons, since they are used to delimit the candidate
// variable from its description (passed to _describe).
var displaySanitizer = strings.NewReplacer(
	"\n", ``,
	"\r", ``,
	"\t", ``,
	`:`, `\:`,
)

// quoter is used on all completion candidates.
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
	// `?`, `\?`,
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
