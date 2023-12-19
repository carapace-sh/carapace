package cmd

import (
	"testing"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/sandbox"
)

func TestFlagDisabled(t *testing.T) {
	sandbox.Package(t, "github.com/rsteube/carapace/example")(func(s *sandbox.Sandbox) {
		s.Run("flag", "disabled", "").
			Expect(carapace.ActionValues(
				"-p1",
				"positional1",
			))

		s.Run("flag", "disabled", "-p1", "").
			Expect(carapace.ActionValues(
				"p2",
				"--positional2",
			))

		s.Run("flag", "disabled", "-p1", "p2", "").
			Expect(carapace.ActionValues())
	})
}
