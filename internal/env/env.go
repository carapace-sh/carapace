package env

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/rsteube/carapace/internal/common"
)

func ColorDisabled() bool {
	return os.Getenv("NO_COLOR") != "" || os.Getenv("CLICOLOR") == "0"
}

func Lenient() bool {
	return os.Getenv("CARAPACE_LENIENT") != ""
}

func Hashdirs() string {
	return os.Getenv("CARAPACE_ZSH_HASH_DIRS")
}

func Sandbox() (m *common.Mock, err error) {
	sandbox := os.Getenv("CARAPACE_SANDBOX")
	if sandbox == "" || !isGoRun() {
		return nil, errors.New("no sandbox")
	}

	err = json.Unmarshal([]byte(sandbox), &m)
	return
}

func Log() bool {
	return os.Getenv("CARAPACE_LOG") != ""
}

func Hidden() bool {
	return os.Getenv("CARAPACE_HIDDEN") != ""
}

func CoverDir() string {
	return os.Getenv("CARAPACE_COVERDIR") // custom env for GOCOVERDIR so that it works together with `-coverprofile`
}

func isGoRun() bool { return strings.HasPrefix(os.Args[0], os.TempDir()+"/go-build") }
