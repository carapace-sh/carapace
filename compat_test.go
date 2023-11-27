package carapace

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRegisterValidArgsFunction(t *testing.T) {
	cmd := &cobra.Command{}

	Gen(cmd).PositionalCompletion(
		ActionValues("1"),
		ActionValuesDescribed("2", "second"),
	)

	Gen(cmd).PositionalAnyCompletion(
		ActionValues("any"),
	)

	registerValidArgsFunction(cmd)

	if vals, _ := cmd.ValidArgsFunction(cmd, []string{}, ""); vals[0] != "1" {
		t.Error("first position wrong")
	}

	if vals, _ := cmd.ValidArgsFunction(cmd, []string{""}, ""); vals[0] != "2\tsecond" {
		t.Error("second position wrong")
	}

	if vals, _ := cmd.ValidArgsFunction(cmd, []string{"", ""}, ""); vals[0] != "any" {
		t.Error("third position wrong")
	}
}

func TestRegisterFlagCompletion(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("flag", "", "")

	Gen(cmd).FlagCompletion(ActionMap{
		"flag": ActionValuesDescribed("1", "one"),
	})

	registerFlagCompletion(cmd)

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	os.Args = []string{"", "__complete", "--flag", ""}
	_ = cmd.Execute()

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout

	if lines := strings.Split(string(out), "\n"); lines[0] != "1\tone" {
		t.Error("flag wrong")
	}
}
