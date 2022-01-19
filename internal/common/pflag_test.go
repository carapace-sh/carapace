package common

import (
	"testing"

	"github.com/spf13/pflag"
)

func TestIsShorthandOnly(t *testing.T) {
	if IsShorthandOnly(&pflag.Flag{}) {
		t.Error("pflag error")
	}
}
