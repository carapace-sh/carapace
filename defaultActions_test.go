package carapace

import (
	"strings"
	"testing"

	"github.com/carapace-sh/carapace/internal/common"
	"github.com/carapace-sh/carapace/pkg/assert"
	"github.com/spf13/cobra"
)

func TestActionImport(t *testing.T) {
	s := `
{
  "version": "unknown",
  "nospace": "",
  "values": [
    {
      "value": "positional1",
      "display": "positional1",
      "description": "",
      "style": "",
	  "tag": "first"
    },
    {
      "value": "p1",
      "display": "p1",
      "description": "",
      "style": "",
	  "tag": "first"
    }
  ]
}`
	assert.Equal(t, ActionValues("positional1", "p1").Tag("first").Invoke(Context{}), ActionImport([]byte(s)).Invoke(Context{}))
}

func TestActionFlags(t *testing.T) {
	cmd := &cobra.Command{Use: "actionFlags"}
	cmd.Flags().BoolP("alpha", "a", false, "")
	cmd.Flags().BoolP("beta", "b", false, "")

	cmd.Flag("alpha").Changed = true
	a := actionFlags(cmd).Invoke(Context{Value: "-a"})
	expected := Action{rawValues: common.RawValues{
		{Value: "-a", Display: "-a", Tag: "shorthand flags", Uid: "cmd://actionFlags?flag=alpha"},
		{Value: "-ab", Display: "b", Tag: "shorthand flags", Uid: "cmd://actionFlags?flag=beta"},
		{Value: "-ah", Display: "h", Description: "help for actionFlags", Tag: "shorthand flags", Uid: "cmd://actionFlags?flag=help"},
	}}
	expected.meta.Nospace.Add('b', 'h')
	assert.Equal(
		t,
		expected.Invoke(Context{}),
		a,
	)
}

func TestActionExecCommandEnv(t *testing.T) {
	ActionExecCommand("env")(func(output []byte) Action {
		lines := strings.SplitSeq(string(output), "\n")
		for line := range lines {
			if strings.Contains(line, "carapace_TestActionExecCommand") {
				t.Error("should not contain env carapace_TestActionExecCommand")
			}
		}
		return ActionValues()
	}).Invoke(Context{})

	c := Context{}
	c.Setenv("carapace_TestActionExecCommand", "test")
	ActionExecCommand("env")(func(output []byte) Action {
		lines := strings.SplitSeq(string(output), "\n")
		for line := range lines {
			if line == "carapace_TestActionExecCommand=test" {
				return ActionValues()
			}
		}
		t.Error("should contain env carapace_TestActionExecCommand=test")
		return ActionValues()
	}).Invoke(c)
}
