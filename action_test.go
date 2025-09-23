package carapace

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/carapace-sh/carapace/internal/common"
	"github.com/carapace-sh/carapace/pkg/assert"
	"github.com/carapace-sh/carapace/pkg/style"
	"github.com/carapace-sh/carapace/pkg/uid"
)

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
	assert.Equal(t, expected, actual)
}

func TestCache(t *testing.T) {
	f := func() Action {
		return ActionCallback(func(c Context) Action {
			return ActionValues(time.Now().String())
		}).Cache(15 * time.Millisecond)
	}

	a1 := f().Invoke(Context{})
	a2 := f().Invoke(Context{})
	assert.Equal(t, a1, a2)

	time.Sleep(16 * time.Millisecond)
	a3 := f().Invoke(Context{})
	assert.NotEqual(t, a1, a3)
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
	assert.Equal(t,
		ActionStyledValues(
			"example/", style.Of(style.Blue, style.Bold),
			"example-nonposix/", style.Of(style.Blue, style.Bold),
			"docs/", style.Of(style.Blue, style.Bold),
			"internal/", style.Of(style.Blue, style.Bold),
			"pkg/", style.Of(style.Blue, style.Bold),
			"third_party/", style.Of(style.Blue, style.Bold),
		).NoSpace('/').Tag("directories").Invoke(Context{}).UidF(uid.Map(
			"example/", "file://"+wd("")+"/example/",
			"example-nonposix/", "file://"+wd("")+"/example-nonposix/",
			"docs/", "file://"+wd("")+"/docs/",
			"internal/", "file://"+wd("")+"/internal/",
			"pkg/", "file://"+wd("")+"/pkg/",
			"third_party/", "file://"+wd("")+"/third_party/",
		)).Query(fmt.Sprintf("file://%v/?directories=true", wd(""))),
		ActionDirectories().Invoke(Context{Value: ""}).Filter("vendor/"),
	)

	assert.Equal(t,
		ActionStyledValues(
			"example/", style.Of(style.Blue, style.Bold),
			"example-nonposix/", style.Of(style.Blue, style.Bold),
			"docs/", style.Of(style.Blue, style.Bold),
			"internal/", style.Of(style.Blue, style.Bold),
			"pkg/", style.Of(style.Blue, style.Bold),
			"third_party/", style.Of(style.Blue, style.Bold),
		).NoSpace('/').Tag("directories").Invoke(Context{}).Prefix("./").UidF(uid.Map(
			"./example/", "file://"+wd("")+"/example/",
			"./example-nonposix/", "file://"+wd("")+"/example-nonposix/",
			"./docs/", "file://"+wd("")+"/docs/",
			"./internal/", "file://"+wd("")+"/internal/",
			"./pkg/", "file://"+wd("")+"/pkg/",
			"./third_party/", "file://"+wd("")+"/third_party/",
		)).Query(fmt.Sprintf("file://%v/?directories=true", wd(""))),
		ActionDirectories().Invoke(Context{Value: "./"}).Filter("./vendor/"),
	)

	assert.Equal(t,
		ActionStyledValues(
			"cmd/", style.Of(style.Blue, style.Bold),
		).NoSpace('/').Tag("directories").Invoke(Context{}).Prefix("example/").UidF(uid.Map(
			"example/cmd/", "file://"+wd("")+"/example/cmd/",
		)).Query(fmt.Sprintf("file://%v/?directories=true", wd(""))),
		ActionDirectories().Invoke(Context{Value: "example/"}),
	)

	assert.Equal(t,
		ActionStyledValues(
			"cmd/", style.Of(style.Blue, style.Bold),
		).NoSpace('/').Tag("directories").Invoke(Context{}).Prefix("example/").UidF(uid.Map(
			"example/cmd/", "file://"+wd("")+"/example/cmd/",
		)).Query(fmt.Sprintf("file://%v/?directories=true", wd(""))),
		ActionDirectories().Invoke(Context{Value: "example/cm"}),
	)
}

func TestActionFiles(t *testing.T) {
	assert.Equal(t,
		ActionStyledValues(
			"README.md", style.Default,
			"example/", style.Of(style.Blue, style.Bold),
			"example-nonposix/", style.Of(style.Blue, style.Bold),
			"docs/", style.Of(style.Blue, style.Bold),
			"internal/", style.Of(style.Blue, style.Bold),
			"pkg/", style.Of(style.Blue, style.Bold),
			"third_party/", style.Of(style.Blue, style.Bold),
		).NoSpace('/').Tag("files").Invoke(Context{}).UidF(uid.Map(
			"README.md", "file://"+wd("")+"/README.md",
			"example/", "file://"+wd("")+"/example/",
			"example-nonposix/", "file://"+wd("")+"/example-nonposix/",
			"docs/", "file://"+wd("")+"/docs/",
			"internal/", "file://"+wd("")+"/internal/",
			"pkg/", "file://"+wd("")+"/pkg/",
			"third_party/", "file://"+wd("")+"/third_party/",
		)).Query(fmt.Sprintf("file://%v/?suffix=.md", wd(""))),
		ActionFiles(".md").Invoke(Context{Value: ""}).Filter("vendor/"),
	)

	assert.Equal(t,
		ActionStyledValues(
			"README.md", style.Default,
			"cmd/", style.Of(style.Blue, style.Bold),
			"main.go", style.Default,
			"main_test.go", style.Default,
		).NoSpace('/').Tag("files").Invoke(Context{}).Prefix("example/").UidF(uid.Map(
			"example/README.md", "file://"+wd("example")+"/README.md",
			"example/cmd/", "file://"+wd("example")+"/cmd/",
			"example/main.go", "file://"+wd("example")+"/main.go",
			"example/main_test.go", "file://"+wd("example")+"/main_test.go",
		)).Query(fmt.Sprintf("file://%v/", wd(""))),
		ActionFiles().Invoke(Context{Value: "example/"}).Filter("example/example"),
	)
}

func TestActionFilesChdir(t *testing.T) {
	oldWd, _ := os.Getwd()

	assert.Equal(t,
		ActionMessage(fmt.Sprintf("stat %v: no such file or directory", wd("nonexistent"))).Invoke(Context{}),
		ActionFiles(".md").Chdir("nonexistent").Invoke(Context{}),
	)

	assert.Equal(t,
		ActionMessage(fmt.Sprintf("not a directory: %v/go.mod", wd(""))).Invoke(Context{}),
		ActionFiles(".md").Chdir("go.mod").Invoke(Context{Value: ""}),
	)

	assert.Equal(t,
		ActionStyledValues(
			"action.go", style.Default,
			"snippet.go", style.Default,
		).NoSpace('/').Tag("files").Invoke(Context{}).Prefix("elvish/").UidF(uid.Map(
			"elvish/action.go", "file://"+wd("internal/shell")+"/elvish/action.go",
			"elvish/snippet.go", "file://"+wd("internal/shell")+"/elvish/snippet.go",
		)).Query(fmt.Sprintf("file://%v/internal/shell", wd(""))), // TODO uid shouldn't include internal/shell
		ActionFiles().Chdir("internal/shell").Invoke(Context{Value: "elvish/"}),
	)

	if newWd, _ := os.Getwd(); oldWd != newWd {
		t.Error("workdir should not be changed")
	}
}

func TestActionMessage(t *testing.T) {
	expected := ActionValues()
	expected.meta.Messages.Add("example message")

	assert.Equal(t,
		expected.Invoke(Context{}),
		ActionMessage("example message").Invoke(Context{Value: "docs/"}),
	)
}

func TestActionMessageSuppress(t *testing.T) {
	assert.Equal(t,
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
	assert.Equal(t,
		ActionMessage("go unknown: unknown command").Invoke(NewContext()).Prefix("docs/"),
		ActionExecCommand("go", "unknown")(func(output []byte) Action { return ActionValues() }).Invoke(context),
	)

	assert.Equal(t,
		ActionValues("module github.com/carapace-sh/carapace\n").Invoke(Context{}),
		ActionExecCommand("head", "-n1", "go.mod")(func(output []byte) Action { return ActionValues(string(output)) }).Invoke(Context{}),
	)
}
