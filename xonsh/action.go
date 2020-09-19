package xonsh

import (
	"fmt"
	"strings"
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

func ActionExecute(command string) string {
	return `{}`
	//return fmt.Sprintf(`"%v" | Out-String | InvokeExpression`, strings.Replace(command, "\n", "`n", -1))
}

func ActionDirectories() string {
	return `{ RichCompletion(f, display=pathlib.PurePath(f).name, description='', prefix_len=0) for f in complete_dir(prefix, line, begidx, endidx, ctx, True)[0]}`
}

// TODO add endswith filter function
func ActionFiles(suffix string) string {
	return `{ RichCompletion(f, display=pathlib.PurePath(f).name, description='', prefix_len=0) for f in complete_path(prefix, line, begidx, endidx, ctx)[0]}`
}

func ActionNetInterfaces() string {
	return `{}`
}

func ActionUsers() string {
	return `{}`
}

func ActionGroups() string {
	return `{}`
}

func ActionHosts() string {
	return `{}`
}

func ActionValues(values ...string) string {
	vals := make([]string, len(values)*2)
	for index, val := range values {
		vals[index*2] = val
		vals[(index*2)+1] = ""
	}
	return ActionValuesDescribed(vals...)
}

func ActionValuesDescribed(values ...string) string {
	sanitized := Sanitize(values...)
	// TODO verify length (description always exists)
	vals := make([]string, len(sanitized)/2)
	for index, val := range sanitized {
		if index%2 == 0 {
			vals[index/2] = fmt.Sprintf(`  RichCompletion('%v', display='%v', description='%v', prefix_len=0),`, val, val, sanitized[index+1])
		}
	}
	return fmt.Sprintf("{\n%v\n}", strings.Join(vals, "\n"))
}

func ActionMessage(msg string) string {
	return ActionValuesDescribed("_", msg, "ERR", msg)
}

func ActionPrefixValues(prefix string, values ...string) string {
	sanitized := Sanitize(values...)
	if len(strings.TrimSpace(strings.Join(sanitized, ""))) == 0 {
		return ActionMessage("no values to complete")
	}

	vals := make([]string, len(sanitized))
	for index, val := range sanitized {
		vals[index] = fmt.Sprintf(`  RichCompletion('%v', display='%v', description='%v', prefix_len=0),`, prefix+val, val, "")
	}
	return fmt.Sprintf("{\n%v\n}", strings.Join(vals, "\n"))
}
