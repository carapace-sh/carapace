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

func TestExecutable(t *testing.T) {
	if !strings.HasSuffix(Executable(), ".test") {
		t.Fail()
	}
}
