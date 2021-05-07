package carapace

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/rsteube/carapace/internal/assert"
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
