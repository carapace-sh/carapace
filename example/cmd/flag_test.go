package cmd

import (
	"testing"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace/pkg/sandbox"
)

func TestFlag(t *testing.T) {
	sandbox.Package(t, "github.com/carapace-sh/carapace/example")(func(s *sandbox.Sandbox) {
		// TODO test flags

		s.Run("flag", "--", "").
			Expect(carapace.ActionValues(
				"d1",
				"dash1",
				"-d1",
				"--dash1",
			))

		s.Run("flag", "--", "-").
			Expect(carapace.ActionValues(
				"-d1",
				"--dash1",
			))

		s.Run("flag", "--", "--").
			Expect(carapace.ActionValues(
				"--dash1",
			))

		s.Run("flag", "--", "d1", "").
			Expect(carapace.ActionValues(
				"d2",
				"dash2",
				"-d2",
				"--dash2",
			))

		s.Run("flag", "--", "d1", "-").
			Expect(carapace.ActionValues(
				"-d2",
				"--dash2",
			))

		s.Run("flag", "--", "d1", "--").
			Expect(carapace.ActionValues(
				"--dash2",
			))

		s.Run("flag", "--", "d1", "d2", "").
			Expect(carapace.ActionValues(
				"dAny",
				"dashAny",
				"-dAny",
				"--dashAny",
			))

		s.Run("flag", "--", "d1", "d2", "-").
			Expect(carapace.ActionValues(
				"-dAny",
				"--dashAny",
			))

		s.Run("flag", "--", "d1", "d2", "--").
			Expect(carapace.ActionValues(
				"--dashAny",
			))
	})
}
