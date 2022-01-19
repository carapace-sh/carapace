package carapace

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/rsteube/carapace/internal/assert"
	"github.com/rsteube/carapace/internal/uid"
	"github.com/spf13/cobra"
)

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

	os.Args = append([]string{"root", "_carapace", "elvish", "_", "root"}, args...)
	rootCmd.Execute()
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
		os.Stdout = sOut
		os.Stderr = sErr

		e, _ := json.Marshal(expected)
		a, _ := json.Marshal(actual)
		assert.Equal(t, string(e), string(a))
	})
}

func TestContext(t *testing.T) {
	testContext(t, Context{
		CallbackValue: "",
		Args:          []string{},
		Parts:         []string{},
	},
		"")

	testContext(t, Context{
		CallbackValue: "",
		Args:          []string{"pos1"},
		Parts:         []string{},
	},
		"pos1", "")

	testContext(t, Context{
		CallbackValue: "po",
		Args:          []string{"pos1", "pos2"},
		Parts:         []string{},
	},
		"pos1", "pos2", "po")

	testContext(t, Context{
		CallbackValue: "",
		Args:          []string{},
		Parts:         []string{},
	},
		"--multiparts", "")

	testContext(t, Context{
		CallbackValue: "fir",
		Args:          []string{},
		Parts:         []string{},
	},
		"--multiparts", "fir")

	testContext(t, Context{
		CallbackValue: "seco",
		Args:          []string{"pos1"},
		Parts:         []string{"first"},
	},
		"pos1", "--multiparts", "first,seco")

	testContext(t, Context{
		CallbackValue: "pos",
		Args:          []string{},
		Parts:         []string{},
	},
		"pos")

	testContext(t, Context{
		CallbackValue: "sec",
		Args:          []string{},
		Parts:         []string{"first"},
	},
		"first:sec")

	testContext(t, Context{
		CallbackValue: "thi",
		Args:          []string{"first:second"},
		Parts:         []string{},
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

func TestInitLogger(t *testing.T) {
	initLogger()
	tmpdir := fmt.Sprintf("%v/carapace", os.TempDir())
	if _, err := os.Stat(fmt.Sprintf("%v/%v.log", tmpdir, uid.Executable())); os.IsNotExist(err) {

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
}
