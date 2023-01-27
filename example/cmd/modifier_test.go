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
			Expect(carapace.ActionValues().
				Usage("explicit flag usage"))
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

func TestToMultiPartsA(t *testing.T) {
	sandbox.Package(t, "github.com/rsteube/carapace/example")(func(s *sandbox.Sandbox) {
		s.Run("modifier", "--tomultiparts", "").
			Expect(carapace.ActionValues("1=", "2=").
				NoSpace('/', '=').
				Usage("ToMultiPartsA()"))

		s.Run("modifier", "--tomultiparts", "1=").
			Expect(carapace.ActionValues("1==", "2==").
				Prefix("1=").
				NoSpace('/', '=').
				Usage("ToMultiPartsA()"))

		s.Run("modifier", "--tomultiparts", "1=1==").
			Expect(carapace.ActionValues("1/", "2/").
				Prefix("1=1==").
				NoSpace('/', '=').
				Usage("ToMultiPartsA()"))

		s.Run("modifier", "--tomultiparts", "1=1==1/").
			Expect(carapace.ActionValuesDescribed(
				"1", "one",
				"2", "two").
				Prefix("1=1==1/").
				NoSpace('/', '=').
				Usage("ToMultiPartsA()"))
	})
}
