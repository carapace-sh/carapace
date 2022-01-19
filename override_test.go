package carapace

import (
	"os"
	"strings"
	"testing"
)

func TestOverride(t *testing.T) {
	os.Args = []string{"command", "--optarg:value", "-long"}
	Override(Opts{
		LongShorthand:   true,
		OptArgDelimiter: ":",
	})
	if strings.Join(os.Args, "|") != "command|--optarg=value|--long" {
		t.Error("override failed")
	}
}
