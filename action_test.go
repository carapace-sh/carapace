package carapace

import (
	"encoding/json"
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

	e, _ := json.MarshalIndent(expected.rawValues, "", "  ")
	a, _ := json.MarshalIndent(actual.rawValues, "", "  ")
	assert.Equal(t, string(e), string(a))

	eMeta, _ := json.MarshalIndent(expected.meta, "", "  ")
	aMeta, _ := json.MarshalIndent(actual.meta, "", "  ")
	assert.Equal(t, string(eMeta), string(aMeta))
}

func assertNotEqual(t *testing.T, expected, actual InvokedAction) {
	sort.Sort(common.ByValue(expected.rawValues))
	sort.Sort(common.ByValue(actual.rawValues))

	e, _ := json.MarshalIndent(expected.rawValues, "", "  ")
	a, _ := json.MarshalIndent(actual.rawValues, "", "  ")

	if string(e) == string(a) {
		t.Errorf("should differ:\n%v", a)
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
		},
	}
	actual := a.Invoke(Context{})
	assertEqual(t, expected, actual)
}

func TestCache(t *testing.T) {
	f := func() Action {
		return ActionCallback(func(c Context) Action {
			return ActionValues(time.Now().String())
		}).Cache(15 * time.Millisecond)
	}

	a1 := f().Invoke(Context{})
	a2 := f().Invoke(Context{})
	assertEqual(t, a1, a2)

	time.Sleep(16 * time.Millisecond)
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
	if !a.meta.Messages.IsEmpty() {
		t.Fatal("uninvoked action should not contain messages")
	}
	if a.Invoke(Context{}).meta.Messages.IsEmpty() {
		t.Fatal("invoked action should contain messages")
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
	if a.meta.Nospace.Matches("x") {
		t.Fatal("uninvoked nospace should not match")
	}
	if !a.Invoke(Context{}).meta.Nospace.Matches("x") {
		t.Fatal("invoked nospace should match")
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
		).NoSpace('/').Invoke(Context{}),
		ActionDirectories().Invoke(Context{CallbackValue: ""}).Filter([]string{"vendor/"}),
	)

	assertEqual(t,
		ActionStyledValues(
			"example/", style.Of("fg-default", "bg-default", style.Blue, style.Bold),
			"docs/", style.Of("fg-default", "bg-default", style.Blue, style.Bold),
			"internal/", style.Of("fg-default", "bg-default", style.Blue, style.Bold),
			"pkg/", style.Of("fg-default", "bg-default", style.Blue, style.Bold),
			"third_party/", style.Of("fg-default", "bg-default", style.Blue, style.Bold),
		).NoSpace('/').Invoke(Context{}).Prefix("./"),
		ActionDirectories().Invoke(Context{CallbackValue: "./"}).Filter([]string{"./vendor/"}),
	)

	assertEqual(t,
		ActionStyledValues(
			"_test/", style.Of("fg-default", "bg-default", style.Blue, style.Bold),
			"cmd/", style.Of("fg-default", "bg-default", style.Blue, style.Bold),
		).NoSpace('/').Invoke(Context{}).Prefix("example/"),
		ActionDirectories().Invoke(Context{CallbackValue: "example/"}),
	)

	assertEqual(t,
		ActionStyledValues(
			"cmd/", style.Of("fg-default", "bg-default", style.Blue, style.Bold),
		).NoSpace('/').Invoke(Context{}).Prefix("example/"),
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
		).NoSpace('/').Invoke(Context{}),
		ActionFiles(".md").Invoke(Context{CallbackValue: ""}).Filter([]string{"vendor/"}),
	)

	assertEqual(t,
		ActionStyledValues(
			"README.md", style.Of("fg-default", "bg-default"),
			"_test/", style.Of("fg-default", "bg-default", style.Blue, style.Bold),
			"cmd/", style.Of("fg-default", "bg-default", style.Blue, style.Bold),
			"main.go", style.Of("fg-default", "bg-default"),
			"main_test.go", style.Of("fg-default", "bg-default"),
		).NoSpace('/').Invoke(Context{}).Prefix("example/"),
		ActionFiles().Invoke(Context{CallbackValue: "example/"}).Filter([]string{"example/example"}),
	)
}

func TestActionFilesChdir(t *testing.T) {
	oldWd, _ := os.Getwd()

	assertEqual(t,
		ActionMessage(fmt.Sprintf("stat %v: no such file or directory", wd("nonexistent"))).NoSpace('/').Invoke(Context{}),
		ActionFiles(".md").Chdir("nonexistent").Invoke(Context{}),
	)

	assertEqual(t,
		ActionMessage(fmt.Sprintf("not a directory: %v/go.mod", wd(""))).Invoke(Context{}),
		ActionFiles(".md").Chdir("go.mod").Invoke(Context{CallbackValue: ""}),
	)

	assertEqual(t,
		ActionStyledValues(
			"action.go", style.Of("fg-default", "bg-default"),
			"snippet.go", style.Of("fg-default", "bg-default"),
		).NoSpace('/').Invoke(Context{}).Prefix("elvish/"),
		ActionFiles().Chdir("internal/shell").Invoke(Context{CallbackValue: "elvish/"}),
	)

	if newWd, _ := os.Getwd(); oldWd != newWd {
		t.Error("workdir should not be changed")
	}
}

func TestActionMessage(t *testing.T) {
	expected := ActionValues().NoSpace()
	expected.meta.Messages.Add("example message")

	assertEqual(t,
		expected.Invoke(Context{}),
		ActionMessage("example message").Invoke(Context{CallbackValue: "docs/"}),
	)
}

func TestActionMessageSuppress(t *testing.T) {
	assertEqual(t,
		Batch(
			ActionMessage("example message").Suppress("example"),
			ActionValues("test"),
		).ToA().Invoke(Context{}),
		ActionValues("test").NoSpace('*').Invoke(Context{}), // TODO suppress does not reset nospace (is that even possible?)
	)
}

func TestActionExecCommand(t *testing.T) {
	assertEqual(t,
		ActionMessage("go unknown: unknown command").NoSpace('/').Invoke(Context{}).Prefix("docs/"),
		ActionExecCommand("go", "unknown")(func(output []byte) Action { return ActionValues() }).Invoke(Context{CallbackValue: "docs/"}),
	)

	assertEqual(t,
		ActionValues("module github.com/rsteube/carapace\n").Invoke(Context{}),
		ActionExecCommand("head", "-n1", "go.mod")(func(output []byte) Action { return ActionValues(string(output)) }).Invoke(Context{}),
	)
}
