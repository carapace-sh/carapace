package bash

import (
	"fmt"
	"strings"
)

var sanitizer = strings.NewReplacer(
	`$`, ``,
	"`", ``,
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
	return fmt.Sprintf(`eval $(_%v_callback '%v')`, prefix, uid)
}

func ActionExecute(command string) string {
	return fmt.Sprintf(`$(%v)`, command)
}

func ActionDirectories() string {
	return `compgen -S / -d -- "$last"`
}

func ActionFiles(suffix string) string {
	return fmt.Sprintf(`compgen -S / -d -- "$last"; compgen -f -X '!*%v' -- "$last"`, suffix)
}

func ActionNetInterfaces() string {
	return `compgen -W "$(ifconfig -a | grep -o '^[^ :]\+')" -- "$last"`
}

func ActionUsers() string {
	return `compgen -u -- "${last//[\"\|\']/}"`
}

func ActionGroups() string {
	return `compgen -g -- "${last//[\"\|\']/}"`
}

func ActionHosts() string {
	return `compgen -W "$(cut -d ' ' -f1 < ~/.ssh/known_hosts | cut -d ',' -f1)" -- "$last"`
}

func ActionValues(values ...string) string {
	sanitized := Sanitize(values...)
	if len(strings.TrimSpace(strings.Join(sanitized, ""))) == 0 {
		return ActionMessage("no values to complete")
	}

	vals := make([]string, len(sanitized))
	for index, val := range sanitized {
		// TODO escape special characters
		//vals[index] = strings.Replace(val, " ", `\ `, -1)
		vals[index] = strings.Replace(val, ` `, `\\\ `, -1)
	}
	return fmt.Sprintf(`compgen -W $'%v' -- "$last"`, strings.Join(vals, `\n`))
}

func ActionValuesDescribed(values ...string) string {
	// TODO verify length (description always exists)
	vals := make([]string, len(values)/2)
	for index, val := range values {
		if index%2 == 0 {
			vals[index/2] = val
		}
	}
	return ActionValues(vals...)
}

func ActionMessage(msg string) string {
	return ActionValues("ERR", Sanitize(msg)[0])
}

func ActionPrefixValues(prefix string, values ...string) string {
	sanitized := Sanitize(values...)
	if len(strings.TrimSpace(strings.Join(sanitized, ""))) == 0 {
		return ActionMessage("no values to complete")
	}

	vals := make([]string, len(sanitized))
	for index, val := range sanitized {
		// TODO escape special characters
		vals[index] = strings.Replace(val, ` `, `\\\ `, -1)
	}

	if index := strings.LastIndexAny(prefix, ":="); index > -1 {
		// COMP_WORD will split on these characters, so $last does not contain the full argument
		prefix = prefix[index:]
	}
	return fmt.Sprintf(`compgen -W $'%v' -- "${last/%v/}"`, strings.Join(vals, `\n`), prefix)
}
