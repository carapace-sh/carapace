package cmd

import (
	"testing"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/sandbox"
)

func TestModifier(t *testing.T) {
	sandbox.Package(t, "github.com/rsteube/carapace/example")(func(s *sandbox.Sandbox) {
		s.Run("modifier", "--timeout", "1s:").
			Expect(carapace.ActionValues("within timeout").
				Prefix("1s:").
				NoSpace(':').
				Usage("Timeout()"))

		s.Run("modifier", "--timeout", "3s:").
			Expect(carapace.ActionMessage("timeout exceeded").
				Usage("Timeout()"))
	})
}
