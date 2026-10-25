package carapace

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/carapace-sh/carapace/pkg/assert"
	"github.com/carapace-sh/carapace/pkg/uid"
	"github.com/spf13/cobra"
)

func TestMain(m *testing.M) {
	os.Unsetenv("LS_COLORS")
	os.Exit(m.Run())
}

func execCompletion(args ...string) (context Context) {
	rootCmd := &cobra.Command{
		Use: "root",
		Run: func(cmd *cobra.Command, args []string) {},
	}
	rootCmd.Flags().String("multiparts", "", "")

	Gen(rootCmd).FlagCompletion(ActionMap{
		"multiparts": ActionMultiParts(",", func(c Context) Action {
			context = c
			return ActionValues()
		}),
	})

	Gen(rootCmd).PositionalAnyCompletion(
		ActionMultiParts(":", func(c Context) Action {
			context = c
			return ActionValues()
		}),
	)

	subCmd := &cobra.Command{
		Use: "sub",
		Run: func(cmd *cobra.Command, args []string) {},
	}

	Gen(subCmd).PositionalAnyCompletion(
		ActionCallback(func(c Context) Action {
			context = c
			return ActionValues()
		}),
	)

	rootCmd.AddCommand(subCmd)

	os.Args = append([]string{"root", "_carapace", "elvish", "root"}, args...)
	_ = rootCmd.Execute()
	return
}

// osGetwd returns the working directory as reported by NewContext (unmodified)
func osGetwd() string {
	wd, _ := os.Getwd()
	return wd
}

func testContext(t *testing.T, expected Context, args ...string) {
	t.Run(strings.Join(args, " "), func(t *testing.T) {
		null, _ := os.Open(os.DevNull)
		defer null.Close()

		sOut := os.Stdout
		sErr := os.Stderr

		os.Stdout = null
		os.Stderr = null
		actual := execCompletion(args...)
		actual.Env = []string{} // skip env
		os.Stdout = sOut
		os.Stderr = sErr

		e, _ := json.MarshalIndent(expected, "", "  ")
		a, _ := json.MarshalIndent(actual, "", "  ")
		assert.Equal(t, string(e), string(a))
	})
}

func TestContext(t *testing.T) {
	testContext(t, Context{
		Value: "",
		Args:  []string{},
		Parts: []string{},
		Env:   []string{},
		Dir:   osGetwd(),
	},
		"")

	testContext(t, Context{
		Value: "",
		Args:  []string{"pos1"},
		Parts: []string{},
		Env:   []string{},
		Dir:   osGetwd(),
	},
		"pos1", "")

	testContext(t, Context{
		Value: "po",
		Args:  []string{"pos1", "pos2"},
		Parts: []string{},
		Env:   []string{},
		Dir:   osGetwd(),
	},
		"pos1", "pos2", "po")

	testContext(t, Context{
		Value: "",
		Args:  []string{},
		Parts: []string{},
		Env:   []string{},
		Dir:   osGetwd(),
	},
		"--multiparts", "")

	testContext(t, Context{
		Value: "fir",
		Args:  []string{},
		Parts: []string{},
		Env:   []string{},
		Dir:   osGetwd(),
	},
		"--multiparts", "fir")

	testContext(t, Context{
		Value: "seco",
		Args:  []string{"pos1"},
		Parts: []string{"first"},
		Env:   []string{},
		Dir:   osGetwd(),
	},
		"pos1", "--multiparts", "first,seco")

	testContext(t, Context{
		Value: "pos",
		Args:  []string{},
		Parts: []string{},
		Env:   []string{},
		Dir:   osGetwd(),
	},
		"pos")

	testContext(t, Context{
		Value: "sec",
		Args:  []string{},
		Parts: []string{"first"},
		Env:   []string{},
		Dir:   osGetwd(),
	},
		"first:sec")

	testContext(t, Context{
		Value: "thi",
		Args:  []string{"first:second"},
		Parts: []string{},
		Env:   []string{},
		Dir:   osGetwd(),
	},
		"first:second", "thi")
}

func TestStandalone(t *testing.T) {
	cmd := &cobra.Command{}
	if cmd.CompletionOptions.DisableDefaultCmd == true {
		t.Fail()
	}

	Gen(cmd).Standalone()

	if cmd.CompletionOptions.DisableDefaultCmd == false {
		t.Fail()
	}
}

func TestIsCallback(t *testing.T) {
	os.Args = []string{uid.Executable(), "subcommand"}
	if IsCallback() {
		t.Fail()
	}

	os.Args = []string{uid.Executable(), "_carapace"}
	if !IsCallback() {
		t.Fail()
	}
}

func TestSnippet(t *testing.T) {
	// snippets need a completer name; an unnamed command yields none
	cmd := &cobra.Command{}
	for _, shell := range []string{"bash", "elvish", "fish", "oil", "powershell", "xonsh", "zsh"} {
		if s, err := Gen(cmd).Snippet(shell); s != "" || err != nil {
			t.Errorf("expected empty snippet for unnamed command (%v)", shell)
		}
	}

	if _, err := Gen(cmd).Snippet("unknown"); err == nil {
		t.Error("unknown shell should fail")
	}
}

func TestTest(t *testing.T) {
	Test(t)
}

func TestComplete(t *testing.T) {
	cmd := &cobra.Command{
		Use: "test",
	}
	cmd.Flags().BoolP("a", "1", false, "")
	cmd.Flags().BoolP("b", "2", false, "")

	if s, err := complete(cmd, []string{"elvish", "_", "test", "-1"}); err != nil || s != `{"Usage":"","Messages":[],"DescriptionStyle":"dim","Candidates":[{"Value":"-12","Display":"2","Description":"","CodeSuffix":"","Style":"default"},{"Value":"-1h","Display":"h","Description":"help for test","CodeSuffix":"","Style":"default"}]}` {
		t.Error(s)
	}
}

func TestCompleteOptarg(t *testing.T) {
	cmd := &cobra.Command{
		Use: "test",
	}
	cmd.Flags().String("opt", "", "")
	cmd.Flag("opt").NoOptDefVal = " "

	Gen(cmd).FlagCompletion(ActionMap{
		"opt": ActionValuesDescribed("value", "description"),
	})

	if s, err := complete(cmd, []string{"elvish", "_", "test", "--opt="}); err != nil || s != `{"Usage":"","Messages":[],"DescriptionStyle":"dim","Candidates":[{"Value":"--opt=value","Display":"value","Description":"description","CodeSuffix":" ","Style":"default"}]}` {
		t.Error(s)
	}
}

func TestCompleteSnippet(t *testing.T) {
	cmd := &cobra.Command{
		Use: "test",
	}

	if s, err := complete(cmd, []string{"bash"}); err != nil || !strings.Contains(s, "#!/bin/bash") {
		t.Error(s)
	}
}

func TestCompletePositionalWithSpace(t *testing.T) {
	cmd := &cobra.Command{
		Use: "test",
	}

	Gen(cmd).PositionalCompletion(
		ActionValues("positional with space"),
	)

	if s, err := complete(cmd, []string{"elvish", "_", "positional "}); err != nil || s != `{"Usage":"","Messages":[],"DescriptionStyle":"dim","Candidates":[{"Value":"positional with space","Display":"positional with space","Description":"","CodeSuffix":" ","Style":"default"}]}` {
		t.Error(s)
	}
}

func TestGenWithOptions(t *testing.T) {
	sub1 := &cobra.Command{Use: "sub1"}
	sub2 := &cobra.Command{Use: "sub2"}
	root := &cobra.Command{Use: "root"}

	Gen(root, WithSubcommands(sub1, sub2))

	if s, err := Gen(root).Snippet("bash"); err != nil {
		t.Error(err)
	} else if !strings.Contains(s, "#!/bin/bash") {
		t.Error("bash multi snippet missing shebang")
	}
}

func TestGenWithSubcommandsSnippet(t *testing.T) {
	sub1 := &cobra.Command{Use: "sub1"}
	sub2 := &cobra.Command{Use: "sub2"}
	root := &cobra.Command{Use: "root"}

	Gen(root, WithSubcommands(sub1, sub2))

	shells := []string{"bash", "zsh", "fish", "elvish", "nushell", "powershell", "xonsh", "oil", "tcsh", "bash-ble"}
	for _, sh := range shells {
		s, err := Gen(root).Snippet(sh)
		if err != nil {
			t.Errorf("%v: %v", sh, err)
		}
		if s == "" {
			t.Errorf("%v: empty snippet", sh)
		}
	}

	_, err := Gen(root).Snippet("cmd-clink")
	if err == nil {
		t.Error("expected error for cmd-clink")
	}
}

func TestGenWithSnippetFuncs(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	Gen(cmd, WithSnippetFuncs(map[string]string{
		"bash": "declare -x CARAPACE_TEST=1",
	}))

	s, err := Gen(cmd).Snippet("bash")
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(s, "declare -x CARAPACE_TEST=1") {
		t.Error("bash snippet missing enrichment code")
	}
	if !strings.Contains(s, "#!/bin/bash") {
		t.Error("bash snippet missing shebang")
	}

	// zsh snippet should not contain bash enrichment
	s, err = Gen(cmd).Snippet("zsh")
	if err != nil {
		t.Error(err)
	}
	if strings.Contains(s, "CARAPACE_TEST=1") {
		t.Error("zsh snippet should not contain bash enrichment")
	}
}

func TestGenWithDefault(t *testing.T) {
	sub1 := &cobra.Command{Use: "sub1"}
	sub2 := &cobra.Command{Use: "sub2"}
	root := &cobra.Command{Use: "root"}

	Gen(root, WithSubcommands(sub1, sub2), WithDefault("sub2"))

	entry := storage.get(root)
	if entry.defaultName != "sub2" {
		t.Errorf("expected defaultName 'sub2', got '%v'", entry.defaultName)
	}
}

func TestGenWithSubcommandsDefaultName(t *testing.T) {
	sub1 := &cobra.Command{Use: "sub1"}
	sub2 := &cobra.Command{Use: "sub2"}
	root := &cobra.Command{Use: "root"}

	Gen(root, WithSubcommands(sub1, sub2))

	entry := storage.get(root)
	exe := uid.Executable()
	if entry.defaultName != exe {
		t.Errorf("expected defaultName '%v' (executable), got '%v'", exe, entry.defaultName)
	}
}

func TestGenWithDefaultNoop(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	Gen(cmd, WithDefault("sub1"))
	// WithDefault is no-op at runtime without WithSubcommands:
	// Execute() does nothing special, Snippet() produces standard snippet
	s, err := Gen(cmd).Snippet("bash")
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(s, "#!/bin/bash") {
		t.Error("should still produce standard bash snippet")
	}
}
