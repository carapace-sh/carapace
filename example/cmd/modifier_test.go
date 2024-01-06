package cmd

import (
	"os"
	"testing"
	"time"

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

func TestCache(t *testing.T) {
	sandbox.Package(t, "github.com/rsteube/carapace/example")(func(s *sandbox.Sandbox) {
		cached := s.Run("modifier", "--cache", "").Output()
		time.Sleep(1 * time.Second)
		s.Run("modifier", "--cache", "").
			Expect(cached.
				Usage("Cache()"))

		s.ClearCache()
		s.Run("modifier", "--cache", "").
			ExpectNot(cached.
				Usage("Cache()"))
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
		s.Run("modifier", "--timeout", "").
			Expect(carapace.ActionMessage("timeout exceeded").
				Usage("Timeout()"))
	})
}

func TestUsage(t *testing.T) {
	sandbox.Package(t, "github.com/rsteube/carapace/example")(func(s *sandbox.Sandbox) {
		s.Run("modifier", "--usage", "").
			Expect(carapace.ActionValues().
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

func TestChdirF(t *testing.T) {
	os.Unsetenv("LS_COLORS")
	sandbox.Package(t, "github.com/rsteube/carapace/example")(func(s *sandbox.Sandbox) {
		s.Files(
			".git/config", "",
			"file1.txt", "",
			"subdirA/file2.txt", "",
			"subdirB/file3.txt", "",
		)

		s.Run("modifier", "--chdirf", "").Expect(
			carapace.ActionValues(
				"subdirA/",
				"subdirB/",
				"file1.txt",
			).StyleF(style.ForPath).
				NoSpace('/').
				Tag("files").
				Usage("ChdirF()"))

		s.Env("GIT_WORK_TREE", "subdirB/") // TODO should also work for subdirB
		s.Run("modifier", "--chdirf", "").Expect(
			carapace.ActionValues(
				"file3.txt",
			).StyleF(style.ForPath).
				NoSpace('/').
				Tag("files").
				Usage("ChdirF()"))
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
	os.Unsetenv("LS_COLORS")
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

func TestFilterArgs(t *testing.T) {
	sandbox.Package(t, "github.com/rsteube/carapace/example")(func(s *sandbox.Sandbox) {
		s.Run("modifier", "--filterargs", "").
			Expect(carapace.ActionValues(
				"one",
				"two",
				"three",
			).Usage("FilterArgs()"))

		s.Run("modifier", "one", "--filterargs", "").
			Expect(carapace.ActionValues(
				"two",
				"three",
			).Usage("FilterArgs()"))

		s.Run("modifier", "one", "three", "--filterargs", "").
			Expect(carapace.ActionValues(
				"two",
			).Usage("FilterArgs()"))
	})
}

func TestFilterParts(t *testing.T) {
	sandbox.Package(t, "github.com/rsteube/carapace/example")(func(s *sandbox.Sandbox) {
		s.Run("modifier", "--filterparts", "").
			Expect(carapace.ActionValues(
				"one",
				"two",
				"three",
			).Suffix(",").
				NoSpace(',').
				Usage("FilterParts()"))

		s.Run("modifier", "--filterparts", "one,").
			Expect(carapace.ActionValues(
				"two",
				"three",
			).Suffix(",").
				NoSpace(',').
				Prefix("one,").
				Usage("FilterParts()"))

		s.Run("modifier", "--filterparts", "one,three,").
			Expect(carapace.ActionValues(
				"two",
			).Suffix(",").
				NoSpace(',').
				Prefix("one,three,").
				Usage("FilterParts()"))
	})
}

func TestMultiPartsP(t *testing.T) {
	sandbox.Package(t, "github.com/rsteube/carapace/example")(func(s *sandbox.Sandbox) {
		s.Run("modifier", "--multipartsp", "").
			Expect(carapace.ActionStyledValuesDescribed(
				"keys/", "", style.Default,
				"styles", "list", style.Yellow,
				"styles/", "", style.Default,
			).NoSpace('/').
				Usage("MultiPartsP()"))

		s.Run("modifier", "--multipartsp", "keys/").
			Expect(carapace.ActionValues(
				"key1",
				"key1/",
				"key2",
				"key2/",
			).Prefix("keys/").
				NoSpace('/').
				Usage("MultiPartsP()"))

		s.Run("modifier", "--multipartsp", "keys/key1/").
			Expect(carapace.ActionValues(
				"val1",
				"val2",
			).Prefix("keys/key1/").
				NoSpace('/').
				Usage("MultiPartsP()"))

		s.Run("modifier", "--multipartsp", "keys/key2/").
			Expect(carapace.ActionValues(
				"val3",
				"val4",
			).Prefix("keys/key2/").
				NoSpace('/').
				Usage("MultiPartsP()"))

		s.Run("modifier", "--multipartsp", "styles/c").
			Expect(carapace.Batch(
				carapace.ActionStyledValues(
					"color", style.Default,
					"cyan", style.Cyan,
				).Tag("styles"),
				carapace.ActionStyledValuesDescribed(
					"custom", "custom style", style.Of(style.Blue, style.Blink),
				),
			).ToA().
				Prefix("styles/").
				NoSpace('/', 'r').
				Usage("MultiPartsP()"))
	})
}

func TestSplit(t *testing.T) {
	os.Unsetenv("LS_COLORS")
	sandbox.Package(t, "github.com/rsteube/carapace/example")(func(s *sandbox.Sandbox) {
		s.Files("subdir/file1.txt", "")

		s.Run("modifier", "--split", "").
			Expect(carapace.ActionValues(
				"pos1",
				"positional1",
			).NoSpace('*').
				Suffix(" ").
				Usage("Split()"))

		s.Run("modifier", "--split", "pos1 ").
			Expect(carapace.ActionValues(
				"subdir/",
			).StyleF(style.ForPathExt).
				Prefix("pos1 ").
				NoSpace('*').
				Usage("Split()").
				Tag("files"))

		s.Run("modifier", "--split", "pos1 --").
			Expect(carapace.ActionStyledValuesDescribed(
				"--bool", "bool flag", style.Default,
				"--string", "string flag", style.Blue,
			).Prefix("pos1 ").
				Suffix(" ").
				NoSpace('*').
				Usage("Split()").
				Tag("flags"))

		s.Run("modifier", "--split", "pos1 --bool=").
			Expect(carapace.ActionStyledValues(
				"true", style.Green,
				"false", style.Red,
			).Prefix("pos1 --bool=").
				Suffix(" ").
				NoSpace('*').
				Usage("bool flag"))

		s.Run("modifier", "--split", "pos1 \"--bool=").
			Expect(carapace.ActionStyledValues(
				"true", style.Green,
				"false", style.Red,
			).Prefix("pos1 \"--bool=").
				Suffix("\" ").
				NoSpace('*').
				Usage("bool flag"))

		s.Run("modifier", "--split", "pos1 '--bool=").
			Expect(carapace.ActionStyledValues(
				"true", style.Green,
				"false", style.Red,
			).Prefix("pos1 '--bool=").
				Suffix("' ").
				NoSpace('*').
				Usage("bool flag"))

		t.Skip("skipping test that don't work yet") // TODO these need to work
		s.Run("modifier", "--split", "pos1 \"").
			Expect(carapace.ActionValues(
				"subdir/",
			).StyleF(style.ForPathExt).
				Prefix("pos1 \"").
				Suffix("\"").
				NoSpace('*').
				Usage("Split()").
				Tag("files"))

		s.Run("modifier", "--split", "pos1 '").
			Expect(carapace.ActionValues(
				"subdir/",
			).StyleF(style.ForPathExt).
				Prefix("pos1 '").
				Suffix("'").
				NoSpace('*').
				Usage("Split()").
				Tag("files"))
	})
}

func TestSplitP(t *testing.T) {
	os.Unsetenv("LS_COLORS")
	sandbox.Package(t, "github.com/rsteube/carapace/example")(func(s *sandbox.Sandbox) {
		s.Files("subdir/file1.txt", "")

		s.Run("modifier", "--splitp", "pos1>").
			Expect(carapace.ActionValues(
				"subdir/",
			).NoSpace('*').
				StyleF(style.ForPath).
				Prefix("pos1>").
				Tag("files").
				Usage("SplitP()"))

		s.Run("modifier", "--splitp", "pos1>subdir/").
			Expect(carapace.ActionValues(
				"file1.txt",
			).NoSpace('*').
				StyleF(style.ForPath).
				Prefix("pos1>subdir/").
				Suffix(" ").
				Tag("files").
				Usage("SplitP()"))

		s.Run("modifier", "--splitp", "pos1>subdir/file1.txt --b").
			Expect(carapace.ActionValuesDescribed(
				"--bool", "bool flag",
			).NoSpace('*').
				Prefix("pos1>subdir/file1.txt ").
				Suffix(" ").
				Tag("flags").
				Usage("SplitP()"))

		s.Run("modifier", "--splitp", "pos1 1>").
			Expect(carapace.ActionValues(
				"subdir/",
			).NoSpace('*').
				StyleF(style.ForPath).
				Prefix("pos1 1>").
				Tag("files").
				Usage("SplitP()"))

		s.Run("modifier", "--splitp", "<> subdir/file1.txt ").
			Expect(carapace.ActionValues(
				"pos1",
				"positional1",
			).NoSpace('*').
				Prefix("<> subdir/file1.txt ").
				Suffix(" ").
				Usage("SplitP()"))

		s.Run("modifier", "--splitp", "pos1|").
			Expect(carapace.ActionValues(
				"pos1",
				"positional1",
			).NoSpace('*').
				Prefix("pos1|").
				Suffix(" ").
				Usage("SplitP()"))

		s.Run("modifier", "--splitp", "pos1|&").
			Expect(carapace.ActionValues(
				"pos1",
				"positional1",
			).NoSpace('*').
				Prefix("pos1|&").
				Suffix(" ").
				Usage("SplitP()"))

		s.Run("modifier", "--splitp", "pos1 ;").
			Expect(carapace.ActionValues(
				"pos1",
				"positional1",
			).NoSpace('*').
				Prefix("pos1 ;").
				Suffix(" ").
				Usage("SplitP()"))

		s.Run("modifier", "--splitp", "pos1 | ").
			Expect(carapace.ActionValues(
				"pos1",
				"positional1",
			).NoSpace('*').
				Prefix("pos1 | ").
				Suffix(" ").
				Usage("SplitP()"))
	})
}
func TestUnless(t *testing.T) {
	sandbox.Package(t, "github.com/rsteube/carapace/example")(func(s *sandbox.Sandbox) {
		s.Run("modifier", "--unless", "").
			Expect(carapace.ActionValues(
				"./local",
				"~/home",
				"/abs",
				"one",
				"two",
				"three",
			).Usage("Unless()"))

		s.Run("modifier", "--unless", "t").
			Expect(carapace.ActionValues(
				"two",
				"three",
			).Usage("Unless()"))

		s.Run("modifier", "--unless", ".").
			Expect(carapace.ActionValues().Usage("Unless()"))

		s.Run("modifier", "--unless", "~").
			Expect(carapace.ActionValues().Usage("Unless()"))

		s.Run("modifier", "--unless", "/").
			Expect(carapace.ActionValues().Usage("Unless()"))
	})
}

func TestUniqueList(t *testing.T) {
	sandbox.Package(t, "github.com/rsteube/carapace/example")(func(s *sandbox.Sandbox) {
		s.Run("modifier", "--uniquelist", "").
			Expect(carapace.ActionValues(
				"one",
				"two",
				"three",
			).NoSpace().
				Usage("UniqueList()"))

		s.Run("modifier", "--uniquelist", "two,").
			Expect(carapace.ActionValues(
				"one",
				"three",
			).Prefix("two,").
				NoSpace().
				Usage("UniqueList()"))
	})
}

func TestUniqueListF(t *testing.T) {
	sandbox.Package(t, "github.com/rsteube/carapace/example")(func(s *sandbox.Sandbox) {
		s.Run("modifier", "--uniquelistf", "").
			Expect(carapace.ActionValues(
				"one",
				"two",
				"three",
			).NoSpace().
				Usage("UniqueListF()"))

		s.Run("modifier", "--uniquelistf", "two,").
			Expect(carapace.ActionValues(
				"one",
				"three",
			).Prefix("two,").
				NoSpace().
				Usage("UniqueListF()"))

		s.Run("modifier", "--uniquelistf", "two:").
			Expect(carapace.ActionValues(
				"1",
				"2",
				"3",
			).Prefix("two:").
				NoSpace().
				Usage("UniqueListF()"))

		s.Run("modifier", "--uniquelistf", "two:1,").
			Expect(carapace.ActionValues(
				"one",
				"three",
			).Prefix("two:1,").
				NoSpace().
				Usage("UniqueListF()"))
	})
}
