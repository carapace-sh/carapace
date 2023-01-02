package cmd

import (
	"testing"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/sandbox"
	"github.com/rsteube/carapace/pkg/style"
)

func TestAction(t *testing.T) {
	sandbox.Run(t, "github.com/rsteube/carapace/example")(func(s *sandbox.Sandbox) {
		s.Files(
			"dirA/file1.txt", "",
			"dirA/file2.png", "",
			"dirB/dirC/file3.go", "",
			"dirB/file4.md", "",
			"file5.go", "",
		)

		s.Reply("git", "remote").With("origin\nfork")

		s.Run("action", "--callback", "").Expect(carapace.ActionMessage("values flag is not set"))
		s.Run("action", "--values", "first", "--callback", "").Expect(carapace.ActionMessage("values flag is set to: 'first'"))
		s.Run("action", "--directories", "").Expect(carapace.ActionValues("dirA/", "dirB/").Tag("directories").StyleF(style.ForPath))
		s.Run("action", "--directories", "dirB/").Expect(carapace.ActionValues("dirC/").Prefix("dirB/").Tag("directories").StyleF(style.ForPath))
		s.Run("action", "--exec-command", "").Expect(carapace.ActionValues("origin", "fork"))
		s.Run("action", "--files", "").Expect(carapace.ActionValues("dirA/", "dirB/", "file5.go").Tag("files").StyleF(style.ForPath))
		s.Run("action", "--files-filtered", "").Expect(carapace.ActionValues("dirA/", "dirB/").Tag("files").StyleF(style.ForPath))
		s.Run("action", "--files-filtered", "dirB/").Expect(carapace.ActionValues("dirC/", "file4.md").Tag("files").Prefix("dirB/").StyleF(style.ForPath))
		s.Run("action", "--import", "").Expect(carapace.ActionValues("first", "second", "third"))
		s.Run("action", "--import", "s").Expect(carapace.ActionValues("second"))
		s.Run("action", "--message", "").Expect(carapace.ActionMessage("example message"))
		s.Run("action", "--message-multiple", "t").Expect(carapace.Batch(carapace.ActionMessage("first message"), carapace.ActionMessage("second message"), carapace.ActionMessage("third message"), carapace.ActionValues("one", "two", "three")).ToA())
		s.Run("action", "--multiparts", "").Expect(carapace.ActionValues("userA", "userB").Suffix(":"))
		s.Run("action", "--multiparts", "userA:").Expect(carapace.ActionValues("groupA", "groupB").Prefix("userA:"))
		s.Run("action", "--multiparts-nested", "").Expect(carapace.ActionValues("DIRECTORY", "FILE", "VALUE").Suffix("="))
		s.Run("action", "--multiparts-nested", "VALUE=").Expect(carapace.ActionValues("one", "two", "three").Prefix("VALUE="))
		s.Run("action", "--multiparts-nested", "VALUE=two,").Expect(carapace.ActionValues("DIRECTORY", "FILE").Prefix("VALUE=two,").Suffix("="))
		s.Run("action", "--multiparts-nested", "VALUE=two,DIRECTORY=").Expect(carapace.ActionValues("dirA/", "dirB/").Tag("directories").StyleF(style.ForPath).Prefix("VALUE=two,DIRECTORY="))
		s.Run("action", "--styled-values", "s").Expect(carapace.ActionStyledValues("second", style.Blue))
		s.Run("action", "--styled-values-described", "t").Expect(carapace.ActionStyledValuesDescribed("third", "description of third", style.Of("#112233", style.Italic), "thirdalias", "description of third", style.BgBrightMagenta))
		s.Run("action", "--values", "sec").Expect(carapace.ActionValues("second"))
		s.Run("action", "--values-described", "third").Expect(carapace.ActionValuesDescribed("third", "description of third"))
		s.Run("action", "pos").Expect(carapace.ActionValues("positional1", "positional1 with space"))
		s.Run("action", "p1", "positional2 ").Expect(carapace.ActionValues("positional2 with space"))
	})
}
