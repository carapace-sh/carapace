package cmd_clink

import (
	"os"

	shlex "github.com/carapace-sh/carapace-shlex/v2"
)

func Patch(args []string) ([]string, error) {
	compline, ok := os.LookupEnv("CARAPACE_COMPLINE")
	if !ok {
		return args, nil
	}
	os.Unsetenv("CARAPACE_COMPLINE")

	if compline == "" {
		return args, nil
	}

	ctx := shlex.Complete(compline, shlex.Cmd)
	args = append(args[:1], ctx.Words...)
	return args, nil
}
