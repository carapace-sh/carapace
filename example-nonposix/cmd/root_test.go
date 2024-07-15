package cmd

import (
	"testing"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace/pkg/sandbox"
	"github.com/carapace-sh/carapace/pkg/style"
)

func TestStandalone(t *testing.T) {
	sandbox.Package(t, "github.com/carapace-sh/carapace/example-nonposix")(func(s *sandbox.Sandbox) {
		s.Run("--h").
			Expect(carapace.ActionValues().
				NoSpace('.'))

		s.Run("hel").
			Expect(carapace.ActionValues())
	})
}

func TestInterspersed(t *testing.T) {
	sandbox.Package(t, "github.com/carapace-sh/carapace/example-nonposix")(func(s *sandbox.Sandbox) {
		s.Run("-delim-colon:d1", "-d").
			Expect(carapace.ActionValuesDescribed(
				"-delim-slash", "OptargDelimiter '/'",
			).NoSpace('.').
				Style(style.Yellow).
				Tag("shorthand flags"))

		s.Run("-delim-colon:d1", "positional1", "-d").
			Expect(carapace.ActionValues())
	})
}

func TestRoot(t *testing.T) {
	sandbox.Package(t, "github.com/carapace-sh/carapace/example-nonposix")(func(s *sandbox.Sandbox) {
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
			Expect(carapace.Batch(
				carapace.ActionValuesDescribed(
					"-c", "CountN",
				).Tag("shorthand flags"),
				carapace.ActionValuesDescribed(
					"-count", "CountN",
				).Tag("longhand flags"),
			).ToA().
				NoSpace('.'))
	})
}

func TestNargs(t *testing.T) {
	sandbox.Package(t, "github.com/carapace-sh/carapace/example-nonposix")(func(s *sandbox.Sandbox) {
		s.Run("--nargs-any", "").
			Expect(carapace.ActionValues("na1", "na2", "na3").
				Usage("Nargs"))

		s.Run("--nargs-any", "na1", "").
			Expect(carapace.ActionValues("na2", "na3").
				Usage("Nargs"))

		s.Run("--nargs-any", "na2", "-c").
			Expect(carapace.Batch(
				carapace.ActionValuesDescribed(
					"-c", "CountN",
				).Tag("shorthand flags"),
				carapace.ActionValuesDescribed(
					"-count", "CountN",
				).Tag("longhand flags"),
			).ToA().
				NoSpace('.'))

		s.Run("--nargs-any", "na1", "na2", "").
			Expect(carapace.ActionValues("na3").
				Usage("Nargs"))

		s.Run("--nargs-two", "").
			Expect(carapace.ActionValues("nt1", "nt2", "nt3").
				Usage("Nargs"))

		s.Run("--nargs-two", "nt1", "").
			Expect(carapace.ActionValues("nt4", "nt5", "nt6").
				Usage("Nargs"))

		s.Run("--nargs-two", "nt1", "-").
			Expect(carapace.ActionValues().
				Usage("Nargs"))

		s.Run("--nargs-two", "nt1", "nt4", "").
			Expect(carapace.ActionValues("p1", "positional1"))

		s.Run("--nargs-two", "nt1", "nt4", "--nargs-").
			Expect(carapace.ActionValuesDescribed(
				"--nargs-any", "Nargs",
				"--nargs-two", "Nargs").
				Style(style.Magenta).
				NoSpace('.').
				Tag("longhand flags"))
	})
}
