package uid

import (
	"strings"
	"testing"

	"github.com/rsteube/carapace/internal/assert"
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

	assert.Equal(t, "_root", Command(root))
	assert.Equal(t, "_root__sub1", Command(sub1))
	assert.Equal(t, "_root__sub1__sub2", Command(sub2))
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

	assert.Equal(t, "_root##rootFlag", Flag(root, root.Flag("rootFlag")))
	assert.Equal(t, "_root__sub1##sub1Flag", Flag(sub1, sub1.Flag("sub1Flag")))
}

func TestUidPositional(t *testing.T) {
	root := &cobra.Command{
		Use: "root",
	}
	sub1 := &cobra.Command{
		Use: "sub1",
	}
	root.AddCommand(sub1)

	assert.Equal(t, "_root#1", Positional(root, 1))
	assert.Equal(t, "_root__sub1#3", Positional(sub1, 3))
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

func TestExecutable(t *testing.T) {
	if !strings.HasSuffix(Executable(), ".test") {
		t.Fail()
	}
}
