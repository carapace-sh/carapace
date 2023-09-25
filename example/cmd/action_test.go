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
				Usage("ActionCallback()"))

		s.Run("action", "--cobra", "").
			Expect(carapace.ActionValues(
				"one",
				"two",
			).NoSpace().
				Usage("ActionCobra()"))

		s.Run("action", "--commands", "s").
			Expect(carapace.ActionValuesDescribed(
				"special", "",
				"subcommand", "subcommand example",
			).Suffix(" ").
				NoSpace().
				Tag("other commands").
				Usage("ActionCommands()"))

		s.Run("action", "--commands", "subcommand ").
			Expect(carapace.Batch(
				carapace.ActionValuesDescribed(
					"a1", "subcommand with alias",
					"a2", "subcommand with alias",
					"alias", "subcommand with alias",
				).Tag("other commands"),
				carapace.ActionValuesDescribed(
					"group", "subcommand with group",
				).Style(style.Blue).Tag("group commands"),
			).ToA().
				Prefix("subcommand ").
				Suffix(" ").
				NoSpace().
				Usage("ActionCommands()"))

		s.Run("action", "--commands", "subcommand unknown ").
			Expect(carapace.ActionMessage(`unknown subcommand "unknown" for "subcommand"`).NoSpace().
				Usage("ActionCommands()"))

		s.Run("action", "--commands", "subcommand hidden ").
			Expect(carapace.ActionValuesDescribed(
				"visible", "visible subcommand of a hidden command",
			).Prefix("subcommand hidden ").
				Suffix(" ").
				NoSpace().
				Tag("commands").
				Usage("ActionCommands()"))

		s.Run("action", "--values", "first", "--callback", "").
			Expect(carapace.ActionMessage("values flag is set to: 'first'").
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

		s.Run("action", "--execcommand", "").
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
				Usage("ActionMessage()"))

		s.Run("action", "--message-multiple", "t").
			Expect(carapace.Batch(
				carapace.ActionMessage("first message"),
				carapace.ActionMessage("second message"),
				carapace.ActionMessage("third message"),
				carapace.ActionValues("one", "two", "three")).
				ToA().
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

		s.Run("action", "embe").
			Expect(carapace.ActionValues("embeddedP1", "embeddedPositional1").
				Usage("action [pos1] [pos2] [--] [dashAny]..."))

		s.Run("action", "embeddedP1", "embeddedP2 ").
			Expect(carapace.ActionValues("embeddedP2 with space").
				Usage("action [pos1] [pos2] [--] [dashAny]..."))

		s.Run("action", "--unknown", "").
			Expect(carapace.ActionMessage("unknown flag: --unknown"))
	})
}

func TestDash(t *testing.T) {
	sandbox.Package(t, "github.com/rsteube/carapace/example")(func(s *sandbox.Sandbox) {
		s.Run("action", "--", "").
			Expect(carapace.ActionValues("embeddedP1", "embeddedPositional1").
				Usage("action [pos1] [pos2] [--] [dashAny]..."))

		s.Run("action", "--", "-").
			Expect(carapace.ActionStyledValuesDescribed(
				"--embedded-bool", "embedded bool flag", style.Default,
				"--embedded-optarg", "embedded optarg flag", style.Yellow,
				"--embedded-string", "embedded string flag", style.Blue,
				"-h", "help for embedded", style.Default,
				"--help", "help for embedded", style.Default).
				NoSpace('.').
				Usage("action [pos1] [pos2] [--] [dashAny]...").
				Tag("flags"))

		s.Run("action", "--", "--").
			Expect(carapace.ActionStyledValuesDescribed(
				"--embedded-bool", "embedded bool flag", style.Default,
				"--embedded-optarg", "embedded optarg flag", style.Yellow,
				"--embedded-string", "embedded string flag", style.Blue,
				"--help", "help for embedded", style.Default).
				NoSpace('.').
				Usage("action [pos1] [pos2] [--] [dashAny]...").
				Tag("flags"))

		s.Run("action", "--", "embeddedP1", "--embedded-optarg=").
			Expect(carapace.ActionValues("eo1", "eo2", "eo3").
				Prefix("--embedded-optarg=").
				Usage("embedded optarg flag"))

		s.Run("action", "--", "embeddedP1", "--embedded-string", "").
			Expect(carapace.ActionValues("es1", "es2", "es3").
				Usage("embedded string flag"))

		s.Run("action", "embeddedP1", "--styled-values", "second", "--", "--embedded-string", "es1", "").
			Expect(carapace.ActionValues("embeddedP2 with space", "embeddedPositional2 with space").
				Usage("action [pos1] [pos2] [--] [dashAny]..."))
	})
}

func TestUnknownFlag(t *testing.T) {
	sandbox.Package(t, "github.com/rsteube/carapace/example")(func(s *sandbox.Sandbox) {
		s.Run("action", "--unknown", "").
			Expect(carapace.ActionMessage("unknown flag: --unknown"))

		s.Env("CARAPACE_LENIENT", "1")
		s.Run("action", "--unknown", "").
			Expect(carapace.ActionValues("embeddedP1", "embeddedPositional1").
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

func TestAttached(t *testing.T) {
	sandbox.Package(t, "github.com/rsteube/carapace/example")(func(s *sandbox.Sandbox) {
		s.Files(
			"dirA/file1.txt", "",
			"dirA/file2.png", "",
			"dirB/dirC/file3.go", "",
			"dirB/file4.md", "",
			"file5.go", "",
		)

		s.Run("action", "--values=").
			Expect(carapace.ActionValues(
				"first",
				"second",
				"third",
			).Prefix("--values=").
				Usage("ActionValues()"))

		s.Run("action", "--values=f").
			Expect(carapace.ActionValues(
				"first",
			).Prefix("--values=").
				Usage("ActionValues()"))

		s.Run("action", "--values=first", "").
			Expect(carapace.ActionValues(
				"embeddedP1",
				"embeddedPositional1",
			).Usage("action [pos1] [pos2] [--] [dashAny]..."))

		s.Run("action", "--multiparts-nested=VALUE=").
			Expect(carapace.ActionValues("one", "two", "three").
				Prefix("--multiparts-nested=VALUE=").
				NoSpace().
				Usage("ActionMultiParts(...ActionMultiParts...)"))

		s.Run("action", "--multiparts-nested=VALUE=two,DIRECTORY=").
			Expect(carapace.ActionValues("dirA/", "dirB/").
				Tag("directories").
				StyleF(style.ForPath).
				Prefix("--multiparts-nested=VALUE=two,DIRECTORY=").
				NoSpace().
				Usage("ActionMultiParts(...ActionMultiParts...)"))
	})
}

func TestActionMultipartsN(t *testing.T) {
	sandbox.Package(t, "github.com/rsteube/carapace/example")(func(s *sandbox.Sandbox) {
		s.Run("action", "--multipartsn", "").
			Expect(carapace.ActionValues("one", "two").
				Suffix("=").
				NoSpace('=').
				Usage("ActionMultiPartsN()"))

		s.Run("action", "--multipartsn", "o").
			Expect(carapace.ActionValues("one").
				Suffix("=").
				NoSpace('=').
				Usage("ActionMultiPartsN()"))

		s.Run("action", "--multipartsn", "one=").
			Expect(carapace.ActionValues("three", "four").
				Prefix("one=").
				Suffix("=").
				NoSpace('=').
				Usage("ActionMultiPartsN()"))

		s.Run("action", "--multipartsn", "one=three=").
			Expect(carapace.ActionValues("five", "six").
				Prefix("one=three=").
				NoSpace('=').
				Usage("ActionMultiPartsN()"))
	})
}
