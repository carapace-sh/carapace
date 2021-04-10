package action

import (
	"os/exec"
	"regexp"
	"strings"

	"github.com/rsteube/carapace"
)

func ActionNetInterfaces() carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		if output, err := exec.Command("ifconfig").Output(); err != nil {
			return carapace.ActionMessage(err.Error())
		} else {
			interfaces := []string{}
			r := regexp.MustCompile("^[0-9a-zA-Z]")
			for _, line := range strings.Split(string(output), "\n") {
				if r.MatchString(line) {
					interfaces = append(interfaces, strings.Split(line, ":")[0])
				}
			}
			return carapace.ActionValues(interfaces...)
		}
	})
}
