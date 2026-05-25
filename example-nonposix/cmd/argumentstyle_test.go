package cmd

import (
	"testing"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace/pkg/sandbox"
)

func TestAcceptNext(t *testing.T) {
	sandbox.Package(t, "github.com/carapace-sh/carapace/example-nonposix")(func(s *sandbox.Sandbox) {
		s.Run("argumentstyle", "--accept-next", "").
			Expect(carapace.ActionValues("val1", "val2").
				Usage("Accept next arg only"))

		s.Run("argumentstyle", "--accept-next=").
			Expect(carapace.ActionValues().
				NoSpace('.'))

		s.Run("argumentstyle", "-1=").
			Expect(carapace.ActionValues().
				NoSpace('.'))

		s.Run("argumentstyle", "-1").
			Expect(carapace.ActionValues().
				NoSpace('.'))
	})
}
