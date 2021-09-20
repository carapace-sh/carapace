package net

import (
	exec "golang.org/x/sys/execabs"
	"regexp"
	"strings"

	"github.com/rsteube/carapace"
)

// ActionNetInterfaces completes net interfaces
func ActionNetInterfaces() carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) (result carapace.Action) {
		if output, err := exec.Command("ifconfig").Output(); err != nil {
			result = carapace.ActionMessage(err.Error())
		} else {
			interfaces := []string{}
			r := regexp.MustCompile("^[0-9a-zA-Z]")
			for _, line := range strings.Split(string(output), "\n") {
				if r.MatchString(line) {
					interfaces = append(interfaces, strings.Split(line, ":")[0])
				}
			}
			result = carapace.ActionValues(interfaces...)
		}
		return
	})
}
