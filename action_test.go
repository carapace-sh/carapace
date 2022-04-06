package carapace

import (
	"fmt"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/rsteube/carapace/internal/assert"
	"github.com/rsteube/carapace/internal/common"
	"github.com/rsteube/carapace/pkg/style"
)

func init() {
	os.Unsetenv("LS_COLORS")
}

func assertEqual(t *testing.T, expected, actual InvokedAction) {
	sort.Sort(common.ByValue(expected.rawValues))
	sort.Sort(common.ByValue(actual.rawValues))

	assert.Equal(t, fmt.Sprintf("%+v\n", expected), fmt.Sprintf("%+v\n", actual))
}

func assertNotEqual(t *testing.T, expected, actual InvokedAction) {
	sort.Sort(common.ByValue(expected.rawValues))
	sort.Sort(common.ByValue(actual.rawValues))

	if fmt.Sprintf("%+v\n", expected) == fmt.Sprintf("%+v\n", actual) {
		t.Error("should be different")
	}
}

func TestActionCallback(t *testing.T) {
	a := ActionCallback(func(c Context) Action {
		return ActionCallback(func(c Context) Action {
			return ActionCallback(func(c Context) Action {
				return ActionValues("a", "b", "c")
			})
		})
	})
	expected := InvokedAction{
		Action{
			rawValues: common.RawValuesFrom("a", "b", "c"),
			nospace:   false,
			skipcache: false,
		},
	}
	actual := a.Invoke(Context{})
	assertEqual(t, expected, actual)
}

func TestCache(t *testing.T) {
	f := func() Action {
		return ActionCallback(func(c Context) Action {
			return ActionValues(time.Now().String())
		}).Cache(10 * time.Millisecond)
	}

	a1 := f().Invoke(Context{})
	a2 := f().Invoke(Context{})
	assertEqual(t, a1, a2)

	time.Sleep(12 * time.Millisecond)
	a3 := f().Invoke(Context{})
	assertNotEqual(t, a1, a3)
}

func TestSkipCache(t *testing.T) {
	a := ActionCallback(func(c Context) Action {
		return ActionValues().Invoke(c).Merge(
			ActionCallback(func(c Context) Action {
				return ActionMessage("skipcache")
			}).Invoke(c)).
			Filter([]string{""}).
			Prefix("").
			Suffix("").
			ToA()
	})
	if a.skipcache {
		t.Fatal("uninvoked skipcache should be false")
	}
	if !a.Invoke(Context{}).skipcache {
		t.Fatal("invoked skipcache should be true")
	}
}

func TestNoSpace(t *testing.T) {
	a := ActionCallback(func(c Context) Action {
		return ActionValues().Invoke(c).Merge(
			ActionMultiParts("", func(c Context) Action {
				return ActionMessage("nospace")
			}).Invoke(c)).
			Filter([]string{""}).
			Prefix("").
			Suffix("").
			ToA()
	})
	if a.nospace {
		t.Fatal("uninvoked nospace should be false")
	}
	if !a.Invoke(Context{}).nospace {
		t.Fatal("invoked nospace should be true")
	}
}

func TestActionDirectories(t *testing.T) {
	assertEqual(t,
		ActionStyledValues(
			"example/", style.Of("fg-default", "bg-default", style.Blue, style.Bold),
			"docs/", style.Of("fg-default", "bg-default", style.Blue, style.Bold),
			"internal/", style.Of("fg-default", "bg-default", style.Blue, style.Bold),
			"pkg/", style.Of("fg-default", "bg-default", style.Blue, style.Bold),
			"third_party/", style.Of("fg-default", "bg-default", style.Blue, style.Bold),
		).noSpace(true).Invoke(Context{}),
		ActionDirectories().Invoke(Context{CallbackValue: ""}).Filter([]string{"vendor/"}),
	)

	assertEqual(t,
		ActionStyledValues(
			"example/", style.Of("fg-default", "bg-default", style.Blue, style.Bold),
			"docs/", style.Of("fg-default", "bg-default", style.Blue, style.Bold),
			"internal/", style.Of("fg-default", "bg-default", style.Blue, style.Bold),
			"pkg/", style.Of("fg-default", "bg-default", style.Blue, style.Bold),
			"third_party/", style.Of("fg-default", "bg-default", style.Blue, style.Bold),
		).noSpace(true).Invoke(Context{}).Prefix("./"),
		ActionDirectories().Invoke(Context{CallbackValue: "./"}).Filter([]string{"./vendor/"}),
	)

	assertEqual(t,
		ActionStyledValues(
			"_test/", style.Of("fg-default", "bg-default", style.Blue, style.Bold),
			"cmd/", style.Of("fg-default", "bg-default", style.Blue, style.Bold),
		).noSpace(true).Invoke(Context{}).Prefix("example/"),
		ActionDirectories().Invoke(Context{CallbackValue: "example/"}),
	)

	assertEqual(t,
		ActionStyledValues(
			"_test/", style.Of("fg-default", "bg-default", style.Blue, style.Bold),
			"cmd/", style.Of("fg-default", "bg-default", style.Blue, style.Bold),
		).noSpace(true).Invoke(Context{}).Prefix("example/"),
		ActionDirectories().Invoke(Context{CallbackValue: "example/cm"}),
	)
}

func TestActionFiles(t *testing.T) {
	assertEqual(t,
		ActionStyledValues(
			"README.md", style.Of("fg-default", "bg-default"),
			"example/", style.Of("fg-default", "bg-default", style.Blue, style.Bold),
			"docs/", style.Of("fg-default", "bg-default", style.Blue, style.Bold),
			"internal/", style.Of("fg-default", "bg-default", style.Blue, style.Bold),
			"pkg/", style.Of("fg-default", "bg-default", style.Blue, style.Bold),
			"third_party/", style.Of("fg-default", "bg-default", style.Blue, style.Bold),
		).noSpace(true).Invoke(Context{}),
		ActionFiles(".md").Invoke(Context{CallbackValue: ""}).Filter([]string{"vendor/"}),
	)

	assertEqual(t,
		ActionStyledValues(
			"_test/", style.Of("fg-default", "bg-default", style.Blue, style.Bold),
			"cmd/", style.Of("fg-default", "bg-default", style.Blue, style.Bold),
			"main.go", style.Of("fg-default", "bg-default"),
			"main_test.go", style.Of("fg-default", "bg-default"),
		).noSpace(true).Invoke(Context{}).Prefix("example/"),
		ActionFiles().Invoke(Context{CallbackValue: "example/"}).Filter([]string{"example/example"}),
	)
}

func TestActionFilesChdir(t *testing.T) {
	oldWd, _ := os.Getwd()

	assertEqual(t,
		ActionValuesDescribed("ERR", "stat nonexistent: no such file or directory", "_", "").noSpace(true).skipCache(true).Invoke(Context{}),
		ActionFiles(".md").Chdir("nonexistent").Invoke(Context{CallbackValue: ""}),
	)

	assertEqual(t,
		ActionValuesDescribed("ERR", "go.mod is not a directory", "_", "").noSpace(true).skipCache(true).Invoke(Context{}),
		ActionFiles(".md").Chdir("go.mod").Invoke(Context{CallbackValue: ""}),
	)

	assertEqual(t,
		ActionStyledValues(
			"action.go", style.Of("fg-default", "bg-default"),
			"snippet.go", style.Of("fg-default", "bg-default"),
		).noSpace(true).Invoke(Context{}).Prefix("elvish/"),
		ActionFiles().Chdir("internal/shell").Invoke(Context{CallbackValue: "elvish/"}),
	)

	if newWd, _ := os.Getwd(); oldWd != newWd {
		t.Error("workdir should not be changed")
	}
}

func TestActionMessage(t *testing.T) {
	assertEqual(t,
		ActionValuesDescribed("_", "", "ERR", "example message").noSpace(true).skipCache(true).Invoke(Context{}).Prefix("docs/"),
		ActionMessage("example message").Invoke(Context{CallbackValue: "docs/"}),
	)
}

func TestActionMessageSuppress(t *testing.T) {
	assertEqual(t,
		Batch(
			ActionMessage("example message").Suppress("example"),
			ActionValues("test"),
		).ToA().Invoke(Context{}),
		ActionValues("test").noSpace(true).skipCache(true).Invoke(Context{}),
	)
}

func TestActionExecCommand(t *testing.T) {
	assertEqual(t,
		ActionMessage("go unknown: unknown command").noSpace(true).skipCache(true).Invoke(Context{}).Prefix("docs/"),
		ActionExecCommand("go", "unknown")(func(output []byte) Action { return ActionValues() }).Invoke(Context{CallbackValue: "docs/"}),
	)

	assertEqual(t,
		ActionValues("module github.com/rsteube/carapace\n").Invoke(Context{}),
		ActionExecCommand("head", "-n1", "go.mod")(func(output []byte) Action { return ActionValues(string(output)) }).Invoke(Context{}),
	)
}
