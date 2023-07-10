package cmd

import (
	"testing"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/sandbox"
	"github.com/rsteube/carapace/pkg/style"
)

func TestInterspersed(t *testing.T) {
	sandbox.Package(t, "github.com/rsteube/carapace/example")(func(s *sandbox.Sandbox) {
		s.Run("interspersed", "--s").
			Expect(carapace.ActionValuesDescribed(
				"--string", "string flag",
			).
				StyleR(&style.Carapace.FlagArg).
				NoSpace('.').
				Tag("flags"))

		s.Run("interspersed", "--bool", "--s").
			Expect(carapace.ActionValuesDescribed(
				"--string", "string flag",
			).
				StyleR(&style.Carapace.FlagArg).
				NoSpace('.').
				Tag("flags"))

		s.Run("interspersed", "--bool", "").
			Expect(carapace.ActionValues(
				"p1", "positional1",
			))

		s.Run("interspersed", "--bool", "p1", "-").
			Expect(carapace.ActionValues())
	})
}
