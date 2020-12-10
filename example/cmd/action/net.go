package action

import (
	"os/exec"
	"regexp"
	"strings"

	"github.com/rsteube/carapace"
)

func ActionNetInterfaces() carapace.Action {
	return carapace.ActionCallback(func(args []string) carapace.Action {
		if output, err := exec.Command("ifconfig").Output(); err != nil {
			return carapace.ActionMessage(err.Error())
		} else {
			interfaces := []string{}
			for _, line := range strings.Split(string(output), "\n") {
				if matches, _ := regexp.MatchString("^[0-9a-zA-Z]", line); matches {
					interfaces = append(interfaces, strings.Split(line, ":")[0])
				}
			}
			return carapace.ActionValues(interfaces...)
		}
	})
}
