package cmd

import (
	"testing"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/sandbox"
	"github.com/rsteube/carapace/pkg/style"
)

func TestHelp(t *testing.T) {
	sandbox.Package(t, "github.com/rsteube/carapace/example")(func(s *sandbox.Sandbox) {
		s.Run("help", "a").
			Expect(carapace.ActionValuesDescribed(
				"action", "action example",
				"alias", "action example",
			).Style(style.Blue).Tag("main commands").
				Usage("help [command]"))
	})
}
