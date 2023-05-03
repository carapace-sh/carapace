package carapace

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/rsteube/carapace/internal/assert"
	"github.com/rsteube/carapace/internal/uid"
	"github.com/spf13/cobra"
)

func init() {
	os.Unsetenv("LS_COLORS")
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
		Dir:   wd(""),
	},
		"")

	testContext(t, Context{
		Value: "",
		Args:  []string{"pos1"},
		Parts: []string{},
		Env:   []string{},
		Dir:   wd(""),
	},
		"pos1", "")

	testContext(t, Context{
		Value: "po",
		Args:  []string{"pos1", "pos2"},
		Parts: []string{},
		Env:   []string{},
		Dir:   wd(""),
	},
		"pos1", "pos2", "po")

	testContext(t, Context{
		Value: "",
		Args:  []string{},
		Parts: []string{},
		Env:   []string{},
		Dir:   wd(""),
	},
		"--multiparts", "")

	testContext(t, Context{
		Value: "fir",
		Args:  []string{},
		Parts: []string{},
		Env:   []string{},
		Dir:   wd(""),
	},
		"--multiparts", "fir")

	testContext(t, Context{
		Value: "seco",
		Args:  []string{"pos1"},
		Parts: []string{"first"},
		Env:   []string{},
		Dir:   wd(""),
	},
		"pos1", "--multiparts", "first,seco")

	testContext(t, Context{
		Value: "pos",
		Args:  []string{},
		Parts: []string{},
		Env:   []string{},
		Dir:   wd(""),
	},
		"pos")

	testContext(t, Context{
		Value: "sec",
		Args:  []string{},
		Parts: []string{"first"},
		Env:   []string{},
		Dir:   wd(""),
	},
		"first:sec")

	testContext(t, Context{
		Value: "thi",
		Args:  []string{"first:second"},
		Parts: []string{},
		Env:   []string{},
		Dir:   wd(""),
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
	cmd := &cobra.Command{}
	if s, _ := Gen(cmd).Snippet("bash"); !strings.Contains(s, "#!/bin/bash") {
		t.Error("bash failed")
	}

	if s, _ := Gen(cmd).Snippet("elvish"); !strings.Contains(s, "edit:completion") {
		t.Error("elvish failed")
	}

	if s, _ := Gen(cmd).Snippet("fish"); !strings.Contains(s, "commandline") {
		t.Error("fish failed")
	}

	if s, _ := Gen(cmd).Snippet("oil"); !strings.Contains(s, "#!/bin/osh") {
		t.Error("oil failed")
	}

	if s, _ := Gen(cmd).Snippet("powershell"); !strings.Contains(s, "System.Management.Automation") {
		t.Error("powershell failed")
	}

	if s, _ := Gen(cmd).Snippet("xonsh"); !strings.Contains(s, "@contextual_command_completer") {
		t.Error("xonsh failed")
	}

	if s, _ := Gen(cmd).Snippet("zsh"); !strings.Contains(s, "compdef") {
		t.Error("zsh")
	}

	if _, err := Gen(cmd).Snippet("unknown"); err == nil {
		t.Error("zsh")
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
