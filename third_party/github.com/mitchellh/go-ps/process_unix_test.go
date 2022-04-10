//go:build linux || solaris
// +build linux solaris

package ps

import (
	"testing"
)

func TestUnixProcess_impl(t *testing.T) {
	var _ Process = new(UnixProcess)
}
