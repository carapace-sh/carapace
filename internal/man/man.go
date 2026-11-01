package man

import (
	"os/exec"
	"regexp"
	"strings"
)

// Descriptions returns manpage descriptions for commands matching given prefix.
func Descriptions(s string) (descriptions map[string]string) {
	descriptions = make(map[string]string)
	if strings.HasPrefix(s, ".") || strings.HasPrefix(s, "~") || strings.HasPrefix(s, "/") {
		return
	}

	output, err := exec.Command("man", "--names-only", "-k", "^"+s).Output()
	if err != nil {
		return
	}

	r := regexp.MustCompile(`^(?P<name>[^ ]+) [^-]+- (?P<description>.*)$`)
	lines := strings.SplitSeq(string(output), "\n")
	for line := range lines {
		if matches := r.FindStringSubmatch(line); len(matches) > 2 {
			descriptions[matches[1]] = matches[2]
		}
	}
	return
}
