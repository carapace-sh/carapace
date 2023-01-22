package cmd

import (
	"testing"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/sandbox"
)

func TestRoot(t *testing.T) {
	sandbox.Package(t, "github.com/rsteube/carapace/example-nonposix")(func(s *sandbox.Sandbox) {
		s.Run("-delim-colon:").
			Expect(carapace.ActionValues("d1", "d2", "d3").
				Prefix("-delim-colon:").
				Usage("OptargDelimiter ':'"))

		s.Run("-delim-colon", "").
			Expect(carapace.ActionValues("p1", "positional1"))

		s.Run("-delim-slash/").
			Expect(carapace.ActionValues("d1", "d2", "d3").
				Prefix("-delim-slash/").
				Usage("OptargDelimiter '/'"))

		s.Run("-c").
			Expect(carapace.ActionValuesDescribed(
				"-c", "CountN",
				"-count", "CountN").
				NoSpace('.').
				Tag("flags"))
	})
}
