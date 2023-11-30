package cmd

import (
	"testing"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/sandbox"
	"github.com/rsteube/carapace/pkg/style"
)

func TestCompat(t *testing.T) {
	sandbox.Package(t, "github.com/rsteube/carapace/example")(func(s *sandbox.Sandbox) {
		s.Files(
			"subdir/file1.txt", "",
			"subdir/subdir2/file2.txt", "",
			"go.mod", "",
			"go.sum", "",
			"README.md", "",
		)

		s.Run("compat", "--error", "").
			Expect(carapace.ActionMessage("an error occurred").
				Usage("ShellCompDirectiveError"))

		s.Run("compat", "--nospace", "").
			Expect(carapace.ActionValues(
				"one",
				"two",
			).NoSpace().
				Usage("ShellCompDirectiveNoSpace"))

		s.Run("compat", "--nofilecomp", "").
			Expect(carapace.ActionValues().
				Usage("ShellCompDirectiveNoFileComp"))

		s.Run("compat", "--filterfileext", "").
			Expect(carapace.ActionValues(
				"subdir/",
				"go.mod",
				"go.sum",
			).NoSpace('/').
				Tag("files").
				StyleF(style.ForPath).
				Usage("ShellCompDirectiveFilterFileExt"))

		s.Run("compat", "--filterdirs", "").
			Expect(carapace.ActionValues(
				"subdir/",
			).NoSpace('/').
				Tag("directories").
				StyleF(style.ForPath).
				Usage("ShellCompDirectiveFilterDirs"))

		s.Run("compat", "--filterdirs-chdir", "").
			Expect(carapace.ActionValues(
				"subdir2/",
			).NoSpace('/').
				Tag("directories").
				StyleF(style.ForPathExt).
				Usage("ShellCompDirectiveFilterDirs"))

		s.Run("compat", "--keeporder", "").
			Expect(carapace.ActionValues(
				"one",
				"two",
				"three",
			).Usage("ShellCompDirectiveKeepOrder"))

		s.Run("compat", "--default", "").
			Expect(carapace.ActionValues(
				"subdir/",
				"go.mod",
				"go.sum",
				"README.md",
			).NoSpace('/').
				Tag("files").
				StyleF(style.ForPath).
				Usage("ShellCompDirectiveDefault"))

		s.Run("compat", "--unset", "").
			Expect(carapace.ActionValues().
				Usage("no completions defined"))

		s.Run("compat", "").
			Expect(carapace.Batch(
				carapace.ActionValues(
					"p1",
					"positional1",
				),
				carapace.ActionValues(
					"sub",
				).Tag("commands"),
			).ToA().Usage(""))

		s.Run("compat", "positional1", "").
			Expect(carapace.ActionValues(
				"subdir/",
				"go.mod",
				"go.sum",
				"README.md",
			).NoSpace('/').
				StyleF(style.ForPath).
				Tag("files").
				Usage(""))

		s.Run("compat", "positional1", "main.go", "").
			Expect(carapace.ActionValues(
				`args: []string{"positional1", "main.go"} toComplete: ""`,
				"alternative",
			).Usage(""))

		s.Run("compat", "positional1", "main.go", "a").
			Expect(carapace.ActionValues(
				`args: []string{"positional1", "main.go"} toComplete: "a"`,
				"alternative",
			).Usage(""))

		s.Run("compat", "--nospace", "one", "positional1", "--", "main.go", "a").
			Expect(carapace.ActionValues(
				`args: []string{"positional1", "main.go"} toComplete: "a"`,
				"alternative",
			).Usage(""))
	})
}
