package cmd

import (
	"testing"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace/pkg/sandbox"
	"github.com/carapace-sh/carapace/pkg/style"
)

func TestCustomPrefixFlagCompletion(t *testing.T) {
	sandbox.Package(t, "github.com/carapace-sh/carapace/example-nonposix")(func(s *sandbox.Sandbox) {
		s.Run("customprefix", "&").
			Expect(carapace.Batch(
				carapace.ActionValuesDescribed(
					"&count", "count flag",
				).Tag("longhand flags"),
				carapace.ActionValuesDescribed(
					"&name", "name flag",
				).Style(style.Blue).Tag("longhand flags"),
				carapace.ActionValuesDescribed(
					"&verbose", "verbose output",
				).Tag("longhand flags"),
				carapace.ActionValuesDescribed(
					"&c", "count flag",
				).Tag("shorthand flags"),
				carapace.ActionValuesDescribed(
					"&n", "name flag",
				).Style(style.Blue).Tag("shorthand flags"),
				carapace.ActionValuesDescribed(
					"&v", "verbose output",
				).Tag("shorthand flags"),
			).ToA().
				NoSpace('.'))
	})
}

func TestCustomPrefixLonghand(t *testing.T) {
	sandbox.Package(t, "github.com/carapace-sh/carapace/example-nonposix")(func(s *sandbox.Sandbox) {
		s.Run("customprefix", "&na").
			Expect(carapace.ActionValuesDescribed(
				"&name", "name flag",
			).Style(style.Blue).Tag("longhand flags").
				NoSpace('.'))
	})
}

func TestCustomPrefixFlagArg(t *testing.T) {
	sandbox.Package(t, "github.com/carapace-sh/carapace/example-nonposix")(func(s *sandbox.Sandbox) {
		s.Run("customprefix", "&n", "").
			Expect(carapace.ActionValues("alpha", "beta", "gamma").
				Usage("name flag"))
	})
}

func TestCustomPrefixFlagArgLonghand(t *testing.T) {
	sandbox.Package(t, "github.com/carapace-sh/carapace/example-nonposix")(func(s *sandbox.Sandbox) {
		s.Run("customprefix", "&name", "").
			Expect(carapace.ActionValues("alpha", "beta", "gamma").
				Usage("name flag"))
	})
}

func TestCustomPrefixDelimited(t *testing.T) {
	sandbox.Package(t, "github.com/carapace-sh/carapace/example-nonposix")(func(s *sandbox.Sandbox) {
		s.Run("customprefix", "&n=").
			Expect(carapace.ActionValues("alpha", "beta", "gamma").
				Prefix("&n=").
				Usage("name flag"))
	})
}

func TestCustomPrefixDelimitedLonghand(t *testing.T) {
	sandbox.Package(t, "github.com/carapace-sh/carapace/example-nonposix")(func(s *sandbox.Sandbox) {
		s.Run("customprefix", "&name=").
			Expect(carapace.ActionValues("alpha", "beta", "gamma").
				Prefix("&name=").
				Usage("name flag"))
	})
}

func TestCustomPrefixPositional(t *testing.T) {
	sandbox.Package(t, "github.com/carapace-sh/carapace/example-nonposix")(func(s *sandbox.Sandbox) {
		s.Run("customprefix", "").
			Expect(carapace.ActionValues("pos1", "positional1"))
	})
}
