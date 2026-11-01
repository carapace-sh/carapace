package carapace

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
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

func TestNoPrefix(t *testing.T) {
	a := ActionCallback(func(c Context) Action {
		return ActionValues().Invoke(c).Merge(
			ActionValues("one", "two").NoPrefix('-').Invoke(c)).
			ToA()
	})
	if a.meta.NoPrefix.Matches("-x") {
		t.Fatal("uninvoked noprefix should not match")
	}
	if !a.Invoke(Context{}).action.meta.NoPrefix.Matches("-x") {
		t.Fatal("invoked noprefix should match")
	}
}

func TestNoPrefixAll(t *testing.T) {
	a := ActionValues("--", "---", "----").NoPrefix()
	if !a.Invoke(Context{}).action.meta.NoPrefix.Matches("x") {
		t.Fatal("NoPrefix() without args should match all prefixes")
	}
}

func TestPrefixAutoNoPrefix(t *testing.T) {
	a := ActionValues("a", "b").Prefix("file://").NoPrefix('f')
	if !a.Invoke(Context{}).action.meta.NoPrefix.Matches("file://test") {
		t.Fatal("NoPrefix('f') should match file:// prefix")
	}
}

func TestActionDirectories(t *testing.T) {
	assert.Equal(t,
		ActionStyledValues(
			"example/", style.Of(style.Blue, style.Bold),
			"example-nonposix/", style.Of(style.Blue, style.Bold),
			"example-multi/", style.Of(style.Blue, style.Bold),
			"docs/", style.Of(style.Blue, style.Bold),
			"internal/", style.Of(style.Blue, style.Bold),
			"pkg/", style.Of(style.Blue, style.Bold),
			"skills/", style.Of(style.Blue, style.Bold),
			"third_party/", style.Of(style.Blue, style.Bold),
		).NoSpace('/').Tag("directories").Invoke(Context{}).UidF(uid.Map(
			"example/", "file://"+wd("")+"/example/",
			"example-nonposix/", "file://"+wd("")+"/example-nonposix/",
			"example-multi/", "file://"+wd("")+"/example-multi/",
			"docs/", "file://"+wd("")+"/docs/",
			"internal/", "file://"+wd("")+"/internal/",
			"pkg/", "file://"+wd("")+"/pkg/",
			"skills/", "file://"+wd("")+"/skills/",
			"third_party/", "file://"+wd("")+"/third_party/",
		)).QueryS(fmt.Sprintf("file://directories?C_DIR=%v",
			url.QueryEscape(wd("")+"/"),
		)),
		ActionDirectories().Invoke(Context{Value: ""}).Filter("vendor/"),
	)

	assert.Equal(t,
		ActionStyledValues(
			"example/", style.Of(style.Blue, style.Bold),
			"example-nonposix/", style.Of(style.Blue, style.Bold),
			"example-multi/", style.Of(style.Blue, style.Bold),
			"docs/", style.Of(style.Blue, style.Bold),
			"internal/", style.Of(style.Blue, style.Bold),
			"pkg/", style.Of(style.Blue, style.Bold),
			"skills/", style.Of(style.Blue, style.Bold),
			"third_party/", style.Of(style.Blue, style.Bold),
		).NoSpace('/').Tag("directories").Invoke(Context{}).Prefix("./").UidF(uid.Map(
			"./example/", "file://"+wd("")+"/example/",
			"./example-nonposix/", "file://"+wd("")+"/example-nonposix/",
			"./example-multi/", "file://"+wd("")+"/example-multi/",
			"./docs/", "file://"+wd("")+"/docs/",
			"./internal/", "file://"+wd("")+"/internal/",
			"./pkg/", "file://"+wd("")+"/pkg/",
			"./skills/", "file://"+wd("")+"/skills/",
			"./third_party/", "file://"+wd("")+"/third_party/",
		)).QueryS(fmt.Sprintf("file://directories?C_DIR=%v&C_VALUE=%v",
			url.QueryEscape(wd("")+"/"),
			url.QueryEscape("./"),
		)),
		ActionDirectories().Invoke(Context{Value: "./"}).Filter("./vendor/"),
	)

	assert.Equal(t,
		ActionStyledValues(
			"cmd/", style.Of(style.Blue, style.Bold),
		).NoSpace('/').Tag("directories").Invoke(Context{}).Prefix("example/").UidF(uid.Map(
			"example/cmd/", "file://"+wd("")+"/example/cmd/",
		)).QueryS(fmt.Sprintf("file://directories?C_DIR=%v&C_VALUE=%v",
			url.QueryEscape(wd("")+"/"),
			url.QueryEscape("example/"),
		)),
		ActionDirectories().Invoke(Context{Value: "example/"}),
	)

	assert.Equal(t,
		ActionStyledValues(
			"cmd/", style.Of(style.Blue, style.Bold),
		).NoSpace('/').Tag("directories").Invoke(Context{}).Prefix("example/").UidF(uid.Map(
			"example/cmd/", "file://"+wd("")+"/example/cmd/",
		)).QueryS(fmt.Sprintf("file://directories?C_DIR=%v&C_VALUE=%v",
			url.QueryEscape(wd("")+"/"),
			url.QueryEscape("example/cm"),
		)),
		ActionDirectories().Invoke(Context{Value: "example/cm"}),
	)
}

func TestActionFiles(t *testing.T) {
	assert.Equal(t,
		ActionStyledValues(
			"AGENTS.md", style.Default,
			"README.md", style.Default,
			"example/", style.Of(style.Blue, style.Bold),
			"example-nonposix/", style.Of(style.Blue, style.Bold),
			"example-multi/", style.Of(style.Blue, style.Bold),
			"docs/", style.Of(style.Blue, style.Bold),
			"internal/", style.Of(style.Blue, style.Bold),
			"pkg/", style.Of(style.Blue, style.Bold),
			"skills/", style.Of(style.Blue, style.Bold),
			"third_party/", style.Of(style.Blue, style.Bold),
		).NoSpace('/').Tag("files").Invoke(Context{}).UidF(uid.Map(
			"AGENTS.md", "file://"+wd("")+"/AGENTS.md",
			"README.md", "file://"+wd("")+"/README.md",
			"example/", "file://"+wd("")+"/example/",
			"example-nonposix/", "file://"+wd("")+"/example-nonposix/",
			"example-multi/", "file://"+wd("")+"/example-multi/",
			"docs/", "file://"+wd("")+"/docs/",
			"internal/", "file://"+wd("")+"/internal/",
			"pkg/", "file://"+wd("")+"/pkg/",
			"skills/", "file://"+wd("")+"/skills/",
			"third_party/", "file://"+wd("")+"/third_party/",
		)).QueryS(fmt.Sprintf("file://files?C_DIR=%v&suffix=.md",
			url.QueryEscape(wd("")+"/"),
		)),
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
		)).QueryS(fmt.Sprintf("file://files?C_DIR=%v&C_VALUE=%v",
			url.QueryEscape(wd("")+"/"),
			url.QueryEscape("example/"),
		)),
		ActionFiles().Invoke(Context{Value: "example/"}).Filter("example/example"),
	)
}

func TestActionFilesChdir(t *testing.T) {
	oldWd, _ := os.Getwd()

	var expectedErrorMsg string
	if runtime.GOOS == "windows" {
		expectedErrorMsg = fmt.Sprintf("CreateFile %v: The system cannot find the file specified.", wd("nonexistent"))
	} else {
		expectedErrorMsg = fmt.Sprintf("stat %v: no such file or directory", wd("nonexistent"))
	}

	assert.Equal(t,
		ActionMessage(expectedErrorMsg).Invoke(Context{}),
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
		)).QueryS(fmt.Sprintf("file://files?C_DIR=%v&C_VALUE=%v",
			url.QueryEscape(wd("internal/shell")),
			url.QueryEscape("elvish/"),
		)),
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
		ActionExecCommand("head", "-n1", "go.mod")(func(output []byte) Action {
			return ActionValues(strings.ReplaceAll(string(output), "\r\n", "\n"))
		}).Invoke(Context{}),
	)
}

func TestActionSplitMultibytePrefix(t *testing.T) {
	// Span offsets are rune-based; the completion prefix must slice the
	// line by runes, not bytes (café is 4 runes, 5 bytes).
	a := ActionValues().Split("bash")
	invoked := a.Invoke(Context{
		Value: `echo café | grep $HO`,
		Env:   []string{"HOME=/tmp"},
	})
	if len(invoked.action.rawValues) != 1 || invoked.action.rawValues[0].Value != "echo café | grep $HOME" {
		t.Errorf("rawValues = %v, want [echo café | grep $HOME]", invoked.action.rawValues)
	}
}

func TestActionFilesHomedir(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("homedir expansion is not supported on windows")
	}

	var username, userHome string
	passwd, err := os.ReadFile("/etc/passwd")
	if err != nil {
		t.Skipf("user database not readable: %v", err)
	}
	for line := range strings.SplitSeq(string(passwd), "\n") {
		fields := strings.Split(line, ":")
		if len(fields) != 7 {
			continue
		}
		home := filepath.ToSlash(strings.TrimSuffix(fields[5], "/"))
		if home == "" || home == "/" {
			continue // system entries (e.g. nobody/daemon) aren't useful for homedir completion
		}
		if entries, err := os.ReadDir(home); err != nil || len(entries) == 0 {
			continue
		}
		username, userHome = fields[0], home
		break
	}
	if username == "" {
		t.Skip("no user with readable home directory found")
	}
	usernamePrefix := "~" + username + "/"

	entries, err := os.ReadDir(userHome)
	if err != nil {
		t.Skipf("home directory not readable: %v", err)
	}
	vals := make([]string, 0, len(entries)*2)
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".") {
			continue // hidden files are not completed
		}
		info, infoErr := e.Info()
		if infoErr != nil {
			t.Skipf("home directory entry unreadable: %v", infoErr)
		}
		if info.IsDir() {
			vals = append(vals, usernamePrefix+e.Name()+"/")
		} else {
			vals = append(vals, usernamePrefix+e.Name())
		}
	}

	// path completion within `~user/`: every visible entry is provided as
	// `~user/<name>` (with trailing slash for directories)
	invoked := ActionFiles().Invoke(Context{Value: usernamePrefix})
	if len(invoked.action.rawValues) == 0 {
		t.Fatalf("%v completion returned no candidates", usernamePrefix)
	}
	for _, actual := range invoked.action.rawValues {
		if !strings.HasPrefix(actual.Value, usernamePrefix) {
			t.Errorf("candidate = %v, want prefix %v", actual.Value, usernamePrefix)
		}
	}
	for _, e := range vals {
		if !slices.ContainsFunc(invoked.action.rawValues, func(r common.RawValue) bool {
			return r.Value == e || r.Value == e+"/"
		}) {
			provided := make([]string, 0, len(invoked.action.rawValues))
			for _, r := range invoked.action.rawValues {
				provided = append(provided, r.Value)
			}
			t.Errorf("%v completion missing %v\nprovided: %v", usernamePrefix, e, strings.Join(provided, " "))
		}
	}

	// usernames: candidates are shaped `~name/`
	users := ActionFiles().Invoke(Context{Value: "~"}).action.rawValues
	if len(users) == 0 {
		t.Fatal("no usernames provided")
	}
	for _, v := range users {
		if !strings.HasPrefix(v.Value, "~") || !strings.HasSuffix(v.Value, "/") {
			t.Errorf("username candidate = %v, want ~/name/", v.Value)
		}
	}

	// full username without slash yet should still provide `~user/`
	completed := ActionFiles().Invoke(Context{Value: "~" + username})
	if len(completed.action.rawValues) == 0 {
		t.Errorf("~%v completion = %v, want at least [%v/]", username, completed.action.rawValues, usernamePrefix)
	}
	for _, v := range completed.action.rawValues {
		if !strings.HasPrefix(v.Value, "~"+username) {
			t.Errorf("~%v completion candidate = %v, want prefix ~%v", username, v.Value, username)
		}
	}
}
