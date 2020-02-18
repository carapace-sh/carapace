package zsh

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestUid(t *testing.T) {
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

  t.Log(uidCommand(sub2))
  t.Log(find(sub1, uidCommand(sub2)).Name())
}
