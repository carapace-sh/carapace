package carapace

import (
	"fmt"
	"sort"
	"testing"

	"github.com/rsteube/carapace/internal/assert"
	"github.com/rsteube/carapace/internal/common"
)

func assertEqual(t *testing.T, expected, actual InvokedAction) {
	sort.Sort(common.ByValue(expected.rawValues))
	sort.Sort(common.ByValue(actual.rawValues))

	assert.Equal(t, fmt.Sprintf("%+v\n", expected), fmt.Sprintf("%+v\n", actual))
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
		rawValues: common.RawValuesFrom("a", "b", "c"),
		nospace:   false,
		skipcache: false,
	}
	actual := a.Invoke(Context{})
	assertEqual(t, expected, actual)
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
		ActionValues("example/", "docs/", "internal/", "pkg/").noSpace(true).Invoke(Context{}),
		ActionDirectories().Invoke(Context{CallbackValue: ""}),
	)

	assertEqual(t,
		ActionValues("example/", "docs/", "internal/", "pkg/").noSpace(true).Invoke(Context{}).Prefix("./"),
		ActionDirectories().Invoke(Context{CallbackValue: "./"}),
	)

	assertEqual(t,
		ActionValues("_test/", "cmd/").noSpace(true).Invoke(Context{}).Prefix("example/"),
		ActionDirectories().Invoke(Context{CallbackValue: "example/"}),
	)

	assertEqual(t,
		ActionValues("_test/", "cmd/").noSpace(true).Invoke(Context{}).Prefix("example/"),
		ActionDirectories().Invoke(Context{CallbackValue: "example/cm"}),
	)
}

func TestActionFiles(t *testing.T) {
	assertEqual(t,
		ActionValues("README.md", "example/", "docs/", "internal/", "pkg/").noSpace(true).Invoke(Context{}),
		ActionFiles(".md").Invoke(Context{CallbackValue: ""}),
	)

	assertEqual(t,
		ActionValues("_test/", "cmd/", "main.go", "main_test.go").noSpace(true).Invoke(Context{}).Prefix("example/"),
		ActionFiles().Invoke(Context{CallbackValue: "example/"}),
	)
}

func TestActionFilesChdir(t *testing.T) {
	assertEqual(t,
		ActionValuesDescribed("ERR", "stat nonexistent: no such file or directory", "_", "").noSpace(true).skipCache(true).Invoke(Context{}),
		ActionFiles(".md").Chdir("nonexistent").Invoke(Context{CallbackValue: ""}),
	)

	assertEqual(t,
		ActionValuesDescribed("ERR", "go.mod is not a directory", "_", "").noSpace(true).skipCache(true).Invoke(Context{}),
		ActionFiles(".md").Chdir("go.mod").Invoke(Context{CallbackValue: ""}),
	)

	assertEqual(t,
		ActionValues("action.go", "snippet.go").noSpace(true).Invoke(Context{}).Prefix("elvish/"),
		ActionFiles().Chdir("internal").Invoke(Context{CallbackValue: "elvish/"}),
	)
}

func TestActionMessage(t *testing.T) {
	assertEqual(t,
		ActionValuesDescribed("_", "", "ERR", "example message").noSpace(true).skipCache(true).Invoke(Context{}).Prefix("docs/"),
		ActionMessage("example message").Invoke(Context{CallbackValue: "docs/"}),
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
