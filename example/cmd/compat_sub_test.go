package cmd

import (
	"testing"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/sandbox"
)

func TestCompatPersistent(t *testing.T) {
	sandbox.Package(t, "github.com/rsteube/carapace/example")(func(s *sandbox.Sandbox) {
		s.Run("compat", "sub", "--persistent-compat", "").
			Expect(carapace.ActionValues(
				`args: []string(nil) toComplete: ""`,
				"alternative",
			).Usage("persistent flag defined with cobra"))

		s.Run("compat", "sub", "one", "--persistent-compat", "").
			Expect(carapace.ActionValues(
				`args: []string{"one"} toComplete: ""`,
				"alternative",
			).Usage("persistent flag defined with cobra"))

		s.Run("compat", "sub", "one", "two", "--persistent-compat", "a").
			Expect(carapace.ActionValues(
				`args: []string{"one", "two"} toComplete: "a"`,
				"alternative",
			).Usage("persistent flag defined with cobra"))
	})
}
