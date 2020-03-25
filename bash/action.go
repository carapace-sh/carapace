package bash

import (
	"fmt"
	"strings"
)

func Callback(uid string) string {
	return fmt.Sprintf(`eval $(_callback '%v')`, uid) // TODO update and use ActionExecute for eval?
}

func ActionExecute(command string) string {
	return fmt.Sprintf(`$(%v)`, command)
}

func ActionPathFiles(suffix string) string {
	return ""
}

func ActionFiles(suffix string) string {
	return fmt.Sprintf(`compgen -f -o plusdirs -X "!*%v" -- $last`, suffix)
}

func ActionNetInterfaces() string {
	return ActionValues(ActionExecute(`ifconfig -a | grep -o '^[^ :]\+' | tr '\n' ' '`))
}

func ActionUsers() string {
	return `compgen -u -- $last`
}

func ActionGroups() string {
	return `compgen -g -- $last`
}

func ActionHosts() string {
	return ActionValues(ActionExecute(`cat ~/.ssh/known_hosts | cut -d ' ' -f1 | cut -d ',' -f1`))
}

func ActionOptions() string {
	return ""
}

func ActionValues(values ...string) string {
	if len(strings.TrimSpace(strings.Join(values, ""))) == 0 {
		return ActionMessage("no values to complete")
	}

	vals := make([]string, len(values))
	for index, val := range values {
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
	return ActionValues("ERR", strings.Replace(msg, " ", "_", -1)) // TODO escape characters
}

func ActionMultiParts(separator rune, values ...string) string {
	return ActionValues(values...)
}
