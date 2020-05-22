package fish

import (
	"fmt"
	"strings"
)

var sanitizer = strings.NewReplacer(
	`$`, ``,
	"`", ``,
	`\`, ``,
	`"`, `'`,
	`(`, `[`,
	`)`, `]`,
)

func Sanitize(values ...string) []string {
	sanitized := make([]string, len(values))
	for index, value := range values {
		sanitized[index] = sanitizer.Replace(value)
	}
	return sanitized
}

func Callback(prefix string, uid string) string {
	return ActionExecute(fmt.Sprintf(`_%v_callback %v`, prefix, uid))
}

func ActionExecute(command string) string {
	return fmt.Sprintf(`%v`, command)
}

func ActionDirectories() string {
	return `__fish_complete_directories`
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

func ActionValues(values ...string) string {
	sanitized := Sanitize(values...)
	if len(strings.TrimSpace(strings.Join(sanitized, ""))) == 0 {
		return ActionMessage("no values to complete")
	}

	vals := make([]string, len(sanitized))
	for index, val := range sanitized {
		// TODO escape special characters
		//vals[index] = strings.Replace(val, " ", `\ `, -1)
		vals[index] = val
	}
	return ActionExecute(fmt.Sprintf(`echo -e "%v"`, strings.Join(vals, `\n`)))
}

func ActionValuesDescribed(values ...string) string {
	sanitized := Sanitize(values...)
	// TODO verify length (description always exists)
	vals := make([]string, len(sanitized))
	for index, val := range sanitized {
		if index%2 == 0 {
			vals[index/2] = fmt.Sprintf(`%v\t%v`, val, values[index+1])
		}
	}
	return ActionExecute(fmt.Sprintf(`echo -e "%v"`, strings.Join(vals, `\n`)))
}

func ActionMessage(msg string) string {
	return ActionExecute(fmt.Sprintf(`echo -e "ERR\t%v\n_"`, Sanitize(msg)[0]))
}

func ActionPrefixValues(prefix string, values ...string) string {
	sanitized := Sanitize(values...)
	if len(strings.TrimSpace(strings.Join(sanitized, ""))) == 0 {
		return ActionMessage("no values to complete")
	}

	vals := make([]string, len(sanitized))
	for index, val := range sanitized {
		// TODO escape special characters
		//vals[index] = strings.Replace(val, " ", `\ `, -1)
		vals[index] = prefix + val
	}
	return ActionExecute(fmt.Sprintf(`echo -e "%v"`, strings.Join(vals, `\n`)))
}
