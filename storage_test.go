package carapace

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestGetFlag(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.PersistentFlags().String("flag", "", "")
	subcmd := &cobra.Command{}
	cmd.AddCommand(subcmd)

	Gen(cmd).FlagCompletion(ActionMap{
		"flag": ActionValues("a", "b"),
	})

	assertEqual(t, ActionValues("a", "b").Invoke(Context{}), storage.getFlag(subcmd, "flag").Invoke(Context{}))
}

func TestGetPositional(t *testing.T) {
	cmd := &cobra.Command{}

	Gen(cmd).PositionalCompletion(
		ActionValues("pos", "1"),
		ActionValues("pos", "2"),
	)

	Gen(cmd).PositionalAnyCompletion(
		ActionValues("pos", "any"),
	)

	assertEqual(t, ActionValues("pos", "1").Invoke(Context{}), storage.getPositional(cmd, 0).Invoke(Context{}))
	assertEqual(t, ActionValues("pos", "2").Invoke(Context{}), storage.getPositional(cmd, 1).Invoke(Context{}))
	assertEqual(t, ActionValues("pos", "any").Invoke(Context{}), storage.getPositional(cmd, 2).Invoke(Context{}))
}

func TestGetDash(t *testing.T) {
	cmd := &cobra.Command{}

	Gen(cmd).DashCompletion(
		ActionValues("dash", "1"),
		ActionValues("dash", "2"),
	)

	Gen(cmd).DashAnyCompletion(
		ActionValues("dash", "any"),
	)

	_ = cmd.Flags().Parse([]string{"--", ""})

	assertEqual(t, ActionValues("dash", "1").Invoke(Context{}), storage.getPositional(cmd, 0).Invoke(Context{}))
	assertEqual(t, ActionValues("dash", "2").Invoke(Context{}), storage.getPositional(cmd, 1).Invoke(Context{}))
	assertEqual(t, ActionValues("dash", "any").Invoke(Context{}), storage.getPositional(cmd, 2).Invoke(Context{}))
}

func TestCheck(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("flag", "", "")

	Gen(cmd).FlagCompletion(ActionMap{
		"flag": ActionValues("a", "b"),
	})

	if len(storage.check()) != 0 {
		t.Error("check should succeed")
	}

	Gen(cmd).FlagCompletion(ActionMap{
		"unknown-flag": ActionValues("a", "b"),
	})

	if len(storage.check()) != 1 {
		t.Error("check should fail")
	}
}

// BenchmarkStorage tests for concurrent map read/write.
func BenchmarkStorage(b *testing.B) {
	cmd := &cobra.Command{}
	cmd2 := &cobra.Command{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Gen(cmd).FlagCompletion(ActionMap{
				"flag1": ActionValues("a", "b"),
			})
			Gen(cmd).PositionalCompletion(ActionValues("a", "b"))

			Gen(cmd2).FlagCompletion(ActionMap{
				"flag2": ActionValues("a", "b"),
			})
			Gen(cmd2).PositionalCompletion(ActionValues("a", "b"))

			storage.getFlag(cmd, "flag1")
			storage.getPositional(cmd, 0)
			storage.getFlag(cmd2, "flag2")
			storage.getPositional(cmd2, 0)
		}
	})
}
