package cmd

import (
	"testing"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/sandbox"
	"github.com/rsteube/carapace/pkg/style"
)

// TODO rename.
func TestMultiparts(t *testing.T) {
	sandbox.Package(t, "github.com/rsteube/carapace/example")(func(s *sandbox.Sandbox) {
		s.Files(
			"dirA/file1.txt", "",
			"dirA/file2.png", "",
			"dirB/dirC/file3.go", "",
			"dirB/file4.md", "",
			"file5.go", "",
		)
		s.Run("multiparts", "").
			Expect(carapace.ActionValues("DIRECTORY", "FILE", "VALUE").
				Suffix("=").
				NoSpace(',', '='))

		s.Run("multiparts", "D").
			Expect(carapace.ActionValues("DIRECTORY").
				Suffix("=").
				NoSpace(',', '='))

		s.Run("multiparts", "DIRECTORY").
			Expect(carapace.ActionValues("DIRECTORY").
				Suffix("=").
				NoSpace(',', '='))

		s.Run("multiparts", "DIRECTORY=").
			Expect(carapace.ActionValues("dirA/", "dirB/").
				Tag("directories").
				StyleF(style.ForPath).
				Prefix("DIRECTORY=").
				NoSpace(',', '/', '='))

		s.Run("multiparts", "VALUE=").
			Expect(carapace.ActionValues("one", "two", "three").
				Prefix("VALUE=").
				NoSpace(',', '='))

		s.Run("multiparts", "VALUE=o").
			Expect(carapace.ActionValues("one").
				Prefix("VALUE=").
				NoSpace(',', '='))

		s.Run("multiparts", "VALUE=one,").
			Expect(carapace.ActionValues("DIRECTORY", "FILE").
				Prefix("VALUE=one,").
				Suffix("=").
				NoSpace(',', '='))

		s.Run("multiparts", "VALUE=one,F").
			Expect(carapace.ActionValues("FILE").
				Prefix("VALUE=one,").
				Suffix("=").
				NoSpace(',', '='))

		s.Run("multiparts", "VALUE=one,FILE=").
			Expect(carapace.ActionValues("dirA/", "dirB/", "file5.go").
				Tag("files").
				StyleF(style.ForPath).
				Prefix("VALUE=one,FILE=").
				NoSpace(',', '/', '='))

		s.Run("multiparts", "VALUE=one,FILE=dirB/").
			Expect(carapace.ActionValues("dirC/", "file4.md").
				Tag("files").
				Prefix("dirB/").
				StyleF(style.ForPath).
				Prefix("VALUE=one,FILE=").
				NoSpace(',', '/', '='))

		s.Run("multiparts", "--none-zero", "").
			Expect(carapace.ActionMessage("invalid value for n [ActionValuesDescribed]: 0").
				Usage("multiparts without divider limited to 0"))

		s.Run("multiparts", "--none-one", "").
			Expect(carapace.ActionValues("a", "b").
				Usage("multiparts without divider limited to 1"))

		s.Run("multiparts", "--none-one", "a").
			Expect(carapace.ActionValues("a").
				Usage("multiparts without divider limited to 1"))

		s.Run("multiparts", "--none-two", "").
			Expect(carapace.ActionValuesDescribed(
				"a", "zero",
				"b", "zero",
			).
				NoSpace().
				Style(style.Blue).
				Usage("multiparts without divider limited to 2"))

		s.Run("multiparts", "--none-two", "a").
			Expect(carapace.ActionValuesDescribed(
				"a", "default",
				"b", "default",
				"c", "default",
			).
				Prefix("a").
				NoSpace().
				Style(style.Red).
				Usage("multiparts without divider limited to 2"))

		s.Run("multiparts", "--none-two", "ab").
			Expect(carapace.ActionValuesDescribed(
				"a", "default",
				"c", "default",
			).
				Prefix("ab").
				NoSpace().
				Style(style.Red).
				Usage("multiparts without divider limited to 2"))

		s.Run("multiparts", "--none-two", "abc").
			Expect(carapace.ActionValuesDescribed(
				"a", "default",
			).
				Prefix("abc").
				NoSpace().
				Style(style.Red).
				Usage("multiparts without divider limited to 2"))

		s.Run("multiparts", "--none-three", "a").
			Expect(carapace.ActionValuesDescribed(
				"a", "one",
				"b", "one",
				"c", "one",
			).
				Prefix("a").
				NoSpace().
				Style(style.Red).
				Usage("multiparts without divider limited to 3"))

		s.Run("multiparts", "--none-three", "ab").
			Expect(carapace.ActionValuesDescribed(
				"a", "default",
				"b", "default",
				"c", "default",
				"d", "default",
			).
				Prefix("ab").
				NoSpace().
				Style(style.Green).
				Usage("multiparts without divider limited to 3"))

		s.Run("multiparts", "--none-three", "abc").
			Expect(carapace.ActionValuesDescribed(
				"a", "default",
				"b", "default",
				"d", "default",
			).
				Prefix("abc").
				NoSpace().
				Style(style.Green).
				Usage("multiparts without divider limited to 3"))

		s.Run("multiparts", "--none-three", "abcd").
			Expect(carapace.ActionValuesDescribed(
				"a", "default",
				"b", "default",
			).
				Prefix("abcd").
				NoSpace().
				Style(style.Green).
				Usage("multiparts without divider limited to 3"))
	})
}
