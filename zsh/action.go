package zsh

import (
	"fmt"
	"github.com/rsteube/carapace/common"
	"strings"
)

var sanitizer = strings.NewReplacer(
	`$`, ``,
	"`", ``,
	"\n", ``,
	`\`, ``,
	`"`, ``,
	`'`, ``,
	"`", ``,
	`|`, ``,
	`>`, ``,
	`<`, ``,
	`&`, ``,
	`(`, ``,
	`)`, ``,
	`;`, ``,
	`#`, ``,
	`[`, `\[`,
	`]`, `\]`,
)

func Sanitize(values ...string) []string {
	sanitized := make([]string, len(values))
	for index, value := range values {
		sanitized[index] = sanitizer.Replace(value)
	}
	return sanitized
}

func Callback(prefix string, cuid string) string {
	return fmt.Sprintf(`{_%v_callback '%v'}`, prefix, cuid)
}

func ActionDirectories() string {
	return `_files -/`
}

// ActionFiles _path_files with all options except -g and -/. These options depend on file-patterns style setting. // TODO fix doc
// [http://zsh.sourceforge.net/Doc/Release/Completion-System.html#index-_005ffiles]
func ActionFiles(pattern string) string {
	if pattern == "" {
		return "_files"
	} else {
		return fmt.Sprintf("_files -g '%v'", pattern)
	}
}

func ActionCandidates(values ...common.Candidate) string {
	vals := make([]string, len(values))
	displays := make([]string, len(values))
	for index, val := range values {
		// TODO sanitize
		vals[index] = fmt.Sprintf("'%v'", sanitizer.Replace(val.Value))
		if strings.TrimSpace(val.Description) == "" {
			displays[index] = fmt.Sprintf("'%v'", sanitizer.Replace(val.Display))
		} else {
			displays[index] = fmt.Sprintf("'%v (%v)'", sanitizer.Replace(val.Display), sanitizer.Replace(val.Description))
		}
	}
	return fmt.Sprintf("{local _comp_desc=(%v);compadd -S '' -d _comp_desc %v}", strings.Join(displays, " "), strings.Join(vals, " "))
}
