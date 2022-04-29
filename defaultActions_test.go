package carapace

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestActionImport(t *testing.T) {
	s := `
{
  "Version": "unknown",
  "Nospace": true,
  "RawValues": [
    {
      "Value": "positional1",
      "Display": "positional1",
      "Description": "",
      "Style": ""
    },
    {
      "Value": "p1",
      "Display": "p1",
      "Description": "",
      "Style": ""
    }
  ]
}`
	assertEqual(t, ActionValues("positional1", "p1").NoSpace().Invoke(Context{}), ActionImport([]byte(s)).Invoke(Context{}))
}

func TestActionFlags(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().BoolP("alpha", "a", false, "")
	cmd.Flags().BoolP("beta", "b", false, "")

	cmd.Flag("alpha").Changed = true
	a := actionFlags(cmd).Invoke(Context{CallbackValue: "-a"})
	assertEqual(t, ActionValuesDescribed("b", "").NoSpace().Invoke(Context{}).Prefix("-a"), a)
}

func TestActionExecCommandEnv(t *testing.T) {
	ActionExecCommand("env")(func(output []byte) Action {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "carapace_TestActionExecCommand") {
				t.Error("should not contain env carapace_TestActionExecCommand")
			}
		}
		return ActionValues()
	}).Invoke(Context{})

	c := Context{}
	c.Setenv("carapace_TestActionExecCommand", "test")
	ActionExecCommand("env")(func(output []byte) Action {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if line == "carapace_TestActionExecCommand=test" {
				return ActionValues()
			}
		}
		t.Error("should contain env carapace_TestActionExecCommand=test")
		return ActionValues()
	}).Invoke(c)
}

func TestActionExecute(t *testing.T) {
	root := &cobra.Command{Use: "pass"}
	rootGitCmd := &cobra.Command{Use: "git"}
	rootShowCmd := &cobra.Command{Use: "show"}

	root.AddCommand(rootGitCmd)
	root.AddCommand(rootShowCmd)

	git := &cobra.Command{Use: "git"}
	gitShowCmd := &cobra.Command{Use: "show"}
	gitStashCmd := &cobra.Command{Use: "stash"}
	gitStashShowCmd := &cobra.Command{Use: "show"}
	git.AddCommand(gitShowCmd)
	git.AddCommand(gitStashCmd)
	gitStashCmd.AddCommand(gitStashShowCmd)

	Gen(root)

	Gen(rootShowCmd).PositionalCompletion(
		ActionValues("rootShowCmd"),
	)

	Gen(rootGitCmd).PositionalAnyCompletion(
		ActionExecute(git).Chdir("/tmp"),
	)

	Gen(gitShowCmd).PositionalCompletion(
		ActionValues("gitShowCmd"),
	)
	
    Gen(gitStashShowCmd).PositionalCompletion(
		ActionValues("gitStashShowCmd"),
	)

	var stdout bytes.Buffer
	root.SetOut(&stdout)
	os.Args = []string{"pass", "_carapace", "export", "pass", "git", "show", ""}
	root.Execute()

	if !strings.Contains(stdout.String(), "gitShowCmd") {
		t.Error("should be gitShowCmd")
	}
}
