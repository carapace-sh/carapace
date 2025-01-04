package main

import (
	"testing"

	"github.com/carapace-sh/carapace/pkg/sandbox"
)

func TestVhs(t *testing.T) {
	sandbox.Package(t, "github.com/carapace-sh/carapace/example")(func(s *sandbox.Sandbox) {
		s.Record(``)
	})
}
