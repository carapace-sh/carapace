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
	sort.Sort(common.ByValue(expected.action.rawValues))
	sort.Sort(common.ByValue(actual.action.rawValues))

	e, _ := json.MarshalIndent(expected.action.rawValues, "", "  ")
	a, _ := json.MarshalIndent(actual.action.rawValues, "", "  ")
	assert.Equal(t, string(e), string(a))

	eMeta, _ := json.MarshalIndent(expected.action.meta, "", "  ")
	aMeta, _ := json.MarshalIndent(actual.action.meta, "", "  ")
	assert.Equal(t, string(eMeta), string(aMeta))
}

func assertNotEqual(t *testing.T, expected, actual InvokedAction) {
	sort.Sort(common.ByValue(expected.action.rawValues))
	sort.Sort(common.ByValue(actual.action.rawValues))

	e, _ := json.MarshalIndent(expected.action.rawValues, "", "  ")
	a, _ := json.MarshalIndent(actual.action.rawValues, "", "  ")

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
			Filter("").
			Prefix("").
			Suffix("").
			ToA()
	})
	if !a.meta.Messages.IsEmpty() {
		t.Fatal("uninvoked action should not contain messages")
	}
	if a.Invoke(Context{}).action.meta.Messages.IsEmpty() {
		t.Fatal("invoked action should contain messages")
	}
}

func TestNoSpace(t *testing.T) {
	a := ActionCallback(func(c Context) Action {
		return ActionValues().Invoke(c).Merge(
			ActionMultiParts("", func(c Context) Action {
				return ActionMessage("nospace")
			}).Invoke(c)).
			Filter("").
			Prefix("").
			Suffix("").
			ToA()
	})
	if a.meta.Nospace.Matches("x") {
		t.Fatal("uninvoked nospace should not match")
	}
	if !a.Invoke(Context{}).action.meta.Nospace.Matches("x") {
		t.Fatal("invoked nospace should match")
	}
}

func TestActionDirectories(t *testing.T) {
	assertEqual(t,
		ActionStyledValues(
			"example/", style.Of(style.Blue, style.Bold),
			"example-nonposix/", style.Of(style.Blue, style.Bold),
			"docs/", style.Of(style.Blue, style.Bold),
			"internal/", style.Of(style.Blue, style.Bold),
			"pkg/", style.Of(style.Blue, style.Bold),
			"third_party/", style.Of(style.Blue, style.Bold),
		).NoSpace('/').Tag("directories").Invoke(Context{}),
		ActionDirectories().Invoke(Context{Value: ""}).Filter("vendor/"),
	)

	assertEqual(t,
		ActionStyledValues(
			"example/", style.Of(style.Blue, style.Bold),
			"example-nonposix/", style.Of(style.Blue, style.Bold),
			"docs/", style.Of(style.Blue, style.Bold),
			"internal/", style.Of(style.Blue, style.Bold),
			"pkg/", style.Of(style.Blue, style.Bold),
			"third_party/", style.Of(style.Blue, style.Bold),
		).NoSpace('/').Tag("directories").Invoke(Context{}).Prefix("./"),
		ActionDirectories().Invoke(Context{Value: "./"}).Filter("./vendor/"),
	)

	assertEqual(t,
		ActionStyledValues(
			"_test/", style.Of(style.Blue, style.Bold),
			"cmd/", style.Of(style.Blue, style.Bold),
		).NoSpace('/').Tag("directories").Invoke(Context{}).Prefix("example/"),
		ActionDirectories().Invoke(Context{Value: "example/"}),
	)

	assertEqual(t,
		ActionStyledValues(
			"cmd/", style.Of(style.Blue, style.Bold),
		).NoSpace('/').Tag("directories").Invoke(Context{}).Prefix("example/"),
		ActionDirectories().Invoke(Context{Value: "example/cm"}),
	)
}

func TestActionFiles(t *testing.T) {
	assertEqual(t,
		ActionStyledValues(
			"README.md", style.Default,
			"example/", style.Of(style.Blue, style.Bold),
			"example-nonposix/", style.Of(style.Blue, style.Bold),
			"docs/", style.Of(style.Blue, style.Bold),
			"internal/", style.Of(style.Blue, style.Bold),
			"pkg/", style.Of(style.Blue, style.Bold),
			"third_party/", style.Of(style.Blue, style.Bold),
		).NoSpace('/').Tag("files").Invoke(Context{}),
		ActionFiles(".md").Invoke(Context{Value: ""}).Filter("vendor/"),
	)

	assertEqual(t,
		ActionStyledValues(
			"README.md", style.Default,
			"_test/", style.Of(style.Blue, style.Bold),
			"cmd/", style.Of(style.Blue, style.Bold),
			"main.go", style.Default,
			"main_test.go", style.Default,
		).NoSpace('/').Tag("files").Invoke(Context{}).Prefix("example/"),
		ActionFiles().Invoke(Context{Value: "example/"}).Filter("example/example"),
	)
}

func TestActionFilesChdir(t *testing.T) {
	oldWd, _ := os.Getwd()

	assertEqual(t,
		ActionMessage(fmt.Sprintf("stat %v: no such file or directory", wd("nonexistent"))).Invoke(Context{}),
		ActionFiles(".md").Chdir("nonexistent").Invoke(Context{}),
	)

	assertEqual(t,
		ActionMessage(fmt.Sprintf("not a directory: %v/go.mod", wd(""))).Invoke(Context{}),
		ActionFiles(".md").Chdir("go.mod").Invoke(Context{Value: ""}),
	)

	assertEqual(t,
		ActionStyledValues(
			"action.go", style.Default,
			"snippet.go", style.Default,
		).NoSpace('/').Tag("files").Invoke(Context{}).Prefix("elvish/"),
		ActionFiles().Chdir("internal/shell").Invoke(Context{Value: "elvish/"}),
	)

	if newWd, _ := os.Getwd(); oldWd != newWd {
		t.Error("workdir should not be changed")
	}
}

func TestActionMessage(t *testing.T) {
	expected := ActionValues()
	expected.meta.Messages.Add("example message")

	assertEqual(t,
		expected.Invoke(Context{}),
		ActionMessage("example message").Invoke(Context{Value: "docs/"}),
	)
}

func TestActionMessageSuppress(t *testing.T) {
	assertEqual(t,
		Batch(
			ActionMessage("example message").Suppress("example"),
			ActionValues("test"),
		).ToA().Invoke(Context{}),
		ActionValues("test").Invoke(Context{}), // TODO suppress does not reset nospace (is that even possible?)
	)
}

func TestActionExecCommand(t *testing.T) {
	context := NewContext()
	context.Value = "docs/"
	assertEqual(t,
		ActionMessage("go unknown: unknown command").Invoke(NewContext()).Prefix("docs/"),
		ActionExecCommand("go", "unknown")(func(output []byte) Action { return ActionValues() }).Invoke(context),
	)

	assertEqual(t,
		ActionValues("module github.com/rsteube/carapace\n").Invoke(Context{}),
		ActionExecCommand("head", "-n1", "go.mod")(func(output []byte) Action { return ActionValues(string(output)) }).Invoke(Context{}),
	)
}
