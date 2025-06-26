package cmd

import (
	"testing"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace/pkg/sandbox"
)

func TestSubcommand(t *testing.T) {
	sandbox.Package(t, "github.com/carapace-sh/carapace/example")(func(s *sandbox.Sandbox) {
		s.Run("subcommand", "").
			Expect(carapace.Batch(
				carapace.ActionValuesDescribed(
					"a1", "subcommand with alias",
					"a2", "subcommand with alias",
					"alias", "subcommand with alias",
				).Tag("other commands"),
				carapace.ActionValuesDescribed(
					"group", "subcommand with group",
				).Tag("group commands").Style("blue"),
			).ToA())

		s.Env("CARAPACE_HIDDEN", "1")
		s.Run("subcommand", "").
			Expect(carapace.Batch(
				carapace.ActionValuesDescribed(
					"a1", "subcommand with alias",
					"a2", "subcommand with alias",
					"alias", "subcommand with alias",
					"hidden", "hidden subcommand",
				).Tag("other commands"),
				carapace.ActionValuesDescribed(
					"group", "subcommand with group",
				).Tag("group commands").Style("blue"),
			).ToA())

		s.Env("CARAPACE_HIDDEN", "2")
		s.Run("subcommand", "").
			Expect(carapace.Batch(
				carapace.ActionValuesDescribed(
					"a1", "subcommand with alias",
					"a2", "subcommand with alias",
					"alias", "subcommand with alias",
					"hidden", "hidden subcommand",
					"_carapace", "",
				).Tag("other commands"),
				carapace.ActionValuesDescribed(
					"group", "subcommand with group",
				).Tag("group commands").Style("blue"),
			).ToA())
	})
}
