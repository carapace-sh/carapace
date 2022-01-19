package carapace

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestActionFlags(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().BoolP("alpha", "a", false, "")
	cmd.Flags().BoolP("beta", "b", false, "")

	cmd.Flag("alpha").Changed = true
	a := actionFlags(cmd).Invoke(Context{CallbackValue: "-a"})
	assertEqual(t, ActionValuesDescribed("b", "").NoSpace().Invoke(Context{}).Prefix("-a"), a)
}
