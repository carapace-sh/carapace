package cmd

import (
	"testing"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace/pkg/sandbox"
	"github.com/carapace-sh/carapace/pkg/style"
)

func TestHelp(t *testing.T) {
	sandbox.Package(t, "github.com/carapace-sh/carapace/example")(func(s *sandbox.Sandbox) {
		s.Run("help", "a").
			Expect(carapace.ActionValuesDescribed(
				"action", "action example",
				"alias", "action example",
			).Style(style.Blue).Tag("main commands").
				Usage("help [command]"))
	})
}
