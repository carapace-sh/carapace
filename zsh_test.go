package zsh

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/alecthomas/chroma/quick"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
)

func highlight(s string) string {
	buf := bytes.NewBufferString("")
	if err := quick.Highlight(buf, s, "bash", "terminal256", "monokai"); err != nil {
		log.Fatal(err)
	}
	return buf.String()
}

func assertEqual(t *testing.T, expected string, actual string) {
	if expected == actual {
		t.Log(highlight(actual))
	} else {
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(expected, actual, false)
		t.Errorf("\nexpected: %v\nactual  : %v", expected, dmp.DiffPrettyText(diffs))
	}
}

func TestZshComp(t *testing.T) {
	root := cobra.Command{
		Use: "test",
	}
	sub1 := cobra.Command{
		Use: "sub1",
	}
	sub2 := cobra.Command{
		Use: "sub2",
	}

	sub2.Flags()

	root.AddCommand(&sub1)
	sub1.AddCommand(&sub2)

	root.PersistentFlags().String("license", "", "name of license for the project")
	sub1.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")

	// TODO stuff
	Gen(&sub2)

	os.Args = []string{"test", "_zsh_completion"}
	//root.Execute() // TODO

	action := ActionCallback(func(args []string) Action {
		return ActionValues("A", "B")
	})
	uid := uidFlag(&sub1, sub1.Flag("author"))
	completions.actions[uid] = action
	t.Log("\n" + highlight(completions.Generate(&root)))
}

func TestSubcommandsSnippet(t *testing.T) {
	root := &cobra.Command{
		Use: "test",
	}
	sub1 := &cobra.Command{
		Use: "sub1",
	}
	sub2 := &cobra.Command{
		Use:   "sub2",
		Short: "some description",
	}

	sub2.Flags().String("file", "f", "file flag")
	sub2.Flags().String("user", "u", "user flag")
	sub2.Flags().String("conditional", "c", "conditional flag")

	Gen(sub2).FlagCompletion(ActionMap{
		"file": ActionFiles("*.go"),
		"user": ActionUsers(),
		"conditional": ActionCallback(func(args []string) Action {
			if sub2.Flag("user").Value.String() == "bob" {
				return ActionValues("bob1", "bob2", "bob3")
			} else {
				return ActionGroups()
			}
		})})

	root.AddCommand(sub1)
	root.AddCommand(sub2)

	t.Log(highlight(subcommands(root)))
	t.Log(highlight(completions.Generate(root)))

	t.Log(completions.actions)
}

func TestTraverse(t *testing.T) {
	root := &cobra.Command{
		Use: "test",
	}
	sub1 := &cobra.Command{
		Use: "sub1",
	}
	sub2 := &cobra.Command{
		Use: "sub2",
	}

	sub2.Flags()

	root.AddCommand(sub1)
	sub1.AddCommand(sub2)

	root.PersistentFlags().String("license", "", "name of license for the project")
	sub2.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")

	args := []string{"--license", "MIT", "sub1", "sub2", "positionalARG", "--author", "bob"}
	targetArgs := traverse(sub1, args)

	if len(targetArgs) != 1 || targetArgs[0] != "positionalARG" {
		t.Error("traverse should return positionalARG")
	}
	if sub2.Flag("license").Value.String() != "MIT" {
		t.Error("flag license should be MIT")
	}
	if sub2.Flag("author").Value.String() != "bob" {
		t.Error("flag author should be bob")
	}
}
