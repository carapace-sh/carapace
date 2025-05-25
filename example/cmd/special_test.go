package cmd

import (
	"testing"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace/pkg/sandbox"
)

func TestMutuallyExclusive(t *testing.T) {
	sandbox.Package(t, "github.com/carapace-sh/carapace/example")(func(s *sandbox.Sandbox) {
		s.Run("special", "--e").
			Expect(carapace.ActionValuesDescribed(
				"--exclusive", "mutually exclusive flag",
				"--exclusiveRepeatable", "mutually exclusive repeatable flag",
			).NoSpace('.').
				Tag("longhand flags"))

		s.Run("special", "--exclusive", "--e").
			Expect(carapace.ActionValues().NoSpace('.'))

		s.Run("special", "--exclusiveRepeatable", "--e").
			Expect(carapace.ActionValuesDescribed(
				"--exclusiveRepeatable", "mutually exclusive repeatable flag",
			).NoSpace('.').
				Tag("longhand flags"))

	})
}
