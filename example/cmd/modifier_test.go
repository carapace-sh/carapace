package cmd

import (
	"testing"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/sandbox"
	"github.com/rsteube/carapace/pkg/style"
)

func TestBatch(t *testing.T) {
	sandbox.Package(t, "github.com/rsteube/carapace/example")(func(s *sandbox.Sandbox) {
		s.Run("modifier", "--batch", "").
			Expect(carapace.ActionValuesDescribed(
				"A", "description of A",
				"B", "description of second B",
				"C", "description of second C",
				"D", "description of D",
			).
				Usage("Batch()"))
	})
}

func TestFilter(t *testing.T) {
	sandbox.Package(t, "github.com/rsteube/carapace/example")(func(s *sandbox.Sandbox) {
		s.Run("modifier", "--filter", "").
			Expect(carapace.ActionValuesDescribed(
				"1", "one",
				"3", "three",
			).Usage("Filter()"))
	})
}

func TestRetain(t *testing.T) {
	sandbox.Package(t, "github.com/rsteube/carapace/example")(func(s *sandbox.Sandbox) {
		s.Run("modifier", "--retain", "").
			Expect(carapace.ActionValuesDescribed(
				"2", "two",
				"4", "four",
			).Usage("Retain()"))
	})
}

func TestShift(t *testing.T) {
	sandbox.Package(t, "github.com/rsteube/carapace/example")(func(s *sandbox.Sandbox) {
		s.Run("modifier", "one", "--shift", "").
			Expect(carapace.ActionMessage(`[]string{}`).Usage("Shift()"))

		s.Run("modifier", "one", "two", "three", "--shift", "").
			Expect(carapace.ActionMessage(`[]string{"two", "three"}`).Usage("Shift()"))
	})
}

func TestTimeout(t *testing.T) {
	sandbox.Package(t, "github.com/rsteube/carapace/example")(func(s *sandbox.Sandbox) {
		s.Run("modifier", "--timeout", "1s:").
			Expect(carapace.ActionValues("within timeout").
				Prefix("1s:").
				NoSpace(':').
				Usage("Timeout()"))

		s.Run("modifier", "--timeout", "3s:").
			Expect(carapace.ActionMessage("timeout exceeded").
				Usage("Timeout()"))
	})
}

func TestUsage(t *testing.T) {
	sandbox.Package(t, "github.com/rsteube/carapace/example")(func(s *sandbox.Sandbox) {
		s.Run("modifier", "--usage", "").
			Expect(carapace.ActionValues("explicit", "implicit").
				Suffix(":").
				NoSpace(':').
				Usage("Usage()"))

		s.Run("modifier", "--usage", "explicit:").
			Expect(carapace.ActionValues().
				NoSpace(':').
				Usage("explicit usage"))
	})
}

func TestChdir(t *testing.T) {
	sandbox.Action(t, func() carapace.Action {
		return carapace.ActionFiles().Chdir("subdir")
	})(func(s *sandbox.Sandbox) {
		s.Files("subdir/file1.txt", "")

		s.Run("").Expect(
			carapace.ActionValues("file1.txt").
				StyleF(func(s string, sc style.Context) string {
					return style.ForPath("subdir/file1.txt", sc)
				}).
				NoSpace('/').
				Tag("files"))
	})
}

func TestMultiParts(t *testing.T) {
	sandbox.Package(t, "github.com/rsteube/carapace/example")(func(s *sandbox.Sandbox) {
		s.Run("modifier", "--multiparts", "").
			Expect(carapace.ActionValues("dir/").
				NoSpace('/').
				Usage("MultiParts()"))

		s.Run("modifier", "--multiparts", "dir/").
			Expect(carapace.ActionValues("subdir1/", "subdir2/").
				Prefix("dir/").
				NoSpace('/').
				Usage("MultiParts()"))

		s.Run("modifier", "--multiparts", "dir/subdir1/").
			Expect(carapace.ActionValues("fileA.txt", "fileB.txt").
				Prefix("dir/subdir1/").
				NoSpace('/').
				Usage("MultiParts()"))

		s.Run("modifier", "--multiparts", "dir/subdir2/").
			Expect(carapace.ActionValues("fileC.txt").
				Prefix("dir/subdir2/").
				NoSpace('/').
				Usage("MultiParts()"))
	})
}

func TestPrefix(t *testing.T) {
	sandbox.Package(t, "github.com/rsteube/carapace/example")(func(s *sandbox.Sandbox) {
		s.Files("subdir/file1.txt", "")

		s.Run("modifier", "--prefix", "").
			Expect(carapace.ActionValues("subdir/").
				StyleF(style.ForPath).
				Prefix("file://").
				NoSpace('/').
				Usage("Prefix()").
				Tag("files"))

		s.Run("modifier", "--prefix", "file").
			Expect(carapace.ActionValues("subdir/").
				StyleF(style.ForPath).
				Prefix("file://").
				NoSpace('/').
				Usage("Prefix()").
				Tag("files"))

		s.Run("modifier", "--prefix", "file://subdir/f").
			Expect(carapace.ActionValues("file1.txt").
				StyleF(style.ForPath).
				Prefix("file://subdir/").
				NoSpace('/').
				Usage("Prefix()").
				Tag("files"))
	})
}
