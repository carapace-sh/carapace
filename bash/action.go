package bash

import (
	"fmt"
	"strings"
)

var sanitizer = strings.NewReplacer(
	`$`, ``,
	"`", ``,
	`\`, ``,
	`"`, `'`,
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

func ActionFiles(suffix string) string {
	return fmt.Sprintf(`compgen -f -o plusdirs -X "!*%v" -- $last`, suffix)
}

func ActionNetInterfaces() string {
	return `compgen -W "$(ifconfig -a | grep -o '^[^ :]\+' | tr '\n' ' ')" -- $last`
}

func ActionUsers() string {
	return `compgen -u -- $last`
}

func ActionGroups() string {
	return `compgen -g -- $last`
}

func ActionHosts() string {
	return `compgen -W "$(cat ~/.ssh/known_hosts | cut -d ' ' -f1 | cut -d ',' -f1)" -- $last`
}

func ActionValues(values ...string) string {
	if len(strings.TrimSpace(strings.Join(values, ""))) == 0 {
		return ActionMessage("no values to complete")
	}

	vals := make([]string, len(values))
	for index, val := range Sanitize(values...) {
		// TODO escape special characters
		//vals[index] = strings.Replace(val, " ", `\ `, -1)
		vals[index] = val
	}
	return fmt.Sprintf(`compgen -W "%v" -- $last`, strings.Join(vals, ` `))
}

func ActionValuesDescribed(values ...string) string {
	// TODO verify length (description always exists)
	vals := make([]string, len(values))
	for index, val := range values {
		if index%2 == 0 {
			vals[index/2] = val
		}
	}
	return ActionValues(vals...)
}

func ActionMessage(msg string) string {
	return ActionValues("ERR", strings.Replace(Sanitize(msg)[0], " ", "_", -1)) // TODO escape characters
}

func ActionMultiParts(separator rune, values ...string) string {
	return ActionValues(values...)
}
