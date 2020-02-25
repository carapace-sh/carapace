package zsh

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestUidCommand(t *testing.T) {
	root := &cobra.Command{
		Use: "root",
	}
	sub1 := &cobra.Command{
		Use: "sub1",
	}
	sub2 := &cobra.Command{
		Use: "sub2",
	}

	root.AddCommand(sub1)
	sub1.AddCommand(sub2)

	assertEqual(t, "_root", uidCommand(root))
	assertEqual(t, "_root__sub1", uidCommand(sub1))
	assertEqual(t, "_root__sub1__sub2", uidCommand(sub2))
}

func TestUidFlag(t *testing.T) {
	root := &cobra.Command{
		Use: "root",
	}
	sub1 := &cobra.Command{
		Use: "sub1",
	}
	root.Flags().Bool("rootFlag", false, "test")
	sub1.Flags().Bool("sub1Flag", false, "test")

	root.AddCommand(sub1)

	assertEqual(t, "_root##rootFlag", uidFlag(root, root.Flag("rootFlag")))
	assertEqual(t, "_root__sub1##sub1Flag", uidFlag(sub1, sub1.Flag("sub1Flag")))
}

func TestUidPositional(t *testing.T) {
	root := &cobra.Command{
		Use: "root",
	}
	sub1 := &cobra.Command{
		Use: "sub1",
	}
	root.AddCommand(sub1)

	assertEqual(t, "_root#1", uidPositional(root, 1))
	assertEqual(t, "_root__sub1#3", uidPositional(sub1, 3))
}

func TestFind(t *testing.T) {
	root := &cobra.Command{
		Use: "root",
	}
	sub1 := &cobra.Command{
		Use: "sub1",
	}
	sub2 := &cobra.Command{
		Use: "sub2",
	}

	root.AddCommand(sub1)
	sub1.AddCommand(sub2)

	if sub2 != find(sub1, "_root__sub1__sub2#5") {
		t.Error(`"_root__sub1__sub2#5" should be sub2`)
	}

	if root != find(sub2, "_root##someflag") {
		t.Error(`"_root##someflag" should be root`)
	}

	if sub1 != find(root, "_root__sub1##someflag") {
		t.Error(`"_root__sub1##someflag" should be sub1`)
	}
}
