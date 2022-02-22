package carapace

import (
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
      "Description": ""
    },
    {
      "Value": "p1",
      "Display": "p1",
      "Description": ""
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
