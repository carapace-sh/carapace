package cmd

import (
	"testing"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/sandbox"
	"github.com/rsteube/carapace/pkg/style"
)

func TestAction(t *testing.T) {
	sandbox.Package(t, "github.com/rsteube/carapace/example")(func(s *sandbox.Sandbox) {
		s.Files(
			"dirA/file1.txt", "",
			"dirA/file2.png", "",
			"dirB/dirC/file3.go", "",
			"dirB/file4.md", "",
			"file5.go", "",
		)

		s.Reply("git", "remote").With("origin\nfork")

		s.Run("action", "--callback", "").
			Expect(carapace.ActionMessage("values flag is not set").
				NoSpace().
				Usage("ActionCallback()"))

		s.Run("action", "--values", "first", "--callback", "").
			Expect(carapace.ActionMessage("values flag is set to: 'first'").
				NoSpace().
				Usage("ActionCallback()"))

		s.Run("action", "--directories", "").
			Expect(carapace.ActionValues("dirA/", "dirB/").
				Tag("directories").
				StyleF(style.ForPath).
				NoSpace('/').
				Usage("ActionDirectories()"))

		s.Run("action", "--directories", "dirB/").
			Expect(carapace.ActionValues("dirC/").
				Prefix("dirB/").
				Tag("directories").
				StyleF(style.ForPath).
				NoSpace('/').
				Usage("ActionDirectories()"))

		s.Run("action", "--exec-command", "").
			Expect(carapace.ActionValues("origin", "fork").
				Usage("ActionExecCommand()"))

		s.Run("action", "--files", "").
			Expect(carapace.ActionValues("dirA/", "dirB/", "file5.go").
				Tag("files").
				StyleF(style.ForPath).
				NoSpace('/').
				Usage("ActionFiles()"))

		s.Run("action", "--files-filtered", "").
			Expect(carapace.ActionValues("dirA/", "dirB/").
				Tag("files").
				StyleF(style.ForPath).
				NoSpace('/').
				Usage("ActionFiles(\".md\", \"go.mod\", \"go.sum\")"))

		s.Run("action", "--files-filtered", "dirB/").
			Expect(carapace.ActionValues("dirC/", "file4.md").
				Tag("files").
				Prefix("dirB/").
				StyleF(style.ForPath).
				NoSpace('/').
				Usage("ActionFiles(\".md\", \"go.mod\", \"go.sum\")"))

		s.Run("action", "--import", "").
			Expect(carapace.ActionValues("first", "second", "third").
				Usage("ActionImport()"))

		s.Run("action", "--import", "s").
			Expect(carapace.ActionValues("second").
				Usage("ActionImport()"))

		s.Run("action", "--message", "").
			Expect(carapace.ActionMessage("example message").
				NoSpace().
				Usage("ActionMessage()"))

		s.Run("action", "--message-multiple", "t").
			Expect(carapace.Batch(
				carapace.ActionMessage("first message"),
				carapace.ActionMessage("second message"),
				carapace.ActionMessage("third message"),
				carapace.ActionValues("one", "two", "three")).
				ToA().
				NoSpace().
				Usage("ActionMessage()"))

		s.Run("action", "--multiparts", "").
			Expect(carapace.ActionValues("userA", "userB").
				Suffix(":").
				NoSpace(':').
				Usage("ActionMultiParts()"))

		s.Run("action", "--multiparts", "userA:").
			Expect(carapace.ActionValues("groupA", "groupB").
				Prefix("userA:").
				NoSpace(':').
				Usage("ActionMultiParts()"))

		s.Run("action", "--multiparts-nested", "").
			Expect(carapace.ActionValues("DIRECTORY", "FILE", "VALUE").
				Suffix("=").
				NoSpace(',', '=').
				Usage("ActionMultiParts(...ActionMultiParts...)"))

		s.Run("action", "--multiparts-nested", "VALUE=").
			Expect(carapace.ActionValues("one", "two", "three").
				Prefix("VALUE=").
				NoSpace().
				Usage("ActionMultiParts(...ActionMultiParts...)"))

		s.Run("action", "--multiparts-nested", "VALUE=two,").
			Expect(carapace.ActionValues("DIRECTORY", "FILE").
				Prefix("VALUE=two,").
				Suffix("=").
				NoSpace(',', '=').
				Usage("ActionMultiParts(...ActionMultiParts...)"))

		s.Run("action", "--multiparts-nested", "VALUE=two,DIRECTORY=").
			Expect(carapace.ActionValues("dirA/", "dirB/").
				Tag("directories").
				StyleF(style.ForPath).
				Prefix("VALUE=two,DIRECTORY=").
				NoSpace().
				Usage("ActionMultiParts(...ActionMultiParts...)"))

		s.Run("action", "--styled-values", "s").
			Expect(carapace.ActionStyledValues("second", style.Blue).
				Usage("ActionStyledValues()"))

		s.Run("action", "--styled-values-described", "t").
			Expect(carapace.ActionStyledValuesDescribed(
				"third", "description of third", style.Of("#112233", style.Italic),
				"thirdalias", "description of third", style.BgBrightMagenta).
				Usage("ActionStyledValuesDescribed()"))

		s.Run("action", "--values", "sec").
			Expect(carapace.ActionValues("second").
				Usage("ActionValues()"))

		s.Run("action", "--values-described", "third").
			Expect(carapace.ActionValuesDescribed("third", "description of third").
				Usage("ActionValuesDescribed()"))

		s.Run("action", "pos").
			Expect(carapace.ActionValues("positional1", "positional1 with space").
				Usage("action [pos1] [pos2] [--] [dashAny]..."))

		s.Run("action", "p1", "positional2 ").
			Expect(carapace.ActionValues("positional2 with space").
				Usage("action [pos1] [pos2] [--] [dashAny]..."))

		s.Run("action", "--unknown", "").
			Expect(carapace.ActionMessage("unknown flag: --unknown"))
	})
}

func TestUnknownFlag(t *testing.T) {
	sandbox.Package(t, "github.com/rsteube/carapace/example")(func(s *sandbox.Sandbox) {
		s.Run("action", "--unknown", "").
			Expect(carapace.ActionMessage("unknown flag: --unknown").NoSpace())

		s.Env("CARAPACE_LENIENT", "1")
		s.Run("action", "--unknown", "").
			Expect(carapace.ActionValues("p1", "positional1", "positional1 with space").
				Usage("action [pos1] [pos2] [--] [dashAny]..."))
	})
}

func TestPersistentFlag(t *testing.T) {
	sandbox.Package(t, "github.com/rsteube/carapace/example")(func(s *sandbox.Sandbox) {
		s.Run("action", "--persistentFlag=").
			Expect(carapace.ActionValues("p1", "p2", "p3").
				Prefix("--persistentFlag=").
				Usage("Help message for persistentFlag"))

		s.Run("action", "--persistentFlag2", "").
			Expect(carapace.ActionValues("p4", "p5", "p6").
				Usage("Help message for persistentFlag2"))
	})
}
