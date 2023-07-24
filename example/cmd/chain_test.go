package cmd

import (
	"testing"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/sandbox"
	"github.com/rsteube/carapace/pkg/style"
)

func TestShorthandChain(t *testing.T) {
	sandbox.Package(t, "github.com/rsteube/carapace/example")(func(s *sandbox.Sandbox) {
		s.Run("chain", "-b").
			Expect(carapace.ActionStyledValues(
				"c", style.Default,
				"o", style.Yellow,
				"v", style.Blue,
			).Prefix("-b").
				NoSpace('c', 'o').
				Tag("flags"))

		s.Run("chain", "-bc").
			Expect(carapace.ActionStyledValues(
				"c", style.Default,
				"o", style.Yellow,
				"v", style.Blue,
			).Prefix("-bc").
				NoSpace('c', 'o').
				Tag("flags"))

		s.Run("chain", "-bcc").
			Expect(carapace.ActionStyledValues(
				"c", style.Default,
				"o", style.Yellow,
				"v", style.Blue,
			).Prefix("-bcc").
				NoSpace('c', 'o').
				Tag("flags"))

		s.Run("chain", "-bcco").
			Expect(carapace.ActionStyledValues(
				"c", style.Default,
				"v", style.Blue,
			).Prefix("-bcco").
				NoSpace('c').
				Tag("flags"))

		s.Run("chain", "-bcco", "").
			Expect(carapace.ActionValues(
				"p1",
				"positional1",
			))

		s.Run("chain", "-bcco=").
			Expect(carapace.ActionValues(
				"opt1",
				"opt2",
			).Prefix("-bcco="))

		s.Run("chain", "-bccv", "").
			Expect(carapace.ActionValues(
				"val1",
				"val2",
			))

		s.Run("chain", "-bccv=").
			Expect(carapace.ActionValues(
				"val1",
				"val2",
			).Prefix("-bccv="))

		s.Run("chain", "-bccv", "val1", "-c").
			Expect(carapace.ActionStyledValues(
				"c", style.Default,
				"o", style.Yellow,
			).Prefix("-c").
				NoSpace('c', 'o').
				Tag("flags"))
	})
}
