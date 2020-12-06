package xonsh

import (
	"fmt"
	"strings"

	"github.com/rsteube/carapace/internal/common"
)

var sanitizer = strings.NewReplacer( // TODO
	`$`, ``,
	"`", ``,
	"\n", ``,
	`\`, ``,
	`"`, ``,
	`'`, ``,
	`|`, ``,
	`>`, ``,
	`<`, ``,
	`&`, ``,
	`(`, ``,
	`)`, ``,
	`;`, ``,
	`#`, ``,
)

func Sanitize(values ...string) []string {
	sanitized := make([]string, len(values))
	for index, value := range values {
		sanitized[index] = sanitizer.Replace(value)
	}
	return sanitized
}

func Callback(prefix string, uid string) string {
	return fmt.Sprintf("_%v_callback('%v')", strings.Replace(prefix, "-", "__", -1), uid)
}

func ActionDirectories() string {
	return `{ RichCompletion(f, display=pathlib.PurePath(f).name, description='', prefix_len=0) for f in complete_dir(prefix, line, begidx, endidx, ctx, True)[0]}`
}

// TODO add endswith filter function
func ActionFiles(suffix string) string {
	return `{ RichCompletion(f, display=pathlib.PurePath(f).name, description='', prefix_len=0) for f in complete_path(prefix, line, begidx, endidx, ctx)[0]}`
}

func ActionCandidates(values ...common.Candidate) string {
	vals := make([]string, len(values))
	for index, val := range values {
		vals[index] = fmt.Sprintf(`  RichCompletion('%v', display='%v', description='%v', prefix_len=0),`, sanitizer.Replace(val.Value), sanitizer.Replace(val.Display), sanitizer.Replace(val.Description))
	}
	return fmt.Sprintf("{\n%v\n}", strings.Join(vals, "\n"))
}
