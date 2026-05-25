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
			Expect(carapace.ActionValues())

		s.Run("argumentstyle", "-1").
			Expect(carapace.ActionValues())

		s.Run("argumentstyle", "-1", "").
			Expect(carapace.ActionValues("val1", "val2").
				Usage("Accept next arg only"))
	})
}

func TestAcceptDelimited(t *testing.T) {
	sandbox.Package(t, "github.com/carapace-sh/carapace/example-nonposix")(func(s *sandbox.Sandbox) {
		s.Run("argumentstyle", "--accept-delimited", "").
			Expect(carapace.ActionMessage("flag needs an argument: --accept-delimited"))

		s.Run("argumentstyle", "--accept-delimited=").
			Expect(carapace.ActionValues("val1", "val2").
				Prefix("--accept-delimited=").
				Usage("Accept delimited only"))

		s.Run("argumentstyle", "-2=").
			Expect(carapace.ActionValues("val1", "val2").
				Prefix("-2=").
				Usage("Accept delimited only"))

		s.Run("argumentstyle", "-2").
			Expect(carapace.ActionValues())

		s.Run("argumentstyle", "-2", "").
			Expect(carapace.ActionMessage("flag needs an argument: '2' in -2"))
	})
}

func TestAcceptAttached(t *testing.T) {
	sandbox.Package(t, "github.com/carapace-sh/carapace/example-nonposix")(func(s *sandbox.Sandbox) {
		s.Run("argumentstyle", "--accept-attached", "").
			Expect(carapace.ActionMessage("flag needs an argument: --accept-attached"))

		// s.Run("argumentstyle", "--accept-attached=").
		// 	Expect(carapace.ActionValues())

		s.Run("argumentstyle", "-3=").
			Expect(carapace.ActionValues())

		s.Run("argumentstyle", "-3").
			Expect(carapace.ActionValues("val1", "val2").
				Prefix("-3").
				Usage("Accept attached only"))

		s.Run("argumentstyle", "-3", "").
			Expect(carapace.ActionMessage("flag needs an argument: '3' in -3"))
	})
}
