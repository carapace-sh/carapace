package cmd

import (
	"testing"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/sandbox"
	"github.com/rsteube/carapace/pkg/style"
)

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
	})
}
