package fish

import (
	"fmt"
	"strings"
)

func Callback(uid string) string {
	return ActionExecute(fmt.Sprintf(`_callback %v`, uid))
}

func ActionExecute(command string) string {
	return fmt.Sprintf(`%v`, command)
}

func ActionPathFiles(suffix string) string {
	return ActionExecute(fmt.Sprintf(`__fish_complete_suffix "%v"`, suffix))
}

func ActionFiles(suffix string) string {
	return ActionExecute(fmt.Sprintf(`__fish_complete_suffix "%v"`, suffix))
}

func ActionNetInterfaces() string {
	return ActionExecute("__fish_print_interfaces")
}

func ActionUsers() string {
	return ActionExecute("__fish_complete_users")
}

func ActionGroups() string {
	return ActionExecute("__fish_complete_groups")
}

func ActionHosts() string {
	return ActionExecute("__fish_print_hostnames")
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
	return ActionExecute(fmt.Sprintf(`echo -e %v`, strings.Join(vals, `\n`)))
}

func ActionValuesDescribed(values ...string) string {
	// TODO verify length (description always exists)
	vals := make([]string, len(values))
	for index, val := range values {
		if index%2 == 0 {
			vals[index/2] = fmt.Sprintf(`%v\t%v`, val, values[index+1])
		}
	}
	return ActionValues(vals...)
}

func ActionMessage(msg string) string {
	return ActionValuesDescribed("ERR", msg, "_", "")
}

func ActionMultiParts(separator rune, values ...string) string {
	return ActionValues(values...)
}
